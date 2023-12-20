package typechecker

import (
	"fmt"
	"strings"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func resolve_name(name parser.Expression) string {
	switch name.Kind() {
	case parser.IdentifierExpressionKind:
		return name.(parser.IdentifierExpression).Value
	default:
		panic("cannot handle names other than identifiers yet")
	}
}

const (
	ReturnTypeOffset int = iota + 1000
	InferredWarningOffset
	FunctionParameterOffset
	MethodOffset int = iota + 3997
)

type type_param struct {
	index int
	typ   *Type
}

type type_param_stack map[string]type_param

func (p type_param_stack) AddGeneric(name string, index int, generic *Type) {
	p[name] = type_param{index: index, typ: generic}
}

func (p type_param_stack) GetGenericByIndex(index int) *Type {
	for _, param := range p {
		if param.index == index {
			return param.typ
		}
	}

	return nil
}

func (p type_param_stack) GetGenericByName(name string) *Type {
	return p[name].typ
}

func (p type_param_stack) AddArgument(index int, param *Type) {
	p[prefixer("param_", index)] = type_param{index: index + FunctionParameterOffset, typ: param}
}

func (p type_param_stack) GetArgument(index int) *Type {
	return p[prefixer("param_", index)].typ
}

func (p type_param_stack) AddMethod(index int, param *Type) {
	p[prefixer("method_", index)] = type_param{index: index + MethodOffset, typ: param}
}

func (p type_param_stack) GetMethod(index int) *Type {
	return p[prefixer("method_", index)].typ
}

func (p type_param_stack) GetMethods(index int) map[string]*Type {
	result := map[string]*Type{}

	for name, property := range p {
		if strings.HasPrefix(name, prefixer("method_")) {
			result[name] = property.typ
		}
	}

	return result
}

func (p type_param_stack) AddProperty(name string, index int, param *Type) {
	property := fmt.Sprintf("property_%s", name)
	p[prefixer(property)] = type_param{index: index + MethodOffset, typ: param}
}

func (p type_param_stack) GetProperty(name string) *Type {
	property := fmt.Sprintf("property_%s", name)
	return p[prefixer(property)].typ
}

func (p type_param_stack) GetProperties(index int) map[string]*Type {
	result := map[string]*Type{}

	for name, property := range p {
		if strings.HasPrefix(name, prefixer("property_")) {
			result[name] = property.typ
		}
	}

	return result
}

func (p type_param_stack) SetReturnType(typ *Type) {
	p[prefixer("return")] = type_param{index: ReturnTypeOffset, typ: typ}
}

func (p type_param_stack) GetReturnType() *Type {
	return p[prefixer("return")].typ
}

const meta_prefix = "#"

func prefixer(message string, indecies ...int) string {
	if len(indecies) > 0 {
		return fmt.Sprintf("%s%s%d", meta_prefix, message, indecies[0])
	}

	return fmt.Sprintf("%s%s", meta_prefix, message)
}
