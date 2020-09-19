package boolean

import (
	"errors"
	"go/ast"
	"go/token"
	"strconv"
)

// ErrUnsupportedOperator is error unsupported operator
var ErrUnsupportedOperator = errors.New("unsupported operator")

// Visitor is boolean visitor interface
type Visitor interface {
	Visit(node ast.Node) ast.Visitor
	Result() (bool, error)
}

// NewVisitor creates new boolean visitor
func NewVisitor() Visitor {
	return &visitor{}
}

type visitor struct {
	res string
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
		case token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ, token.LAND, token.LOR:
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
		case token.EQL:
			v.res = strconv.FormatBool(x == y)
		case token.NEQ:
			v.res = strconv.FormatBool(x != y)
		case token.GTR:
			v.res = strconv.FormatBool(x > y)
		case token.GEQ:
			v.res = strconv.FormatBool(x >= y)
		case token.LSS:
			v.res = strconv.FormatBool(x < y)
		case token.LEQ:
			v.res = strconv.FormatBool(x <= y)
		case token.LAND:
			xbool, _ := strconv.ParseBool(x)
			ybool, _ := strconv.ParseBool(y)
			v.res = strconv.FormatBool(xbool && ybool)
		case token.LOR:
			xbool, _ := strconv.ParseBool(x)
			ybool, _ := strconv.ParseBool(y)
			v.res = strconv.FormatBool(xbool || ybool)
		default:
			v.err = ErrUnsupportedOperator
		}

		return nil
	case *ast.BasicLit:
		v.res = d.Value
		return nil
	case *ast.Ident:
		v.res = d.String()
		return nil
	}

	return v
}

func (v *visitor) Result() (bool, error) {
	if v.err != nil {
		return false, v.err
	}
	var res bool
	res, v.err = strconv.ParseBool(v.res)
	return res, v.err
}
