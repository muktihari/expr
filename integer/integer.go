package integer

import (
	"go/ast"
	"go/token"
	"strconv"
)

// Visitor is integer visitor interface
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() int
}

// NewVisitor creates new integer visitor
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res int
}

func (v *visitor) visit(x ast.Expr) int {
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
			value, _ := strconv.Atoi(d.Value)
			return value
		case token.FLOAT:
			value, _ := strconv.ParseFloat(d.Value, 32) // we expect int, so we go with small float32
			return int(value)
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

func (v *visitor) Result() int {
	return v.res
}
