package sitter

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var input = `
// the ancient pond
// a frog jumps in
// the sound of water
`

func evalQuery(query string) int {
	parser := NewParser()
	grammar := getTestGrammar()
	parser.SetLanguage(grammar)

	tree := parser.Parse(nil, []byte(input))
	root := tree.RootNode()
	// log.Print("Root: ", root)

	q, err := NewQuery([]byte(query), grammar)
	if err != nil {
		log.Panic("Failed to parse query", err)
	}

	qc := NewQueryCursor()
	qc.Exec(q, root)

	count := 0

	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}

		// Change to `true` to enable predicate check (doesn't affect number of
		// matches, however)
		if false {
			filteredMatch := qc.FilterPredicates(match, []byte(input))
			isMatch := len(filteredMatch.Captures) > 0
			if !isMatch {
				count--
				continue
			}
		}

		count++
	}

	return count
}

func TestWTF(t *testing.T) {
	log.Printf("input = \"%s\"", input)

	assert.Equal( // PASSES
		t,
		3,
		evalQuery(`
		(expression (comment) @foo)
	`),
	)

	assert.Equal( // FAILS
		t,
		3,
		evalQuery(`
		(expression (comment) @foo
			(#match? @foo "^// the")
		)
	`),
	)

	assert.Equal( // FAILS
		t,
		3,
		evalQuery(`
			(expression (comment) @foo)
			(#match? @foo "^// the")
		`),
	)

	assert.Equal( // FAILS
		t,
		3,
		evalQuery(`
			(expression (comment) @foo)
				(#match? @foo "^// the")
				(#match? @foo "water$")
		`),
	)
}
