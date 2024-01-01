// Copyright 2023 The Expr Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package explain

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/muktihari/expr"
	"github.com/muktihari/expr/internal/conv"
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
	Explaination   string
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

		xValue := evx.Value()
		yValue := evy.Value()

		v.value = fmt.Sprintf("%s %s %s", xValue, d.Op, yValue)

		transform := Transform{
			Segmented:      fmt.Sprintf("%s %s %s", vx.value, d.Op, vy.value),
			EquivalentForm: v.value,
			Evaluated:      ev.Value(),
		}

		// Special case for explaining bitwise
		switch d.Op {
		case token.AND, token.OR, token.XOR, token.AND_NOT, token.SHL, token.SHR:
			explainBitwise(&transform, xValue, yValue, d.Op)
		}

		v.transforms = append(v.transforms, transform)
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

var operatorStringMap = map[token.Token]string{
	token.AND:     "AND",
	token.OR:      "OR",
	token.XOR:     "XOR",
	token.AND_NOT: "AND NOT",
}

func explainBitwise(transform *Transform, xValue, yValue string, op token.Token) {
	fx, _ := strconv.ParseFloat(xValue, 64)
	fy, _ := strconv.ParseFloat(yValue, 64)

	var result int64
	switch op {
	case token.AND:
		result = int64(fx) & int64(fy)
	case token.OR:
		result = int64(fx) | int64(fy)
	case token.XOR:
		result = int64(fx) ^ int64(fy)
	case token.AND_NOT:
		result = int64(fx) &^ int64(fy)
	}

	size := len(fmt.Sprintf("%.8b", int64(fx)))
	sizeY := len(fmt.Sprintf("%.8b", int64(fy)))
	if sizeY > size {
		size = sizeY
	}
	formatter := fmt.Sprintf("0b%%.%db", size)

	xbits := fmt.Sprintf(formatter, int64(fx))
	ybits := fmt.Sprintf(formatter, int64(fy))

	if op != token.SHL && op != token.SHR {
		transform.Explaination = fmt.Sprintf("%s\n%s\n%s %s\n%s",
			xbits, ybits,
			strings.Repeat("-", (size*2)-(size*2/10)), operatorStringMap[op],
			fmt.Sprintf(formatter, result))
		return
	}

	var shiftDirection string
	switch op {
	case token.SHL:
		shiftDirection = "left"
		result = int64(fx) << int64(fy)
	case token.SHR:
		shiftDirection = "right"
		result = int64(fx) >> int64(fy)
	}

	transform.Explaination = fmt.Sprintf("%s %s-shifted by %d = %s",
		xbits, shiftDirection, int64(fy), fmt.Sprintf(formatter, result))
}
