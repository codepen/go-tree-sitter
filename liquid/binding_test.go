package liquid_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/liquid"

	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	content := `{% render "folder-2/header.liquid" %}`
	rootNode, err := sitter.ParseCtx(context.Background(), []byte(content), liquid.GetLanguage())

	assert.NoError(err)
	expected := `(program (render_statement file: (string)))`
	assert.Equal(expected, rootNode.String())
}
