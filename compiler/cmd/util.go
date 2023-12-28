package cmd

import (
	"os"

	"github.com/moonbite-org/moonbite/common"
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

type DummyTypeChecker struct{}

func (t DummyTypeChecker) GetDefault(typ parser.TypeLiteral) common.Object {
	switch typ.TypeKind() {
	case parser.TypeIdentifierKind:
		switch typ.(parser.TypeIdentifier).Name.(parser.IdentifierExpression).Value {
		case "string", "String":
			return common.StringObject{
				Value: "",
			}
		case "bool", "Bool":
			return common.BoolObject{
				Value: false,
			}
		case "int", "Int",
			"int8", "Int8",
			"int16", "Int16",
			"int32", "Int32",
			"int64", "Int64",
			"uint8", "Uint8",
			"uint16", "Uint16",
			"uint32", "Uint32",
			"uint64", "Uint64":
			return common.Uint32Object{
				Value: 0,
			}
		default:
			return common.StringObject{
				Value: "not implemented",
			}
		}
	default:
		return common.StringObject{
			Value: "not implemented",
		}
	}
}
