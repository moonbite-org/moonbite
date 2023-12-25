package cmd

import (
	"fmt"

	"github.com/moonbite-org/moonbite/common"
	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type PackageCompiler struct {
	EventTarget
	Definitions  []parser.Definition
	SymbolTable  *SymbolTable
	ConstantPool common.ConstantPool
	Typechecker  DummyTypeChecker
	Instructions common.InstructionSet
}

func (c PackageCompiler) Compile() errors.Error {
	for _, builtin := range common.Builtins {
		c.SymbolTable.Define(builtin.Name, parser.ConstantKind)
	}

	for _, definition := range c.Definitions {
		instructions, err := c.compile_statement(definition)
		if err.Exists {
			return err
		}
		c.Instructions = append(c.Instructions, instructions...)
		fmt.Println(instructions)
	}

	fmt.Println(c.Instructions)
	fmt.Println(c.ConstantPool)
	fmt.Println(c.SymbolTable)

	return errors.EmptyError
}

func (c PackageCompiler) GetBytes() []byte {
	result := []byte{}
	result = append(result, c.ConstantPool.Serialize()...)
	result = append(result, c.Instructions.GetBytes()...)

	return result
}

func NewPackageCompiler(definitions []parser.Definition) PackageCompiler {
	return PackageCompiler{
		Definitions: definitions,
		SymbolTable: NewSymbolTable(),
		ConstantPool: common.ConstantPool{
			Values: [1024]common.Object{},
		},
	}
}
