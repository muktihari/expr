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
	"go/parser"
	"testing"
)

func TestVisit(t *testing.T) {
	tt := []struct {
		in            string
		expectedValue value
		expectedKind  Kind
		expectedErr   error
	}{
		{
			in:           "<-a", // unary operator
			expectedKind: KindIllegal,
			expectedErr:  ErrUnsupportedOperator,
		},
		{
			in:            "'a' == 'a'",
			expectedValue: boolValue(true),
			expectedKind:  KindBoolean,
		},
		{
			in:            "-(20+10i)",
			expectedValue: complex128Value(-(20 + 10i)),
			expectedKind:  KindImag,
		},
		{
			in:            "true && 20 > 9",
			expectedValue: boolValue(true),
			expectedKind:  KindBoolean,
		},
		{
			in:            "1 + 1",
			expectedValue: float64Value(2),
			expectedKind:  KindFloat,
		},
		{
			in:            "1 + 2 * 10",
			expectedValue: float64Value(21),
			expectedKind:  KindFloat,
		},
		{
			in:            "2.5 * 2.1",
			expectedValue: float64Value(5.25),
			expectedKind:  KindFloat,
		},
		{
			in:            "2.50 > 2.4",
			expectedValue: boolValue(true),
			expectedKind:  KindBoolean,
		},
		{
			in:           "true & 2.4",
			expectedKind: KindIllegal,
			expectedErr:  ErrBitwiseOperation,
		},
		{
			in:            "4 == 0b0100",
			expectedValue: boolValue(true),
			expectedKind:  KindBoolean,
		},
		{
			in:           "true || !(!(10 * 100 %2))",
			expectedKind: KindIllegal,
			expectedErr:  ErrUnaryOperation,
		},
		{
			in:           "!(!(10 * 100 %2)) || true ",
			expectedKind: KindIllegal,
			expectedErr:  ErrUnaryOperation,
		},
		{
			in:            "expr == expr",
			expectedValue: boolValue(true),
			expectedKind:  KindBoolean,
		},
	}

	for i, tc := range tt {
		if i < 2 {
			continue
		}
		tc := tc
		t.Run(fmt.Sprintf("[%d] %s", i, tc.in), func(t *testing.T) {
			e, err := parser.ParseExpr(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			v := NewVisitor(WithNumericType(NumericTypeAuto))
			ast.Walk(v, e)

			if err := v.Err(); !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, err)
			}
			if val := v.ValueAny(); val != tc.expectedValue.Any() {
				t.Fatalf("expected val: %v (%T), got: %v (%T)", tc.expectedValue.Any(), tc.expectedValue.Any(), val, val)
			}
			if kind := v.Kind(); kind != tc.expectedKind {
				t.Fatalf("expected kind: %v, got: %v", tc.expectedKind, kind)
			}
		})
	}

	tt2 := []struct {
		in       ast.Expr
		expected *Visitor
	}{
		{
			in:       nil,
			expected: &Visitor{},
		},
		{
			in:       &ast.BadExpr{},
			expected: &Visitor{},
		},
	}

	for _, tc := range tt2 {
		tc := tc
		t.Run("", func(t *testing.T) {
			v := &Visitor{}
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
		KindIllegal: "KindIllegal",
		KindBoolean: "KindBoolean",
		KindInt:     "KindInt",
		KindFloat:   "KindFloat",
		KindImag:    "KindImag",
		KindString:  "KindString",
	}

	for kind, expected := range kinds {
		kind, expected := kind, expected
		t.Run(expected, func(t *testing.T) {
			kind := Kind(kind)
			if kind.String() != expected {
				t.Fatalf("expected kind string: %s, got: %s", expected, kind.String())
			}
		})
	}

	// unsupported kinds
	tt := []struct {
		kind     Kind
		expected string
	}{
		{kind: 255, expected: "kind(255)"},
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
