package main

import (
	"os"

	formatter "github.com/moonbite-org/moonbite/formatter/cmd"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func main() {
	if len(os.Args) < 2 {
		os.Stdout.WriteString("no file provided.")
		os.Exit(1)
	}

	input, err := os.ReadFile(os.Args[1])

	if err != nil {
		os.Stdout.WriteString(err.Error())
		os.Exit(1)
	}

	ast, p_err := parser.Parse(input, os.Args[1])
	if p_err.Exists {
		os.Stdout.WriteString(p_err.String())
		os.Exit(1)
	}

	os.Stdout.WriteString(formatter.Format(ast))
}
