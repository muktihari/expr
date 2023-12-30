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

package explain_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/muktihari/expr"
	"github.com/muktihari/expr/exp/explain"
)

func TestExplain(t *testing.T) {
	tt := []struct {
		str      string
		explains []explain.Step
		err      error
	}{
		{
			str: "1 + 2",
			explains: []explain.Step{
				{[]string{"1 + 2"}, "3"},
			},
		},
		{
			str: "1 + 2 + 3",
			explains: []explain.Step{
				{[]string{"1 + 2"}, "3"},
				{[]string{"(1 + 2) + 3", "3 + 3"}, "6"},
			},
		},
		{
			str: "!true || ((5 > 3) && 1 == 1)",
			explains: []explain.Step{
				{[]string{"!true"}, "false"},
				{[]string{"5 > 3"}, "true"},
				{[]string{"1 == 1"}, "true"},
				{[]string{"(5 > 3) && (1 == 1)", "true && true"}, "true"},
				{[]string{"false || (true && true)", "false || true"}, "true"},
			},
		},
		{
			str: "!(true) && !7",
			err: expr.ErrUnaryOperation,
		},
	}

	t.Run("parser error", func(t *testing.T) {
		_, err := explain.Explain("1 +")
		if err == nil {
			t.Fatalf("expected error, got: %v", err)
		}
	})

	for _, tc := range tt {
		tc := tc
		t.Run(tc.str, func(t *testing.T) {
			explains, err := explain.Explain(tc.str)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected err: %v, got: %v", tc.err, err)
			}

			if diff := cmp.Diff(explains, tc.explains); diff != "" {
				t.Fatal(diff)
			}
		})
	}

}
