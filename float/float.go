// Deprecated: This package is no longer maintained and might be deleted in the future, use expr.Visitor instead.
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
//
// Deprecated: use expr.Visitor instead.
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() (float64, error)
}

// NewVisitor creates new float visitor
//
// Deprecated: use this instead:
//
// v := expr.NewVisitor(
//
//	expr.WithNumericType(visitor.NumericTypeFloat)
//
// ).
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res float64
	err error
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

func (v *visitor) visitBinary(binaryExpr *ast.BinaryExpr) ast.Visitor {
	switch binaryExpr.Op {
	case token.ADD, token.SUB, token.MUL, token.QUO:
	default:
		v.err = ErrUnsupportedOperator
		return nil
	}

	x := &visitor{}
	ast.Walk(x, binaryExpr.X)
	if x.err != nil {
		v.err = x.err
		return nil
	}

	y := &visitor{}
	ast.Walk(y, binaryExpr.Y)
	if y.err != nil {
		v.err = y.err
		return nil
	}

	switch binaryExpr.Op {
	case token.ADD:
		v.res = x.res + y.res
	case token.SUB:
		v.res = x.res - y.res
	case token.MUL:
		v.res = x.res * y.res
	case token.QUO:
		v.res = x.res / y.res
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
