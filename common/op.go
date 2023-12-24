package common

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

type Op byte

const (
	OpNoop Op = iota + 10
	OpConstant
	OpSet
	OpGet
	OpCall
	OpPop
	OpReturn
	OpTrue
	OpFalse
	OpJump
	OpJumpIfFalse
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNegate
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEqual
)

var op_map = map[Op]string{
	OpNoop:               "Noop",
	OpConstant:           "Constant",
	OpSet:                "Set",
	OpGet:                "Get",
	OpCall:               "Call",
	OpPop:                "Pop",
	OpReturn:             "Return",
	OpTrue:               "True",
	OpFalse:              "False",
	OpJump:               "Jump",
	OpJumpIfFalse:        "JumpIfFalse",
	OpAdd:                "Add",
	OpSub:                "Sub",
	OpMul:                "Mul",
	OpDiv:                "Div",
	OpMod:                "Mod",
	OpNegate:             "Negate",
	OpEqual:              "Equal",
	OpNotEqual:           "NotEqual",
	OpGreaterThan:        "GreaterThan",
	OpGreaterThanOrEqual: "GreaterThanOrEqual",
}

func (op Op) String() string {
	if len(op_map[op]) > 1 {
		return op_map[op]
	}

	return fmt.Sprintf("%d", op)
}

func IntToBytes[T uint32 | uint64](value T) []byte {
	var result []byte

	switch reflect.TypeOf(value).Kind() {
	case reflect.Uint32:
		result = make([]byte, 4)
		binary.LittleEndian.PutUint32(result, uint32(value))
	case reflect.Uint64:
		result = make([]byte, 4)
		binary.LittleEndian.PutUint64(result, uint64(value))
	}

	return result
}

func BytesToInt[T uint32 | uint64](data []byte) T {
	switch len(data) {
	case 8:
		return T(binary.LittleEndian.Uint64(data))
	case 4:
		return T(binary.LittleEndian.Uint32(data))
	default:
		panic("unsupported int size")
	}
}

func MakeOp(op Op, operands ...uint32) []byte {
	result := []byte{byte(op)}

	for _, operand := range operands {
		result = append(result, IntToBytes(operand)...)
	}

	return result
}
