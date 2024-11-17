# go tree-sitter

[![Build Status](https://github.com/codepen/go-tree-sitter/workflows/Test/badge.svg?branch=master)](https://github.com/codepen/go-tree-sitter/actions/workflows/test.yml?query=branch%3Amaster)
[![GoDoc](https://godoc.org/github.com/codepen/go-tree-sitter?status.svg)](https://godoc.org/github.com/codepen/go-tree-sitter)

Golang bindings for [tree-sitter](https://github.com/tree-sitter/tree-sitter)

## How to Update Pug Grammar

### 1. Update Automation Grammars Configuration File
Update the grammar to the `codepen/go-tree-sitter/_automation/grammars.json` file. The format is:

```json
{
"language": "pug",
"url": "https://github.com/codepen/tree-sitter-pug",
"files": ["parser.c", "scanner.cc"],
"reference": "master",
"revision": "757e95a5fbf26058e38f9beb1fd2f05c140410a7"
},
```

### 2. Run Script to Update Grammar

This will pull the latest generated C parser files from the GitHub repo.

```bash
go run _automation/main.go update pug -force
```

### 3. Add Go Binding

Manually add a file called `binding.go` to the `codepen/go-tree-sitter/pug` directory. This file should contain the following code, specific to your grammar:

```go
package pug

//#include "parser.h"
//TSLanguage *tree_sitter_pug();
import "C"

import (
	"unsafe"

	sitter "github.com/codepen/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_pug())
	return sitter.NewLanguage(ptr)
}
```

This logic adds a `GetLanguage` method that allows you to call the underlying C code from Go.

Be sure to set the correct language name in two places in this file:
Line 4: `//TSLanguage *tree_sitter_<language>();`
Line 14: `ptr := unsafe.Pointer(C.tree_sitter_<language>())`

### 4. Add Binding Test
Test the syntax tree your newly added parser generates by creating a test file in the `codepen/go-tree-sitter/<language>` directory. This file should contain the following code, specific to your grammar:

```go
package pug_test

import (
	"context"
	"testing"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/pug"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte(`extends layout.pug`), pug.GetLanguage())
	assert.NoError(err)
	assert.Equal("(source_file (extends_statement (path)))", n.String())
}
```

Run the test from the root of the repo with the following command:
```bash
go test -v pug/binding_test.go
```

## Usage

Create a parser with a grammar:

```go
import (
	"context"
	"fmt"

	sitter "github.com/codepen/go-tree-sitter"
	"github.com/codepen/go-tree-sitter/javascript"
)

parser := sitter.NewParser()
parser.SetLanguage(javascript.GetLanguage())
```

Parse some code:

```go
sourceCode := []byte("let a = 1")
tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
```

Inspect the syntax tree:

```go
n := tree.RootNode()

fmt.Println(n) // (program (lexical_declaration (variable_declarator (identifier) (number))))

child := n.NamedChild(0)
fmt.Println(child.Type()) // lexical_declaration
fmt.Println(child.StartByte()) // 0
fmt.Println(child.EndByte()) // 9
```

### Custom grammars

This repository provides grammars for many common languages out of the box.

But if you need support for any other language you can keep it inside your own project or publish it as a separate repository to share with the community.

See explanation on how to create a grammar for go-tree-sitter [here](https://github.com/codepen/go-tree-sitter/issues/57).

Known external grammars:

- [Salesforce grammars](https://github.com/aheber/tree-sitter-sfapex) - including Apex, SOQL, and SOSL languages.
- [Ruby](https://github.com/shagabutdinov/go-tree-sitter-ruby) - Deprecated, grammar is provided by main repo instead
- [Go Template](https://github.com/mrjosh/helm-ls/tree/master/internal/tree-sitter/gotemplate) - Used for helm

### Editing

If your source code changes, you can update the syntax tree. This will take less time than the first parse.

```go
// change 1 -> true
newText := []byte("let a = true")
tree.Edit(sitter.EditInput{
    StartIndex:  8,
    OldEndIndex: 9,
    NewEndIndex: 12,
    StartPoint: sitter.Point{
        Row:    0,
        Column: 8,
    },
    OldEndPoint: sitter.Point{
        Row:    0,
        Column: 9,
    },
    NewEndPoint: sitter.Point{
        Row:    0,
        Column: 12,
    },
})

// check that it changed tree
assert.True(n.HasChanges())
assert.True(n.Child(0).HasChanges())
assert.False(n.Child(0).Child(0).HasChanges()) // left side of the tree didn't change
assert.True(n.Child(0).Child(1).HasChanges())

// generate new tree
newTree := parser.Parse(tree, newText)
```

### Predicates

You can filter AST by using [predicate](https://tree-sitter.github.io/tree-sitter/using-parsers#predicates) S-expressions.

Similar to [Rust](https://github.com/tree-sitter/tree-sitter/tree/master/lib/binding_rust) or [WebAssembly](https://github.com/tree-sitter/tree-sitter/blob/master/lib/binding_web) bindings we support filtering on a few common predicates:
- `eq?`, `not-eq?`
- `match?`, `not-match?`

Usage [example](./_examples/predicates/main.go):

```go
func main() {
	// Javascript code
	sourceCode := []byte(`
		const camelCaseConst = 1;
		const SCREAMING_SNAKE_CASE_CONST = 2;
		const lower_snake_case_const = 3;`)
	// Query with predicates
	screamingSnakeCasePattern := `(
		(identifier) @constant
		(#match? @constant "^[A-Z][A-Z_]+")
	)`

	// Parse source code
	lang := javascript.GetLanguage()
	n, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)
	// Execute the query
	q, _ := sitter.NewQuery([]byte(screamingSnakeCasePattern), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)
	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			fmt.Println(c.Node.Content(sourceCode))
		}
	}
}

// Output of this program:
// SCREAMING_SNAKE_CASE_CONST
```

## Development

### Updating a grammar

Check if any updates for vendored files are available:

```
go run _automation/main.go check-updates
```

Update vendor files:

- open `_automation/grammars.json`
- modify `reference` (for tagged grammars) or `revision` (for grammars from a branch)
- run `go run _automation/main.go update <grammar-name>`

It is also possible to update all grammars in one go using

```
go run _automation/main.go update-all
```
