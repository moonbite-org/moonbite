package abi

import "github.com/moonbite-org/moonbite/common"

var NativeABI = ABI{
	Builtins: []Builtin{
		CreateBuiltinMap("syscall", map[string]Builtin{
			"write": BuiltinFun{
				name: "write",
				Fun: func(params ...common.Object) common.Object {
					return nil
				},
			},
		}),
	},
}
