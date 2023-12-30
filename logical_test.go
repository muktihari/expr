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

func TestLogical(t *testing.T) {
	tt := []struct {
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []interface{}
		expectedErrs   []error
	}{
		{
			v:              &Visitor{},
			vx:             &Visitor{value: true, kind: KindBoolean},
			ops:            []token.Token{token.LAND, token.NEQ},
			vy:             &Visitor{value: false, kind: KindBoolean},
			expectedValues: []interface{}{false, true},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "1", kind: KindInt},
			ops:            []token.Token{token.LAND},
			vy:             &Visitor{value: "false", kind: KindBoolean},
			expectedValues: []interface{}{nil},
			expectedErrs:   []error{ErrLogicalOperation},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "false", kind: KindBoolean},
			ops:            []token.Token{token.LAND},
			vy:             &Visitor{value: "1", kind: KindInt},
			expectedValues: []interface{}{nil},
			expectedErrs:   []error{ErrLogicalOperation},
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
				logical(tc.v, tc.vx, tc.vy, be)
				if !errors.Is(tc.v.err, expectedErr) {
					t.Fatalf("expected err: %v, got: %v", expectedErr, tc.v.err)
				}
				if tc.v.value != expectedValue {
					t.Fatalf("expected value: %v, got: %v", expectedValue, tc.v.value)
				}
			})
		}
	}
}
