package hooks

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func ExecuteValidationHook(hook, name, value string) (string, error) {
	out, err := exec.Command("/bin/sh", hook, name, value).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute validation hook %s on %s with value %s: %w", hook, name, value, err)
	}
	return string(out), nil
}

func ExecuteHooks(repoDir string, hooks []string) error {
	for _, hook := range hooks {
		if err := executeHook(filepath.Join(repoDir, hook)); err != nil {
			return err
		}
	}
	return nil
}

func executeHook(hook string) error {
	err := exec.Command("/bin/sh", hook).Run()
	if err != nil {
		return fmt.Errorf("failed to execute hook %s: %w", hook, err)
	}
	return nil
}
