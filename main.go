package main

import (
	"fmt"
	"os"

	"github.com/moonbite-org/moonbite/parser"
)

func main() {
	input := []byte(`package main
	fun main() {
		match (data) {
			(.) {
				a = b
			}
		}
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	if err.Exists {
		fmt.Println(err)
		os.Exit(1)
	}

	f := ast.Definitions[0].(*parser.UnboundFunDefinitionStatement).Body
	m := f[0].(parser.ExpressionStatement).Expression.(parser.MatchExpression)

	b := m.Blocks[0].Predicate
	fmt.Println(b.Kind())
}
