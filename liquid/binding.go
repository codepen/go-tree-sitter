package liquid

//#include "tree_sitter/parser.h"
//TSLanguage *tree_sitter_liquid();
import "C"

import (
	"unsafe"

	sitter "github.com/codepen/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_liquid())
	return sitter.NewLanguage(ptr)
}
