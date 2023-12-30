package cmd

import (
	"fmt"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type Typechecker struct {
	TypeTable *TypeTable
}

func (c Typechecker) IsAssignable(typ parser.TypeLiteral, value parser.Expression) {
	fmt.Println(typ, value)
}

func (c Typechecker) IsCallable() {}

func (c Typechecker) IsCoercible() {}

func (c Typechecker) IsIndexable() {}

func New() Typechecker {
	return Typechecker{
		TypeTable: NewTypeTable(),
	}
}
