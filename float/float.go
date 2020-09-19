package float

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
)

// ErrUnsupportedOperator is error unsupported operator
var ErrUnsupportedOperator = errors.New("unsupported operator")

// Visitor is float visitor interface
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() (float64, error)
}

// NewVisitor creates new float visitor
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res float64
	err error
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil || v.err != nil {
		return nil
	}

	switch d := node.(type) {
	case *ast.ParenExpr:
		return v.Visit(d.X)
	case *ast.BinaryExpr:
		switch d.Op { // early validate operator
		case token.ADD, token.SUB, token.MUL, token.QUO:
		default:
			v.err = ErrUnsupportedOperator
			return nil
		}

		xVisitor := &visitor{}
		ast.Walk(xVisitor, d.X)
		if xVisitor.err != nil {
			v.err = xVisitor.err
			return nil
		}
		x := xVisitor.res

		yVisitor := &visitor{}
		ast.Walk(yVisitor, d.Y)
		if yVisitor.err != nil {
			v.err = yVisitor.err
			return nil
		}
		y := yVisitor.res

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
	case *ast.BasicLit:
		switch d.Kind {
		case token.INT:
			fallthrough
		case token.FLOAT:
			v.res, _ = strconv.ParseFloat(d.Value, 64)
		}
		return nil
	}

	return v
}

func (v *visitor) Result() (float64, error) {
	return v.res, v.err
}
