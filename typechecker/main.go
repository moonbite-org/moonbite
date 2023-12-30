package main

import (
	"fmt"

	checker "github.com/moonbite-org/moonbite/typechecker/cmd"
)

var builtin_map = map[string]int{
	"any":      1,
	"bool":     5,
	"iterable": 7,
	"uint8":    11,
	"uint16":   13,
	"uint32":   17,
	"uint64":   19,
}

// type TypeLiteral interface {
// 	Name() string
// 	Kind() parser.TypeKind
// }

// type Int32TypeLiteral struct{}

// func (t Int32TypeLiteral) Name() string {
// 	return "uint32"
// }

// func (t Int32TypeLiteral) Kind() parser.TypeKind {
// 	return parser.TypeIdentifierKind
// }

func main() {
	checker := checker.New()

	checker.TypeTable.Define("i32", false, 29)
	checker.TypeTable.Extend("rune", "i32", false)
	checker.TypeTable.Define("iterable", false, 7)
	checker.TypeTable.Extend("List", "iterable", false)
	checker.TypeTable.Extend("string", "List", false)
	// checker.TypeTable.Define("rune", false)

	fmt.Println(checker)
}
