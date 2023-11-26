package compiler

import "github.com/moonbite-org/moonbite/serde"

var InstructionMap = map[string]byte{
	"define_module": 10,
	"save":          20,
	"load_pointer":  21,
	"load_value":    22,
	"load_param":    23,
	"call":          24,
}

func define_module(name string, is_entry bool) []byte {
	result := []byte{InstructionMap["define_module"]}
	result = append(result, serde.SerializeBool(is_entry)...)
	result = append(result, serde.SerializeString(name)...)

	return result
}

func save(pointer int, value interface{}) ([]byte, error) {
	result := []byte{InstructionMap["save"]}
	p, err := serde.SerializeValue(pointer)

	if err != nil {
		return result, err
	}

	result = append(result, p...)

	v, err := serde.SerializeValue(value)

	if err != nil {
		return result, err
	}

	result = append(result, v...)

	return result, nil
}

func call(pointer int, args []int) ([]byte, error) {
	result := []byte{InstructionMap["call"]}
	p, err := serde.SerializeValue(pointer)

	if err != nil {
		return result, err
	}

	result = append(result, p...)

	arg_list := []byte{}
	for _, arg := range args {
		val, err := serde.SerializeValue(arg)

		if err != nil {
			return result, err
		}

		arg_list = append(arg_list, val...)
	}

	length, err := serde.SerializeValue(len(args))

	if err != nil {
		return result, err
	}

	result = append(result, length...)
	result = append(result, arg_list...)

	return result, nil
}

func load_pointer(pointer int) ([]byte, error) {
	result := []byte{InstructionMap["load_pointer"]}

	p, err := serde.SerializeValue(pointer)

	if err != nil {
		return result, err
	}

	result = append(result, p...)

	return result, nil
}

func load_value(value interface{}) ([]byte, error) {
	result := []byte{InstructionMap["load_value"]}

	v, err := serde.SerializeValue(value)

	if err != nil {
		return result, err
	}

	result = append(result, v...)

	return result, nil
}

func load_param(index int) ([]byte, error) {
	result := []byte{InstructionMap["load_index"]}

	i, err := serde.SerializeValue(index)

	if err != nil {
		return result, err
	}

	result = append(result, i...)

	return result, nil
}
