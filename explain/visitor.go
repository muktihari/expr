package explain

import (
	"fmt"
	"go/ast"

	"github.com/muktihari/expr"
	"github.com/muktihari/expr/conv"
)

type exprType int

const (
	illegalExpr exprType = iota
	parentExpr
	unaryExpr
	binaryExpr
	basicLit
	ident
)

// Transform holds transformation of operations result.
type Transform struct {
	Segmented      string
	EquivalentForm string
	Evaluated      string
}

// Visitor satisfies ast.Visitor and it will evaluate given expression in step-by-step transformation values.
type Visitor struct {
	transforms []Transform
	depth      int
	exprType   exprType
	value      string
	err        error
}

var _ ast.Visitor = &Visitor{}

// Value return []Transform
func (v *Visitor) Value() []Transform { return v.transforms }

// Err returns visitor's error
func (v *Visitor) Err() error { return v.err }

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil || v.err != nil {
		return nil
	}

	switch d := node.(type) {
	case *ast.ParenExpr:
		v.exprType = parentExpr

		vx := &Visitor{depth: v.depth + 1}
		ast.Walk(vx, d.X)
		if vx.err != nil {
			v.err = vx.err
			return nil
		}

		v.value = "(" + vx.value + ")"
		v.transforms = append(v.transforms, vx.transforms...)
	case *ast.UnaryExpr:
		v.exprType = unaryExpr

		vx := &Visitor{depth: v.depth + 1}
		ast.Walk(vx, d.X)
		if vx.err != nil {
			v.err = vx.err
			return nil
		}
		v.transforms = append(v.transforms, vx.transforms...)

		ev := newExprVisitor()
		ast.Walk(ev, d)
		if err := ev.Err(); err != nil {
			v.err = err
			return nil
		}

		rv := conv.FormatExpr(d)
		v.value = ev.Value()

		v.transforms = append(v.transforms, Transform{
			Segmented:      rv,
			EquivalentForm: rv,
			Evaluated:      ev.Value(),
		})
	case *ast.BinaryExpr:
		v.exprType = binaryExpr

		vx := &Visitor{depth: v.depth + 1}
		ast.Walk(vx, d.X)
		if vx.err != nil {
			v.err = vx.err
			return nil
		}
		v.transforms = append(v.transforms, vx.transforms...)

		vy := &Visitor{depth: v.depth + 1}
		ast.Walk(vy, d.Y)
		if vy.err != nil {
			v.err = vy.err
			return nil
		}
		v.transforms = append(v.transforms, vy.transforms...)

		if vx.exprType == binaryExpr {
			vx.value = "(" + conv.FormatExpr(d.X) + ")"
		}
		if vy.exprType == binaryExpr {
			vy.value = "(" + conv.FormatExpr(d.Y) + ")"
		}

		// evaluated values
		ev := newExprVisitor()
		ast.Walk(ev, d)
		if err := ev.Err(); err != nil {
			v.err = err
			return nil
		}

		evx := newExprVisitor()
		ast.Walk(evx, d.X)

		evy := newExprVisitor()
		ast.Walk(evy, d.Y)

		v.value = fmt.Sprintf("%s %s %s", evx.Value(), d.Op, evy.Value())

		v.transforms = append(v.transforms, Transform{
			Segmented:      fmt.Sprintf("%s %s %s", vx.value, d.Op, vy.value),
			EquivalentForm: v.value,
			Evaluated:      ev.Value(),
		})
	case *ast.BasicLit:
		v.value, v.exprType = d.Value, basicLit
	case *ast.Ident:
		v.value, v.exprType = d.Name, ident
	}

	return nil
}

func newExprVisitor() *expr.Visitor {
	return expr.NewVisitor(
		expr.WithNumericType(expr.NumericTypeAuto),
		expr.WithAllowIntegerDividedByZero(true),
	)
}
