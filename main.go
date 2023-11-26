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
}
