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

package expr_test

import (
	"errors"
	"go/ast"
	"go/parser"
	"testing"

	"github.com/muktihari/expr"
)

func TestVisit(t *testing.T) {
	tt := []struct {
		in            string
		expectedValue interface{}
		expectedKind  expr.Kind
		expectedErr   error
	}{
		{
			in:           "<-a", // unary operator
			expectedKind: expr.KindIllegal,
			expectedErr:  expr.ErrUnsupportedOperator,
		},
		{
			in:            "'a' == 'a'",
			expectedValue: true,
			expectedKind:  expr.KindBoolean,
		},
		{
			in:            "-(20+10i)",
			expectedValue: -(20 + 10i),
			expectedKind:  expr.KindImag,
		},
		{
			in:            "true && 20 > 9",
			expectedValue: true,
			expectedKind:  expr.KindBoolean,
		},
		{
			in:            "1 + 1",
			expectedValue: float64(2),
			expectedKind:  expr.KindFloat,
		},
		{
			in:            "1 + 2 * 10",
			expectedValue: float64(21),
			expectedKind:  expr.KindFloat,
		},
		{
			in:            "2.5 * 2.1",
			expectedValue: float64(5.25),
			expectedKind:  expr.KindFloat,
		},
		{
			in:            "2.50 > 2.4",
			expectedValue: true,
			expectedKind:  expr.KindBoolean,
		},
		{
			in:           "true & 2.4",
			expectedKind: expr.KindIllegal,
			expectedErr:  expr.ErrBitwiseOperation,
		},
		{
			in:            "4 == 0b0100",
			expectedValue: true,
			expectedKind:  expr.KindBoolean,
		},
		{
			in:           "true || !(!(10 * 100 %2))",
			expectedKind: expr.KindIllegal,
			expectedErr:  expr.ErrUnaryOperation,
		},
		{
			in:           "!(!(10 * 100 %2)) || true ",
			expectedKind: expr.KindIllegal,
			expectedErr:  expr.ErrUnaryOperation,
		},
		{
			in:            "expr == expr",
			expectedValue: true,
			expectedKind:  expr.KindBoolean,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			e, err := parser.ParseExpr(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			v := expr.NewVisitor(expr.WithNumericType(expr.NumericTypeAuto))
			ast.Walk(v, e)

			if err := v.Err(); !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, err)
			}
			if val := v.ValueAny(); val != tc.expectedValue {
				t.Fatalf("expected val: %v (%T), got: %v (%T)", tc.expectedValue, tc.expectedValue, val, val)
			}
			if kind := v.Kind(); kind != tc.expectedKind {
				t.Fatalf("expected kind: %v, got: %v", tc.expectedKind, kind)
			}
		})
	}

	tt2 := []struct {
		in       ast.Expr
		expected *expr.Visitor
	}{
		{
			in:       nil,
			expected: &expr.Visitor{},
		},
		{
			in:       &ast.BadExpr{},
			expected: &expr.Visitor{},
		},
	}

	for _, tc := range tt2 {
		tc := tc
		t.Run("", func(t *testing.T) {
			v := &expr.Visitor{}
			ast.Walk(v, tc.in)
			if v.Err() != tc.expected.Err() {
				t.Fatalf("expected err: %v, got: %v", tc.expected.Err(), v.Err())
			}
			if v.Kind() != tc.expected.Kind() {
				t.Fatalf("expected kind: %s, got: %s", tc.expected.Kind(), v.Kind())
			}
			if v.Value() != tc.expected.Value() {
				t.Fatalf("expected value: %v, got: %v", tc.expected.Value(), v.Value())
			}
		})
	}
}

func TestKindString(t *testing.T) {
	kinds := [...]string{
		expr.KindIllegal: "KindIllegal",
		expr.KindBoolean: "KindBoolean",
		expr.KindInt:     "KindInt",
		expr.KindFloat:   "KindFloat",
		expr.KindImag:    "KindImag",
		expr.KindString:  "KindString",
	}

	for kind, expected := range kinds {
		kind, expected := kind, expected
		t.Run(expected, func(t *testing.T) {
			kind := expr.Kind(kind)
			if kind.String() != expected {
				t.Fatalf("expected kind string: %s, got: %s", expected, kind.String())
			}
		})
	}

	// unsupported kinds
	tt := []struct {
		kind     expr.Kind
		expected string
	}{
		{kind: -100, expected: "kind(-100)"},
		{kind: 1000, expected: "kind(1000)"},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.expected, func(t *testing.T) {
			if tc.kind.String() != tc.expected {
				t.Fatalf("expected kind string: %s, got: %s", tc.expected, tc.kind.String())
			}
		})
	}
}
