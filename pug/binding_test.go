package pug_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/pug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte(`extends layout.pug`), pug.GetLanguage())
	assert.NoError(err)
	assert.Equal("(source_file (extends (keyword) (filename)))", n.String())
}

func TestCaptureIncludeWithoutSpaces(t *testing.T) {
	assert := assert.New(t)

	content := []byte(`include folder/file.pug`)
	node, err := sitter.ParseCtx(context.Background(), content, pug.GetLanguage())

	assert.NoError(err)
	assert.Equal("(source_file (include (keyword) (filename)))", node.String())

	// Execute query and ensure values match
	query := "(include (keyword) (filename) @reference ) @node"
	tsitterQuery, err := sitter.NewQuery([]byte(query), pug.GetLanguage())
	assert.NoError(err)

	queryCursor := sitter.NewQueryCursor()
	queryCursor.Exec(tsitterQuery, node)

	// Iterate over the query results.
	indexToCapture := map[uint32]string{
		1: "include folder/file.pug",
		0: "folder/file.pug",
	}

	for {
		match, ok := queryCursor.NextMatch()
		if !ok {
			break
		}

		for _, cap := range match.Captures {
			expected := indexToCapture[cap.Index]
			actual := string(content[cap.Node.StartByte():cap.Node.EndByte()])
			require.Equal(t, expected, actual)
		}
	}
}

func TestCaptureExtendWithoutSpaces(t *testing.T) {
	assert := assert.New(t)

	content := []byte(`extends layout.pug`)
	node, err := sitter.ParseCtx(context.Background(), content, pug.GetLanguage())

	assert.NoError(err)
	assert.Equal("(source_file (extends (keyword) (filename)))", node.String())

	// Execute query and ensure values match
	query := "(extends (keyword) (filename) @reference ) @node"
	tsitterQuery, err := sitter.NewQuery([]byte(query), pug.GetLanguage())
	assert.NoError(err)

	queryCursor := sitter.NewQueryCursor()
	queryCursor.Exec(tsitterQuery, node)

	// Iterate over the query results.
	indexToCapture := map[uint32]string{
		1: "extends layout.pug",
		0: "layout.pug",
	}

	for {
		match, ok := queryCursor.NextMatch()
		if !ok {
			break
		}

		for _, cap := range match.Captures {
			expected := indexToCapture[cap.Index]
			actual := string(content[cap.Node.StartByte():cap.Node.EndByte()])
			require.Equal(t, expected, actual)
		}
	}
}
