package cmd

import (
	"github.com/moonbite-org/moonbite/common"
	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type PackageCompiler struct {
	IsRoot               bool
	Definitions          []parser.Definition
	SymbolTable          *SymbolTable
	ConstantPool         common.ConstantPool
	Typechecker          DummyTypeChecker
	Instructions         common.InstructionSet
	current_match_target common.InstructionSet
}

func (c *PackageCompiler) Compile() errors.Error {
	for _, builtin := range common.Builtins {
		c.SymbolTable.Define(builtin.Name, parser.ConstantKind, false)
	}

	for _, definition := range c.Definitions {
		instructions, err := c.compile_statement(definition)
		if err.Exists {
			return err
		}
		c.Instructions = append(c.Instructions, instructions...)
	}

	if c.IsRoot {
		instructions, err := c.compile_call_expression(parser.CallExpression{
			Callee:    parser.IdentifierExpression{Value: "main"},
			Arguments: []parser.Expression{},
		})
		if err.Exists {
			return err
		}

		c.Instructions = append(c.Instructions, instructions...)
	}

	return errors.EmptyError
}

func (c PackageCompiler) GetBytes() []byte {
	result := []byte{}
	result = append(result, c.ConstantPool.Serialize()...)
	result = append(result, c.Instructions.GetBytes()...)

	return result
}

func NewPackageCompiler(definitions []parser.Definition, is_root bool) PackageCompiler {
	return PackageCompiler{
		Definitions: definitions,
		SymbolTable: NewSymbolTable(),
		ConstantPool: common.ConstantPool{
			Values: [1024]common.Object{},
		},
		IsRoot: is_root,
	}
}
