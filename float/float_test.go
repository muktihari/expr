package float

import (
	"errors"
	"go/ast"
	"go/token"
	"math"
	"testing"
)

func TestVisitUnary(t *testing.T) {
	tt := []struct {
		name           string
		unaryExpr      *ast.UnaryExpr
		expectedResult float64
		expectedErr    error
	}{
		{
			name: "+2.7",
			unaryExpr: &ast.UnaryExpr{
				Op: token.ADD,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2.7",
				},
			},
			expectedResult: 2.7,
		},
		{
			name: "-2.7",
			unaryExpr: &ast.UnaryExpr{
				Op: token.SUB,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2.7",
				},
			},
			expectedResult: -2.7,
		},
		{
			name: "-(-2.7)",
			unaryExpr: &ast.UnaryExpr{
				Op: token.SUB,
				X: &ast.ParenExpr{
					X: &ast.UnaryExpr{
						Op: token.SUB,
						X: &ast.BasicLit{
							Kind:  token.INT,
							Value: "2.7",
						},
					},
				},
			},
			expectedResult: 2.7,
		},
		{
			name: "!2",
			unaryExpr: &ast.UnaryExpr{
				Op: token.NOT,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2",
				},
			},
			expectedErr: ErrUnsupportedOperator,
		},
		{
			name: "-(!2)",
			unaryExpr: &ast.UnaryExpr{
				Op: token.SUB,
				X: &ast.ParenExpr{
					X: &ast.UnaryExpr{
						Op: token.NOT,
						X: &ast.BasicLit{
							Kind:  token.INT,
							Value: "2",
						},
					},
				},
			},
			expectedErr: ErrUnsupportedOperator,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := visitor{}
			_ = v.visitUnary(tc.unaryExpr)
			if !errors.Is(v.err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, v.err)
			}
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %f, got: %f", tc.expectedResult, v.res)
			}
		})
	}
}

func TestVisitBinary(t *testing.T) {
	tt := []struct {
		name           string
		x              ast.Expr
		y              ast.Expr
		op             token.Token
		expectedResult float64
		expectedErr    error
	}{
		{
			name:           "7.1 + 3 = 10.1",
			x:              &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
			y:              &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			op:             token.ADD,
			expectedResult: 10.1,
		},
		{
			name:           "7.1 - 3 = 4.1",
			x:              &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
			y:              &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			op:             token.SUB,
			expectedResult: 4.1,
		},
		{
			name:           "7.1 * 3 = 21.3",
			x:              &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
			y:              &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			op:             token.MUL,
			expectedResult: 21.3,
		},
		{
			name:           "7.1 / 3 = 2.37",
			x:              &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
			y:              &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			op:             token.QUO,
			expectedResult: 2.37,
		},
		{
			name:        "40.0 % 2 = ErrUnsupportedOperator",
			x:           &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
			y:           &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			op:          token.REM,
			expectedErr: ErrUnsupportedOperator,
		},
		{
			name: "10 % 2 * 2 = ErrUnsupportedOperator",
			x: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
				Y:  &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
				Op: token.REM,
			},
			y:           &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			op:          token.MUL,
			expectedErr: ErrUnsupportedOperator,
		},
		{
			name: "10 * 2 % 2 = ErrUnsupportedOperator",
			x:    &ast.BasicLit{Value: "3", Kind: token.FLOAT},
			y: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
				Y:  &ast.BasicLit{Value: "7.1", Kind: token.FLOAT},
				Op: token.REM,
			},
			op:          token.MUL,
			expectedErr: ErrUnsupportedOperator,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			binaryExpr := &ast.BinaryExpr{X: tc.x, Y: tc.y, Op: tc.op}
			v := &visitor{}
			_ = v.visitBinary(binaryExpr)
			if !errors.Is(v.err, tc.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, v.err)
			}

			// avoid floating point precision problem, only compare up to 2 decimal should be ok
			if math.Round(v.res*100)/100 != math.Round(tc.expectedResult*100)/100 {
				t.Fatalf("expected result: %f, got: %f", tc.expectedResult, v.res)
			}
		})
	}
}

func TestVisit(t *testing.T) {
	tt := []struct {
		name           string
		node           ast.Node
		expectedResult float64
	}{
		{
			name: "*ast.ParenExpr (1.1) = 1.1",
			node: &ast.ParenExpr{
				X: &ast.BasicLit{Value: "1.1", Kind: token.FLOAT},
			},
			expectedResult: 1.1,
		},
		{
			name: "*ast.BinaryExpr 1.6 + 2 = 3.6",
			node: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1.6", Kind: token.FLOAT},
				Y:  &ast.BasicLit{Value: "2", Kind: token.FLOAT},
				Op: token.ADD,
			},
			expectedResult: 3.6,
		},
		{
			name:           "*ast.BasicLit 1",
			node:           &ast.BasicLit{Value: "1", Kind: token.INT},
			expectedResult: 1,
		},
		{
			name:           "not supported, e.g. *ast.ArrayType",
			node:           &ast.ArrayType{},
			expectedResult: 0,
		},
		{
			name:           "nil node",
			node:           nil,
			expectedResult: 0,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := NewVisitor().(*visitor)
			_ = v.Visit(tc.node)

			// avoid floating point precision problem, only compare up to 2 decimal should be sufficient
			if math.Round(v.res*100)/100 != math.Round(tc.expectedResult*100)/100 {
				t.Fatalf("expected result: %f, got: %f", tc.expectedResult, v.res)
			}
		})
	}
}

func TestResult(t *testing.T) {
	tt := []struct {
		name           string
		v              *visitor
		expectedResult float64
		expectedErr    error
	}{
		{name: "9.9", v: &visitor{res: 9.9}, expectedResult: 9.9},
		{name: "ErrUnsupportedOperator", v: &visitor{err: ErrUnsupportedOperator}, expectedErr: ErrUnsupportedOperator},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.v.Result()
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, err)
			}
			if res != tc.expectedResult {
				t.Fatalf("expected result: %v, got: %v", tc.expectedErr, err)
			}
		})
	}
}
