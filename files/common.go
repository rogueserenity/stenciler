package files

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/bmatcuk/doublestar/v4"
)

// ensureDirExists ensures that the directory containing the file at relFilePath exists in destRootPath with the same
// permissions as the directory containing the file at relFilePath in srcRootPath.
func ensureDirExists(srcRootPath, destRootPath, relFilePath string) error {
	if relFilePath == "." {
		return nil
	}

	srcPath := filepath.Join(srcRootPath, relFilePath)
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", srcPath, err)
	}

	if !srcInfo.IsDir() {
		return ensureDirExists(srcRootPath, destRootPath, filepath.Dir(relFilePath))
	}

	err = ensureDirExists(srcRootPath, destRootPath, filepath.Dir(relFilePath))
	if err != nil {
		return fmt.Errorf("failed to ensure directory exists: %w", err)
	}
	perms := srcInfo.Mode().Perm()
	destPath := filepath.Join(destRootPath, relFilePath)
	err = os.Mkdir(destPath, perms)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to create directory %s: %w", destPath, err)
		}
	}

	return nil
}

// openSourceFile opens the file at relFilePath in srcRootPath and returns the file and its FileInfo.
func openSourceFile(srcRootPath string, relFilePath string) (*os.File, fs.FileInfo, error) {
	srcPath := filepath.Join(srcRootPath, relFilePath)
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stat %s: %w", srcPath, err)
	}
	if !srcInfo.Mode().IsRegular() {
		return nil, nil, fmt.Errorf("failed to copy %s: not a regular file", relFilePath)
	}
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open %s: %w", srcPath, err)
	}
	return file, srcInfo, nil
}

// createDestFile creates the file at relFilePath in destRootPath with the same permissions as the file at relFilePath
// in srcRootPath.
func createDestFile(destRootPath, relFilePath string, perm fs.FileMode) (*os.File, error) {
	destPath := filepath.Join(destRootPath, relFilePath)
	file, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s: %w", destPath, err)
	}
	err = file.Chmod(perm)
	if err != nil {
		return nil, fmt.Errorf("failed to set permissions on %s: %w", destPath, err)
	}
	return file, nil
}

// isRegularFile returns true if the file at relFilePath in srcRootPath is a regular file.
func isRegularFile(srcRootPath, relFilePath string) bool {
	srcPath := filepath.Join(srcRootPath, relFilePath)
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return false
	}
	return srcInfo.Mode().IsRegular()
}

// createFileList generates a list of files in the specified root that match the list of glob patterns.
func createFileList(root string, globPatterns []string) ([]string, error) {
	var copyList []string
	srcRoot := os.DirFS(root)
	for _, rawGlob := range globPatterns {
		rawFiles, err := doublestar.Glob(srcRoot, rawGlob)
		if err != nil {
			return nil, fmt.Errorf("failed to glob %s: %w", rawGlob, err)
		}
		copyList = append(copyList, rawFiles...)
	}

	slog.Debug("copy list", slog.Any("files", copyList))

	return copyList, nil
}

func removeFromFileList(fileList, removeList []string) []string {
	var resultList []string
	for _, f := range fileList {
		if !slices.Contains(removeList, f) {
			resultList = append(resultList, f)
		}
	}
	return resultList
}
