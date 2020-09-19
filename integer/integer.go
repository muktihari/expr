package integer

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
)

// ErrUnsupportedOperator is error unsupported operator
var ErrUnsupportedOperator = errors.New("unsupported operator")

// Visitor is integer visitor interface
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() (int, error)
}

// NewVisitor creates new integer visitor
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res int
	err error
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch d := node.(type) {
	case *ast.ParenExpr:
		return v.Visit(d.X)
	case *ast.BinaryExpr:
		switch d.Op { // early validate operator
		case token.ADD, token.SUB, token.MUL, token.QUO, token.REM, // arithmatic operators
			token.AND, token.OR, token.XOR, token.AND_NOT, token.SHL, token.SHR: // bitwise operators
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
		case token.REM:
			v.res = x % y
		case token.AND:
			v.res = x & y
		case token.OR:
			v.res = x | y
		case token.XOR:
			v.res = x ^ y
		case token.AND_NOT:
			v.res = x &^ y
		case token.SHL:
			v.res = x << y
		case token.SHR:
			v.res = x >> y
		}

		return nil
	case *ast.BasicLit:
		switch d.Kind {
		case token.INT:
			v.res, _ = strconv.Atoi(d.Value)
		case token.FLOAT:
			floatValue, _ := strconv.ParseFloat(d.Value, 32)
			v.res = int(floatValue)
		}
		return nil
	}

	return v
}

func (v *visitor) Result() (int, error) {
	return v.res, v.err
}
