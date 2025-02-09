package finder

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func (f *FileFinder) FindNextFile(rootDir string) (string, error) {
	slog.Debug("starting file search", "rootDir", rootDir)
	var foundFile string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("error accessing path", "path", path, "error", err)
			return err
		}

		// Skip if already visited
		if f.visited[path] {
			slog.Debug("skipping visited file", "path", path)
			return nil
		}

		// Skip directories and non-typescript files
		if info.IsDir() {
			slog.Debug("skipping directory", "path", path)
			return nil
		}
		if !strings.HasSuffix(path, f.language.GetFileExtension()) {
			slog.Debug("skipping non-target file", "path", path, "extension", filepath.Ext(path))
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, f.language.GetTestFilePattern()) {
			slog.Debug("skipping test file", "path", path)
			return nil
		}

		// Skip if test file already exists
		testPath := f.GetTestPath(path)
		if _, err := os.Stat(testPath); err == nil {
			slog.Debug("skipping file with existing test", "path", path, "testPath", testPath)
			return nil
		}

		// Found a file that needs tests
		slog.Info("found file needing tests", "path", path)
		foundFile = path
		f.visited[path] = true
		return filepath.SkipAll
	})

	if err != nil {
		slog.Error("error walking directory", "error", err)
		return "", err
	}

	if foundFile == "" {
		slog.Info("no files found needing tests")
		return "", nil
	}

	return foundFile, nil
}
