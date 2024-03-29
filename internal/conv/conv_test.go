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

package conv_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/muktihari/expr/internal/conv"
)

func TestFormatExpr(t *testing.T) {
	tt := []struct {
		in  ast.Expr
		val string
	}{
		{
			in: &ast.BinaryExpr{
				X: &ast.ParenExpr{
					Lparen: 10,
					X: &ast.BinaryExpr{
						X: &ast.BasicLit{
							Kind:     token.INT,
							Value:    "1234",
							ValuePos: 11,
						},
						Op:    token.MUL,
						OpPos: 16,
						Y: &ast.BasicLit{
							Kind:     token.INT,
							Value:    "1",
							ValuePos: 18,
						},
					},
					Rparen: 19,
				},
				Op:    token.ADD,
				OpPos: 21,
				Y: &ast.BasicLit{
					Kind:     token.INT,
					Value:    "20",
					ValuePos: 23,
				},
			},
			val: "(1234 * 1) + 20",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.val, func(t *testing.T) {
			s := conv.FormatExpr(tc.in)
			if s != tc.val {
				t.Fatalf("expected value: %s, got: %s", tc.val, s)
			}
		})

	}
}
