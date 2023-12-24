package cmd

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
