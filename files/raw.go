package files

import (
	"fmt"
	"io"
	"log/slog"
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
		_, err = copyRawFile(srcRootPath, destRootPath, f)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", f, err)
		}
	}

	return nil
}

// generateRawCopyList generates a list of files to copy from the template directory without template processing.
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

	slog.Debug("copy list", slog.Any("files", copyList))

	return copyList, nil
}

// copyRawFile copies the file at relFilePath in srcRootPath to destRootPath. It skips any non-regular files and will
// ensure that the directory containing the file exists in destRootPath.
func copyRawFile(srcRootPath, destRootPath, relFilePath string) (int64, error) {
	if !isRegularFile(srcRootPath, relFilePath) {
		return 0, nil
	}

	err := ensureDirExists(srcRootPath, destRootPath, relFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to ensure directory exists: %w", err)
	}

	srcFile, srcInfo, err := openSourceFile(srcRootPath, relFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := createDestFile(destRootPath, relFilePath, srcInfo.Mode().Perm())
	if err != nil {
		return 0, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	b, err := io.Copy(destFile, srcFile)
	if err != nil {
		return 0, fmt.Errorf("failed to copy %s: %w", relFilePath, err)
	}

	slog.Debug("copied file",
		slog.String("src", srcFile.Name()),
		slog.String("dest", destFile.Name()),
	)
	return b, nil
}
