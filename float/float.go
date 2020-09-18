package float

import (
	"go/ast"
	"go/token"
	"strconv"
)

// Visitor is float visitor interface
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() float64
}

// NewVisitor creates new float visitor
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res float64
}

func (v *visitor) visit(x ast.Expr) float64 {
	switch d := x.(type) {
	case *ast.BinaryExpr:
		dv := &visitor{}
		ast.Walk(dv, d)
		return dv.res
	case *ast.ParenExpr:
		dv := &visitor{}
		ast.Walk(dv, d.X)
		return dv.res
	case *ast.BasicLit:
		switch d.Kind {
		case token.INT:
			fallthrough
		case token.FLOAT:
			value, _ := strconv.ParseFloat(d.Value, 64)
			return value
		}
	}
	return 0
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
		case token.ADD:
			v.res = x + y
		case token.SUB:
			v.res = x - y
		case token.MUL:
			v.res = x * y
		case token.QUO:
			v.res = x / y
		}
		return nil
	}
	return v
}

func (v *visitor) Result() float64 {
	return v.res
}
