package common

import (
	"bytes"
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
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpAssign
	OpAssignLocal
	OpSetItem
	OpCall
	OpPop
	OpReturn
	OpReturnEmpty
	OpBreak
	OpBreakEmpty
	OpDefer
	OpContinue
	OpYield
	OpIndex
	OpTrue
	OpFalse
	OpJump
	OpJumpIfFalse
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpAnd
	OpOr
	OpArray
	OpMap
	OpNegate
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEqual
	OpInstanceof
	OpCast
	OpExit
)

var op_map = map[Op]string{
	OpNoop:               "Noop",
	OpConstant:           "Constant",
	OpSet:                "Set",
	OpGet:                "Get",
	OpSetLocal:           "SetLocal",
	OpGetLocal:           "GetLocal",
	OpGetBuiltin:         "GetBuiltin",
	OpAssign:             "Assign",
	OpAssignLocal:        "AssignLocal",
	OpSetItem:            "SetItem",
	OpCall:               "Call",
	OpPop:                "Pop",
	OpReturn:             "Return",
	OpReturnEmpty:        "ReturnEmpty",
	OpBreak:              "Break",
	OpDefer:              "Defer",
	OpBreakEmpty:         "BreakEmpty",
	OpContinue:           "Continue",
	OpYield:              "Yield",
	OpIndex:              "Index",
	OpTrue:               "True",
	OpFalse:              "False",
	OpJump:               "Jump",
	OpJumpIfFalse:        "JumpIfFalse",
	OpAdd:                "Add",
	OpSub:                "Sub",
	OpMul:                "Mul",
	OpDiv:                "Div",
	OpMod:                "Mod",
	OpAnd:                "And",
	OpOr:                 "Or",
	OpArray:              "Array",
	OpMap:                "Map",
	OpNegate:             "Negate",
	OpEqual:              "Equal",
	OpNotEqual:           "NotEqual",
	OpGreaterThan:        "GreaterThan",
	OpGreaterThanOrEqual: "GreaterThanOrEqual",
	OpInstanceof:         "Instanceof",
	OpExit:               "Exit",
}

func (op Op) String() string {
	if len(op_map[op]) > 1 {
		return op_map[op]
	}

	return fmt.Sprintf("%d", op)
}

func NumberToBytes[T uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64](value T) []byte {
	var result []byte

	switch reflect.TypeOf(value).Kind() {
	case reflect.Uint8:
		result = []byte{uint8(value)}
	case reflect.Uint16:
		result = make([]byte, 2)
		binary.LittleEndian.PutUint16(result, uint16(value))
	case reflect.Uint32:
		result = make([]byte, 4)
		binary.LittleEndian.PutUint32(result, uint32(value))
	case reflect.Uint64:
		result = make([]byte, 8)
		binary.LittleEndian.PutUint64(result, uint64(value))
	case reflect.Int8:
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.LittleEndian, int8(value))
		result = buffer.Bytes()
	case reflect.Int16:
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.LittleEndian, int16(value))
		result = buffer.Bytes()
	case reflect.Int32:
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.LittleEndian, int32(value))
		result = buffer.Bytes()
	case reflect.Int64:
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.LittleEndian, int64(value))
		result = buffer.Bytes()
	case reflect.Float32:
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.LittleEndian, float32(value))
		result = buffer.Bytes()
	case reflect.Float64:
		buffer := bytes.Buffer{}
		binary.Write(&buffer, binary.LittleEndian, float64(value))
		result = buffer.Bytes()
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
		result = append(result, NumberToBytes(operand)...)
	}

	return result
}

func (i Instruction) GetSize() int {
	return len(i.GetBytes())
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

func (i InstructionSet) GetSize() int {
	return len(i.GetBytes())
}
