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
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"testing"
)

func TestArithmetic(t *testing.T) {
	newVisitor := func(numericType NumericType) *Visitor {
		return &Visitor{options: options{numericType: numericType}}
	}
	tt := []struct {
		name          string
		v, vx, vy     *Visitor
		op            token.Token
		expectedValue value
		expectedErr   error
	}{
		{
			name:          "arithmetic non numeric error: x numeric y boolean",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: int64Value(1)},
			vy:            &Visitor{value: boolValue(true)},
			op:            token.ADD,
			expectedValue: value{},
			expectedErr:   ErrArithmeticOperation,
		},
		{
			name:          "arithmetic non numeric error: x boolean y numeric",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: boolValue(false)},
			vy:            &Visitor{value: int64Value(10)},
			op:            token.ADD,
			expectedValue: value{},
			expectedErr:   ErrArithmeticOperation,
		},
		{
			name:          "arithmetic numeric complex",
			v:             newVisitor(NumericTypeComplex),
			vx:            &Visitor{value: complex128Value(1 + 0i)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: complex128Value(3 + 0i),
		},
		{
			name:          "arithmetic numeric float",
			v:             newVisitor(NumericTypeFloat),
			vx:            &Visitor{value: float64Value(1.5)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: float64Value(3.5),
		},
		{
			name:          "arithmetic numeric int",
			v:             newVisitor(NumericTypeInt),
			vx:            &Visitor{value: float64Value(1.5)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: int64Value(3),
		},
		{
			name:          "arithmetic numeric auto: x imag y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: complex128Value(1.5 + 1i)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: complex128Value(3.5 + 1i),
		},
		{
			name:          "arithmetic numeric auto: x int y imag",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: int64Value(2)},
			vy:            &Visitor{value: complex128Value(1.5 + 1i)},
			op:            token.ADD,
			expectedValue: complex128Value(3.5 + 1i),
		},
		{
			name:          "arithmetic numeric auto: x float y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: float64Value(1.5)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: float64Value(3.5),
		},
		{
			name:          "arithmetic numeric auto: x int y float",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: float64Value(1.5)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: float64Value(3.5),
		},
		{
			name:          "arithmetic numeric auto: x int y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: int64Value(1)},
			vy:            &Visitor{value: int64Value(2)},
			op:            token.ADD,
			expectedValue: float64Value(3),
		},
	}

	for i, tc := range tt {
		if i < 2 {
			continue
		}
		tc := tc
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			be := &ast.BinaryExpr{Op: tc.op}
			arithmetic(tc.v, tc.vx, tc.vy, be)
			if !errors.Is(tc.v.err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, tc.v.err)
			}
			if tc.v.value.Any() != tc.expectedValue.Any() {
				t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValue.Any(), tc.expectedValue.Any(),
					tc.v.value, tc.v.value)
			}
		})
	}
}

func TestCalculateComplex(t *testing.T) {
	newComplexVisitor := func() *Visitor {
		return &Visitor{options: options{numericType: NumericTypeComplex}}
	}
	tt := []struct {
		name           string
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []value
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: int64Value(1)},
			vy:             &Visitor{value: int64Value(2)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []value{complex128Value(3 + 0i), complex128Value(-1 + 0i), complex128Value(2 + 0i), complex128Value(0.5 + 0i)},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "calculate floats",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: float64Value(1.0)},
			vy:             &Visitor{value: float64Value(2.0)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []value{complex128Value(3 + 0i), complex128Value(-1 + 0i), complex128Value(2 + 0i), complex128Value(0.5 + 0i)},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: complex128Value(1 + 1i)},
			vy:             &Visitor{value: complex128Value(2 + 1i)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []value{complex128Value(3 + 2i), complex128Value(-1 + 0i), complex128Value(1 + 3i), complex128Value(0.6 + 0.2i)},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "unsupported complex operation",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: complex128Value(1 + 1i)},
			vy:             &Visitor{value: complex128Value(2 + 1i)},
			ops:            []token.Token{token.REM},
			expectedValues: []value{{}},
			expectedErrs:   []error{ErrArithmeticOperation},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := fmt.Sprintf("%v%s%v", tc.vx.value.Any(), op, tc.vy.value.Any())
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateComplex(tc.v, parseComplex(tc.vx.value), parseComplex(tc.vy.value), be.Op, be.OpPos)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %v, got: %v", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value.Any() != tc.expectedValues[i].Any() {
						t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValues[i].Any(), tc.expectedValues[i].Any(),
							tc.v.value, tc.v.value)
					}
				})
			}
		})
	}
}

func TestCalculateFloat(t *testing.T) {
	newFloatVisitor := func() *Visitor {
		return &Visitor{options: options{numericType: NumericTypeFloat}}
	}

	tt := []struct {
		name           string
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []value
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: int64Value(10)},
			vy:             &Visitor{value: int64Value(3)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []value{float64Value(13), float64Value(7), float64Value(30), float64Value(3.3333333333333335), float64Value(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate floats",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: float64Value(10.0)},
			vy:             &Visitor{value: float64Value(3.0)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []value{float64Value(13), float64Value(7), float64Value(30), float64Value(3.3333333333333335), float64Value(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: complex128Value(10 + 1i)},
			vy:             &Visitor{value: complex128Value(3 + 1i)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []value{float64Value(13), float64Value(7), float64Value(30), float64Value(3.3333333333333335), float64Value(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := fmt.Sprintf("%v%s%v", tc.vx.value.Any(), op, tc.vy.value.Any())
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateFloat(tc.v, parseFloat(tc.vx.value), parseFloat(tc.vy.value), be.Op)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %v, got: %v", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value.Any() != tc.expectedValues[i].Any() {
						t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValues[i].Any(), tc.expectedValues[i].Any(),
							tc.v.value, tc.v.value)
					}
				})
			}
		})
	}
}

func TestCalculateInt(t *testing.T) {
	newIntVisitor := func() *Visitor {
		return &Visitor{options: options{
			numericType:               NumericTypeInt,
			allowIntegerDividedByZero: true,
		}}
	}

	tt := []struct {
		name           string
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []value
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newIntVisitor(),
			vx:             &Visitor{value: int64Value(10)},
			vy:             &Visitor{value: int64Value(3)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []value{int64Value(13), int64Value(7), int64Value(30), int64Value(3), int64Value(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate integers allowIntegerDividedByZero == true",
			v:              &Visitor{options: options{numericType: NumericTypeInt, allowIntegerDividedByZero: true}},
			vx:             &Visitor{value: int64Value(10)},
			vy:             &Visitor{value: int64Value(0)},
			ops:            []token.Token{token.QUO},
			expectedValues: []value{int64Value(0)},
			expectedErrs:   []error{nil},
		},
		{
			name:           "calculate integers allowIntegerDividedByZero == false",
			v:              &Visitor{options: options{numericType: NumericTypeInt, allowIntegerDividedByZero: false}},
			vx:             &Visitor{value: int64Value(10)},
			vy:             &Visitor{value: int64Value(0)},
			ops:            []token.Token{token.QUO},
			expectedValues: []value{{}},
			expectedErrs:   []error{ErrIntegerDividedByZero},
		},
		{
			name:           "calculate floats",
			v:              newIntVisitor(),
			vx:             &Visitor{value: float64Value(10.0)},
			vy:             &Visitor{value: float64Value(3.0)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []value{int64Value(13), int64Value(7), int64Value(30), int64Value(3), int64Value(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newIntVisitor(),
			vx:             &Visitor{value: complex128Value(10 + 1i)},
			vy:             &Visitor{value: complex128Value(3 + 1i)},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []value{int64Value(13), int64Value(7), int64Value(30), int64Value(3), int64Value(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := fmt.Sprintf("%v%s%v", tc.vx.value.Any(), op, tc.vy.value.Any())
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateInt(tc.v, parseInt(tc.vx.value), parseInt(tc.vy.value), tc.vy.pos, be.Op)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %v, got: %v", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value.Any() != tc.expectedValues[i].Any() {
						t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValues[i].Any(), tc.expectedValues[i].Any(),
							tc.v.value, tc.v.value)
					}
				})
			}
		})
	}
}

func TestParseInvalidValue(t *testing.T) {
	i64 := parseInt(stringValue("invalid"))
	if i64 != 0 {
		t.Fatalf("expected 0, got: %v", i64)
	}
	f64 := parseFloat(boolValue(true))
	if f64 != 0 {
		t.Fatalf("expected 0, got: %v", f64)
	}
	c128 := parseComplex(boolValue(false))
	if c128 != 0 {
		t.Fatalf("expected 0, got: %v", c128)
	}
}
