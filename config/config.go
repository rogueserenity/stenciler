package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Param holds all of the values for a parameter. THe validation hook path is relative to the repository root.
type Param struct {
	Name string `yaml:"name"`

	Prompt         string `yaml:"prompt"`
	Default        string   `yaml:"default"`
	ValidationHook string `yaml:"validation-hook"`

	Value string `yaml:"value"`
}

// Template holds all of the values for a template configuration. The paths defined by init-only and raw-copy are
// relative to directory. The directory and hook paths are all relative to the repository root.
type Template struct {
	Repository string `yaml:"repository"`
	Directory  string `yaml:"directory"`

	Params []Param `yaml:"params"`

	InitOnlyPaths []string `yaml:"init-only"`
	RawCopyPaths  []string `yaml:"raw-copy"`

	PreInitHookPaths    []string `yaml:"pre-init-hooks"`
	PostInitHookPaths   []string `yaml:"post-init-hooks"`
	PreUpdateHookPaths  []string `yaml:"pre-update-hooks"`
	PostUpdateHookPaths []string `yaml:"post-update-hooks"`
}

// Config holds the contents of a configuration file
type Config struct {
	Templates []Template `yaml:"templates"`
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
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	return c.Write(file)
}

// Write attempts to write the config to the specified writer
func (c *Config) Write(out io.Writer) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	return err
}
