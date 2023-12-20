package typechecker

import (
	"fmt"

	"github.com/moonbite-org/moonbite/common"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type TypeChecker struct {
	SymbolTable SymbolTable
}

func (c *TypeChecker) SetComptimeBuiltins() {
	for name, prime := range builtin_primes {
		typ := NewType(prime, parser.TypeIdentifierKind, []*Type{}, map[string]type_param{})
		c.SymbolTable.Set(name, &typ)
	}
}

func (c *TypeChecker) SetRuntimeBuiltins() {
	defer (func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	})()

	resolve := func(name string) *Type {
		t, err := c.ResolveTypeLiteral(parser.TypeIdentifier{
			Name: parser.IdentifierExpression{Value: name},
		}, nil)

		if err != nil {
			panic(err)
		}

		return t
	}

	t_uint8 := resolve("uint8")

	t_Uint8 := t_uint8.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Uint8", &t_Uint8)

	t_uint16 := resolve("uint16")

	t_Uint16 := t_uint16.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Uint16", &t_Uint16)

	t_uint32 := resolve("uint32")

	t_Uint32 := t_uint32.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Uint32", &t_Uint32)

	t_uint64 := resolve("uint64")

	t_Uint64 := t_uint64.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Uint64", &t_Uint64)

	t_int8 := resolve("int8")

	t_Int8 := t_int8.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Int8", &t_Int8)

	t_int16 := resolve("int16")

	t_Int16 := t_int16.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Int16", &t_Int16)

	t_int32 := resolve("int32")

	t_Int32 := t_int32.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Int32", &t_Int32)

	t_int64 := resolve("int64")

	t_Int64 := t_int64.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Int64", &t_Int64)

	t_Rune := t_Uint32.Extend([]*Type{}, map[string]type_param{})
	c.SymbolTable.Set("Rune", &t_Rune)

	t_iterable := resolve("iterable")

	t_List := t_iterable.Extend([]*Type{}, map[string]type_param{"T": {0, &t_any}})
	c.SymbolTable.Set("List", &t_List)

	t_String := t_List.Extend([]*Type{&t_Rune}, map[string]type_param{})
	c.SymbolTable.Set("String", &t_String)

	var ref_string parser.TypeLiteral = parser.TypeIdentifier{
		Name: parser.IdentifierExpression{Value: "String"},
	}
	t_Printable, err := c.ResolveTypeLiteral(parser.TraitDefinitionStatement{
		Name:     parser.IdentifierExpression{Value: "Printable"},
		Generics: []parser.ConstrainedType{},
		Mimics:   []parser.TypeIdentifier{},
		Definition: []parser.UnboundFunctionSignature{
			{
				Name:       parser.IdentifierExpression{Value: "string"},
				Parameters: []parser.TypedParameter{},
				Generics:   []parser.ConstrainedType{},
				ReturnType: &ref_string,
			},
		},
	}, nil)

	if err != nil {
		panic(err)
	}

	c.SymbolTable.Set("Printable", t_Printable)
}

func (c TypeChecker) ResolveTypeLiteral(literal parser.TypeLiteral, frame *SymbolTable) (*Type, error) {
	switch literal.TypeKind() {
	case parser.TypeIdentifierKind:
		name := resolve_name(literal.(parser.TypeIdentifier).Name)
		symbol := c.SymbolTable.GetByName(name)

		if symbol == nil {
			if frame != nil {
				symbol = frame.GetByName(name)
			}

			if symbol == nil {
				return nil, fmt.Errorf("cannot find type %s", name)
			}
		}

		return symbol, nil
	case parser.OperatedTypeKind:
		operated := literal.(parser.OperatedType)
		left, err := c.ResolveTypeLiteral(operated.LeftHandSide, nil)
		if err != nil {
			return nil, err
		}

		right, err := c.ResolveTypeLiteral(operated.RightHandSide, nil)
		if err != nil {
			return nil, err
		}

		switch operated.Operator {
		case "|":
			new_type := NewType(left.Base*right.Base, parser.OperatedTypeKind, []*Type{}, map[string]type_param{})
			return &new_type, nil
		case "&":
			new_type := NewType(left.Base*right.Base, parser.OperatedTypeKind, []*Type{}, map[string]type_param{})
			return &new_type, nil
		default:
			return nil, fmt.Errorf("unknown operator %s", operated.Operator)
		}
	case parser.GroupTypeKind:
		return c.ResolveTypeLiteral(literal.(parser.GroupType).Type, nil)
	case parser.TypedLiteralKind:
		typ, err := c.ResolveTypeLiteral(literal.(parser.TypedLiteral).Type, nil)
		if err != nil {
			return nil, err
		}

		result := typ.Extend([]*Type{}, map[string]type_param{})
		// TODO: create object from literal
		result.SetMetaData(common.StringObject{Value: "debug"})

		return &result, nil
	case parser.FunTypeKind:
		fun_symbol := c.SymbolTable.GetByName("function")
		signature := literal.(parser.FunctionSignature)

		parameters := type_param_stack{}
		for _, generic := range signature.GetGenerics() {
			if generic.Constraint == nil {
				parameters.AddGeneric(generic.Name.Value, generic.Index, &t_any)
			} else {
				param, err := c.ResolveTypeLiteral(*generic.Constraint, nil)
				if err != nil {
					return nil, err
				}
				parameters.AddGeneric(generic.Name.Value, generic.Index, param)
			}
		}

		result := fun_symbol.Extend([]*Type{}, parameters)
		frame := result.GetSymbolTable()

		for i, f_param := range signature.GetParameters() {
			param, err := c.ResolveTypeLiteral(f_param.Type, &frame)
			if err != nil {
				return nil, err
			}
			parameters.AddArgument(i, param)
		}

		return_type := signature.GetReturnType()
		if return_type != nil {
			typ, err := c.ResolveTypeLiteral(*return_type, &frame)
			if err != nil {
				return nil, err
			}
			parameters.SetReturnType(typ)
		}

		result = fun_symbol.Extend([]*Type{}, parameters)
		return &result, nil
	case parser.TraitTypeKind:
		trait := literal.(parser.TraitDefinitionStatement)
		map_symbol := c.SymbolTable.GetByName("map")

		parameters := type_param_stack{}

		for _, generic := range trait.Generics {
			if generic.Constraint == nil {
				parameters.AddGeneric(generic.Name.Value, generic.Index, &t_any)
			} else {
				param, err := c.ResolveTypeLiteral(*generic.Constraint, nil)
				if err != nil {
					return nil, err
				}
				parameters.AddGeneric(generic.Name.Value, generic.Index, param)
			}
		}

		result := map_symbol.Extend([]*Type{}, parameters)
		frame := result.GetSymbolTable()

		for i, definition := range trait.Definition {
			typ, err := c.ResolveTypeLiteral(definition, &frame)
			if err != nil {
				return nil, err
			}

			parameters.AddMethod(i, typ)
		}

		result = map_symbol.Extend([]*Type{}, parameters)

		return &result, nil
	case parser.StructLiteralKind:
		fmt.Println(literal)
		map_symbol := c.SymbolTable.GetByName("map")

		parameters := type_param_stack{}
		for i, pair := range literal.(parser.StructLiteral).Values {
			name := pair.Key.Value

			if parameters.GetProperty(name) != nil {
				return nil, fmt.Errorf("duplicate field '%s'", name)
			}

			typ, err := c.ResolveTypeLiteral(pair.Type, nil)
			if err != nil {
				return nil, err
			}

			parameters.AddProperty(name, i, typ)
		}

		result := map_symbol.Extend([]*Type{}, parameters)
		return &result, nil
	default:
		return nil, fmt.Errorf("unknown type kind %s", literal.TypeKind())
	}
}

func (c TypeChecker) DefineType(statement parser.TypeDefinitionStatement) error {
	type_name := resolve_name(statement.Name)

	if c.SymbolTable.GetByName(type_name) != nil {
		return fmt.Errorf("type %s has already been declared", type_name)
	}

	switch statement.Definition.TypeKind() {
	case parser.TypeIdentifierKind:
		return c.DefineTypeIdentifier(type_name, statement)
	default:
		return fmt.Errorf("cannot define this type yet %s", statement.Definition.TypeKind())
	}
}

func (c *TypeChecker) DefineTypeIdentifier(type_name string, statement parser.TypeDefinitionStatement) error {
	definition := statement.Definition.(parser.TypeIdentifier)
	reference_name := resolve_name(definition.Name)
	reference_symbol := c.SymbolTable.GetByName(reference_name)

	if reference_symbol == nil {
		return fmt.Errorf("cannot find type %s", reference_name)
	}

	if len(definition.Generics) > len(reference_symbol.Parameters) {
		return fmt.Errorf("type %s only accepts %d generic(s)", reference_name, len(reference_symbol.Parameters))
	}

	parameters := type_param_stack{}
	for _, generic := range statement.Generics {
		if parameters.GetGenericByName(generic.Name.Value) != nil {
			return fmt.Errorf("duplicate generic '%s'", generic.Name.Value)
		}

		if generic.Constraint == nil {
			parameters.AddGeneric(generic.Name.Value, generic.Index, &t_any)
		} else {
			constraint, err := c.ResolveTypeLiteral(*generic.Constraint, nil)
			if err != nil {
				return err
			}

			parameters.AddGeneric(generic.Name.Value, generic.Index, constraint)
		}
	}

	result := reference_symbol.Extend([]*Type{}, parameters)
	frame := result.GetSymbolTable()

	for name, param := range reference_symbol.Parameters {
		argument, ok := definition.Generics[param.index]

		if !ok {
			return fmt.Errorf("generic parameter '%s' is required for type '%s'", name, reference_name)
		}

		typ, err := c.ResolveTypeLiteral(argument, &frame)

		if err != nil {
			return err
		}
		result.Arguments = append(result.Arguments, typ)
	}

	c.SymbolTable.Set(type_name, &result)
	return nil
}

func New() TypeChecker {
	return TypeChecker{
		SymbolTable: SymbolTable{
			Definitions: map[string]*Type{},
		},
	}
}

// func (c TypeChecker) Check(a, b parser.TypeLiteral) (bool, error) {
// 	symbol_a, err := c.ResolveTypeLiteral(a)
// 	if err != nil {
// 		return false, err
// 	}

// 	symbol_b, err := c.ResolveTypeLiteral(b)
// 	if err != nil {
// 		return false, err
// 	}

// 	// fmt.Println(symbol_a.Arguments, symbol_b.Arguments, symbol_a.Parameters, symbol_b.Parameters)
// 	// result := true

// 	// if len(symbol_a.Arguments) != len(symbol_b.Parameters) {
// 	// 	return false, nil
// 	// }

// 	// for i := 0; i < len(symbol_a.Parameters); i++ {
// 	// 	parameter := symbol_a.Parameters[i]
// 	// 	argument := symbol_b.Arguments[i]

// 	// 	if argument&parameter != 0 {
// 	// 		result = false
// 	// 		break
// 	// 	}
// 	// }

// 	// if !result {
// 	// 	return result, nil
// 	// }

// 	return c.check(a, b, symbol_a, symbol_b), nil
// }

// func (c TypeChecker) check(a, b parser.TypeLiteral, symbol_a, symbol_b *Type) bool {
// 	if a.TypeKind() == parser.OperatedTypeKind {
// 		if a.(parser.OperatedType).Operator == "|" {
// 			return symbol_a.Base%symbol_b.Base == 0
// 		}
// 	}

// 	if b.TypeKind() == parser.OperatedTypeKind {
// 		if b.(parser.OperatedType).Operator == "|" {
// 			return symbol_a.Base%symbol_b.Base == 0
// 		}
// 	}

// 	return symbol_b.Base%symbol_a.Base == 0
// }
