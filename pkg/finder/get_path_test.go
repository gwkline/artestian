package finder

import (
	"testing"

	"github.com/gwkline/artestian/pkg/golang"
	"github.com/gwkline/artestian/pkg/typescript"
	"github.com/gwkline/artestian/types"
	"github.com/stretchr/testify/assert"
)

func TestGetTestPath(t *testing.T) {
	tests := []struct {
		name         string
		sourcePath   string
		filePattern  string
		expectedPath string
		language     types.ILanguage
	}{
		{
			name:         "converts .go file to _test.go",
			sourcePath:   "/path/to/source.go",
			filePattern:  "_test.go",
			expectedPath: "/path/to/source_test.go",
			language:     golang.NewGoSupport(),
		},
		{
			name:         "converts .ts file to .test.ts",
			sourcePath:   "/path/to/source.ts",
			filePattern:  ".test.ts",
			expectedPath: "/path/to/source.test.ts",
			language:     typescript.NewTypeScriptSupport(),
		},
		{
			name:         "handles paths with no extension",
			sourcePath:   "/path/to/source",
			filePattern:  "_test.go",
			expectedPath: "/path/to/source_test.go",
			language:     golang.NewGoSupport(),
		},
		{
			name:         "handles paths with multiple dots",
			sourcePath:   "/path/to/source.handler.go",
			filePattern:  "_test.go",
			expectedPath: "/path/to/source.handler_test.go",
			language:     golang.NewGoSupport(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := &FileFinder{
				language: tt.language,
			}
			result := finder.GetTestPath(tt.sourcePath)
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}
