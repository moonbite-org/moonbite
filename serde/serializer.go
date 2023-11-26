package serde

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

var TypeMap = map[string]byte{
	"false":      10,
	"true":       11,
	"string":     12,
	"rune":       13,
	"int8":       14,
	"int16":      15,
	"int32":      16,
	"int64":      17,
	"uint8":      18,
	"uint16":     19,
	"uint32":     20,
	"uint64":     21,
	"float32":    22,
	"float64":    23,
	"list":       24,
	"struct":     25,
	"fun":        26,
	"terminator": 0,
}

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
	}

	return []byte{}, nil
}

func SerializeBool(value bool) []byte {
	if value {
		return []byte{TypeMap["true"]}
	}

	return []byte{TypeMap["false"]}
}

func SerializeString(value string) []byte {
	result := []byte{TypeMap["string"]}
	result = append(result, []byte(value)...)
	result = append(result, TypeMap["terminator"])
	return result
}

func SerializeRune(value rune) []byte {
	result := []byte{TypeMap["rune"]}
	result = append(result, []byte(string(value))...)
	result = append(result, TypeMap["terminator"])
	return result
}

func SerializeNumber(value any) ([]byte, error) {
	result := []byte{}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Int8:
		result = append(result, TypeMap["int8"])
	case reflect.Int16:
		result = append(result, TypeMap["int16"])
	case reflect.Int:
		value = int32(value.(int))
		result = append(result, TypeMap["int32"])
	case reflect.Int32:
		result = append(result, TypeMap["int32"])
	case reflect.Int64:
		result = append(result, TypeMap["int64"])
	case reflect.Uint8:
		result = append(result, TypeMap["uint8"])
	case reflect.Uint16:
		result = append(result, TypeMap["uint16"])
	case reflect.Uint32:
		result = append(result, TypeMap["uint32"])
	case reflect.Uint64:
		result = append(result, TypeMap["uint64"])
	case reflect.Float32:
		result = append(result, TypeMap["float32"])
	case reflect.Float64:
		result = append(result, TypeMap["float64"])
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
	result := []byte{TypeMap["list"]}

	for _, v := range value {
		val, err := SerializeValue(v)

		if err != nil {
			return []byte{TypeMap["list"], TypeMap["terminator"]}, err
		}

		result = append(result, val...)
	}

	result = append(result, TypeMap["terminator"])
	return result, nil
}

func SerializeStruct(value map[string]interface{}) ([]byte, error) {
	result := []byte{TypeMap["struct"]}

	for key, v := range value {
		result = append(result, SerializeString(key)...)
		val, err := SerializeValue(v)

		if err != nil {
			return []byte{TypeMap["struct"], TypeMap["terminator"]}, err
		}

		result = append(result, val...)
	}

	result = append(result, TypeMap["terminator"])
	return result, nil
}

func SerializeFunction() ([]byte, error) {
	result := []byte{TypeMap["fun"], 0, TypeMap["terminator"]}

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
