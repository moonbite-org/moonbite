package main

import (
	"fmt"
	"os"

	"github.com/moonbite-org/moonbite/compiler"
)

func main() {
	// p := "/Users/muhammedalican/Documents/projects/moonbite/design/test.mb"
	// input, _ := os.ReadFile(p)

	// start := time.Now()
	// ast, err := parser.Parse(input, p)
	// fmt.Println(time.Since(start))

	// if err.Exists {
	// 	fmt.Println("There is an error", err)
	// }

	// data, _ := json.MarshalIndent(&ast, "", "  ")
	// fmt.Println(string(data))

	bin := compiler.Compiler{}

	mod := compiler.NewModule("main", true)
	mod.Save(10, map[string]interface{}{
		"read":  11,
		"write": 12,
	})
	p1 := mod.Save(mod.NextPointer(), "Hello ")
	p2 := mod.Save(mod.NextPointer(), "World")
	mod.Call(12, []int{p1, p2})
	bin.RegisterModule(mod)

	result, err := bin.Compile()

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	os.WriteFile("mbin", result, 0644)
}
