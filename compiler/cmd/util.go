package cmd

import (
	"github.com/moonbite-org/moonbite/common"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type listener struct {
	event   string
	handler func(event Event[any])
}

type Event[T any] struct {
	Target EventTarget
	Data   T
}

type EventTarget struct {
	listeners []listener
}

func (t *EventTarget) AddListener(event string, callbackfn func(event Event[any])) {
	t.listeners = append(t.listeners, listener{event: event, handler: callbackfn})
}

func (t EventTarget) Dispatch(name string, data any) {
	event := Event[any]{
		Target: t,
		Data:   data,
	}

	for _, listener := range t.listeners {
		if listener.event == name {
			listener.handler(event)
		}
	}
}

type DummyTypeChecker struct{}

func (t DummyTypeChecker) GetDefault(typ parser.TypeLiteral) common.Object {
	switch typ.TypeKind() {
	case parser.TypeIdentifierKind:
		switch typ.(parser.TypeIdentifier).Name.(parser.IdentifierExpression).Value {
		case "string", "String":
			return common.StringObject{
				Value: "",
			}
		case "bool", "Bool":
			return common.BoolObject{
				Value: false,
			}
		case "int", "Int",
			"int8", "Int8",
			"int16", "Int16",
			"int32", "Int32",
			"int64", "Int64",
			"uint8", "Uint8",
			"uint16", "Uint16",
			"uint32", "Uint32",
			"uint64", "Uint64":
			return common.Uint32Object{
				Value: 0,
			}
		default:
			return common.StringObject{
				Value: "not implemented",
			}
		}
	default:
		return common.StringObject{
			Value: "not implemented",
		}
	}
}
