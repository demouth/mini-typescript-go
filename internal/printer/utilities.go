package printer

import (
	"fmt"

	"github.com/demouth/mini-typescript-go/internal/ast"
)

type getLiteralTextFlags int

const (
	getLiteralTextFlagsNone getLiteralTextFlags = 0
)

func isNotPrologueDirective(node *ast.Node) bool {
	return !ast.IsPrologueDirective(node)
}
func greatestEnd(end int, nodes ...interface{ End() int }) int {
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		if nodeEnd, ok := tryGetEnd(node); ok && end < nodeEnd {
			end = nodeEnd
		}
	}
	return end
}
func tryGetEnd(node interface{ End() int }) (int, bool) {
	switch v := node.(type) {
	case (*ast.Node):
		if v != nil {
			return v.End(), true
		}
	default:
		panic(fmt.Sprintf("unhandled type: %T", node))
	}
	return 0, false
}
