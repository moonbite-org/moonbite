package cmd

import (
	"os"

	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type FileCompiler struct {
	EventTarget
	Path        string
	SymbolTable SymbolTable
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

	return errors.EmptyError
}

func new_file_compiler(path string) FileCompiler {
	return FileCompiler{
		Path:        path,
		SymbolTable: *NewSymbolTable(),
	}
}
