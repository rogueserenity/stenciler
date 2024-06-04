package files

import (
	"fmt"
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

	fileList, err := generateFileList(srcRootPath, tmplate)
	if err != nil {
		return fmt.Errorf("failed to generate file list: %w", err)
	}

	for _, f := range fileList {
		err = copyTemplatedFile(srcRootPath, destRootPath, f, params)
		if err != nil {
			return fmt.Errorf("failed to copy %s: %w", f, err)
		}
	}

	return nil
}

func generateFileList(srcRootPath string, tmplate *config.Template) ([]string, error) {
	var fileList []string
	srcRoot := os.DirFS(srcRootPath)
	allFiles, err := doublestar.Glob(srcRoot, "**")
	if err != nil {
		return nil, fmt.Errorf("failed to glob all files: %w", err)
	}
	for _, file := range allFiles {
		isRawCopyFile, err := isRawCopyFile(file, tmplate.RawCopyPaths)
		if err != nil {
			return nil, fmt.Errorf("failed to check if %s is a raw copy file: %w", file, err)
		}
		if isRawCopyFile {
			continue
		}

		fileList = append(fileList, file)
	}

	return fileList, nil
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

	return nil
}

func isRawCopyFile(filePath string, rawCopyPaths []string) (bool, error) {
	for _, rawGlob := range rawCopyPaths {
		match, err := doublestar.Match(rawGlob, filePath)
		if err != nil {
			return false, fmt.Errorf("failed to match %s: %w", rawGlob, err)
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}
