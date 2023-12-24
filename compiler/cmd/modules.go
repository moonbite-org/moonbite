package cmd

import (
	"fmt"
	"os"
	"path"
	"slices"

	errors "github.com/moonbite-org/moonbite/error"
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
	Modules map[string]Module
}

func New(dir string) Compiler {
	return Compiler{
		RootDir: dir,
		Modules: map[string]Module{},
	}
}

var allowed_extensions = []string{".mb"}

func resolve_dir(dir string, is_root bool) (map[string]Module, error) {
	result := map[string]Module{}
	mod := Module{
		Dir: dir,
	}
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entry_path := path.Join(dir, entry.Name())

		if entry.Type().IsDir() {
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
		result["root"] = mod
	} else {
		result[path.Base(dir)] = mod
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

	if len(modules["root"].FilePaths) == 0 {
		return errors.CreateAnonError(errors.CompileError, "no root module found")
	}

	return modules["root"].Compile()
}

type Module struct {
	PackageName string
	Dir         string
	FilePaths   []string
}

func (m Module) Compile() errors.Error {
	var multi_package_error *string

	for _, file_path := range m.FilePaths {
		if multi_package_error != nil {
			break
		}
		compiler := new_file_compiler(file_path)

		compiler.AddListener("announce:package", func(event Event[any]) {
			name := event.Data.(string)
			if len(m.PackageName) == 0 {
				m.PackageName = name
			} else {
				if m.PackageName != name {
					multi_package_error = &name
				}
			}
		})

		if err := compiler.Compile(); err.Exists {
			return err
		}
	}

	if multi_package_error != nil {
		return errors.CreateAnonError(errors.CompileError, fmt.Sprintf("found multiple packages in this directory '%s', '%s' %s", *multi_package_error, m.PackageName, m.Dir))
	}

	return errors.EmptyError
}
