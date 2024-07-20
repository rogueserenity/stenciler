package files

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/rogueserenity/stenciler/config"
)

// CopyTemplated copies all templated files (those no located in raw-copy) and optionaly excluding files listed in
// init-only by passing them through the template engine, then writing the copies into the current working directory.
func CopyTemplated(repoDir string, tmplate *config.Template) error {
	srcRootPath := filepath.Join(repoDir, tmplate.Directory)

	destRootPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	params := make(map[string]string)
	for _, p := range tmplate.Params {
		params[p.Name] = p.Value
	}

	fileList, err := createSourceFileList(srcRootPath)
	if err != nil {
		return fmt.Errorf("failed to generate file list: %w", err)
	}
	fileList = removeFromFileList(fileList, tmplate.RawCopyPaths)

	if tmplate.Update {
		initOnlyList, err := createFileList(srcRootPath, tmplate.InitOnlyPaths)
		if err != nil {
			return fmt.Errorf("failed to generate init-only list: %w", err)
		}
		fileList = removeFromFileList(fileList, initOnlyList)
	}

	for _, f := range fileList {
		err = copyTemplatedFile(srcRootPath, destRootPath, f, params)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", f, err)
		}
	}

	return nil
}

func createSourceFileList(root string) ([]string, error) {
	srcRoot := os.DirFS(root)
	allFiles, err := doublestar.Glob(srcRoot, "**")
	if err != nil {
		return nil, fmt.Errorf("failed to glob all files: %w", err)
	}
	return allFiles, nil
}

func copyTemplatedFile(srcRootPath, destRootPath, relFilePath string, params map[string]string) error {
	if !isRegularFile(srcRootPath, relFilePath) {
		return nil
	}

	err := ensureDirExists(srcRootPath, destRootPath, relFilePath)
	if err != nil {
		return fmt.Errorf("failed to ensure directory exists: %w", err)
	}

	templateFile, err := template.ParseFiles(filepath.Join(srcRootPath, relFilePath))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	srcFilePath := filepath.Join(srcRootPath, relFilePath)
	srcFileInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	destFilePath := filepath.Join(destRootPath, relFilePath)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	err = destFile.Chmod(srcFileInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to set destination file permissions: %w", err)
	}

	err = templateFile.Execute(destFile, params)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	slog.Debug("copied templated file",
		slog.String("src", srcFilePath),
		slog.String("dest", destFilePath),
	)

	return nil
}
