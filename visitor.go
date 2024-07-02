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

package expr

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"sync"

	"github.com/muktihari/expr/internal/conv"
)

var pool = sync.Pool{New: func() interface{} { return new(Visitor) }}

// NumericType determines what type of number represented in the expr string
type NumericType byte

const (
	NumericTypeAuto    NumericType = iota // [1 * 2 = 2]       [1 * 2.5 = 2.5]
	NumericTypeComplex                    // [1 * 2 = 2+0i]    [1 * (2+2i) = (2+2i)]    [(1+2i) * (2+2i) = (-2+6i)]
	NumericTypeFloat                      // [1 * 2 = 2.0]     [1 * 2.5 = 2.5]
	NumericTypeInt                        // [1 * 2 = 2,]      [1 * 2.5 = 2]
)

type options struct {
	allowIntegerDividedByZero bool        // true: 2/0 = 0, false: return error
	numericType               NumericType // treat numeric type as specific type
}

// Option is Visitor's option.
type Option func(o *options)

// WithAllowIntegerDividedByZero allows integer divided by zero operation.
func WithAllowIntegerDividedByZero(v bool) Option {
	return func(o *options) { o.allowIntegerDividedByZero = v }
}

// WithNumericType treats all numeric types as v.
func WithNumericType(v NumericType) Option {
	return func(o *options) { o.numericType = v }
}

var _ ast.Visitor = (*Visitor)(nil)

// Visitor satisfies ast.Visitor interface.
type Visitor struct {
	value   value
	err     error
	pos     int
	options options // Visitor's Option
}

func defaultOptions() options {
	return options{
		allowIntegerDividedByZero: true,
		numericType:               NumericTypeAuto,
	}
}

// NewVisitor create new Visitor. If Option is not specified, these following default options will be set:
//   - allowIntegerDividedByZero: true
//   - numericType:               NumericTypeAuto
func NewVisitor(opts ...Option) *Visitor {
	v := &Visitor{
		options: defaultOptions(),
	}
	for i := range opts {
		opts[i](&v.options)
	}
	return v
}

// Value returns visitor's value in string
func (v *Visitor) Value() string { return fmt.Sprintf("%v", v.value.Any()) }

// ValueAny returns visitor's value interface{}
func (v *Visitor) ValueAny() interface{} { return v.value.Any() }

// Kind returns visitor's kind
func (v *Visitor) Kind() Kind { return v.value.Kind() }

// Err returns visitor's error
func (v *Visitor) Err() error { return v.err }

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil || v.err != nil {
		return nil
	}
	v.pos = int(node.Pos())

	switch d := node.(type) {
	case *ast.ParenExpr:
		return v.Visit(d.X)
	case *ast.UnaryExpr:
		return v.visitUnary(d)
	case *ast.BinaryExpr:
		return v.visitBinary(d)
	case *ast.BasicLit: // handle type: int, float, imag, char, string
		return v.visitBasicLit(d)
	case *ast.Ident: // handle type: bolean, string without quotation
		return v.visitIndent(d)
	}

	return v
}

func (v *Visitor) visitUnary(unaryExpr *ast.UnaryExpr) ast.Visitor {
	switch unaryExpr.Op {
	case token.NOT, token.ADD, token.SUB:
		vx := pool.Get().(*Visitor)
		defer pool.Put(vx)
		vx.reset(v.options)

		vx.Visit(unaryExpr.X)
		if vx.err != nil {
			v.err = vx.err
			return nil
		}

		v.value.SetKind(vx.value.Kind())
		switch unaryExpr.Op {
		case token.NOT: // negation: !true -> false, !false -> true
			if vx.value.Kind() != KindBoolean {
				s := conv.FormatExpr(unaryExpr.X)
				v.err = &SyntaxError{
					Msg: "could not do negation: result of \"" + s + "\" is \"" + fmt.Sprintf("%v", vx.value) + "\" not a boolean",
					Pos: vx.pos,
					Err: ErrUnaryOperation,
				}
				return nil
			}
			v.value = boolValue(!vx.value.Bool())
		case token.ADD:
			v.value = vx.value
		case token.SUB:
			switch vx.value.Kind() {
			case KindInt:
				v.value = int64Value(vx.value.Int64() * -1)
			case KindFloat:
				v.value = float64Value(vx.value.Float64() * -1)
			case KindImag:
				v.value = complex128Value(vx.value.Complex128() * -1)
			}
		}
	default:
		v.err = &SyntaxError{
			Msg: "operator \"" + unaryExpr.Op.String() + "\" is unsupported",
			Pos: int(unaryExpr.OpPos),
			Err: ErrUnsupportedOperator,
		}
	}
	return nil
}

func (v *Visitor) visitBinary(binaryExpr *ast.BinaryExpr) ast.Visitor {
	vx := pool.Get().(*Visitor)
	defer pool.Put(vx)
	vx.reset(v.options)

	vx.Visit(binaryExpr.X)
	if vx.err != nil {
		v.err = vx.err
		return nil
	}

	vy := pool.Get().(*Visitor)
	defer pool.Put(vy)
	vy.reset(v.options)

	vy.Visit(binaryExpr.Y)
	if vy.err != nil {
		v.err = vy.err
		return nil
	}

	switch binaryExpr.Op {
	case token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ:
		comparison(v, vx, vy, binaryExpr)
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		arithmetic(v, vx, vy, binaryExpr)
	case token.AND, token.OR, token.XOR, token.AND_NOT, token.SHL, token.SHR:
		bitwise(v, vx, vy, binaryExpr)
	case token.LAND, token.LOR:
		logical(v, vx, vy, binaryExpr)
	}
	return nil
}

func (v *Visitor) visitBasicLit(basicLit *ast.BasicLit) ast.Visitor {
	switch basicLit.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(basicLit.Value, 0, 64)
		v.value = int64Value(val)
	case token.FLOAT:
		val, _ := strconv.ParseFloat(basicLit.Value, 64)
		v.value = float64Value(val)
	case token.IMAG:
		val, _ := strconv.ParseComplex(basicLit.Value, 128)
		v.value = complex128Value(val)
	case token.CHAR:
		fallthrough // treat as string
	case token.STRING:
		v.value = stringValue(strings.TrimFunc(basicLit.Value, func(r rune) bool { return r == '\'' || r == '`' || r == '"' }))
	}
	return nil
}

func (v *Visitor) visitIndent(indent *ast.Ident) ast.Visitor {
	v.value = stringValue(indent.String())
	vb, err := strconv.ParseBool(indent.String())
	if err != nil {
		return nil
	}
	v.value = boolValue(vb)
	return nil
}

func (v *Visitor) reset(o options) {
	v.value = value{}
	v.err = nil
	v.pos = 0
	v.options = o
}
