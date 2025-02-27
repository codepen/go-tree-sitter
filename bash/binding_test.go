package bash_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/bash"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte("echo 1"), bash.GetLanguage())
	assert.NoError(err)
	assert.Equal(
		"(program (command name: (command_name (word)) argument: (number)))",
		n.String(),
	)
}
