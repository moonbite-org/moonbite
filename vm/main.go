package main

import (
	"fmt"

	compiler "github.com/moonbite-org/moonbite/compiler/cmd"
)

func main() {
	c := compiler.New("/Users/muhammedalican/Documents/projects/a-star")
	if err := c.Compile(); err.Exists {
		fmt.Println(err)
	}
}
