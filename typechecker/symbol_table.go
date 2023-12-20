package typechecker

type SymbolTable struct {
	Definitions map[string]*Type
}

func (table SymbolTable) GetByName(name string) *Type {
	return table.Definitions[name]
}

func (table SymbolTable) GetByBase(base int) *Type {
	for _, typ := range table.Definitions {
		if typ.Base == base {
			return typ
		}
	}

	return nil
}

func (table *SymbolTable) Set(name string, typ *Type) {
	table.Definitions[name] = typ
}
