package cue_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/cue"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	code := `a: {
		F=f: string
		X="my-identifier": bool
	}
`

	n, err := sitter.ParseCtx(context.Background(), []byte(code), cue.GetLanguage())
	assert.NoError(err)
	assert.Equal(
		"(source_file (field (label (identifier)) (value (struct_lit (field (label alias: (identifier) (identifier)) (value (primitive_type))) (field (label alias: (identifier) (string)) (value (primitive_type)))))))",
		n.String(),
	)
}
