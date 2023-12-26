package common

import (
	"os"
)

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
			return nil
		}},
	},
	{
		Name: "exit",
		Builtin: &Builtin{func(args ...Object) Object {
			os.Exit(int(args[0].(Int32Object).Value))

			return nil
		}},
	},
}
