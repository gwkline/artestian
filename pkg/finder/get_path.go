package finder

import (
	"path/filepath"
	"strings"
)

func (f *FileFinder) GetTestPath(sourcePath string) string {
	ext := filepath.Ext(sourcePath)
	return strings.TrimSuffix(sourcePath, ext) + f.language.GetTestFilePattern()
}
