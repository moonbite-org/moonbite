package main

import (
	"fmt"
	"os"

	"github.com/moonbite-org/moonbite/parser"
)

func main() {
	input := []byte(`package main
	fun main() {
		index++ 
		++index --index
		2
		2 + 2
		2 + 2 * 3
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	if err.Exists {
		fmt.Println(err)
		os.Exit(1)
	}

	f := ast.Definitions[0].(*parser.UnboundFunDefinitionStatement).Body
	fmt.Println(f)
	// m := f[0].(parser.ExpressionStatement).Expression.(parser.ArithmeticUnaryExpression)

	// fmt.Println(m)
}
