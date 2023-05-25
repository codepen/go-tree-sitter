package cssimports_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/cssimports"
	"github.com/stretchr/testify/assert"
)

func TestCssImportsGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte("@import \"hello\";"), cssimports.GetLanguage())
	assert.NoError(err)
	assert.Equal(
		"(doc (import_statement (quoted_import_reference (import_reference))))",
		n.String(),
	)
}
