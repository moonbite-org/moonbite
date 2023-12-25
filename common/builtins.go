package common

import "fmt"

type Builtin struct {
	Fun func(args ...Object) Object
}

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "print",
		Builtin: &Builtin{func(args ...Object) Object {
			fmt.Println("hi")

			return nil
		}},
	},
}
