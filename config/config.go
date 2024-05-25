package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Param holds all of the values for a parameter. THe validation hook path is relative to the repository root.
type Param struct {
	Name string `yaml:"name,omitempty"`

	Prompt         string `yaml:"prompt,omitempty"`
	Default        string `yaml:"default,omitempty"`
	ValidationHook string `yaml:"validation-hook,omitempty"`

	Value string `yaml:"value,omitempty"`
}

// Template holds all of the values for a template configuration. The paths defined by init-only and raw-copy are
// relative to directory. The directory and hook paths are all relative to the repository root.
type Template struct {
	Repository string `yaml:"repository"`
	Directory  string `yaml:"directory"`

	Params []Param `yaml:"params,omitempty"`

	InitOnlyPaths []string `yaml:"init-only,omitempty"`
	RawCopyPaths  []string `yaml:"raw-copy,omitempty"`

	PreInitHookPaths    []string `yaml:"pre-init-hooks,omitempty"`
	PostInitHookPaths   []string `yaml:"post-init-hooks,omitempty"`
	PreUpdateHookPaths  []string `yaml:"pre-update-hooks,omitempty"`
	PostUpdateHookPaths []string `yaml:"post-update-hooks,omitempty"`
}

// Config holds the contents of a configuration file
type Config struct {
	Templates []Template `yaml:"templates,omitempty"`
}

// ReadFromFile attempts to read a config from the specified path
func ReadFromFile(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	return Read(file)
}

// Read attempts to read a config from the specified reader
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

// WriteToFile attempts to write the config out to the specified path
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

// Write attempts to write the config to the specified writer
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
