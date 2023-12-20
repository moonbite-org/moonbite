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

func (t SymbolTable) Merge(table *SymbolTable) SymbolTable {
	new_table := SymbolTable{
		Definitions: t.Definitions,
	}

	if table != nil {
		for name, definition := range table.Definitions {
			new_table.Definitions[name] = definition
		}
	}

	return new_table
}
