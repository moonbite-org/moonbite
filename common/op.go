package common

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
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
	OpArray
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
	OpArray:              "Array",
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

type Instruction struct {
	Op       Op
	Operands []uint32
}

func (i Instruction) GetBytes() []byte {
	result := []byte{}
	result = append(result, byte(i.Op))

	for _, operand := range i.Operands {
		result = append(result, IntToBytes(operand)...)
	}

	return result
}

func (i Instruction) String() string {
	operands := []string{}

	for _, operand := range i.Operands {
		operands = append(operands, fmt.Sprintf("%d", operand))
	}

	return fmt.Sprintf("%s(%s)", op_map[i.Op], strings.Join(operands, " "))
}

func NewInstruction(op Op, operands ...int) Instruction {
	u_operands := []uint32{}

	for _, operand := range operands {
		u_operands = append(u_operands, uint32(operand))
	}

	return Instruction{
		Op:       op,
		Operands: u_operands,
	}
}

type InstructionSet []Instruction

func (s InstructionSet) GetBytes() []byte {
	result := []byte{}

	for _, instruction := range s {
		result = append(result, instruction.GetBytes()...)
	}

	return result
}
