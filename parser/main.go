package main

import (
	"encoding/base64"
	"encoding/json"
	"os"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("no input provided\n")
		os.Exit(1)
	}

	var input []byte
	var file_path string

	if os.Args[1] == "--file" {
		if len(os.Args) < 3 {
			os.Stderr.WriteString("not enough arguments\n")
			os.Exit(1)
		}

		var err error
		input, err = os.ReadFile(os.Args[2])
		file_path = os.Args[2]

		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	} else {
		if len(os.Args) < 3 {
			os.Stderr.WriteString("not enough arguments\n")
			os.Exit(1)
		}

		var err error
		input, err = base64.RawStdEncoding.DecodeString(os.Args[1])
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		file_path = os.Args[2]
	}

	ast, p_err := parser.Parse(input, file_path)

	if p_err.Exists {
		message, _ := json.Marshal(p_err)
		os.Stderr.Write(message)
		os.Exit(0)
	}

	data, err := json.Marshal(ast)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	os.Stdout.WriteString(string(data))
}
