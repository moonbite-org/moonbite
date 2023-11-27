package compiler

import (
	"fmt"

	"github.com/moonbite-org/moonbite/common"
)

type Binary struct {
	Modules []Module
}

func (b *Binary) RegisterModule(mod Module) {
	b.Modules = append(b.Modules, mod)
}

func (b *Binary) Compile() ([]byte, error) {
	result := []byte{}

	result = append(result, SerializeString(common.Config.ArchiveName)...)
	result = append(result, SerializeString(common.Config.VersionStamp)...)

	for _, mod := range b.Modules {
		compiled, err := mod.Compile()

		if err != nil {
			return result, err
		}

		result = append(result, compiled...)
	}

	return result, nil
}

type Module struct {
	Name             string
	IsEntry          bool
	context_stack    [][]byte
	pointer          int
	internal_pointer int
}

func NewModule(name string, is_entry bool) Module {
	return Module{
		Name:             name,
		IsEntry:          is_entry,
		pointer:          100,
		internal_pointer: 10,
		context_stack:    [][]byte{},
	}
}

func (m *Module) Compile() ([]byte, error) {
	if len(m.context_stack) > 1 {
		return []byte{}, fmt.Errorf("unclosed context")
	}
	bin := m.PopContext()

	sys := map[string]interface{}{
		"read":  m.next_pointer(),
		"write": m.next_pointer(),
	}
	m.Save(m.next_pointer(), sys)

	result := []byte{common.InstructionMap["define_module"]}
	result = append(result, SerializeBool(m.IsEntry)...)

	rest := []byte{}
	rest = append(rest, SerializeString(m.Name)...)
	rest = append(rest, m.PopContext()...)
	rest = append(rest, bin...)

	size, err := SerializeNumber(len(rest))

	if err != nil {
		return result, err
	}

	result = append(result, size...)
	result = append(result, rest...)

	return result, nil
}

func (m *Module) NextPointer() int {
	pointer := m.pointer
	m.pointer++

	return pointer
}

func (m *Module) next_pointer() int {
	pointer := m.internal_pointer
	m.internal_pointer++

	return pointer
}

func (m *Module) append(value []byte) {
	if len(m.context_stack) == 0 {
		m.context_stack = append(m.context_stack, []byte{})
	}

	m.context_stack[len(m.context_stack)-1] = append(m.context_stack[len(m.context_stack)-1], value...)
}

func (m *Module) Save(pointer int, value interface{}) int {
	result := []byte{common.InstructionMap["save"]}

	p, err := SerializeNumber(pointer)

	if err != nil {
		return pointer
	}

	result = append(result, p...)

	v, err := serialize_value(value)

	if err != nil {
		return pointer
	}

	result = append(result, v...)

	m.append(result)
	return pointer
}

func (m *Module) SaveRaw(pointer int, value []byte) int {
	result := []byte{common.InstructionMap["save"]}

	p, err := SerializeNumber(pointer)

	if err != nil {
		return pointer
	}

	result = append(result, p...)
	result = append(result, value...)

	m.append(result)
	return pointer
}

func (m *Module) Skip(size int) {
	result := []byte{common.InstructionMap["skip"]}

	s, _ := SerializeNumber(size)
	result = append(result, s...)

	m.append(result)
}

func (m *Module) SkipIfZero(pointer, size int) {
	result := []byte{common.InstructionMap["skip_if_zero"]}

	p, _ := SerializeNumber(pointer)
	result = append(result, p...)

	s, _ := SerializeNumber(size)
	result = append(result, s...)

	m.append(result)
}

// func (m *Module) Compare() int {}

// func (m *Module) ComparePointers() int {}

func (m *Module) Arithmetic(operation byte, left, right, output Argument) {
	result := []byte{common.InstructionMap["arithmetic"], operation}

	result = append(result, left.Bin()...)
	result = append(result, right.Bin()...)
	result = append(result, output.Bin()...)

	m.append(result)
}

type ArgumentKind byte

const (
	ValueArgumentKind ArgumentKind = iota
	PointerArgumentKind
	ParamArgumentKind
	ReturnPointerArgumentKind
)

type Argument interface {
	Kind() ArgumentKind
	Bin() []byte
}

type ValueArgument struct {
	Value []byte
}

func (a ValueArgument) Kind() ArgumentKind {
	return ValueArgumentKind
}

func (a ValueArgument) Bin() []byte {
	result := []byte{byte(ValueArgumentKind)}
	result = append(result, a.Value...)

	return result
}

type PointerArgument struct {
	Value int
}

func (a PointerArgument) Kind() ArgumentKind {
	return PointerArgumentKind
}

func (a PointerArgument) Bin() []byte {
	result := []byte{byte(PointerArgumentKind)}

	v, _ := SerializeNumber(a.Value)

	result = append(result, v...)

	return result
}

type ParamArgument struct {
	Value int
}

func (a ParamArgument) Kind() ArgumentKind {
	return ParamArgumentKind
}

func (a ParamArgument) Bin() []byte {
	result := []byte{byte(ParamArgumentKind)}

	v, _ := SerializeNumber(a.Value)

	result = append(result, v...)

	return result
}

type ReturnPointerArgument struct{}

func (a ReturnPointerArgument) Kind() ArgumentKind {
	return ReturnPointerArgumentKind
}

func (a ReturnPointerArgument) Bin() []byte {
	result := []byte{byte(ReturnPointerArgumentKind)}

	return result
}

func (m *Module) Call(pointer int, return_pointer int, args []Argument) int {
	result := []byte{}

	if return_pointer >= 0 {
		result = append(result, common.InstructionMap["call_returned"])
	} else {
		result = append(result, common.InstructionMap["call"])
	}

	p, err := SerializeNumber(pointer)

	if err != nil {
		return pointer
	}

	result = append(result, p...)

	if return_pointer >= 0 {
		r_p, err := SerializeNumber(return_pointer)

		if err != nil {
			return pointer
		}

		result = append(result, r_p...)
	}

	result = append(result, common.TypeMap["list"])
	for _, arg := range args {
		result = append(result, arg.Bin()...)
	}
	result = append(result, common.TypeMap["terminator"])

	m.append(result)
	return pointer
}

func (m *Module) PushContext() {
	m.context_stack = append(m.context_stack, []byte{})
}

func (m *Module) PopContext() []byte {
	if len(m.context_stack) == 0 {
		return []byte{}
	}

	current := m.context_stack[len(m.context_stack)-1]
	m.context_stack = m.context_stack[0 : len(m.context_stack)-1]

	return current
}
