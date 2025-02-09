package finder

import (
	"github.com/gwkline/artestian/types"
)

type FileFinder struct {
	language types.ILanguage
	visited  map[string]bool
}

func NewFileFinder(lang types.ILanguage) *FileFinder {
	return &FileFinder{
		language: lang,
		visited:  make(map[string]bool),
	}
}
