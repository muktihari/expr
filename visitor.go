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

// Kind of value (value's type)
type Kind byte

const (
	KindIllegal Kind = iota
	KindBoolean      // true false

	// Identifiers of numeric type
	numeric_beg
	KindInt   // 12345
	KindFloat // 123.45
	KindImag  // 123.45i
	numeric_end

	KindString // "abc" 'abc' `abc`
)

var kinds = [...]string{
	KindIllegal: "KindIllegal",
	KindBoolean: "KindBoolean",
	KindInt:     "KindInt",
	KindFloat:   "KindFloat",
	KindImag:    "KindImag",
	KindString:  "KindString",
}

func (k Kind) String() string {
	if k < Kind(len(kinds)) {
		return kinds[k]
	}
	return "kind(" + strconv.Itoa(int(k)) + ")"
}

type Option interface{ apply(o *options) }

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

func WithAllowIntegerDividedByZero(v bool) Option {
	return fnApply(func(o *options) { o.allowIntegerDividedByZero = v })
}

func WithNumericType(v NumericType) Option {
	return fnApply(func(o *options) { o.numericType = v })
}

type fnApply func(o *options)

func (f fnApply) apply(o *options) { f(o) }

var _ ast.Visitor = &Visitor{}

// Visitor satisfies ast.Visitor interface.
type Visitor struct {
	value   interface{}
	err     error
	pos     int
	kind    Kind    // used to reduce the need to do type assertion on numeric operation.
	options options // Visitor's Option
}

func defaultOptions() *options {
	return &options{
		allowIntegerDividedByZero: true,
		numericType:               NumericTypeAuto,
	}
}

// NewVisitor create new Visitor. If Option is not specified, these following default options will be set:
//   - allowIntegerDividedByZero: true
//   - numericType:               NumericTypeAuto
func NewVisitor(opts ...Option) *Visitor {
	options := defaultOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	return &Visitor{
		options: *options,
	}
}

// Value returns visitor's value in string
func (v *Visitor) Value() string { return fmt.Sprintf("%v", v.value) }

// ValueAny returns visitor's value interface{}
func (v *Visitor) ValueAny() interface{} { return v.value }

// Kind returns visitor's kind
func (v *Visitor) Kind() Kind { return v.kind }

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

var visitorPool = sync.Pool{
	New: func() interface{} {
		return new(Visitor)
	},
}

func (v *Visitor) visitUnary(unaryExpr *ast.UnaryExpr) ast.Visitor {
	switch unaryExpr.Op {
	case token.NOT, token.ADD, token.SUB:
		vx := visitorPool.Get().(*Visitor)
		defer visitorPool.Put(vx)
		vx.reset()
		vx.options = v.options

		ast.Walk(vx, unaryExpr.X)
		if vx.err != nil {
			v.err = vx.err
			return nil
		}

		v.kind = vx.kind
		switch unaryExpr.Op {
		case token.NOT: // negation: !true -> false, !false -> true
			res, ok := vx.value.(bool)
			if !ok {
				s := conv.FormatExpr(unaryExpr.X)
				v.err = &SyntaxError{
					Msg: "could not do negation: result of \"" + s + "\" is \"" + fmt.Sprintf("%v", vx.value) + "\" not a boolean",
					Pos: vx.pos,
					Err: ErrUnaryOperation,
				}
				return nil
			}
			v.value = !res
		case token.ADD:
			v.value = vx.value
		case token.SUB:
			switch val := vx.value.(type) {
			case complex128:
				v.value = val * -1
			case float64:
				v.value = val * -1
			case int64:
				v.value = val * -1
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
	vx := visitorPool.Get().(*Visitor)
	defer visitorPool.Put(vx)
	vx.reset()
	vx.options = v.options

	ast.Walk(vx, binaryExpr.X)
	if vx.err != nil {
		v.err = vx.err
		return nil
	}

	vy := visitorPool.Get().(*Visitor)
	defer visitorPool.Put(vy)
	vy.reset()
	vy.options = v.options

	ast.Walk(vy, binaryExpr.Y)
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
		v.kind = KindInt
		v.value, _ = strconv.ParseInt(basicLit.Value, 0, 64)
	case token.FLOAT:
		v.kind = KindFloat
		v.value, _ = strconv.ParseFloat(basicLit.Value, 64)
	case token.IMAG:
		v.kind = KindImag
		v.value, _ = strconv.ParseComplex(basicLit.Value, 128)
	case token.CHAR:
		fallthrough // treat as string
	case token.STRING:
		v.kind = KindString
		v.value = strings.TrimFunc(basicLit.Value, func(r rune) bool { return r == '\'' || r == '`' || r == '"' })
	}
	return nil
}

func (v *Visitor) visitIndent(indent *ast.Ident) ast.Visitor {
	v.kind, v.value = KindString, indent.String()
	vb, err := strconv.ParseBool(indent.String())
	if err != nil {
		return nil
	}
	v.kind, v.value = KindBoolean, vb
	return nil
}

func (v *Visitor) reset() {
	v.err = nil
	v.kind = KindIllegal
	v.pos = 0
	v.value = nil
}
