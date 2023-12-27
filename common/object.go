package common

import (
	"fmt"
	"reflect"
	"strings"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type ObjectKind string

const (
	StringObjectKind   ObjectKind = "object:string"
	RuneObjectKind     ObjectKind = "object:rune"
	BoolObjectKind     ObjectKind = "object:bool"
	Uint8ObjectKind    ObjectKind = "object:uint8"
	Uint16ObjectKind   ObjectKind = "object:uint16"
	Uint32ObjectKind   ObjectKind = "object:uint32"
	Uint64ObjectKind   ObjectKind = "object:uint64"
	Int8ObjectKind     ObjectKind = "object:int8"
	Int16ObjectKind    ObjectKind = "object:int16"
	Int32ObjectKind    ObjectKind = "object:int32"
	Int64ObjectKind    ObjectKind = "object:int64"
	Float32ObjectKind  ObjectKind = "object:float32"
	Float64ObjectKind  ObjectKind = "object:float64"
	ListObjectKind     ObjectKind = "object:list"
	MapObjectKind      ObjectKind = "object:map"
	InstanceObjectKind ObjectKind = "object:instance"
	FunObjectKind      ObjectKind = "object:fun"
	terminator_kind    ObjectKind = "object:terminator"
	pool_block_kind    ObjectKind = "object:pool"
)

var type_map = map[ObjectKind]byte{
	StringObjectKind:   12,
	BoolObjectKind:     13,
	Uint8ObjectKind:    15,
	Uint16ObjectKind:   16,
	Uint32ObjectKind:   17,
	Uint64ObjectKind:   18,
	Int8ObjectKind:     19,
	Int16ObjectKind:    20,
	Int32ObjectKind:    21,
	Int64ObjectKind:    22,
	ListObjectKind:     23,
	MapObjectKind:      24,
	InstanceObjectKind: 25,
	FunObjectKind:      26,
	pool_block_kind:    126,
	terminator_kind:    0,
}

type Object interface {
	Kind() ObjectKind
	GetValue() interface{}
	Serialize() []byte
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

func (o StringObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, []byte(o.Value)...)
	result = append(result, type_map[terminator_kind])

	return result
}

type RuneObject struct {
	Value rune
}

func (o RuneObject) Kind() ObjectKind {
	return RuneObjectKind
}

func (o RuneObject) GetValue() interface{} {
	return o.Value
}

func (o RuneObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, []byte(string(o.Value))...)
	result = append(result, type_map[terminator_kind])

	return result
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

func (o BoolObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}

	if o.Value {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}

	return result
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

func (o Uint8Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Uint16Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Uint32Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Uint64Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Int8Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Int16Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Int32Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o Int64Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
}

type Float32Object struct {
	Value float32
}

func (o Float32Object) Kind() ObjectKind {
	return Float32ObjectKind
}

func (o Float32Object) GetValue() interface{} {
	return o.Value
}

func (o Float32Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
}

type Float64Object struct {
	Value float32
}

func (o Float64Object) Kind() ObjectKind {
	return Float64ObjectKind
}

func (o Float64Object) GetValue() interface{} {
	return o.Value
}

func (o Float64Object) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	result = append(result, NumberToBytes(o.Value)...)
	return result
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

func (o ListObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}

	for _, value := range o.Value {
		result = append(result, value.Serialize()...)
	}

	result = append(result, type_map[terminator_kind])

	return result
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

func (o MapObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}

	for _, entry := range o.Value {
		result = append(result, entry.Key.Serialize()...)
		result = append(result, entry.Value.Serialize()...)
	}

	result = append(result, type_map[terminator_kind])

	return result
}

type InstanceObject struct {
	Value []struct {
		Key   Object
		Value Object
	}
}

func (o InstanceObject) Kind() ObjectKind {
	return InstanceObjectKind
}

func (o InstanceObject) GetValue() interface{} {
	return o.Value
}

func (o InstanceObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}

	for _, entry := range o.Value {
		result = append(result, entry.Key.Serialize()...)
		result = append(result, entry.Value.Serialize()...)
	}

	result = append(result, type_map[terminator_kind])

	return result
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

func (o FunctionObject) Serialize() []byte {
	result := []byte{type_map[o.Kind()]}
	value := o.Value.GetBytes()
	result = append(result, NumberToBytes(int32(len(value)))...)
	result = append(result, value...)
	return result
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

func (p ConstantPool) Serialize() []byte {
	result := []byte{type_map[pool_block_kind]}

	for i, constant := range p.Values {
		if constant == nil {
			break
		}

		result = append(result, NumberToBytes(int32(i))...)
		result = append(result, constant.Serialize()...)
	}

	result = append(result, type_map[terminator_kind])

	return result
}

func (p ConstantPool) Has(value Object) int {
	for i, v := range p.Values {
		if v == nil {
			return -1
		}

		if reflect.TypeOf(v.GetValue()) == reflect.TypeOf(InstructionSet{}) {
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
