package generator

import (
	"log/slog"

	"github.com/gwkline/artestian/types"
)

func (g *TestGenerator) findBestExample(sourceCode string) types.TestExample {
	if len(g.examples) == 0 {
		return types.TestExample{}
	}

	example, err := g.ai.PickExample(sourceCode, g.examples)
	if err != nil {
		slog.Error("failed to pick example", "error", err)
		return g.examples[0]
	}

	return example
}
