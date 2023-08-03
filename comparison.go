package expr

import (
	"go/ast"
	"go/token"
	"strconv"
)

func comparison(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindBoolean
	if vx.kind == KindBoolean || vy.kind == KindBoolean {
		compareBoolean(v, vx, vy, binaryExpr)
		return
	}
	if vx.kind == KindString || vy.kind == KindString {
		compareString(v, vx, vy, binaryExpr)
		return
	}
	if vx.kind == KindImag || vy.kind == KindImag {
		compareComplex(v, vx, vy, binaryExpr)
		return
	}
	if vx.kind == KindFloat || vy.kind == KindFloat {
		compareFloat(v, vx, vy, binaryExpr)
		return
	}

	compareInt(v, vx, vy, binaryExpr)
}

func compareBoolean(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	switch binaryExpr.Op {
	case token.EQL:
		v.value = strconv.FormatBool(vx.value == vy.value)
	case token.NEQ:
		v.value = strconv.FormatBool(vx.value != vy.value)
	default:
		v.kind = KindIllegal
		v.err = &SyntaxError{
			Msg: "operator \"" + binaryExpr.Op.String() + "\" is not supported for comparing boolean values",
			Pos: int(binaryExpr.OpPos),
			Err: ErrUnsupportedOperator,
		}
		return
	}
}

func compareString(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	switch binaryExpr.Op {
	case token.EQL:
		v.value = strconv.FormatBool(vx.value == vy.value)
	case token.NEQ:
		v.value = strconv.FormatBool(vx.value != vy.value)
	case token.GTR:
		v.value = strconv.FormatBool(vx.value > vy.value)
	case token.GEQ:
		v.value = strconv.FormatBool(vx.value >= vy.value)
	case token.LSS:
		v.value = strconv.FormatBool(vx.value < vy.value)
	case token.LEQ:
		v.value = strconv.FormatBool(vx.value <= vy.value)
	}
}

// IEEE 754 says that only NaNs satisfy f != f.
func compareComplex(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := parseComplex(vx.value, vx.kind)
	y := parseComplex(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = strconv.FormatBool(x == y)
	case token.NEQ:
		v.value = strconv.FormatBool(x != y)
	default:
		v.kind = KindIllegal
		v.err = &SyntaxError{
			Msg: "operator \"" + binaryExpr.Op.String() + "\" is not supported for comparing complex numbers",
			Pos: int(binaryExpr.OpPos),
			Err: ErrUnsupportedOperator,
		}
		return
	}
}

func compareFloat(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x := parseFloat(vx.value, vx.kind)
	y := parseFloat(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = strconv.FormatBool(x == y)
	case token.NEQ:
		v.value = strconv.FormatBool(x != y)
	case token.GTR:
		v.value = strconv.FormatBool(x > y)
	case token.GEQ:
		v.value = strconv.FormatBool(x >= y)
	case token.LSS:
		v.value = strconv.FormatBool(x < y)
	case token.LEQ:
		v.value = strconv.FormatBool(x <= y)
	}
}

func compareInt(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	x, _ := strconv.ParseInt(vx.value, 0, 64)
	y, _ := strconv.ParseInt(vy.value, 0, 64)

	switch binaryExpr.Op {
	case token.EQL:
		v.value = strconv.FormatBool(x == y)
	case token.NEQ:
		v.value = strconv.FormatBool(x != y)
	case token.GTR:
		v.value = strconv.FormatBool(x > y)
	case token.GEQ:
		v.value = strconv.FormatBool(x >= y)
	case token.LSS:
		v.value = strconv.FormatBool(x < y)
	case token.LEQ:
		v.value = strconv.FormatBool(x <= y)
	}
}
