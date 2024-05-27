package hooks

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// ExecuteValidationHook executes a validation hook with the given name and value. It expects the hook to return an
// updated value written to standard out for the given parameter. No other output should be present on stdout. Standard
// error content is ignored.
func ExecuteValidationHook(hook, name, value string) (string, error) {
	out, err := exec.Command("/bin/sh", hook, name, value).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute validation hook %s on %s with value %s: %w", hook, name, value, err)
	}
	return string(out), nil
}

// ExecuteHooks executes each of the pre/post hooks in the order they were listed. If a hook exits with a non-zero exit
// code, all execution with stop and an error will be returned.
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
