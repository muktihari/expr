package boolean

import (
	"go/ast"
	"go/token"
	"strconv"
)

// Visitor is boolean visitor interface
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() bool
}

// NewVisitor creates new boolean visitor
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res bool
}

func (v *visitor) visit(x ast.Expr) string {
	switch d := x.(type) {
	case *ast.BinaryExpr:
		dv := &visitor{}
		ast.Walk(dv, d)
		return strconv.FormatBool(dv.res)
	case *ast.ParenExpr:
		dv := &visitor{}
		ast.Walk(dv, d.X)
		return strconv.FormatBool(dv.res)
	case *ast.BasicLit:
		return d.Value
	case *ast.Ident:
		return d.String()
	}
	return ""
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch d := node.(type) {
	case *ast.BinaryExpr:
		x := v.visit(d.X)
		y := v.visit(d.Y)

		switch d.Op {
		case token.EQL:
			v.res = x == y
		case token.NEQ:
			v.res = x != y
		case token.GTR:
			v.res = x > y
		case token.GEQ:
			v.res = x >= y
		case token.LSS:
			v.res = x < y
		case token.LEQ:
			v.res = x <= y
		case token.LAND:
			xbool, _ := strconv.ParseBool(x)
			ybool, _ := strconv.ParseBool(y)
			v.res = xbool && ybool
		case token.LOR:
			xbool, _ := strconv.ParseBool(x)
			ybool, _ := strconv.ParseBool(y)
			v.res = xbool || ybool
		}

		return nil
	}
	return v
}

func (v *visitor) Result() bool {
	return v.res
}
