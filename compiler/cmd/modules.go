package cmd

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
	"gopkg.in/yaml.v3"
)

type Dependency struct {
	Name    string
	Version string
}

type WorkspaceEntry struct {
	Path    string
	Module  string
	Version string
	Deps    []Dependency
	DevDeps []Dependency
}

type ModConfig struct {
	Module    string
	Moonbite  string
	Version   string
	Deps      []Dependency
	DevDeps   []Dependency
	Workspace []WorkspaceEntry
}

func ParseConfig(config_path string) (ModConfig, error) {
	config_file, err := os.ReadFile(config_path)
	if err != nil {
		return ModConfig{}, err
	}
	config := ModConfig{}

	err = yaml.Unmarshal(config_file, &config)
	if err != nil {
		return ModConfig{}, err
	}

	return config, nil
}

type Compiler struct {
	RootDir string
	Config  ModConfig
	Modules map[string]*Module
}

func New(dir string) Compiler {
	return Compiler{
		RootDir: dir,
		Modules: map[string]*Module{},
	}
}

var allowed_extensions = []string{".mb"}

func resolve_dir(dir string, is_root bool) (map[string]*Module, error) {
	result := map[string]*Module{}
	mod := Module{
		Dir:    dir,
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

			sub_modules, err := resolve_dir(entry_path, false)

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

	modules, err := resolve_dir(c.RootDir, true)

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
	Dir         string
	FilePaths   []string
	IsRoot      bool
	Compiler    PackageCompiler
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

	m.Compiler = NewPackageCompiler(definitions, m.IsRoot)

	if err := m.Compiler.Compile(); err.Exists {
		return err
	}

	return errors.EmptyError
}
