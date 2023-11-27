package vm

import "fmt"

type Stack struct {
	Scopes []Scope
}

func NewStack() Stack {
	stack := Stack{
		Scopes: []Scope{},
	}

	stack.PushScope("__root__")

	return stack
}

func (s *Stack) PushScope(name string) {
	depth := 0

	if len(s.Scopes) != 0 {
		depth = s.CurrentScope().Depth + 1
	}

	s.Scopes = append(s.Scopes, Scope{
		Name:  name,
		Depth: depth,
		Mem:   map[int][]byte{},
		Args:  map[int][]byte{},
	})
}

func (s *Stack) PopScope() {
	if len(s.Scopes) == 0 {
		return
	}

	s.Scopes = s.Scopes[0 : len(s.Scopes)-1]
}

func (s Stack) CurrentScope() *Scope {
	if len(s.Scopes) == 0 {
		s.PushScope("__root__")
	}

	return &s.Scopes[len(s.Scopes)-1]
}

func (s Stack) ScopeAt(depth int) *Scope {
	if len(s.Scopes) < depth {
		return nil
	}
	return &s.Scopes[len(s.Scopes)-depth]
}

func (s Stack) Save(pointer int, value []byte) {
	s.CurrentScope().Save(pointer, value)
}

func (s Stack) Load(pointer int, source source_type) ([]byte, error) {
	immediate_value := s.CurrentScope().Load(pointer, source)

	if immediate_value != nil {
		return *immediate_value, nil
	}

	if len(s.Scopes) >= 2 {
		immediate_value := s.ScopeAt(2).Load(pointer, source)

		if immediate_value != nil {
			return *immediate_value, nil
		}
	}

	value_stack := map[int][]byte{}

	for _, scope := range s.Scopes {
		value := scope.Load(pointer, source)

		if value != nil {
			value_stack[scope.Depth] = *value
		}
	}

	max := 0
	for depth := range value_stack {
		if depth > max {
			max = depth
		}
	}

	value := value_stack[max]

	if len(value) == 0 {
		return []byte{}, fmt.Errorf("value not found at pointer %d", pointer)
	}

	return value, nil
}

type Scope struct {
	Name  string
	Depth int
	Mem   map[int][]byte
	Args  map[int][]byte
}

func (s *Scope) Save(pointer int, value []byte) {
	s.Mem[pointer] = value
}

func (s Scope) Load(pointer int, source source_type) *[]byte {
	var value []byte

	switch source {
	case Memory:
		value = s.Mem[pointer]
	case Argument:
		value = s.Args[pointer]
	}

	if len(value) == 0 {
		return nil
	}

	return &value
}
