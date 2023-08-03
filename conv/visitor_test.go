package conv

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestVisit(t *testing.T) {
	tt := []struct {
		name string
		in   ast.Expr
		val  string
		pos  int
	}{
		{
			name: "visit nil",
			in:   nil,
			val:  "",
			pos:  0,
		},
		{
			name: "visit nor supported expr",
			in:   &ast.BadExpr{},
			val:  "",
			pos:  0,
		},
		{
			name: "visit parent pos 10",
			in: &ast.ParenExpr{
				Lparen: 10,
				X: &ast.BasicLit{
					Value:    "1234",
					ValuePos: 11,
				},
				Rparen: 15,
			},
			val: "(1234)",
			pos: 10,
		},
		{
			name: "visit unary pos 15",
			in: &ast.UnaryExpr{
				OpPos: 15,
				Op:    token.NOT,
				X: &ast.Ident{
					Name:    "true",
					NamePos: 16,
				},
			},
			val: "!true",
			pos: 15,
		},
		{
			name: "visit binary pos 17",
			in: &ast.BinaryExpr{
				X: &ast.BasicLit{
					Kind:     token.INT,
					Value:    "1",
					ValuePos: 17,
				},
				Op:    token.ADD,
				OpPos: 19,
				Y: &ast.BasicLit{
					Kind:     token.INT,
					Value:    "2",
					ValuePos: 21,
				},
			},
			val: "1 + 2",
			pos: 17,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := &Visitor{}
			ast.Walk(v, tc.in)
			if v.value != v.Value() {
				t.Fatalf("expected Value(): %s, got: %s", v.value, v.Value())
			}
			if v.value != tc.val {
				t.Fatalf("expected value: %s, got: %s", tc.val, v.value)
			}
			if v.pos != tc.pos {
				t.Fatalf("expected pos: %d, got: %d", tc.pos, v.pos)
			}
		})

	}

}
