package main

import (
	"fmt"
	"os"

	compiler "github.com/moonbite-org/moonbite/compiler/cmd"
)

func main() {
	c := compiler.New("/Users/muhammedalican/Documents/projects/a-star")
	if err := c.Compile(); err.Exists {
		fmt.Println(err)
	}

	root := c.Modules["root"].Compiler
	fmt.Println(root.Instructions)
	fmt.Println(root.ConstantPool)
	fmt.Println(root.SymbolTable)
	fmt.Println(root.GetBytes())
	os.WriteFile("/Users/muhammedalican/Documents/projects/a-star/main.mbin", root.GetBytes(), 0644)
}
