package expr

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/muktihari/expr/conv"
)

func logical(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	if vx.kind != KindBoolean {
		v.err = newLogicalNonBooleanError(vx, binaryExpr.X)
		return
	}

	if vy.kind != KindBoolean {
		v.err = newLogicalNonBooleanError(vy, binaryExpr.Y)
		return
	}

	x := vx.value.(bool)
	y := vy.value.(bool)

	v.kind = KindBoolean
	if binaryExpr.Op == token.LAND {
		v.value = x && y
		return
	}

	v.value = x || y // token.LOR
}

func newLogicalNonBooleanError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value) + "\" which is not a boolean",
		Pos: v.pos,
		Err: ErrLogicalOperation,
	}
}
