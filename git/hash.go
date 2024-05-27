package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

// CommitHash returns the git commit hash of the HEAD of repository located at the provided path. Returns an error if
// it cannot locate a git repository or the path or cannot find the HEAD of the repository.
func CommitHash(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	return head.Hash().String(), nil
}
