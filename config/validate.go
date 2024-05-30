package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Validate validates all the hooks in the template exist and are executable.
func (t *Template) Validate(repoPath string) error {
	var errs []error
	hookPaths := t.gatherHookPaths()
	for _, hookPath := range hookPaths {
		if err := validateHook(hookPath, repoPath); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (t *Template) gatherHookPaths() []string {
	var hookPaths []string

	for _, param := range t.Params {
		if param.ValidationHook != "" {
			hookPaths = append(hookPaths, param.ValidationHook)
		}
	}

	hookPaths = append(hookPaths, t.PreInitHookPaths...)
	hookPaths = append(hookPaths, t.PostInitHookPaths...)
	hookPaths = append(hookPaths, t.PreUpdateHookPaths...)
	hookPaths = append(hookPaths, t.PostUpdateHookPaths...)

	return hookPaths
}

func validateHook(hookPath, repoPath string) error {
	path := filepath.Join(repoPath, hookPath)

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("hook %s does not exist", hookPath)
	}

	mode := info.Mode()
	if !mode.IsRegular() || mode.Perm()&0111 == 0 {
		return fmt.Errorf("hook %s is not executable", hookPath)
	}

	return nil
}
