package files

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rogueserenity/stenciler/config"
)

// CopyRaw copies files that do not require template processing into the current working directory.
func CopyRaw(repoDir string, template *config.Template) error {
	srcRootPath := filepath.Join(repoDir, template.Directory)

	destRootPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	copyList, err := createFileList(srcRootPath, template.RawCopyPaths)
	if err != nil {
		return fmt.Errorf("failed to generate copy list: %w", err)
	}

	if template.Update {
		initOnlyList, err := createFileList(srcRootPath, template.InitOnlyPaths)
		if err != nil {
			return fmt.Errorf("failed to generate init-only list: %w", err)
		}
		copyList = removeFromFileList(copyList, initOnlyList)
	}

	for _, f := range copyList {
		_, err = copyRawFile(srcRootPath, destRootPath, f)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", f, err)
		}
	}

	return nil
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
