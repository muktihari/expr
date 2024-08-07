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
	"errors"
	"go/ast"
	"go/parser"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/muktihari/expr"
)

func TestVisit(t *testing.T) {
	tt := []struct {
		str        string
		transforms []Transform
		err        error
	}{
		{
			str: "((1 * 2) + 5) * (3 * 4)",
			transforms: []Transform{
				{Segmented: "1 * 2", EquivalentForm: "1 * 2", Evaluated: "2"},
				{Segmented: "(1 * 2) + 5", EquivalentForm: "2 + 5", Evaluated: "7"},
				{Segmented: "3 * 4", EquivalentForm: "3 * 4", Evaluated: "12"},
				{Segmented: "(2 + 5) * (3 * 4)", EquivalentForm: "7 * 12", Evaluated: "84"},
			},
		},
		{
			str: "(1 * 2) * (3 * 4)",
			transforms: []Transform{
				{Segmented: "1 * 2", EquivalentForm: "1 * 2", Evaluated: "2"},
				{Segmented: "3 * 4", EquivalentForm: "3 * 4", Evaluated: "12"},
				{Segmented: "(1 * 2) * (3 * 4)", EquivalentForm: "2 * 12", Evaluated: "24"},
			},
		},
		{
			str: "2 + 1 + 2 * 3 + 3",
			transforms: []Transform{
				{Segmented: "2 + 1", EquivalentForm: "2 + 1", Evaluated: "3"},
				{Segmented: "2 * 3", EquivalentForm: "2 * 3", Evaluated: "6"},
				{Segmented: "(2 + 1) + (2 * 3)", EquivalentForm: "3 + 6", Evaluated: "9"},
				{Segmented: "(2 + 1 + 2 * 3) + 3", EquivalentForm: "9 + 3", Evaluated: "12"},
			},
		},
		{
			str: "1 * 1 > 1 + 2",
			transforms: []Transform{
				{Segmented: "1 * 1", EquivalentForm: "1 * 1", Evaluated: "1"},
				{Segmented: "1 + 2", EquivalentForm: "1 + 2", Evaluated: "3"},
				{Segmented: "(1 * 1) > (1 + 2)", EquivalentForm: "1 > 3", Evaluated: "false"},
			},
		},
		{
			str: "(1+2)",
			transforms: []Transform{
				{Segmented: "1 + 2", EquivalentForm: "1 + 2", Evaluated: "3"},
			},
		},
		{
			str: "!true || ((5 > 3) && 1==1)",
			transforms: []Transform{
				{Segmented: "!true", EquivalentForm: "!true", Evaluated: "false"},
				{Segmented: "5 > 3", EquivalentForm: "5 > 3", Evaluated: "true"},
				{Segmented: "1 == 1", EquivalentForm: "1 == 1", Evaluated: "true"},
				{
					Segmented:      "(5 > 3) && (1==1)",
					EquivalentForm: "true && true",
					Evaluated:      "true",
				},
				{
					Segmented:      "false || (true && true)",
					EquivalentForm: "false || true",
					Evaluated:      "true",
				},
			},
		},
		{
			str: "true && 1 == 2 && ((!7))",
			err: expr.ErrUnaryOperation,
		},
		{
			str: "!(!9) && (!(7))",
			err: expr.ErrUnaryOperation,
		},
		{
			str: "1.2 & 1",
			err: expr.ErrBitwiseOperation,
		},
		{
			str: "4 << 10",
			transforms: []Transform{
				{Segmented: "4 << 10", EquivalentForm: "4 << 10", Evaluated: "4096"},
			},
		},
		{
			str: "4 >> 10",
			transforms: []Transform{
				{Segmented: "4 >> 10", EquivalentForm: "4 >> 10", Evaluated: "0"},
			},
		},
		{
			str: "4 & 10",
			transforms: []Transform{
				{Segmented: "4 & 10", EquivalentForm: "4 & 10", Evaluated: "0"},
			},
		},
		{
			str: "4 | 1000000000",
			transforms: []Transform{
				{Segmented: "4 | 1000000000", EquivalentForm: "4 | 1000000000", Evaluated: "1000000004"},
			},
		},
		{
			str: "4 ^ 10",
			transforms: []Transform{
				{Segmented: "4 ^ 10", EquivalentForm: "4 ^ 10", Evaluated: "14"},
			},
		},
		{
			str: "4 &^ 10",
			transforms: []Transform{
				{Segmented: "4 &^ 10", EquivalentForm: "4 &^ 10", Evaluated: "4"},
			},
		},
	}

	// test nil node
	t.Run("nil node", func(t *testing.T) { ast.Walk((&Visitor{}), nil) })

	for _, tc := range tt {
		tc := tc
		t.Run(tc.str, func(t *testing.T) {
			expr, err := parser.ParseExpr(tc.str)
			if err != nil {
				t.Fatal(err)
			}
			v := &Visitor{}
			ast.Walk(v, expr)
			if !errors.Is(v.Err(), tc.err) {
				t.Fatalf("expected err: %v, got: %v", tc.err, v.err)
			}
			if v.err != nil {
				return
			}

			transforms := v.Value()

			// Ignore Explanation
			for i := range transforms {
				transforms[i].Explanation = ""
			}

			if diff := cmp.Diff(transforms, tc.transforms); diff != "" {
				t.Fatal(diff)
			}
		})
	}

}
