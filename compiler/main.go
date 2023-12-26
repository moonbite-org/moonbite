package main

import (
	"encoding/json"
	"os"

	compiler "github.com/moonbite-org/moonbite/compiler/cmd"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("no input provided\n")
		os.Exit(1)
	}

	c := compiler.New(os.Args[1])
	if err := c.Compile(); err.Exists {
		message, _ := json.Marshal(err)
		os.Stderr.Write(message)
		os.Exit(0)
	} else {
		os.Stdout.WriteString("{}")
	}
}
