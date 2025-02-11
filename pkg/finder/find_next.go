package finder

import (
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func (f *FileFinder) FindNextFile(rootDir string, excludeDirs []string) (string, error) {
	slog.Debug("starting file search", "rootDir", rootDir, "excludeDirs", excludeDirs)
	var eligibleFiles []string

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
			// Skip excluded directories
			for _, excludeDir := range excludeDirs {
				// Use filepath.Clean to normalize paths for comparison
				cleanPath := filepath.Clean(path)
				cleanExcludeDir := filepath.Clean(excludeDir)

				// Check if the current path matches the excluded directory exactly
				// or if it's a subdirectory of the excluded directory
				rel, err := filepath.Rel(cleanExcludeDir, cleanPath)
				if err == nil && (rel == "." || !strings.HasPrefix(rel, "..")) {
					slog.Debug("skipping excluded directory", "path", path, "excludeDir", excludeDir)
					return filepath.SkipDir
				}
			}
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

		// Add eligible file to the list
		eligibleFiles = append(eligibleFiles, path)
		return nil
	})

	if err != nil {
		slog.Error("error walking directory", "error", err)
		return "", err
	}

	if len(eligibleFiles) == 0 {
		slog.Info("no files found needing tests")
		return "", nil
	}

	slog.Debug("eligible files", "files", eligibleFiles)
	slog.Debug("exclude dirs", "dirs", excludeDirs)

	// Randomly select one of the eligible files
	selectedFile := eligibleFiles[rand.Intn(len(eligibleFiles))]
	f.visited[selectedFile] = true
	slog.Info("randomly selected file for testing", "path", selectedFile)

	return selectedFile, nil
}
