package cmd

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/moonbite-org/moonbite/abi"
	"github.com/moonbite-org/moonbite/common"
	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

type Compiler struct {
	RootDir string
	Config  ModConfig
	Modules map[string]*Module
	ABI     abi.ABI
}

func New(dir string, interface_ abi.ABI) Compiler {
	return Compiler{
		RootDir: dir,
		Modules: map[string]*Module{},
		ABI:     interface_,
	}
}

var allowed_extensions = []string{".mb"}

func resolve_dir(dir string, is_root bool, interface_ abi.ABI) (map[string]*Module, error) {
	result := map[string]*Module{}
	mod := Module{
		Dir:    dir,
		ABI:    interface_,
		IsRoot: is_root,
	}
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entry_path := path.Join(dir, entry.Name())

		if entry.Type().IsDir() {
			if strings.HasPrefix(entry_path, ".") {
				continue
			}

			sub_modules, err := resolve_dir(entry_path, false, interface_)

			if err != nil {
				return result, err
			}

			for name, mod := range sub_modules {
				result[name] = mod
			}
		} else {
			extension := path.Ext(entry_path)

			if slices.Contains(allowed_extensions, extension) {
				mod.FilePaths = append(mod.FilePaths, entry_path)
			}
		}
	}

	if is_root {
		result["root"] = &mod
	} else {
		result[path.Base(dir)] = &mod
	}

	return result, nil
}

func (c *Compiler) Compile() errors.Error {
	config_path := path.Join(c.RootDir, "moon.yml")
	_, err := os.Stat(config_path)

	if err != nil {
		return errors.CreateAnonError(errors.CompileError, "A moon.yml config file is required to compile your module")
	}

	mod_config, err := ParseConfig(config_path)

	if err != nil {
		return errors.CreateAnonError(errors.CompileError, err.Error())
	}

	c.Config = mod_config

	modules, err := resolve_dir(c.RootDir, true, c.ABI)

	if err != nil {
		return errors.CreateAnonError(errors.CompileError, err.Error())
	}

	c.Modules = modules
	root := modules["root"]

	if len(root.FilePaths) == 0 {
		return errors.CreateAnonError(errors.CompileError, "no root module found")
	}

	return root.Compile()
}

type Module struct {
	PackageName string
	ABI         abi.ABI
	Dir         string
	FilePaths   []string
	IsRoot      bool
	Compiler    package_compiler
}

func (m *Module) Compile() errors.Error {
	definitions := []parser.Definition{}

	for _, file_path := range m.FilePaths {
		program, err := os.ReadFile(file_path)
		if err != nil {
			return errors.CreateAnonError(errors.CompileError, err.Error())
		}

		ast, p_err := parser.Parse(program, file_path)
		if p_err.Exists {
			return p_err
		}

		if len(m.PackageName) == 0 {
			m.PackageName = ast.Package.Name.Value
		} else {
			if m.PackageName != ast.Package.Name.Value {
				return errors.CreateCompileError(fmt.Sprintf("found multiple packages in this directory '%s', '%s' %s", ast.Package.Name.Value, m.PackageName, m.Dir), ast.Package.Location())
			}
		}

		definitions = append(definitions, ast.Definitions...)
	}

	m.Compiler = new_package_compiler(m.PackageName, definitions, m.IsRoot, m.ABI)

	if err := m.Compiler.Compile(); err.Exists {
		return err
	}

	return errors.EmptyError
}

var builtins = []string{"exit", "#null"}

type package_compiler struct {
	package_name         string
	ABI                  abi.ABI
	IsRoot               bool
	Definitions          []parser.Definition
	SymbolTable          *SymbolTable
	TypeSymbolTable      map[string]*SymbolTable
	ConstantPool         common.ConstantPool
	Typechecker          DummyTypeChecker
	Instructions         common.InstructionSet
	current_match_target common.InstructionSet
}

func (c *package_compiler) Compile() errors.Error {
	for _, builtin := range builtins {
		c.SymbolTable.DefineBuiltin(builtin)
	}

	for _, builtin := range c.ABI.Builtins {
		c.SymbolTable.DefineBuiltin(builtin.Name())
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

func (c package_compiler) GetBytes() []byte {
	result := []byte{}
	result = append(result, c.ConstantPool.Serialize()...)
	result = append(result, c.Instructions.GetBytes()...)

	return result
}

func new_package_compiler(package_ string, definitions []parser.Definition, is_root bool, interface_ abi.ABI) package_compiler {
	return package_compiler{
		package_name:    package_,
		ABI:             interface_,
		Definitions:     definitions,
		SymbolTable:     NewSymbolTable(),
		TypeSymbolTable: map[string]*SymbolTable{},
		ConstantPool: common.ConstantPool{
			Values: [1024]common.Object{},
		},
		IsRoot:               is_root,
		current_match_target: nil,
	}
}
