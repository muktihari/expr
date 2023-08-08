// explain is a standalone package aimed to explain step by step operation in expr.
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
