package abi

import (
	"github.com/moonbite-org/moonbite/common"
)

type Builtin interface {
	Name() string
}

type BuiltinFun struct {
	name string
	Fun  func(params ...common.Object) common.Object
}

func (b BuiltinFun) Name() string {
	return b.name
}

type BuiltinMap struct {
	name  string
	Value common.MapObject
}

func CreateBuiltinMap(name string, data map[string]Builtin) BuiltinMap {
	return BuiltinMap{
		name:  name,
		Value: common.MapObject{},
	}
}

func (b BuiltinMap) Name() string {
	return b.name
}

type ABI struct {
	Builtins []Builtin
}
