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

func parseUnaryExpr(unaryExpr *ast.UnaryExpr) (string, token.Token, error) {
	switch unaryExpr.Op {
	case token.NOT:
		v := &visitor{}
		ast.Walk(v, unaryExpr.X)
		if v.err != nil {
			return "", v.kind, v.err
		}
		res, _ := strconv.ParseBool(v.res)
		return strconv.FormatBool(!res), 0, nil
	case token.ADD:
		v := &visitor{}
		ast.Walk(v, unaryExpr.X)
		if v.err != nil {
			return "", v.kind, v.err
		}
		return v.res, v.kind, nil
	case token.SUB:
		v := &visitor{}
		ast.Walk(v, unaryExpr.X)
		if v.err != nil {
			return "", v.kind, v.err
		}

		if strings.HasPrefix(v.res, "-") {
			if v.kind == token.FLOAT {
				value, _ := strconv.ParseFloat(v.res, 64)
				return fmt.Sprintf("%f", value*-1), v.kind, nil
			}
			value, _ := strconv.Atoi(v.res)
			return fmt.Sprintf("%d", value*-1), v.kind, nil
		}

		return "-" + v.res, v.kind, nil
	default:
		return "", 0, ErrUnsupportedOperator
	}
}

func tryArithmetic(xVisitor, yVisitor *visitor, op token.Token) (string, token.Token, error) {
	switch op {
	case token.ADD:
		if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
			x, _ := strconv.ParseFloat(xVisitor.res, 64)
			y, _ := strconv.ParseFloat(yVisitor.res, 64)
			return fmt.Sprintf("%f", x+y), token.FLOAT, nil
		}
		x, _ := strconv.Atoi(xVisitor.res)
		y, _ := strconv.Atoi(yVisitor.res)
		return fmt.Sprintf("%d", x+y), token.INT, nil
	case token.SUB:
		if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
			x, _ := strconv.ParseFloat(xVisitor.res, 64)
			y, _ := strconv.ParseFloat(yVisitor.res, 64)
			return fmt.Sprintf("%f", x-y), token.FLOAT, nil
		}
		x, _ := strconv.Atoi(xVisitor.res)
		y, _ := strconv.Atoi(yVisitor.res)
		return fmt.Sprintf("%d", x-y), token.INT, nil
	case token.MUL:
		if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
			x, _ := strconv.ParseFloat(xVisitor.res, 64)
			y, _ := strconv.ParseFloat(yVisitor.res, 64)
			return fmt.Sprintf("%f", x*y), token.FLOAT, nil
		}
		x, _ := strconv.Atoi(xVisitor.res)
		y, _ := strconv.Atoi(yVisitor.res)
		return fmt.Sprintf("%d", x*y), token.INT, nil
	case token.QUO:
		if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
			x, _ := strconv.ParseFloat(xVisitor.res, 64)
			y, _ := strconv.ParseFloat(yVisitor.res, 64)
			return fmt.Sprintf("%f", x/y), token.FLOAT, nil
		}
		x, _ := strconv.Atoi(xVisitor.res)
		y, _ := strconv.Atoi(yVisitor.res)
		return fmt.Sprintf("%d", x/y), token.INT, nil
	case token.REM:
		if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
			return "", 0, fmt.Errorf("operator %s is not supported on untyped float: %w", "%", ErrInvalidOperationOnFloat)
		}
		x, _ := strconv.Atoi(xVisitor.res)
		y, _ := strconv.Atoi(yVisitor.res)
		return fmt.Sprintf("%d", x%y), token.INT, nil
	default:
		return "", 0, ErrUnsupportedOperator
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
		v.res, v.kind, v.err = parseUnaryExpr(d)
		return nil
	case *ast.BinaryExpr:
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
			if xVisitor.kind != token.STRING && yVisitor.kind != token.STRING {
				if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
					x, _ := strconv.ParseFloat(xVisitor.res, 64)
					y, _ := strconv.ParseFloat(yVisitor.res, 64)

					v.res = strconv.FormatBool(x > y)
					return nil
				}
				x, _ := strconv.Atoi(xVisitor.res)
				y, _ := strconv.Atoi(yVisitor.res)
				v.res = strconv.FormatBool(x > y)
				return nil
			}
			v.res = strconv.FormatBool(x > y)
		case token.GEQ:
			if xVisitor.kind != token.STRING && yVisitor.kind != token.STRING {
				if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
					x, _ := strconv.ParseFloat(xVisitor.res, 64)
					y, _ := strconv.ParseFloat(yVisitor.res, 64)

					v.res = strconv.FormatBool(x >= y)
					return nil
				}
				x, _ := strconv.Atoi(xVisitor.res)
				y, _ := strconv.Atoi(yVisitor.res)
				v.res = strconv.FormatBool(x >= y)
				return nil
			}
			v.res = strconv.FormatBool(x >= y)
		case token.LSS:
			if xVisitor.kind != token.STRING && yVisitor.kind != token.STRING {
				if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
					x, _ := strconv.ParseFloat(xVisitor.res, 64)
					y, _ := strconv.ParseFloat(yVisitor.res, 64)

					v.res = strconv.FormatBool(x < y)
					return nil
				}
				x, _ := strconv.Atoi(xVisitor.res)
				y, _ := strconv.Atoi(yVisitor.res)
				v.res = strconv.FormatBool(x < y)
				return nil
			}
			v.res = strconv.FormatBool(x < y)
		case token.LEQ:
			if xVisitor.kind != token.STRING && yVisitor.kind != token.STRING {
				if xVisitor.kind == token.FLOAT || yVisitor.kind == token.FLOAT {
					x, _ := strconv.ParseFloat(xVisitor.res, 64)
					y, _ := strconv.ParseFloat(yVisitor.res, 64)

					v.res = strconv.FormatBool(x <= y)
					return nil
				}
				x, _ := strconv.Atoi(xVisitor.res)
				y, _ := strconv.Atoi(yVisitor.res)
				v.res = strconv.FormatBool(x <= y)
				return nil
			}
			v.res = strconv.FormatBool(x <= y)
		case token.LAND:
			xbool, err := strconv.ParseBool(x)
			if err != nil {
				v.err = err
				return nil
			}
			ybool, err := strconv.ParseBool(y)
			if err != nil {
				v.err = err
				return nil
			}
			v.res = strconv.FormatBool(xbool && ybool)
		case token.LOR:
			xbool, err := strconv.ParseBool(x)
			if err != nil {
				v.err = err
				return nil
			}
			ybool, err := strconv.ParseBool(y)
			if err != nil {
				v.err = err
				return nil
			}
			v.res = strconv.FormatBool(xbool || ybool)
		default:
			v.res, v.kind, v.err = tryArithmetic(xVisitor, yVisitor, d.Op)
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
	var res bool
	res, v.err = strconv.ParseBool(v.res)
	return res, v.err
}
