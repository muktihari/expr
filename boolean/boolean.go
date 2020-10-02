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
	if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
		x, _ := strconv.ParseFloat(xVisitor.res, 64)
		y, _ := strconv.ParseFloat(yVisitor.res, 64)
		switch op {
		case token.ADD:
			v.res, v.kind = fmt.Sprintf("%f", x+y), token.FLOAT
		case token.SUB:
			v.res, v.kind = fmt.Sprintf("%f", x-y), token.FLOAT
		case token.MUL:
			v.res, v.kind = fmt.Sprintf("%f", x*y), token.FLOAT
		case token.QUO:
			v.res, v.kind = fmt.Sprintf("%f", x/y), token.FLOAT
		case token.REM:
			v.res, v.kind = fmt.Sprintf("%f", x+y), token.FLOAT
		}
		return
	}

	x, _ := strconv.Atoi(xVisitor.res)
	y, _ := strconv.Atoi(yVisitor.res)
	switch op {
	case token.ADD:
		v.res, v.kind = fmt.Sprintf("%d", x+y), token.INT
	case token.SUB:
		v.res, v.kind = fmt.Sprintf("%d", x-y), token.INT
	case token.MUL:
		v.res, v.kind = fmt.Sprintf("%d", x*y), token.INT
	case token.QUO:
		v.res, v.kind = fmt.Sprintf("%d", x/y), token.INT
	case token.REM:
		v.res, v.kind = fmt.Sprintf("%d", x%y), token.INT
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
		switch op {
		case token.GTR:
			v.res = strconv.FormatBool(xVisitor.res > yVisitor.res)
		case token.GEQ:
			v.res = strconv.FormatBool(xVisitor.res >= yVisitor.res)
		case token.LSS:
			v.res = strconv.FormatBool(xVisitor.res > yVisitor.res)
		case token.LEQ:
			v.res = strconv.FormatBool(xVisitor.res >= yVisitor.res)
		}
		return
	}

	if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
		x, _ := strconv.ParseFloat(xVisitor.res, 64)
		y, _ := strconv.ParseFloat(yVisitor.res, 64)
		switch op {
		case token.GTR:
			v.res = strconv.FormatBool(x > y)
		case token.GEQ:
			v.res = strconv.FormatBool(x >= y)
		case token.LSS:
			v.res = strconv.FormatBool(x < y)
		case token.LEQ:
			v.res = strconv.FormatBool(x <= y)
		}
		return
	}

	x, _ := strconv.Atoi(xVisitor.res)
	y, _ := strconv.Atoi(yVisitor.res)
	switch op {
	case token.GTR:
		v.res = strconv.FormatBool(x > y)
	case token.GEQ:
		v.res = strconv.FormatBool(x >= y)
	case token.LSS:
		v.res = strconv.FormatBool(x < y)
	case token.LEQ:
		v.res = strconv.FormatBool(x <= y)
	}
}

func (v *visitor) logical(xVisitor, yVisitor *visitor, op token.Token) {
	var x, y bool
	x, v.err = strconv.ParseBool(xVisitor.res)
	if v.err != nil {
		return
	}
	y, v.err = strconv.ParseBool(yVisitor.res)
	if v.err != nil {
		return
	}
	if op == token.LAND {
		v.res = strconv.FormatBool(x && y)
		return
	}
	v.res = strconv.FormatBool(x || y) // token.LOR
}

func (v *visitor) visitBinary(binaryExpr *ast.BinaryExpr) ast.Visitor {
	xVisitor := &visitor{}
	ast.Walk(xVisitor, binaryExpr.X)
	if xVisitor.err != nil {
		v.err = xVisitor.err
		return nil
	}
	yVisitor := &visitor{}
	ast.Walk(yVisitor, binaryExpr.Y)
	if yVisitor.err != nil {
		v.err = yVisitor.err
		return nil
	}

	switch binaryExpr.Op {
	case token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ:
		v.comparison(xVisitor, yVisitor, binaryExpr.Op)
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		v.arithmetic(xVisitor, yVisitor, binaryExpr.Op)
	case token.LAND, token.LOR:
		v.logical(xVisitor, yVisitor, binaryExpr.Op)
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
