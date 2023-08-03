package conv

import (
	"go/ast"
	"strings"
)

// Visitor satisfies ast.Visitor and it will turn given expression type into a string representation.
type Visitor struct {
	value string
	pos   int
}

var _ ast.Visitor = &Visitor{} // satisfies ast.Visitor

// Value returns visitor's value
func (v *Visitor) Value() string { return v.value }

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	v.pos = int(node.Pos())

	switch d := node.(type) {
	case *ast.ParenExpr:
		vx := &Visitor{}
		ast.Walk(vx, d.X)
		spacerX := createSpacer(vx.pos - int(d.Lparen) - 1)
		spacerY := createSpacer(int(d.Rparen) - (int(vx.pos) + len(vx.value)))
		v.value = "(" + spacerX + vx.value + spacerY + ")"
		return nil
	case *ast.UnaryExpr:
		vx := &Visitor{}
		ast.Walk(vx, d.X)
		spacer := createSpacer(vx.pos - int(d.OpPos) - 1)
		v.value = d.Op.String() + spacer + vx.value
		return nil
	case *ast.BinaryExpr:
		vx, vy := &Visitor{}, &Visitor{}
		ast.Walk(vx, d.X)
		ast.Walk(vy, d.Y)
		spacerX := createSpacer(int(d.OpPos) - (int(vx.pos) + len(vx.value)))
		spacerY := createSpacer(int(vy.pos) - (int(d.OpPos) + len(d.Op.String())))
		v.value = vx.value + spacerX + d.Op.String() + spacerY + vy.value
		return nil
	case *ast.BasicLit:
		v.value = d.Value
		return nil
	case *ast.Ident:
		v.value = d.Name
		return nil
	}

	return nil
}

func createSpacer(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}
