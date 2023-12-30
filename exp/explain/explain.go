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
	"go/ast"
	"go/parser"
)

// Step contains smaller operations of given s with its result.
// One step can have multiple equivalent forms, starts with original form until the final form.
type Step struct {
	EquivalentForms []string
	Result          string
}

// Explains explains step-by-step process of evaluating s.
func Explain(s string) ([]Step, error) {
	e, err := parser.ParseExpr(s)
	if err != nil {
		return nil, err
	}

	v := &Visitor{}
	ast.Walk(v, e)
	if err := v.err; err != nil {
		return nil, err
	}

	// sanitize results
	explains := make([]Step, len(v.transforms))
	for i, transform := range v.transforms {
		explains[i] = Step{
			EquivalentForms: []string{transform.Segmented},
			Result:          transform.Evaluated,
		}
		if transform.Segmented != transform.EquivalentForm {
			explains[i].EquivalentForms = append(explains[i].EquivalentForms, transform.EquivalentForm)
		}
	}

	return explains, nil
}
