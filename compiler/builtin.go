package compiler

import "syscall"

var Builtin = map[int8]interface{}{
	11: func(fd int, p []byte) int {
		n, _ := syscall.Read(fd, p)
		return n
	},
	12: func(fd int, p []byte) int {
		n, _ := syscall.Write(fd, p)
		return n
	},
}
