package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Param holds all of the values for a parameter.
type Param struct {
	// Name is the name of the parameter. Required.
	Name string `yaml:"name"`

	// Prompt is the prompt to display to the user when initializing a new repository. Optional. If not provided, the
	// parameter is considered internal only.
	Prompt string `yaml:"prompt,omitempty"`
	// Default is the default value to use if the user does not provide one. Optional. An empty string is used as the
	// default if no default is provided and the user does not set a value.
	Default string `yaml:"default,omitempty"`
	// ValidationHook is the path to a script to run to validate the value. Optional. The path is relative to the
	// repository root.
	ValidationHook string `yaml:"validation-hook,omitempty"`

	// Value is the value of the parameter. For a template, this is ignored if Prompt is set.
	// For a repository, the value is determined by the following rules:
	// 1. If the parameter is internal only, the value is the value from the template.
	// 2. If the parameter has a prompt, the user is prompted for the value. The default is used if the user does not
	//    provide a value.
	// 3. If the parameter has a ValidationHook, then that is executed and the output is the value.
	Value string `yaml:"value,omitempty"`
}

// Template holds all of the values for a template configuration. The paths defined by init-only and raw-copy are
// relative to directory. The directory and hook paths are all relative to the repository root.
type Template struct {
	// Repository is the URL of the repository to clone. Required.
	Repository string `yaml:"repository"`
	// Directory is the directory at the root of the repository that holds the template data. Required.
	Directory string `yaml:"directory"`

	// Params is a list of parameters to prompt the user for when initializing a new repository. Optional.
	Params []Param `yaml:"params,omitempty"`

	// InitOnlyPaths are a list of glob paths that are only copied over during intialization. Optional. The glob paths
	// are relative to Directory.
	InitOnlyPaths []string `yaml:"init-only,omitempty"`
	// RawCopyPaths are a list of glob paths that are copied without being run through the template engine. Optional.
	// The glob paths are relative to Directory.
	RawCopyPaths []string `yaml:"raw-copy,omitempty"`

	// PreInitHookPaths are a list of paths to scripts to run before initializing the repository. Optional. The paths
	// are relative to the repository root. The hooks are run in the order they are defined.
	PreInitHookPaths []string `yaml:"pre-init-hooks,omitempty"`
	// PostInitHookPaths are a list of paths to scripts to run after initializing the repository. Optional. The paths
	// are relative to the repository root. The hooks are run in the order they are defined.
	PostInitHookPaths []string `yaml:"post-init-hooks,omitempty"`
	// PreUpdateHookPaths are a list of paths to scripts to run before updating the repository. Optional. The paths
	// are relative to the repository root. The hooks are run in the order they are defined.
	PreUpdateHookPaths []string `yaml:"pre-update-hooks,omitempty"`
	// PostUpdateHookPaths are a list of paths to scripts to run after updating the repository. Optional. The paths
	// are relative to the repository root. The hooks are run in the order they are defined.
	PostUpdateHookPaths []string `yaml:"post-update-hooks,omitempty"`
}

// Config holds the contents of a configuration file.
type Config struct {
	Templates []Template `yaml:"templates,omitempty"`
}

// ReadFromFile attempts to read a config from the specified path.
func ReadFromFile(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	return Read(file)
}

// Read attempts to read a config from the specified reader.
func Read(in io.Reader) (*Config, error) {
	b, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// WriteToFile attempts to write the config out to the specified path.
func (c *Config) WriteToFile(configPath string) error {
	if len(c.Templates) == 0 {
		return errors.New("unable to write empty config")
	}
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	return c.Write(file)
}

// Write attempts to write the config to the specified writer.
func (c *Config) Write(out io.Writer) error {
	if len(c.Templates) == 0 {
		return errors.New("unable to write empty config")
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	return err
}
