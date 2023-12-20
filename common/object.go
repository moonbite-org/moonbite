package common

import parser "github.com/moonbite-org/moonbite/parser/cmd"

type Object interface {
	Type() string
}

type StringObject struct {
	Value string
}

func (o StringObject) Type() string {
	return "string"
}

type BoolObject struct {
	Value bool
}

func (o BoolObject) Type() string {
	return "bool"
}

type Uint8Object struct {
	Value uint8
}

func (o Uint8Object) Type() string {
	return "uint8"
}

type Uint16Object struct {
	Value uint16
}

func (o Uint16Object) Type() string {
	return "uint16"
}

type Uint32Object struct {
	Value uint32
}

func (o Uint32Object) Type() string {
	return "uint32"
}

type Uint64Object struct {
	Value uint64
}

func (o Uint64Object) Type() string {
	return "uint64"
}

type Int8Object struct {
	Value int8
}

func (o Int8Object) Type() string {
	return "int8"
}

type Int16Object struct {
	Value int16
}

func (o Int16Object) Type() string {
	return "int16"
}

type Int32Object struct {
	Value int32
}

func (o Int32Object) Type() string {
	return "int32"
}

type Int64Object struct {
	Value int64
}

func (o Int64Object) Type() string {
	return "int64"
}

type ListObject struct {
	Value []Object
}

func (o ListObject) Type() string {
	return "list"
}

type MapObject struct {
	Value []struct {
		Key   Object
		Value Object
	}
}

func (o MapObject) Type() string {
	return "map"
}

func ObjectFromLiteral(literal parser.LiteralExpression) Object {
	switch literal.LiteralKind() {
	case parser.StringLiteralKind:
		return StringObject{Value: literal.(parser.StringLiteralExpression).Value}
	case parser.BoolLiteralKind:
		return BoolObject{Value: literal.(parser.BoolLiteralExpression).Value}
	case parser.NumberLiteralKind:
		return StringObject{Value: "not implemented"}
	default:
		return StringObject{Value: "not implemented"}
	}
}
