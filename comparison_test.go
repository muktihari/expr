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

func TestComparison(t *testing.T) {
	tt := []struct {
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []value
		expectedErrs   []error
	}{
		// compareBoolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: boolValue(true)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: boolValue(false)},
			expectedValues: []value{boolValue(false), boolValue(true)},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: boolValue(true)},
			ops:            []token.Token{token.GTR},
			vy:             &Visitor{value: boolValue(false)},
			expectedValues: []value{{}},
			expectedErrs:   []error{ErrUnsupportedOperator},
		},
		// compareString
		{
			v:              &Visitor{},
			vx:             &Visitor{value: stringValue("\"abc\"")},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: stringValue("\"abc\"")},
			expectedValues: []value{boolValue(true), boolValue(false), boolValue(false), boolValue(true), boolValue(false), boolValue(true)},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compareImag
		{
			v:              &Visitor{},
			vx:             &Visitor{value: complex128Value(2 + 0i)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: complex128Value(2 + 0i)},
			expectedValues: []value{boolValue(true), boolValue(false)},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: complex128Value(2 + 0i)},
			ops:            []token.Token{token.GTR},
			vy:             &Visitor{value: complex128Value(2 + 0i)},
			expectedValues: []value{{}},
			expectedErrs:   []error{ErrUnsupportedOperator},
		},
		// compareFloat
		{
			v:              &Visitor{},
			vx:             &Visitor{value: float64Value(2.0)},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: int64Value(2)},
			expectedValues: []value{boolValue(true), boolValue(false), boolValue(false), boolValue(true), boolValue(false), boolValue(true)},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compareInt
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64Value(2)},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: int64Value(2)},
			expectedValues: []value{boolValue(true), boolValue(false), boolValue(false), boolValue(true), boolValue(false), boolValue(true)},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compare imag to float
		{
			v:              &Visitor{},
			vx:             &Visitor{value: complex128Value(2 + 0i)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: float64Value(2.0)},
			expectedValues: []value{boolValue(true), boolValue(false)},
			expectedErrs:   []error{nil, nil},
		},
		// compare imag to int
		{
			v:              &Visitor{},
			vx:             &Visitor{value: complex128Value(2 + 0i)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: int64Value(2)},
			expectedValues: []value{boolValue(true), boolValue(false)},
			expectedErrs:   []error{nil, nil},
		},
		// compare float to imag
		{
			v:              &Visitor{},
			vx:             &Visitor{value: float64Value(2.0)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: complex128Value(2 + 0i)},
			expectedValues: []value{boolValue(true), boolValue(false)},
			expectedErrs:   []error{nil, nil},
		},
		// compare int to complex
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64Value(2)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: complex128Value(2 + 0i)},
			expectedValues: []value{boolValue(true), boolValue(false)},
			expectedErrs:   []error{nil, nil},
		},
		// compare int to float64
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64Value(2)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: float64Value(2.0)},
			expectedValues: []value{boolValue(true), boolValue(false)},
			expectedErrs:   []error{nil, nil},
		},
		// compare int to boolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64Value(10)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: boolValue(true)},
			expectedValues: []value{{}, {}},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
		// compare boolean to int
		{
			v:              &Visitor{},
			vx:             &Visitor{value: boolValue(true)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: int64Value(10)},
			expectedValues: []value{{}, {}},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
		// compare boolean to string
		{
			v:              &Visitor{},
			vx:             &Visitor{value: boolValue(true)},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: stringValue("true")},
			expectedValues: []value{{}, {}},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
		// compare string to boolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: stringValue("true")},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: boolValue(true)},
			expectedValues: []value{{}, {}},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
	}

	for i, tc := range tt {
		tc := tc

		for j, op := range tc.ops {
			var (
				op            = op
				expectedErr   = tc.expectedErrs[j]
				expectedValue = tc.expectedValues[j]
				be            = &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: fmt.Sprintf("%v", tc.vx.value)},
					Op: op,
					Y:  &ast.BasicLit{Value: fmt.Sprintf("%v", tc.vy.value)},
				}
				name = fmt.Sprintf("%v %s %v", tc.vx.value.Any(), op, tc.vy.value.Any())
			)

			t.Run(fmt.Sprintf("[%d][%d] %s", i, j, name), func(t *testing.T) {
				comparison(tc.v, tc.vx, tc.vy, be)
				if !errors.Is(tc.v.err, expectedErr) {
					t.Fatalf("expected err: %v, got: %v", expectedErr, tc.v.err)
				}

				if tc.v.value.Any() != expectedValue.Any() {
					t.Fatalf("expected value: %v (%T), got: %v (%T)", expectedValue.Any(), tc.v.value.Any(),
						tc.v.value, tc.v.value)
				}
			})
		}
	}

}
