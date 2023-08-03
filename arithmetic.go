package expr

import (
	"go/ast"
	"go/token"
	"math"
	"strconv"

	"github.com/muktihari/expr/conv"
)

func arithmetic(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	if vx.kind <= numeric_beg || vx.kind >= numeric_end {
		v.err = newArithmeticNonNumericError(vx, binaryExpr.X)
		return
	}
	if vy.kind <= numeric_beg || vy.kind >= numeric_end {
		v.err = newArithmeticNonNumericError(vy, binaryExpr.Y)
		return
	}

	switch v.options.numericType {
	case NumericTypeComplex:
		calculateComplex(v, vx, vy, binaryExpr)
		return
	case NumericTypeFloat:
		calculateFloat(v, vx, vy, binaryExpr)
		return
	case NumericTypeInt:
		calculateInt(v, vx, vy, binaryExpr)
		return
	}

	// Auto figure out value type
	if vx.kind == KindImag || vy.kind == KindImag {
		calculateComplex(v, vx, vy, binaryExpr)
		return
	}

	if vx.kind == KindFloat || vy.kind == KindFloat {
		calculateFloat(v, vx, vy, binaryExpr)
		return
	}

	calculateInt(v, vx, vy, binaryExpr)
}

func newArithmeticNonNumericError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + v.value + "\" which is not a number",
		Pos: v.pos,
		Err: ErrArithmeticOperation,
	}
}

func calculateComplex(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindImag
	x := parseComplex(vx.value, vx.kind)
	y := parseComplex(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.ADD:
		v.value = strconv.FormatComplex(x+y, 'f', -1, 128)
	case token.SUB:
		v.value = strconv.FormatComplex(x-y, 'f', -1, 128)
	case token.MUL:
		v.value = strconv.FormatComplex(x*y, 'f', -1, 128)
	case token.QUO:
		v.value = strconv.FormatComplex(x/y, 'f', -1, 128)
	case token.REM:
		v.kind = KindIllegal
		v.err = &SyntaxError{
			Msg: "operator \"" + binaryExpr.Op.String() + "\" is not supported to do arithmetic on complex number",
			Pos: int(binaryExpr.OpPos),
			Err: ErrArithmeticOperation,
		}
		return
	}
}

func calculateFloat(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindFloat
	x := parseFloat(vx.value, vx.kind)
	y := parseFloat(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.ADD:
		v.value = strconv.FormatFloat(x+y, 'f', -1, 64)
	case token.SUB:
		v.value = strconv.FormatFloat(x-y, 'f', -1, 64)
	case token.MUL:
		v.value = strconv.FormatFloat(x*y, 'f', -1, 64)
	case token.QUO:
		v.value = strconv.FormatFloat(x/y, 'f', -1, 64)
	case token.REM:
		v.value = strconv.FormatFloat(math.Mod(x, y), 'f', -1, 64)
	}
}

func calculateInt(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	v.kind = KindInt
	x := parseInt(vx.value, vx.kind)
	y := parseInt(vy.value, vy.kind)

	switch binaryExpr.Op {
	case token.ADD:
		v.value = strconv.FormatInt(x+y, 10)
	case token.SUB:
		v.value = strconv.FormatInt(x-y, 10)
	case token.MUL:
		v.value = strconv.FormatInt(x*y, 10)
	case token.QUO:
		if y == 0 {
			if v.options.allowIntegerDividedByZero {
				v.value = "0"
				return
			}
			v.kind = KindIllegal
			v.err = &SyntaxError{
				Msg: "could not divide x with zero y, allowIntegerDividedByZero == false",
				Pos: vy.pos,
				Err: ErrIntegerDividedByZero,
			}
			return
		}
		v.value = strconv.FormatInt(x/y, 10)
	case token.REM:
		v.value = strconv.FormatInt(x%y, 10)
	}
}

func parseComplex(s string, kind Kind) complex128 {
	switch kind {
	case KindImag, KindFloat:
		v, _ := strconv.ParseComplex(s, 128)
		return v
	}

	v, _ := strconv.ParseInt(s, 0, 64)
	return complex(float64(v), 0)
}

func parseFloat(s string, kind Kind) float64 {
	switch kind {
	case KindImag:
		v, _ := strconv.ParseComplex(s, 128)
		return real(v)
	case KindFloat:
		v, _ := strconv.ParseFloat(s, 64)
		return v
	default: // INT: 0xFF, 0b1010, 0o77, 071, 90
		v, _ := strconv.ParseInt(s, 0, 64)
		return float64(v)
	}
}

func parseInt(s string, kind Kind) int64 {
	switch kind {
	case KindImag:
		v, _ := strconv.ParseComplex(s, 128)
		return int64(real(v))
	case KindFloat:
		v, _ := strconv.ParseFloat(s, 64)
		return int64(v)
	default: // INT: 0xFF, 0b1010, 0o77, 071, 90
		v, _ := strconv.ParseInt(s, 0, 64)
		return v
	}
}
