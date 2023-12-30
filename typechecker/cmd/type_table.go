package cmd

import (
	"fmt"
	"strings"
)

type Type struct {
	Base     int
	extender int
	Name     string
	Index    int
	Hidden   bool
}

type TypeTable struct {
	store map[string]Type
	count int
}

func (t TypeTable) Extend(new_type, extended_type string, hidden bool) (Type, error) {
	symbol := t.Resolve(extended_type)

	if symbol == nil {
		return Type{}, fmt.Errorf("no type %s", extended_type)
	}

	new_symbol, err := t.Define(new_type, hidden, symbol.extender*symbol.Base)
	if err != nil {
		return Type{}, err
	}

	symbol.extender++

	return new_symbol, nil
}

func (t *TypeTable) Define(name string, hidden bool, base int) (Type, error) {
	_, exists := t.store[name]

	if exists {
		return Type{}, fmt.Errorf("cannot redeclare variable '%s' on the same block", name)
	}

	symbol := Type{
		Base:     base,
		extender: 2,
		Name:     name,
		Index:    t.count,
		Hidden:   hidden,
	}

	t.count++
	t.store[name] = symbol

	return symbol, nil
}

func (t TypeTable) Resolve(name string) *Type {
	symbol, ok := t.store[name]

	if !ok {
		return nil
	}

	return &symbol
}

func (t TypeTable) String() string {
	result := []string{}

	for name, symbol := range t.store {
		result = append(result, fmt.Sprintf("%s: {index: %d base: %d extender: %d}", name, symbol.Index, symbol.Base, symbol.extender))
	}

	return fmt.Sprintf("[%s]", strings.Join(result, " "))
}

func NewTypeTable() *TypeTable {
	store := make(map[string]Type)
	return &TypeTable{store: store}
}
