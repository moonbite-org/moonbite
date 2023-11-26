package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/moonbite-org/moonbite/common"
)

func SerializeValue(value interface{}) ([]byte, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return SerializeString(value.(string)), nil
	case reflect.Bool:
		return SerializeBool(value.(bool)), nil
	case reflect.Array, reflect.Slice:
		return SerializeList(to_interface_slice(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return SerializeNumber(value)
	case reflect.Map:
		return SerializeStruct(value.(map[string]interface{}))
	default:
		return []byte{}, fmt.Errorf("unsupported type")
	}
}

func SerializeBool(value bool) []byte {
	if value {
		return []byte{common.TypeMap["true"]}
	}

	return []byte{common.TypeMap["false"]}
}

func SerializeString(value string) []byte {
	result := []byte{common.TypeMap["string"]}
	result = append(result, []byte(value)...)
	result = append(result, common.TypeMap["terminator"])
	return result
}

func SerializeRune(value rune) []byte {
	result := []byte{common.TypeMap["rune"]}
	result = append(result, []byte(string(value))...)
	result = append(result, common.TypeMap["terminator"])
	return result
}

func SerializeNumber(value any) ([]byte, error) {
	result := []byte{}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int8:
		result = append(result, common.TypeMap["int8"])
	case reflect.Int16:
		result = append(result, common.TypeMap["int16"])
	case reflect.Int:
		value = int32(value.(int))
		result = append(result, common.TypeMap["int32"])
	case reflect.Int32:
		result = append(result, common.TypeMap["int32"])
	case reflect.Int64:
		result = append(result, common.TypeMap["int64"])
	case reflect.Uint8:
		result = append(result, common.TypeMap["uint8"])
	case reflect.Uint16:
		result = append(result, common.TypeMap["uint16"])
	case reflect.Uint32:
		result = append(result, common.TypeMap["uint32"])
	case reflect.Uint64:
		result = append(result, common.TypeMap["uint64"])
	case reflect.Float32:
		result = append(result, common.TypeMap["float32"])
	case reflect.Float64:
		result = append(result, common.TypeMap["float64"])
	}

	buffer := bytes.Buffer{}
	err := binary.Write(&buffer, binary.BigEndian, value)

	if err != nil {
		return []byte{}, err
	}

	result = append(result, buffer.Bytes()...)

	return result, nil
}

func SerializeList(value []interface{}) ([]byte, error) {
	result := []byte{common.TypeMap["list"]}

	for _, v := range value {
		val, err := SerializeValue(v)

		if err != nil {
			return []byte{common.TypeMap["list"], common.TypeMap["terminator"]}, err
		}

		result = append(result, val...)
	}

	result = append(result, common.TypeMap["terminator"])
	return result, nil
}

func SerializeStruct(value map[string]interface{}) ([]byte, error) {
	result := []byte{common.TypeMap["struct"]}

	for key, v := range value {
		result = append(result, SerializeString(key)...)
		val, err := SerializeValue(v)

		if err != nil {
			return []byte{common.TypeMap["struct"], common.TypeMap["terminator"]}, err
		}

		result = append(result, val...)
	}

	result = append(result, common.TypeMap["terminator"])
	return result, nil
}

func SerializeFunction() ([]byte, error) {
	result := []byte{common.TypeMap["fun"], 0, common.TypeMap["terminator"]}

	return result, nil
}

func to_interface_slice(values interface{}) []interface{} {
	result := []interface{}{}

	v := reflect.ValueOf(values)

	for i := 0; i < v.Len(); i++ {
		result = append(result, v.Index(i).Interface())
	}

	return result
}
