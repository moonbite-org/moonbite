package cmd

import (
	"fmt"
	"strings"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type SymbolScope string

const (
	BuiltinScope SymbolScope = "scope:builtin"
	GlobalScope  SymbolScope = "scope:global"
	LocalScope   SymbolScope = "scope:local"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
	Kind  parser.VarKind
}

type SymbolTable struct {
	Outer *SymbolTable

	store map[string]Symbol
	count int
}

func (t *SymbolTable) Define(name string, kind parser.VarKind) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: t.count,
		Kind:  kind,
	}

	if t.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	t.count++
	t.store[name] = symbol

	return symbol
}

func (t SymbolTable) Resolve(name string) *Symbol {
	symbol, ok := t.store[name]

	if !ok {
		if t.Outer != nil {
			return t.Outer.Resolve(name)
		}

		return nil
	}

	return &symbol
}

func (t SymbolTable) String() string {
	result := []string{}

	for name, symbol := range t.store {
		result = append(result, fmt.Sprintf("%s: {%d %s %s}", name, symbol.Index, symbol.Scope, symbol.Kind))
	}

	return fmt.Sprintf("[%s]", strings.Join(result, " "))
}

func NewSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	return &SymbolTable{store: store, Outer: nil}
}

func NewScopedSymbolTable(outer *SymbolTable) *SymbolTable {
	table := NewSymbolTable()
	table.Outer = outer
	return table
}
