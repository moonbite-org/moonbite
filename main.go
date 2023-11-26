package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/moonbite-org/moonbite/parser"
)

func main() {
	p := "/Users/muhammedalican/Documents/projects/moonbite/design/test.mb"
	input, _ := os.ReadFile(p)

	start := time.Now()
	ast, err := parser.Parse(input, p)
	fmt.Println(time.Since(start))

	if err.Exists {
		fmt.Println("There is an error", err)
	}

	data, _ := json.MarshalIndent(&ast, "", "  ")
	fmt.Println(string(data))

	// bin := compiler.Compiler{}

	// mod := compiler.NewModule("main", true)
	// p0 := mod.Save(mod.NextPointer(), 0)
	// p1 := mod.Save(mod.NextPointer(), 65)
	// mod.Call(mod.Builtins["Syscall.Write"], []int{p0, p1})
	// bin.RegisterModule(mod)

	// result, err := bin.Compile()

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(result)

	// v := vm.NewVm(result)

	// fmt.Println(v)
}
