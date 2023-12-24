package cmd

import (
	"fmt"
	"os"

	"github.com/moonbite-org/moonbite/common"
	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type FileCompiler struct {
	EventTarget
	Path          string
	SymbolTable   SymbolTable
	ConstantPool  common.ConstantPool
	Typechecker   DummyTypeChecker
	Instructions  common.InstructionSet
	current_scope SymbolScope
	last_scope    SymbolScope
}

func (c FileCompiler) Compile() errors.Error {
	program, err := os.ReadFile(c.Path)
	if err != nil {
		return errors.CreateAnonError(errors.CompileError, err.Error())
	}

	ast, p_err := parser.Parse(program, c.Path)
	if p_err.Exists {
		return p_err
	}

	c.Dispatch("announce:package", ast.Package.Name.Value)

	for _, definition := range ast.Definitions {
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

func (c FileCompiler) GetBytes() []byte {
	return c.Instructions.GetBytes()
}

func new_file_compiler(path string) FileCompiler {
	return FileCompiler{
		Path:        path,
		SymbolTable: *NewSymbolTable(),
		ConstantPool: common.ConstantPool{
			Values: [1024]common.Object{},
		},
		current_scope: GlobalScope,
		last_scope:    GlobalScope,
	}
}
