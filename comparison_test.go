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
		expectedValues []interface{}
		expectedErrs   []error
	}{
		// compareBoolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: true, kind: KindBoolean},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: false, kind: KindBoolean},
			expectedValues: []interface{}{false, true},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: true, kind: KindBoolean},
			ops:            []token.Token{token.GTR},
			vy:             &Visitor{value: false, kind: KindBoolean},
			expectedValues: []interface{}{nil},
			expectedErrs:   []error{ErrUnsupportedOperator},
		},
		// compareString
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "\"abc\"", kind: KindString},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: "\"abc\"", kind: KindString},
			expectedValues: []interface{}{true, false, false, true, false, true},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compareImag
		{
			v:              &Visitor{},
			vx:             &Visitor{value: (2 + 0i), kind: KindImag},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: (2 + 0i), kind: KindImag},
			expectedValues: []interface{}{true, false},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: (2 + 0i), kind: KindImag},
			ops:            []token.Token{token.GTR},
			vy:             &Visitor{value: (2 + 0i), kind: KindImag},
			expectedValues: []interface{}{nil},
			expectedErrs:   []error{ErrUnsupportedOperator},
		},
		// compareFloat
		{
			v:              &Visitor{},
			vx:             &Visitor{value: float64(2.0), kind: KindFloat},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: int64(2), kind: KindInt},
			expectedValues: []interface{}{true, false, false, true, false, true},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compareInt
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64(2), kind: KindInt},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: int64(2), kind: KindInt},
			expectedValues: []interface{}{true, false, false, true, false, true},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compare imag to float
		{
			v:              &Visitor{},
			vx:             &Visitor{value: (2 + 0i), kind: KindImag},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: 2.0, kind: KindFloat},
			expectedValues: []interface{}{true, false},
			expectedErrs:   []error{nil, nil},
		},
		// compare imag to int
		{
			v:              &Visitor{},
			vx:             &Visitor{value: (2 + 0i), kind: KindImag},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: int64(2), kind: KindInt},
			expectedValues: []interface{}{true, false},
			expectedErrs:   []error{nil, nil},
		},
		// compare float to imag
		{
			v:              &Visitor{},
			vx:             &Visitor{value: 2.0, kind: KindInt},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: (2 + 0i), kind: KindImag},
			expectedValues: []interface{}{true, false},
			expectedErrs:   []error{nil, nil},
		},
		// compare int to complex
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64(2), kind: KindInt},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: 2 + 0i, kind: KindImag},
			expectedValues: []interface{}{true, false},
			expectedErrs:   []error{nil, nil},
		},
		// compare int to float64
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64(2), kind: KindInt},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: 2.0, kind: KindFloat},
			expectedValues: []interface{}{true, false},
			expectedErrs:   []error{nil, nil},
		},
		// compare int to boolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: int64(10), kind: KindInt},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: true, kind: KindBoolean},
			expectedValues: []interface{}{nil, nil},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
		// compare boolean to int
		{
			v:              &Visitor{},
			vx:             &Visitor{value: true, kind: KindBoolean},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: int64(10), kind: KindInt},
			expectedValues: []interface{}{nil, nil},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
		// compare boolean to string
		{
			v:              &Visitor{},
			vx:             &Visitor{value: true, kind: KindBoolean},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: "true", kind: KindString},
			expectedValues: []interface{}{nil, nil},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
		// compare string to boolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "true", kind: KindString},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: true, kind: KindBoolean},
			expectedValues: []interface{}{nil, nil},
			expectedErrs:   []error{ErrComparisonOperation, ErrComparisonOperation},
		},
	}

	for _, tc := range tt {
		tc := tc

		for i, op := range tc.ops {
			var (
				op            = op
				expectedErr   = tc.expectedErrs[i]
				expectedValue = tc.expectedValues[i]
				be            = &ast.BinaryExpr{
					X:  &ast.BasicLit{Value: fmt.Sprintf("%v", tc.vx.value)},
					Op: op,
					Y:  &ast.BasicLit{Value: fmt.Sprintf("%v", tc.vy.value)},
				}
				name = fmt.Sprintf("%v %s %v", tc.vx.value, op, tc.vy.value)
			)

			t.Run(name, func(t *testing.T) {
				comparison(tc.v, tc.vx, tc.vy, be)
				if !errors.Is(tc.v.err, expectedErr) {
					t.Fatalf("expected err: %v, got: %v", expectedErr, tc.v.err)
				}

				if tc.v.value != expectedValue {
					t.Fatalf("expected value: %v (%T), got: %v (%T)", expectedValue, expectedValue,
						tc.v.value, tc.v.value)
				}
			})
		}
	}

}
