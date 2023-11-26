package compiler

import "syscall"

func generate_builtins(mod *Module) map[int]interface{} {
	mod.Builtins["Syscall"] = mod.next_pointer()
	mod.Builtins["Syscall.Read"] = mod.next_pointer()
	mod.Builtins["Syscall.Write"] = mod.next_pointer()

	mod.Save(mod.Builtins["Syscall"], map[string]interface{}{
		"read":  mod.Builtins["Syscall.Read"],
		"write": mod.Builtins["Syscall.Write"],
	})

	return map[int]interface{}{
		mod.Builtins["Syscall.Read"]: func(fd int, p []byte) int {
			n, _ := syscall.Read(fd, p)
			return n
		},
		mod.Builtins["Syscall.Write"]: func(fd int, p []byte) int {
			n, _ := syscall.Write(fd, p)
			return n
		},
	}
}
