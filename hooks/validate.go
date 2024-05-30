package hooks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rogueserenity/stenciler/config"
)

// Validate validates all the hooks in the template exist and are executable.
func Validate(template config.Template, repoPath string) error {
	var errs []error
	hookPaths := gatherHookPaths(template)
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

func gatherHookPaths(template config.Template) []string {
	var hookPaths []string

	for _, param := range template.Params {
		if param.ValidationHook != "" {
			hookPaths = append(hookPaths, param.ValidationHook)
		}
	}

	hookPaths = append(hookPaths, template.PreInitHookPaths...)
	hookPaths = append(hookPaths, template.PostInitHookPaths...)
	hookPaths = append(hookPaths, template.PreUpdateHookPaths...)
	hookPaths = append(hookPaths, template.PostUpdateHookPaths...)

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
