package generator

import (
	"github.com/gwkline/artestian/types"
)

type TestGenerator struct {
	finder       types.IFileFinder
	ai           types.IAgent
	language     types.ILanguage
	examples     []types.TestExample
	contextFiles []types.ContextFile
}

func NewTestGenerator(
	finder types.IFileFinder,
	ai types.IAgent,
	language types.ILanguage,
	examples []types.TestExample,
	contextFiles []types.ContextFile,
) *TestGenerator {
	return &TestGenerator{
		finder:       finder,
		ai:           ai,
		language:     language,
		examples:     examples,
		contextFiles: contextFiles,
	}
}
