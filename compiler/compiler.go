package compiler

import (
	"fmt"
)

type Compiler struct {
	Modules             []Module
	is_entry_registered bool
}

func (c *Compiler) RegisterModule(mod Module) error {
	if mod.IsEntry {
		if c.is_entry_registered {
			return fmt.Errorf("cannot register more than 1 entry module")
		}

		c.is_entry_registered = true
	}

	c.Modules = append(c.Modules, mod)
	return nil
}

func (c Compiler) Compile() ([]byte, error) {
	result := []byte{}

	for _, mod := range c.Modules {
		if mod.IsEntry {
			compiled, err := mod.Compile()

			if err != nil {
				return result, err
			}

			result = append(result, compiled...)
			break
		}
	}

	for _, mod := range c.Modules {
		if mod.IsEntry {
			continue
		}

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
	Builtins         map[string]int
	IsEntry          bool
	bin              []byte
	error            error
	pointer          int
	internal_pointer int
}

func (m Module) Compile() ([]byte, error) {
	result := []byte{}

	if m.error != nil {
		return result, m.error
	}

	head := m.bin[0:2]
	tail := m.bin[2:]

	result = append(result, head...)
	size, _ := SerializeNumber(len(tail))
	result = append(result, size...)
	result = append(result, tail...)

	return result, nil
}

func (m *Module) NextPointer() int {
	p := m.pointer
	m.pointer++

	return p
}

func (m *Module) next_pointer() int {
	p := m.internal_pointer
	m.internal_pointer++

	return p
}

func (m *Module) Save(pointer int, value interface{}) int {
	if m.error != nil {
		return -1
	}

	i, e := save(pointer, value)
	m.error = e

	m.bin = append(m.bin, i...)

	return pointer
}

func (m *Module) Call(pointer int, args []int) {
	if m.error != nil {
		return
	}

	i, e := call(pointer, args)
	m.error = e

	m.bin = append(m.bin, i...)
}

func (m *Module) LoadPointer(pointer int) {
	if m.error != nil {
		return
	}

	i, e := load_pointer(pointer)
	m.error = e

	m.bin = append(m.bin, i...)
}

func (m *Module) LoadValue(value interface{}) {
	if m.error != nil {
		return
	}

	i, e := load_value(value)
	m.error = e

	m.bin = append(m.bin, i...)
}

func (m *Module) LoadParam(index int) {
	if m.error != nil {
		return
	}

	i, e := load_param(index)
	m.error = e

	m.bin = append(m.bin, i...)
}

func NewModule(name string, is_entry bool) Module {
	mod := Module{
		Name:             name,
		IsEntry:          is_entry,
		bin:              []byte{},
		pointer:          100,
		internal_pointer: 10,
		Builtins:         map[string]int{},
	}

	mod.bin = append(mod.bin, define_module(name, is_entry)...)
	generate_builtins(&mod)

	return mod
}
