package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Clone clones a repository and returns the path to the cloned repository
func Clone(url, authToken string) (string, error) {
	path, err := os.MkdirTemp("", "stenciler-clone-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	opts := &git.CloneOptions{
		URL: url,
	}
	if len(authToken) > 0 {
		opts.Auth = &http.BasicAuth{
			Username: "token",
			Password: authToken,
		}
	}

	_, err = git.PlainClone(path, false, opts)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return path, nil
}
