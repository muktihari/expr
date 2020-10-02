package integer

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
)

// ErrUnsupportedOperator is error unsupported operator
var ErrUnsupportedOperator = errors.New("unsupported operator")

// ErrIntegerDividedByZero occurs when x/y and y equals to 0, Go does not allow integer to be divided by zero
var ErrIntegerDividedByZero = errors.New("integer divided by zero")

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
	res  int
	base int // base10 except stated otherwise
	err  error
}

func (v *visitor) visitUnary(unaryExpr *ast.UnaryExpr) ast.Visitor {
	switch unaryExpr.Op {
	case token.ADD, token.SUB:
		xVisitor := &visitor{}
		ast.Walk(xVisitor, unaryExpr.X)
		if xVisitor.err != nil {
			v.err = xVisitor.err
			return nil
		}
		switch unaryExpr.Op {
		case token.ADD:
			v.res = xVisitor.res
		case token.SUB:
			v.res = xVisitor.res * -1
		}
	default:
		v.err = ErrUnsupportedOperator
	}
	return nil
}

func (v *visitor) arithmetic(binaryExpr *ast.BinaryExpr) {
	x := &visitor{}
	ast.Walk(x, binaryExpr.X)
	if x.err != nil {
		v.err = x.err
		return
	}
	y := &visitor{}
	ast.Walk(y, binaryExpr.Y)
	if y.err != nil {
		v.err = y.err
		return
	}

	switch binaryExpr.Op {
	case token.ADD:
		v.res = x.res + y.res
	case token.SUB:
		v.res = x.res - y.res
	case token.MUL:
		v.res = x.res * y.res
	case token.QUO:
		if y.res == 0 {
			v.err = ErrIntegerDividedByZero
			return
		}
		v.res = x.res / y.res
	case token.REM:
		v.res = x.res % y.res
	}
}

func (v *visitor) bitwise(binaryExpr *ast.BinaryExpr) {
	x := &visitor{base: 2}
	ast.Walk(x, binaryExpr.X)
	if x.err != nil {
		v.err = x.err
		return
	}

	y := &visitor{base: 2}
	ast.Walk(y, binaryExpr.Y)
	if y.err != nil {
		v.err = y.err
		return
	}

	switch binaryExpr.Op {
	case token.AND:
		v.res = x.res & y.res
	case token.OR:
		v.res = x.res | y.res
	case token.XOR:
		v.res = x.res ^ y.res
	case token.AND_NOT:
		v.res = x.res &^ y.res
	case token.SHL:
		v.res = x.res << y.res
	case token.SHR:
		v.res = x.res >> y.res
	}
}

func (v *visitor) visitBinary(binaryExpr *ast.BinaryExpr) ast.Visitor {
	switch binaryExpr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		v.arithmetic(binaryExpr)
	case token.AND, token.OR, token.XOR, token.AND_NOT, token.SHL, token.SHR:
		v.bitwise(binaryExpr)
	default:
		v.err = ErrUnsupportedOperator
	}
	return nil
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil || v.err != nil {
		return nil
	}

	switch d := node.(type) {
	case *ast.ParenExpr:
		return v.Visit(d.X)
	case *ast.UnaryExpr:
		return v.visitUnary(d)
	case *ast.BinaryExpr:
		return v.visitBinary(d)
	case *ast.BasicLit:
		switch d.Kind {
		case token.INT:
			if v.base == 2 {
				var res int64
				res, v.err = strconv.ParseInt(d.Value, 2, 64)
				v.res = int(res)
				return nil
			}
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
