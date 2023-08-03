// conv package contains converter function(s)
package conv

import (
	"go/ast"
)

// FormatExpr formats ast.Expr into expr's string format.
func FormatExpr(e ast.Expr) string {
	v := &Visitor{}
	ast.Walk(v, e)
	return v.Value()
}
