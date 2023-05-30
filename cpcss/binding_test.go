package cpcss_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/cpcss"
	"github.com/stretchr/testify/assert"
)

func TestCPCssImportsGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte("@import \"hello\";"), cpcss.GetLanguage())
	assert.NoError(err)
	assert.Equal(
		"(doc (import_statement (import_reference)))",
		n.String(),
	)
}
