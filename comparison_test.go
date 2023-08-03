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
		expectedValues []string
		expectedErrs   []error
	}{
		// compareBoolean
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "true", kind: KindBoolean},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: "false", kind: KindBoolean},
			expectedValues: []string{"false", "true"},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "true", kind: KindBoolean},
			ops:            []token.Token{token.GTR},
			vy:             &Visitor{value: "false", kind: KindBoolean},
			expectedValues: []string{""},
			expectedErrs:   []error{ErrUnsupportedOperator},
		},
		// compareString
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "\"abc\"", kind: KindString},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: "\"abc\"", kind: KindString},
			expectedValues: []string{"true", "false", "false", "true", "false", "true"},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compareImag
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "(2+0i)", kind: KindImag},
			ops:            []token.Token{token.EQL, token.NEQ},
			vy:             &Visitor{value: "(2+0i)", kind: KindImag},
			expectedValues: []string{"true", "false"},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "(2+0i)", kind: KindImag},
			ops:            []token.Token{token.GTR},
			vy:             &Visitor{value: "(2+0i)", kind: KindImag},
			expectedValues: []string{""},
			expectedErrs:   []error{ErrUnsupportedOperator},
		},
		// compareFloat
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "2.0", kind: KindFloat},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: "2", kind: KindInt},
			expectedValues: []string{"true", "false", "false", "true", "false", "true"},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
		},
		// compareInt
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "2", kind: KindInt},
			ops:            []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			vy:             &Visitor{value: "2", kind: KindInt},
			expectedValues: []string{"true", "false", "false", "true", "false", "true"},
			expectedErrs:   []error{nil, nil, nil, nil, nil, nil},
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
					X:  &ast.BasicLit{Value: tc.vx.value},
					Op: op,
					Y:  &ast.BasicLit{Value: tc.vy.value},
				}
				name = fmt.Sprintf("%s %s %s", tc.vx.value, op, tc.vy.value)
			)

			t.Run(name, func(t *testing.T) {
				comparison(tc.v, tc.vx, tc.vy, be)
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
