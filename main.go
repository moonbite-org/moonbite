package main

import (
	"fmt"
	"os"

	"github.com/moonbite-org/moonbite/parser"
)

func main() {
	filepath := "/Users/muhammedalican/Documents/projects/moonbite/design/test.mb"
	input, _ := os.ReadFile(filepath)
	ast, err := parser.Parse(input, filepath)

	if err.Exists {
		panic(err)
	}

	fmt.Println(ast)
	// analyzer := compiler.Analyzer{Program: ast}
	// fmt.Printf("%+v\n\n", analyzer)

	// analyzer.Analyze()

}
