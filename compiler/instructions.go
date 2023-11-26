package compiler

import "github.com/moonbite-org/moonbite/common"

func define_module(name string, is_entry bool) []byte {
	result := []byte{common.InstructionMap["define_module"]}
	result = append(result, SerializeBool(is_entry)...)
	result = append(result, SerializeString(name)...)

	return result
}

func save(pointer int, value interface{}) ([]byte, error) {
	result := []byte{common.InstructionMap["save"]}
	p, err := SerializeValue(pointer)

	if err != nil {
		return result, err
	}

	result = append(result, p...)

	v, err := SerializeValue(value)

	if err != nil {
		return result, err
	}

	result = append(result, v...)

	return result, nil
}

func call(pointer int, args []int) ([]byte, error) {
	result := []byte{common.InstructionMap["call"]}
	p, err := SerializeValue(pointer)

	if err != nil {
		return result, err
	}

	result = append(result, p...)

	arg_list := []byte{}
	for _, arg := range args {
		val, err := SerializeValue(arg)

		if err != nil {
			return result, err
		}

		arg_list = append(arg_list, val...)
	}

	length, err := SerializeValue(len(args))

	if err != nil {
		return result, err
	}

	result = append(result, length...)
	result = append(result, arg_list...)

	return result, nil
}

func load_pointer(pointer int) ([]byte, error) {
	result := []byte{common.InstructionMap["load_pointer"]}

	p, err := SerializeValue(pointer)

	if err != nil {
		return result, err
	}

	result = append(result, p...)

	return result, nil
}

func load_value(value interface{}) ([]byte, error) {
	result := []byte{common.InstructionMap["load_value"]}

	v, err := SerializeValue(value)

	if err != nil {
		return result, err
	}

	result = append(result, v...)

	return result, nil
}

func load_param(index int) ([]byte, error) {
	result := []byte{common.InstructionMap["load_index"]}

	i, err := SerializeValue(index)

	if err != nil {
		return result, err
	}

	result = append(result, i...)

	return result, nil
}
