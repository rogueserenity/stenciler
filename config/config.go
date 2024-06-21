package config

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

// HookClass is an enumeration type for the different classes of hooks that can be run.
type HookClass int

const (
	// PreInitHook is a hook that runs before initializing the repository.
	PreInitHook HookClass = iota
	// PostInitHook is a hook that runs after initializing the repository.
	PostInitHook
	// PreUpdateHook is a hook that runs before updating the repository.
	PreUpdateHook
	// PostUpdateHook is a hook that runs after updating the repository.
	PostUpdateHook
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
	Params []*Param `yaml:"params,omitempty"`

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
	Templates []*Template `yaml:"templates,omitempty"`
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

// LogValue returns the slog.Value representation of the Config.
func (c *Config) LogValue() slog.Value {
	return slog.AnyValue(c.Templates)
}

// LogValue returns the slog.Value representation of the Template.
func (t *Template) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("repository", t.Repository),
		slog.String("directory", t.Directory),
		slog.Any("params", t.Params),
		slog.Any("init-only", t.InitOnlyPaths),
		slog.Any("raw-copy", t.RawCopyPaths),
		slog.Any("pre-init-hooks", t.PreInitHookPaths),
		slog.Any("post-init-hooks", t.PostInitHookPaths),
		slog.Any("pre-update-hooks", t.PreUpdateHookPaths),
		slog.Any("post-update-hooks", t.PostUpdateHookPaths),
	)
}

// LogValue returns the slog.Value representation of the Param.
func (p *Param) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", p.Name),
		slog.String("prompt", p.Prompt),
		slog.String("default", p.Default),
		slog.String("validation-hook", p.ValidationHook),
		slog.String("value", p.Value),
	)
}

// Validate runs the validation hook for the parameter and updates the value with the output of the hook. If there is no
// validation hook, then the value is unchanged and the function returns no error. It expects the hook to return an
// updated value written to standard out for the given parameter. No other output should be present on stdout. Standard
// error content is ignored.
func (p *Param) Validate(repoDir string) error {
	if len(p.ValidationHook) == 0 {
		return nil
	}
	hook := filepath.Join(repoDir, p.ValidationHook)
	out, err := exec.Command("/bin/sh", hook, p.Name, p.Value).Output()
	if err != nil {
		return fmt.Errorf("failed to execute validation hook %s on %s with value %s: %w",
			p.ValidationHook, p.Name, p.Value, err)
	}
	p.Value = strings.TrimSpace(string(out))

	return nil
}

// ExecuteHooks executes each of the pre/post hooks in the order they were listed. If a hook exits with a non-zero exit
// code, all execution with stop and an error will be returned.
func (t *Template) ExecuteHooks(repoDir string, hookClass HookClass) error {
	var hooks []string
	switch hookClass {
	case PreInitHook:
		hooks = t.PreInitHookPaths
	case PostInitHook:
		hooks = t.PostInitHookPaths
	case PreUpdateHook:
		hooks = t.PreUpdateHookPaths
	case PostUpdateHook:
		hooks = t.PostUpdateHookPaths
	default:
		return fmt.Errorf("unknown hook class %d", hookClass)
	}

	for _, hook := range hooks {
		if err := t.executeHook(filepath.Join(repoDir, hook)); err != nil {
			return err
		}
	}

	return nil
}

func (t *Template) executeHook(hook string) error {
	cmd := exec.Command("/bin/sh", hook)

	env := os.Environ()
	for _, p := range t.Params {
		name := strcase.ToScreamingSnake(p.Name)
		env = append(env, fmt.Sprintf("STENCILER_%s=%s", name, p.Value))
	}
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute hook %s: %w", hook, err)
	}

	return nil
}
