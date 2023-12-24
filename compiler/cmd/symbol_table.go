package cmd

type SymbolScope string

const (
	GlobalScope SymbolScope = "scope:global"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store map[string]Symbol
	count int
}

func (t *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: t.count,
	}
	t.count++
	t.store[name] = symbol

	return symbol
}

func (t SymbolTable) Resolve(name string) *Symbol {
	symbol := t.store[name]
	return &symbol
}

func NewSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	return &SymbolTable{store: store}
}
