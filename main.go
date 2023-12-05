package main

import (
	"fmt"
	"os"

	"github.com/moonbite-org/moonbite/parser"
)

func main() {
	input := []byte(`package main
	const test = read_file or giveup
	`)
	ast, err := parser.Parse(input, "test.mb")

	if err.Exists {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", ast.Definitions[0])
}
