package typechecker

import (
	"fmt"
	"os"
	"path"

	"github.com/moonbite-org/moonbite/common"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type TypeChecker struct {
	File        string
	SymbolTable SymbolTable
	LibPath     string
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

	file_path := path.Join(c.LibPath, "types.mb")
	types, err := os.ReadFile(file_path)
	if err != nil {
		panic(err)
	}

	ast, f_err := parser.Parse(types, path.Join(c.LibPath, "types.mb"))
	if f_err.Exists {
		panic(f_err)
	}

	for _, definition := range ast.Definitions {
		switch definition.Kind() {
		case parser.TypeDefinitionStatementKind:
			if err := c.DefineType(definition.(parser.TypeDefinitionStatement)); err.Exists {
				panic(err)
			}
		case parser.TraitDefinitionStatementKind:
			if err := c.DefineTrait(definition.(parser.TraitDefinitionStatement)); err.Exists {
				panic(err)
			}
		}
	}
}

func (c TypeChecker) ResolveTypeLiteral(literal parser.TypeLiteral, top_frame *SymbolTable) (*Type, parser.Error) {
	switch literal.TypeKind() {
	case parser.TypeIdentifierKind:
		name := resolve_name(literal.(parser.TypeIdentifier).Name)
		symbol := c.SymbolTable.GetByName(name)

		if symbol == nil {
			if top_frame != nil {
				symbol = top_frame.GetByName(name)
			}

			if symbol == nil {
				return nil, parser.CreateTypeError(fmt.Sprintf("cannot find type %s", name), literal.Location())
			}
		}

		return symbol, parser.EmptyError
	case parser.OperatedTypeKind:
		operated := literal.(parser.OperatedType)
		left, err := c.ResolveTypeLiteral(operated.LeftHandSide, nil)
		if err.Exists {
			return nil, err
		}

		right, err := c.ResolveTypeLiteral(operated.RightHandSide, nil)
		if err.Exists {
			return nil, err
		}

		switch operated.Operator {
		case "|":
			new_type := NewType(left.Base*right.Base, parser.OperatedTypeKind, []*Type{}, map[string]type_param{})
			return &new_type, parser.EmptyError
		case "&":
			new_type := NewType(left.Base*right.Base, parser.OperatedTypeKind, []*Type{}, map[string]type_param{})
			return &new_type, parser.EmptyError
		default:
			return nil, parser.CreateTypeError(fmt.Sprintf("unknown operator %s", operated.Operator), operated.LeftHandSide.Location())
		}
	case parser.GroupTypeKind:
		return c.ResolveTypeLiteral(literal.(parser.GroupType).Type, nil)
	case parser.TypedLiteralKind:
		typ, err := c.ResolveTypeLiteral(literal.(parser.TypedLiteral).Type, nil)
		if err.Exists {
			return nil, err
		}

		result := typ.Extend([]*Type{}, map[string]type_param{}, parser.TypedLiteralKind)
		data := common.ObjectFromLiteral(literal.(parser.TypedLiteral).Literal)
		result.SetMetaData(data)

		return &result, parser.EmptyError
	case parser.FunTypeKind:
		fun_symbol := c.SymbolTable.GetByName("function")
		signature := literal.(parser.FunctionSignature)

		parameters := type_param_stack{}
		for _, generic := range signature.GetGenerics() {
			if generic.Constraint == nil {
				parameters.AddGeneric(generic.Name.Value, generic.Index, &t_any)
			} else {
				param, err := c.ResolveTypeLiteral(*generic.Constraint, nil)
				if err.Exists {
					return nil, err
				}
				parameters.AddGeneric(generic.Name.Value, generic.Index, param)
			}
		}

		result := fun_symbol.Extend([]*Type{}, parameters, parser.FunTypeKind)
		frame := result.GetSymbolTable().Merge(top_frame)

		for i, f_param := range signature.GetParameters() {
			param, err := c.ResolveTypeLiteral(f_param.Type, &frame)
			if err.Exists {
				return nil, err
			}
			parameters.AddArgument(i, param)
		}

		return_type := signature.GetReturnType()
		if return_type != nil {
			typ, err := c.ResolveTypeLiteral(*return_type, &frame)
			if err.Exists {
				return nil, err
			}
			parameters.SetReturnType(typ)
		}

		result = fun_symbol.Extend([]*Type{}, parameters, parser.FunTypeKind)
		return &result, parser.EmptyError
	case parser.TraitTypeKind:
		trait := literal.(parser.TraitDefinitionStatement)
		map_symbol := c.SymbolTable.GetByName("map")

		parameters := type_param_stack{}

		for _, generic := range trait.Generics {
			if generic.Constraint == nil {
				parameters.AddGeneric(generic.Name.Value, generic.Index, &t_any)
			} else {
				param, err := c.ResolveTypeLiteral(*generic.Constraint, nil)
				if err.Exists {
					return nil, err
				}
				parameters.AddGeneric(generic.Name.Value, generic.Index, param)
			}
		}

		result := map_symbol.Extend([]*Type{}, parameters, parser.TraitTypeKind)
		frame := result.GetSymbolTable().Merge(top_frame)

		for _, mimick := range trait.Mimics {
			typ, err := c.ResolveTypeLiteral(mimick, &frame)
			if err.Exists {
				return nil, err
			}

			if typ.Kind != parser.TraitTypeKind {
				return nil, parser.CreateTypeError(fmt.Sprintf("cannot mimic non-trait %s kind", typ.Kind), mimick.Location())
			}

			methods := typ.Parameters.GetMethods()

			index := 0
			for name, method := range methods {
				parameters.AddMethod(name, index, method)
				index++
			}
		}

		result = map_symbol.Extend([]*Type{}, parameters, parser.TraitTypeKind)
		frame = result.GetSymbolTable().Merge(top_frame)
		offset := len(parameters.GetMethods())

		for i, definition := range trait.Definition {
			if parameters.GetMethod(definition.Name.Value) != nil {
				return nil, parser.CreateTypeError(fmt.Sprintf("duplicate method '%s'", definition.Name.Value), definition.Name.Location())
			}

			typ, err := c.ResolveTypeLiteral(definition, &frame)
			if err.Exists {
				return nil, err
			}

			parameters.AddMethod(definition.Name.Value, i+offset, typ)
		}

		result = map_symbol.Extend([]*Type{}, parameters, parser.TraitTypeKind)

		return &result, parser.EmptyError
	case parser.StructLiteralKind:
		map_symbol := c.SymbolTable.GetByName("map")

		parameters := type_param_stack{}
		for i, pair := range literal.(parser.StructLiteral).Values {
			name := pair.Key.Value

			if parameters.GetProperty(name) != nil {
				return nil, parser.CreateTypeError(fmt.Sprintf("duplicate field '%s'", name), pair.Key.Location())
			}

			typ, err := c.ResolveTypeLiteral(pair.Type, nil)
			if err.Exists {
				return nil, err
			}

			parameters.AddProperty(name, i, typ)
		}

		result := map_symbol.Extend([]*Type{}, parameters, parser.StructLiteralKind)
		return &result, parser.EmptyError
	default:
		return nil, parser.CreateTypeError(fmt.Sprintf("unknown type kind %s", literal.TypeKind()), literal.Location())
	}
}

func (c TypeChecker) DefineType(statement parser.TypeDefinitionStatement) parser.Error {
	type_name := resolve_name(statement.Name)

	if c.SymbolTable.GetByName(type_name) != nil {
		return parser.CreateTypeError(fmt.Sprintf("type %s has already been declared", type_name), statement.Name.Location())
	}

	switch statement.Definition.TypeKind() {
	case parser.TypeIdentifierKind:
		return c.DefineTypeIdentifier(type_name, statement)
	default:
		return parser.CreateTypeError(fmt.Sprintf("cannot define this type yet %s", statement.Definition.TypeKind()), statement.Location())
	}
}

func (c *TypeChecker) DefineTypeIdentifier(type_name string, statement parser.TypeDefinitionStatement) parser.Error {
	definition := statement.Definition.(parser.TypeIdentifier)
	reference_name := resolve_name(definition.Name)
	reference_symbol := c.SymbolTable.GetByName(reference_name)

	if reference_symbol == nil {
		return parser.CreateTypeError(fmt.Sprintf("cannot find type %s", reference_name), definition.Name.Location())
	}

	if len(definition.Generics) > len(reference_symbol.Parameters) {
		return parser.CreateTypeError(fmt.Sprintf("type %s only accepts %d generic(s)", reference_name, len(reference_symbol.Parameters)), definition.Name.Location())
	}

	parameters := type_param_stack{}
	for _, generic := range statement.Generics {
		if parameters.GetGenericByName(generic.Name.Value) != nil {
			return parser.CreateTypeError(fmt.Sprintf("duplicate generic '%s'", generic.Name.Value), generic.Name.Location())
		}

		if generic.Constraint == nil {
			parameters.AddGeneric(generic.Name.Value, generic.Index, &t_any)
		} else {
			constraint, err := c.ResolveTypeLiteral(*generic.Constraint, nil)
			if err.Exists {
				return err
			}

			parameters.AddGeneric(generic.Name.Value, generic.Index, constraint)
		}
	}

	result := reference_symbol.Extend([]*Type{}, parameters, parser.TypeIdentifierKind)
	frame := result.GetSymbolTable()

	for name, param := range reference_symbol.Parameters {
		argument, ok := definition.Generics[param.index]

		if !ok {
			return parser.CreateTypeError(fmt.Sprintf("generic parameter '%s' is required for type '%s'", name, reference_name), definition.Name.Location())
		}

		typ, err := c.ResolveTypeLiteral(argument, &frame)

		if err.Exists {
			return err
		}
		result.Arguments = append(result.Arguments, typ)
	}

	c.SymbolTable.Set(type_name, &result)
	return parser.EmptyError
}

func (c TypeChecker) DefineTrait(statement parser.TraitDefinitionStatement) parser.Error {
	typ, err := c.ResolveTypeLiteral(statement, nil)

	if err.Exists {
		return err
	}

	c.SymbolTable.Set(statement.Name.Value, typ)

	return parser.EmptyError
}

func New(file string) TypeChecker {
	return TypeChecker{
		File: file,
		SymbolTable: SymbolTable{
			Definitions: map[string]*Type{},
		},
		LibPath: os.Getenv("MOON_LIB"),
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
