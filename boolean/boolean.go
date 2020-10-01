package boolean

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// ErrUnsupportedOperator is error unsupported operator
var ErrUnsupportedOperator = errors.New("unsupported operator")

// ErrInvalidOperationOnFloat is error invalid operation on float
var ErrInvalidOperationOnFloat = errors.New("invalid operation on float")

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
	kind token.Token
	res  string
	err  error
}

func (v *visitor) visitUnary(unaryExpr *ast.UnaryExpr) ast.Visitor {
	switch unaryExpr.Op {
	case token.NOT, token.ADD, token.SUB:
		xVisitor := &visitor{}
		ast.Walk(xVisitor, unaryExpr.X)
		if xVisitor.err != nil {
			v.err = xVisitor.err
			return nil
		}
		switch unaryExpr.Op {
		case token.NOT:
			res, _ := strconv.ParseBool(xVisitor.res)
			v.res = strconv.FormatBool(!res)
		case token.ADD:
			v.res, v.kind = xVisitor.res, xVisitor.kind
		case token.SUB:
			if strings.HasPrefix(xVisitor.res, "-") {
				v.res, v.kind = strings.TrimPrefix(xVisitor.res, "-"), xVisitor.kind
				return nil
			}
			v.res, v.kind = "-"+xVisitor.res, xVisitor.kind
		}
	default:
		v.err = ErrUnsupportedOperator
	}
	return nil
}

func (v *visitor) arithmetic(xVisitor, yVisitor *visitor, op token.Token) {
	var x, y interface{}
	if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
		x, _ = strconv.ParseFloat(xVisitor.res, 64)
		y, _ = strconv.ParseFloat(yVisitor.res, 64)
	} else {
		x, _ = strconv.Atoi(xVisitor.res)
		y, _ = strconv.Atoi(yVisitor.res)
	}

	_, ok := x.(float64)
	switch op {
	case token.ADD:
		if ok {
			v.res, v.kind = fmt.Sprintf("%f", x.(float64)+y.(float64)), token.FLOAT
			return
		}
		v.res, v.kind = fmt.Sprintf("%d", x.(int)+y.(int)), token.INT
	case token.SUB:
		if ok {
			v.res, v.kind = fmt.Sprintf("%f", x.(float64)-y.(float64)), token.FLOAT
			return
		}
		v.res, v.kind = fmt.Sprintf("%d", x.(int)-y.(int)), token.INT
	case token.MUL:
		if ok {
			v.res, v.kind = fmt.Sprintf("%f", x.(float64)*y.(float64)), token.FLOAT
			return
		}
		v.res, v.kind = fmt.Sprintf("%d", x.(int)*y.(int)), token.INT
	case token.QUO:
		if ok {
			v.res, v.kind = fmt.Sprintf("%f", x.(float64)/y.(float64)), token.FLOAT
			return
		}
		v.res, v.kind = fmt.Sprintf("%d", x.(int)/y.(int)), token.INT
	case token.REM:
		if ok {
			v.err = fmt.Errorf("operator %s is not supported on untyped float: %w", "%", ErrInvalidOperationOnFloat)
		}
		v.res, v.kind = fmt.Sprintf("%d", x.(int)%y.(int)), token.INT
	}
}

func (v *visitor) comparison(xVisitor, yVisitor *visitor, op token.Token) {
	switch op {
	case token.EQL:
		v.res = strconv.FormatBool(xVisitor.res == yVisitor.res)
		return
	case token.NEQ:
		v.res = strconv.FormatBool(xVisitor.res != yVisitor.res)
		return
	}

	if xVisitor.kind == token.STRING || yVisitor.kind == token.STRING {
		v.res = strconv.FormatBool(xVisitor.res > yVisitor.res)
		return
	}

	var x, y interface{}
	if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
		x, _ = strconv.ParseFloat(xVisitor.res, 64)
		y, _ = strconv.ParseFloat(yVisitor.res, 64)
	} else {
		x, _ = strconv.Atoi(xVisitor.res)
		y, _ = strconv.Atoi(yVisitor.res)
	}

	_, ok := x.(float64)
	switch op {
	case token.GTR:
		if ok {
			v.res = strconv.FormatBool(x.(float64) > y.(float64))
			return
		}
		v.res = strconv.FormatBool(x.(int) > y.(int))
	case token.GEQ:
		if ok {
			v.res = strconv.FormatBool(x.(float64) >= y.(float64))
			return
		}
		v.res = strconv.FormatBool(x.(int) >= y.(int))
	case token.LSS:
		if ok {
			v.res = strconv.FormatBool(x.(float64) < y.(float64))
			return
		}
		v.res = strconv.FormatBool(x.(int) < y.(int))
	case token.LEQ:
		if ok {
			v.res = strconv.FormatBool(x.(float64) <= y.(float64))
			return
		}
		v.res = strconv.FormatBool(x.(int) <= y.(int))
	}
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
		xVisitor := &visitor{}
		ast.Walk(xVisitor, d.X)
		if xVisitor.err != nil {
			v.err = xVisitor.err
			return nil
		}
		yVisitor := &visitor{}
		ast.Walk(yVisitor, d.Y)
		if yVisitor.err != nil {
			v.err = yVisitor.err
			return nil
		}
		switch d.Op {
		case token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ:
			v.comparison(xVisitor, yVisitor, d.Op)
		case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
			v.arithmetic(xVisitor, yVisitor, d.Op)
		case token.LAND, token.LOR:
			var x, y bool
			x, v.err = strconv.ParseBool(xVisitor.res)
			if v.err != nil {
				return nil
			}
			y, v.err = strconv.ParseBool(yVisitor.res)
			if v.err != nil {
				return nil
			}
			if d.Op == token.LAND {
				v.res = strconv.FormatBool(x && y)
				return nil
			}
			v.res = strconv.FormatBool(x || y)
		default:
			v.err = ErrUnsupportedOperator
		}
		return nil
	case *ast.BasicLit:
		v.res = d.Value
		v.kind = d.Kind
		return nil
	case *ast.Ident:
		v.res = d.String()
		v.kind = token.STRING
		return nil
	}

	return v
}

func (v *visitor) Result() (bool, error) {
	if v.err != nil {
		return false, v.err
	}
	return strconv.ParseBool(v.res)
}
