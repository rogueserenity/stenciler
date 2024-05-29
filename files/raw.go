package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/rogueserenity/stenciler/config"
)

// CopyRaw copies files that do not require template processing into the current working directory.
func CopyRaw(repoDir string, template *config.Template) error {
	srcRootPath := filepath.Join(repoDir, template.Directory)

	destRootPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	copyList, err := generateRawCopyList(srcRootPath, template.RawCopyPaths)
	if err != nil {
		return fmt.Errorf("failed to generate copy list: %w", err)
	}

	for _, f := range copyList {
		_, err = copyFile(srcRootPath, destRootPath, f)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", f, err)
		}
	}

	return nil
}

func generateRawCopyList(srcRootPath string, rawCopyPaths []string) ([]string, error) {
	var copyList []string
	srcRoot := os.DirFS(srcRootPath)
	for _, rawGlob := range rawCopyPaths {
		rawFiles, err := doublestar.Glob(srcRoot, rawGlob)
		if err != nil {
			return nil, fmt.Errorf("failed to glob %s: %w", rawGlob, err)
		}
		copyList = append(copyList, rawFiles...)
	}
	return copyList, nil
}
