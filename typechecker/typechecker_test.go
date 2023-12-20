package typechecker_test

import (
	"fmt"
	"testing"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
	"github.com/moonbite-org/moonbite/typechecker"
)

func parse_type_def(name, def string) parser.TypeDefinitionStatement {
	input := fmt.Sprintf(`package main
	type %s %s
	`, name, def)

	ast, _ := parser.Parse([]byte(input), "")

	return ast.Definitions[0].(parser.TypeDefinitionStatement)
}

func assert(t *testing.T, given, expected any) {
	if given != expected {
		t.Errorf("expected '%t' but got '%t'", expected, given)
	}
}

func TestComptimeBuiltins(t *testing.T) {
	checker := typechecker.New()
	checker.SetComptimeBuiltins()

	assert(t, checker.SymbolTable.GetByName("uint32"), checker.SymbolTable.GetByName("uint32"))
	cases := []struct {
		name string
		base int
	}{
		{
			name: "function",
			base: 11,
		},
		{
			name: "iterable",
			base: 13,
		},
		{
			name: "map",
			base: 17,
		},
		{
			name: "uint8",
			base: 19,
		},
		{
			name: "uint16",
			base: 23,
		},
		{
			name: "uint32",
			base: 29,
		},
		{
			name: "uint64",
			base: 31,
		},
		{
			name: "int8",
			base: 37,
		},
		{
			name: "int16",
			base: 41,
		},
		{
			name: "int32",
			base: 43,
		},
		{
			name: "int64",
			base: 47,
		},
		{
			name: "bool",
			base: 53,
		},
	}

	for _, case_ := range cases {
		typ := checker.SymbolTable.GetByName(case_.name)
		if typ == nil {
			t.Errorf("expected %s but type is not defined in symbol table", case_.name)
		} else {
			assert(t, typ.Base, case_.base)
		}
	}
}

func Test(t *testing.T) {
	checker := typechecker.New()
	checker.SetComptimeBuiltins()
	checker.SetRuntimeBuiltins()

	fmt.Println(checker.DefineType(parse_type_def("Foo<T, K>", "List<K>")))

	// fmt.Println(checker.ResolveTypeLiteral(parser.StructLiteral{
	// 	Values: []parser.ValueTypePair{
	// 		{
	// 			Key: parser.IdentifierExpression{Value: "username"},
	// 			Type: parser.TypeIdentifier{
	// 				Name:     parser.IdentifierExpression{Value: "String"},
	// 				Generics: make(map[int]parser.TypeLiteral),
	// 			},
	// 		},
	// 		{
	// 			Key: parser.IdentifierExpression{Value: "age"},
	// 			Type: parser.TypeIdentifier{
	// 				Name:     parser.IdentifierExpression{Value: "Int32"},
	// 				Generics: make(map[int]parser.TypeLiteral),
	// 			},
	// 		},
	// 		{
	// 			Key: parser.IdentifierExpression{Value: "age"},
	// 			Type: parser.TypeIdentifier{
	// 				Name:     parser.IdentifierExpression{Value: "Uint32"},
	// 				Generics: make(map[int]parser.TypeLiteral),
	// 			},
	// 		},
	// 	},
	// }, nil))
}
