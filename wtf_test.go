package sitter

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input = `
// the ancient pond
// a frog jumps in
// the sound of water
`

func evalQuery(t *testing.T, query string) {

	parser := NewParser()
	grammar := getTestGrammar()
	parser.SetLanguage(grammar)

	tree := parser.Parse(nil, []byte(input))
	root := tree.RootNode()
	// log.Print("Root: ", root)

	q, err := NewQuery([]byte(query), grammar)
	if err != nil {
		assert.FailNow(t, "Failed to parse query", err)
	}

	qc := NewQueryCursor()
	qc.Exec(q, root)

	nMatches := 0
	nPass := 0

	for {
		match, ok := qc.NextMatch()

		if !ok {
			break
		}

		nMatches++
		isMatch := true

		// Enable / disable predicate check
		if true {
			filteredMatch := qc.FilterPredicates(match, []byte(input))
			isMatch = len(filteredMatch.Captures) > 0
		}

		if isMatch {
			nPass++
		}
	}

	log.Printf(`query = "%s"

Predicates passed %d of %d matches
======
`,
		strings.TrimSpace(query),
		nPass,
		nMatches,
	)
}

func TestWTF(t *testing.T) {
	tests := []string{
		`
(expression (comment) @foo
	(#match? @foo "^// the")
)
		`,
		`
(expression (comment) @foo)
(#match? @foo "^// the")
		`,
		`
(expression (comment) @foo)
(#match? @foo "^// the")
(#match? @foo "water$")
		`,
	}

	log.Printf("input = \"%s\"", input)

	for _, test := range tests {
		evalQuery(t, test)
	}
}
