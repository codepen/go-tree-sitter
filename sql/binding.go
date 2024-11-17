package sql

//#include "tree_sitter/parser.h"
//TSLanguage *tree_sitter_sql();
import "C"
import (
	"unsafe"

	sitter "github.com/codepen/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_sql())
	return sitter.NewLanguage(ptr)
}
