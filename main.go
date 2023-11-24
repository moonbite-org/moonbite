package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CanPacis/tokenizer/formatter"
	"github.com/CanPacis/tokenizer/parser"
)

func main() {
	p := "/Users/muhammedalican/Documents/projects/moonlang/design/test.mb"
	input, _ := os.ReadFile(p)

	start := time.Now()
	ast, err := parser.Parse(input, p)
	fmt.Printf("Time: %s\n\n//%s\n", time.Since(start), ast.FileName)

	if err.Exists {
		fmt.Println("There is an error", err)
	}

	data, _ := json.Marshal(ast)
	os.WriteFile("ast.json", data, 0644)
	fmt.Println(formatter.Format(ast))
}
