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
		expectedValues []string
		expectedErrs   []error
	}{
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "true", kind: KindBoolean},
			ops:            []token.Token{token.LAND, token.NEQ},
			vy:             &Visitor{value: "false", kind: KindBoolean},
			expectedValues: []string{"false", "true"},
			expectedErrs:   []error{nil, nil},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "1", kind: KindInt},
			ops:            []token.Token{token.LAND},
			vy:             &Visitor{value: "false", kind: KindBoolean},
			expectedValues: []string{""},
			expectedErrs:   []error{ErrLogicalOperation},
		},
		{
			v:              &Visitor{},
			vx:             &Visitor{value: "false", kind: KindBoolean},
			ops:            []token.Token{token.LAND},
			vy:             &Visitor{value: "1", kind: KindInt},
			expectedValues: []string{""},
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
					X:  &ast.BasicLit{Value: tc.vx.value},
					Op: op,
					Y:  &ast.BasicLit{Value: tc.vy.value},
				}
				name = fmt.Sprintf("%s %s %s", tc.vx.value, op, tc.vy.value)
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
