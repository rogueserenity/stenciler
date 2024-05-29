package files

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

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

func copyFile(srcRootPath, destRootPath, relFilePath string) (int64, error) {
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
