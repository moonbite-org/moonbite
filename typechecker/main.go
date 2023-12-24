package main

import (
	"fmt"
	"strings"

	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
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

type Generic struct {
	Index   int
	Literal TypeLiteral
}

type TypeDefinition struct {
	Name       string
	Definition TypeLiteral
	Generics   map[string]Generic
}

func (d TypeDefinition) String() string {
	return fmt.Sprintf("<type %s %+v %+v>", d.Name, d.Definition, d.Generics)
}

type TypeLiteral struct {
	Base             int
	Literal          parser.TypeLiteral
	GenericArguments map[string]parser.TypeLiteral
	extender         int
}

func (l *TypeLiteral) Extend(literal parser.TypeLiteral) TypeLiteral {
	new_base := l.extender * l.Base
	l.extender++

	return TypeLiteral{
		Base:     new_base,
		Literal:  literal,
		extender: 2,
	}
}

type BuiltinType struct {
	Name string
}

func (t BuiltinType) TypeKind() parser.TypeKind {
	return parser.BuiltinTypeKind
}

func (t BuiltinType) Location() errors.Location {
	return errors.CreateAnonError(0, "").Location
}

type SymbolTable struct {
	Symbols map[string]*TypeDefinition
}

func (t SymbolTable) String() string {
	result := []string{}

	for _, def := range t.Symbols {
		result = append(result, def.String())
	}

	return strings.Join(result, "\n")
}

func (t SymbolTable) Get(name string) *TypeDefinition {
	return t.Symbols[name]
}

func (t SymbolTable) Set(name string, def TypeDefinition) {
	t.Symbols[name] = &def
}

type Typechecker struct {
	SymbolTable SymbolTable
}

func (c Typechecker) ResolveLiteral(literal parser.TypeLiteral) (*TypeLiteral, errors.Error) {
	switch literal.TypeKind() {
	case parser.TypeIdentifierKind:
		name := resolve_name(literal.(parser.TypeIdentifier).Name)
		symbol := c.SymbolTable.Get(name)

		if symbol == nil {
			return nil, errors.CreateTypeError(fmt.Sprintf("cannot find type '%s'", name), literal.Location())
		}

		generics := literal.(parser.TypeIdentifier).Generics

		if len(symbol.Generics) != len(generics) {
			return nil, errors.CreateTypeError(fmt.Sprintf("type '%s' expects %d arguments but %d provided", name, len(symbol.Generics), len(generics)), literal.Location())
		}

		return &symbol.Definition, errors.EmptyError
	case parser.OperatedTypeKind:
		operated := literal.(parser.OperatedType)
		left, err := c.ResolveLiteral(operated.LeftHandSide)
		if err.Exists {
			return nil, err
		}
		right, err := c.ResolveLiteral(operated.RightHandSide)
		if err.Exists {
			return nil, err
		}

		return &TypeLiteral{
			Base:             left.Base * right.Base,
			Literal:          operated,
			GenericArguments: map[string]parser.TypeLiteral{},
			extender:         2,
		}, errors.EmptyError
	case parser.GroupTypeKind:
		return c.ResolveLiteral(literal.(parser.GroupType).Type)
	case parser.TypedLiteralKind:
		return c.ResolveLiteral(literal.(parser.TypedLiteral).Type)
	default:
		return nil, errors.CreateTypeError(fmt.Sprintf("cannot resolve this type %s", literal.TypeKind()), literal.Location())
	}
}

func (c Typechecker) Define(statement parser.TypeDefinitionStatement) errors.Error {
	if c.SymbolTable.Get(statement.Name.Value) != nil {
		return errors.CreateTypeError(fmt.Sprintf("type '%s' is already defined", statement.Name.Value), statement.Name.Location())
	}

	literal, err := c.ResolveLiteral(statement.Definition)
	if err.Exists {
		return err
	}

	generics := map[string]Generic{}
	for _, generic := range statement.Generics {
		generics[generic.Name.Value] = Generic{
			Index:   generic.Index,
			Literal: c.SymbolTable.Get("any").Definition,
		}
	}

	switch statement.Definition.TypeKind() {
	case parser.TypeIdentifierKind:
		definition := statement.Definition.(parser.TypeIdentifier)
		reference_name := resolve_name(definition.Name)
		reference_symbol := c.SymbolTable.Get(reference_name)

		if len(definition.Generics) != len(reference_symbol.Generics) {
			return errors.CreateTypeError(fmt.Sprintf("type '%s' expects %d arguments but %d provided", reference_name, len(reference_symbol.Generics), len(definition.Generics)), definition.Name.Location())
		}

		for index, literal := range definition.Generics {
			generic := find_indexed_generic(reference_symbol.Generics, index)
			argument, err := c.ResolveLiteral(literal)
			if err.Exists {
				return err
			}

			if argument.Base&generic.Literal.Base != 0 {
				return errors.CreateTypeError("cannot assign", literal.Location())
			}
		}
	}

	c.SymbolTable.Set(statement.Name.Value, TypeDefinition{
		Name:       statement.Name.Value,
		Definition: literal.Extend(statement.Definition),
		Generics:   generics,
	})

	return errors.EmptyError
}

func (c Typechecker) Check(a, b parser.TypeLiteral) (bool, errors.Error) {
	left, err := c.ResolveLiteral(a)
	if err.Exists {
		return false, err
	}
	right, err := c.ResolveLiteral(b)
	if err.Exists {
		return false, err
	}

	if right.Base%left.Base != 0 {
		return false, errors.EmptyError
	}

	return true, errors.EmptyError
}

func New() Typechecker {
	checker := Typechecker{
		SymbolTable: SymbolTable{
			Symbols: map[string]*TypeDefinition{},
		},
	}

	for name, base := range builtin_map {
		def := TypeDefinition{
			Name: name,
			Definition: TypeLiteral{
				Base:     base,
				Literal:  BuiltinType{name},
				extender: 2,
			},
			Generics: map[string]Generic{},
		}
		checker.SymbolTable.Set(name, def)
	}

	return checker
}

func define_types(checker Typechecker) {
	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "Uint32"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypeIdentifier{
			Name:     parser.IdentifierExpression{Value: "uint32"},
			Generics: map[int]parser.TypeLiteral{},
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "Uint64"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypeIdentifier{
			Name:     parser.IdentifierExpression{Value: "uint64"},
			Generics: map[int]parser.TypeLiteral{},
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "Bool"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypeIdentifier{
			Name:     parser.IdentifierExpression{Value: "bool"},
			Generics: map[int]parser.TypeLiteral{},
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "Rune"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypeIdentifier{
			Name:     parser.IdentifierExpression{Value: "Uint32"},
			Generics: map[int]parser.TypeLiteral{},
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name: parser.IdentifierExpression{Value: "List"},
		Generics: map[string]parser.ConstrainedType{
			"T": {
				Index:      0,
				Name:       parser.IdentifierExpression{Value: "T"},
				Constraint: nil,
			},
		},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypeIdentifier{
			Name:     parser.IdentifierExpression{Value: "iterable"},
			Generics: map[int]parser.TypeLiteral{},
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "String"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypeIdentifier{
			Name: parser.IdentifierExpression{Value: "List"},
			Generics: map[int]parser.TypeLiteral{
				0: parser.TypeIdentifier{
					Name:     parser.IdentifierExpression{Value: "Rune"},
					Generics: map[int]parser.TypeLiteral{},
				},
			},
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "Number"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.OperatedType{
			LeftHandSide: parser.TypeIdentifier{
				Name:     parser.IdentifierExpression{Value: "Uint32"},
				Generics: map[int]parser.TypeLiteral{},
			},
			RightHandSide: parser.TypeIdentifier{
				Name:     parser.IdentifierExpression{Value: "Uint64"},
				Generics: map[int]parser.TypeLiteral{},
			},
			Operator: "|",
		},
	}); err.Exists {
		fmt.Println(err)
	}

	if err := checker.Define(parser.TypeDefinitionStatement{
		Name:            parser.IdentifierExpression{Value: "Admin"},
		Generics:        map[string]parser.ConstrainedType{},
		Implementations: []parser.TypeIdentifier{},
		Definition: parser.TypedLiteral{
			Type: parser.TypeIdentifier{
				Name:     parser.IdentifierExpression{Value: "String"},
				Generics: map[int]parser.TypeLiteral{},
			},
			Literal: parser.StringLiteralExpression{
				Value: "admin",
			},
		},
	}); err.Exists {
		fmt.Println(err)
	}
}

func main() {
	checker := New()
	define_types(checker)

	fmt.Println(checker.Check(parser.TypeIdentifier{
		Name:     parser.IdentifierExpression{Value: "Uint32"},
		Generics: map[int]parser.TypeLiteral{
			// 0: parser.TypeIdentifier{
			// 	Name:     parser.IdentifierExpression{Value: "Rune"},
			// 	Generics: map[int]parser.TypeLiteral{},
			// },
		},
	}, parser.TypeIdentifier{
		Name:     parser.IdentifierExpression{Value: "Rune"},
		Generics: map[int]parser.TypeLiteral{},
	}))

	fmt.Println(checker.SymbolTable)
}

func resolve_name(name parser.Expression) string {
	switch name.Kind() {
	case parser.IdentifierExpressionKind:
		return name.(parser.IdentifierExpression).Value
	case parser.MemberExpressionKind:
		left := resolve_name(name.(parser.MemberExpression).LeftHandSide)
		right := name.(parser.MemberExpression).RightHandSide.Value

		return fmt.Sprintf("%s.%s", left, right)
	default:
		return "cannot resolve"
	}
}

func find_indexed_generic(generics map[string]Generic, index int) Generic {
	for _, g := range generics {
		if g.Index == index {
			return g
		}
	}

	return Generic{}
}
