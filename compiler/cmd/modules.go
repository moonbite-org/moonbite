package cmd

import (
	"os"

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
