package common

import (
	"fmt"
	"strings"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type ObjectKind string

const (
	StringObjectKind ObjectKind = "object:string"
	BoolObjectKind   ObjectKind = "object:bool"
	Uint8ObjectKind  ObjectKind = "object:uint8"
	Uint16ObjectKind ObjectKind = "object:uint16"
	Uint32ObjectKind ObjectKind = "object:uint32"
	Uint64ObjectKind ObjectKind = "object:uint64"
	Int8ObjectKind   ObjectKind = "object:int8"
	Int16ObjectKind  ObjectKind = "object:int16"
	Int32ObjectKind  ObjectKind = "object:int32"
	Int64ObjectKind  ObjectKind = "object:int64"
	ListObjectKind   ObjectKind = "object:list"
	MapObjectKind    ObjectKind = "object:map"
	FunObjectKind    ObjectKind = "object:fun"
)

type Object interface {
	Kind() ObjectKind
	GetValue() interface{}
}

type StringObject struct {
	Value string
}

func (o StringObject) Kind() ObjectKind {
	return StringObjectKind
}

func (o StringObject) GetValue() interface{} {
	return o.Value
}

type BoolObject struct {
	Value bool
}

func (o BoolObject) Kind() ObjectKind {
	return BoolObjectKind
}

func (o BoolObject) GetValue() interface{} {
	return o.Value
}

type Uint8Object struct {
	Value uint8
}

func (o Uint8Object) Kind() ObjectKind {
	return Uint8ObjectKind
}

func (o Uint8Object) GetValue() interface{} {
	return o.Value
}

type Uint16Object struct {
	Value uint16
}

func (o Uint16Object) Kind() ObjectKind {
	return Uint16ObjectKind
}

func (o Uint16Object) GetValue() interface{} {
	return o.Value
}

type Uint32Object struct {
	Value uint32
}

func (o Uint32Object) Kind() ObjectKind {
	return Uint32ObjectKind
}

func (o Uint32Object) GetValue() interface{} {
	return o.Value
}

type Uint64Object struct {
	Value uint64
}

func (o Uint64Object) Kind() ObjectKind {
	return Uint64ObjectKind
}

func (o Uint64Object) GetValue() interface{} {
	return o.Value
}

type Int8Object struct {
	Value int8
}

func (o Int8Object) Kind() ObjectKind {
	return Int8ObjectKind
}

func (o Int8Object) GetValue() interface{} {
	return o.Value
}

type Int16Object struct {
	Value int16
}

func (o Int16Object) Kind() ObjectKind {
	return Int16ObjectKind
}

func (o Int16Object) GetValue() interface{} {
	return o.Value
}

type Int32Object struct {
	Value int32
}

func (o Int32Object) Kind() ObjectKind {
	return Int32ObjectKind
}

func (o Int32Object) GetValue() interface{} {
	return o.Value
}

type Int64Object struct {
	Value int64
}

func (o Int64Object) Kind() ObjectKind {
	return Int64ObjectKind
}

func (o Int64Object) GetValue() interface{} {
	return o.Value
}

type ListObject struct {
	Value []Object
}

func (o ListObject) Kind() ObjectKind {
	return ListObjectKind
}

func (o ListObject) GetValue() interface{} {
	return o.Value
}

type MapObject struct {
	Value []struct {
		Key   Object
		Value Object
	}
}

func (o MapObject) Kind() ObjectKind {
	return MapObjectKind
}

func (o MapObject) GetValue() interface{} {
	return o.Value
}

type FunctionObject struct {
	Value InstructionSet
}

func (o FunctionObject) Kind() ObjectKind {
	return FunObjectKind
}

func (o FunctionObject) GetValue() interface{} {
	return o.Value
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

type ConstantPool struct {
	Values  [1024]Object
	pointer int
}

func (p *ConstantPool) Add(value Object) int {
	has := p.Has(value)

	if has >= 0 {
		return has
	}

	p.Values[p.pointer] = value
	p.pointer++

	return p.pointer - 1
}

func (p *ConstantPool) Get(index int) Object {
	return p.Values[index]
}

func (p ConstantPool) Has(value Object) int {
	for i, v := range p.Values {
		if v == nil {
			return -1
		}
		if v.GetValue() == value.GetValue() {
			return i
		}
	}

	return -1
}

func (p ConstantPool) String() string {
	result := []string{}

	for i, constant := range p.Values {
		if constant == nil {
			break
		}

		result = append(result, fmt.Sprintf("%d: %+v", i, constant))
	}

	return fmt.Sprintf("[%s]", strings.Join(result, " "))
}
