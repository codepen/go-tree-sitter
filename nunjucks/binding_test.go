package nunjucks_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/nunjucks"

	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	content := `{% include "folder-2/header.njk" %}`
	rootNode, err := sitter.ParseCtx(context.Background(), []byte(content), nunjucks.GetLanguage())

	assert.NoError(err)
	expected := `(fragment (fragment_statement (statement (include_statement template: (expression (string))))))`
	assert.Equal(expected, rootNode.String())
}
