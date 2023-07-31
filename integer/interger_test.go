package integer

import (
	"errors"
	"go/ast"
	"go/token"
	"testing"
)

func TestVisitUnary(t *testing.T) {
	tt := []struct {
		name           string
		unaryExpr      *ast.UnaryExpr
		expectedResult int
		expectedErr    error
	}{
		{
			name: "+2",
			unaryExpr: &ast.UnaryExpr{
				Op: token.ADD,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2",
				},
			},
			expectedResult: 2,
		},
		{
			name: "-2",
			unaryExpr: &ast.UnaryExpr{
				Op: token.SUB,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2",
				},
			},
			expectedResult: -2,
		},
		{
			name: "-(-2)",
			unaryExpr: &ast.UnaryExpr{
				Op: token.SUB,
				X: &ast.ParenExpr{
					X: &ast.UnaryExpr{
						Op: token.SUB,
						X: &ast.BasicLit{
							Kind:  token.INT,
							Value: "2",
						},
					},
				},
			},
			expectedResult: 2,
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
				t.Fatalf("expected result: %d, got: %d", tc.expectedResult, v.res)
			}
		})
	}
}

func TestArithmetic(t *testing.T) {
	tt := []struct {
		name           string
		x              ast.Expr
		y              ast.Expr
		op             token.Token
		expectedResult int
		expectedErr    error
	}{
		{
			name:           "7 + 3 = 10",
			x:              &ast.BasicLit{Value: "7", Kind: token.FLOAT},
			y:              &ast.BasicLit{Value: "3", Kind: token.INT},
			op:             token.ADD,
			expectedResult: 10,
		},
		{
			name:           "7 - 3 = 4",
			x:              &ast.BasicLit{Value: "7", Kind: token.INT},
			y:              &ast.BasicLit{Value: "3", Kind: token.INT},
			op:             token.SUB,
			expectedResult: 4,
		},
		{
			name:           "7 * 3 = 21",
			x:              &ast.BasicLit{Value: "7", Kind: token.INT},
			y:              &ast.BasicLit{Value: "3", Kind: token.INT},
			op:             token.MUL,
			expectedResult: 21,
		},
		{
			name:           "7 / 3 = 2",
			x:              &ast.BasicLit{Value: "7", Kind: token.INT},
			y:              &ast.BasicLit{Value: "3", Kind: token.INT},
			op:             token.QUO,
			expectedResult: 2,
		},
		{
			name:           "7 % 3 = 1",
			x:              &ast.BasicLit{Value: "7", Kind: token.INT},
			y:              &ast.BasicLit{Value: "3", Kind: token.INT},
			op:             token.REM,
			expectedResult: 1,
		},
		{
			name:        "7 / 0 = ErrIntegerDividedByZero",
			x:           &ast.BasicLit{Value: "7", Kind: token.INT},
			y:           &ast.BasicLit{Value: "0", Kind: token.INT},
			op:          token.QUO,
			expectedErr: ErrIntegerDividedByZero,
		},
		{
			name:        "!7 / 3 = 2",
			x:           &ast.UnaryExpr{X: &ast.BasicLit{Value: "7"}, Op: token.MUL},
			y:           &ast.BasicLit{Value: "3", Kind: token.INT},
			op:          token.QUO,
			expectedErr: ErrUnsupportedOperator,
		},
		{
			name:        "7 / !3 = 2",
			x:           &ast.BasicLit{Value: "7", Kind: token.INT},
			y:           &ast.UnaryExpr{X: &ast.BasicLit{Value: "3"}, Op: token.MUL},
			op:          token.QUO,
			expectedErr: ErrUnsupportedOperator,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			binaryExpr := &ast.BinaryExpr{X: tc.x, Y: tc.y, Op: tc.op}
			v := &visitor{}
			v.arithmetic(binaryExpr)
			if !errors.Is(v.err, tc.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, v.err)
			}

			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %d, got: %d", tc.expectedResult, v.res)
			}
		})
	}
}

func TestBitwise(t *testing.T) {
	tt := []struct {
		name           string
		x              ast.Expr
		y              ast.Expr
		op             token.Token
		expectedResult int
		expectedErr    error
	}{
		{
			name:           "1101 & 1011 = 1001 ≈ 9",
			x:              &ast.BasicLit{Value: "1101", Kind: token.INT},
			y:              &ast.BasicLit{Value: "1011", Kind: token.INT},
			op:             token.AND,
			expectedResult: 9,
		},
		{
			name:           "1101 | 1011 = 1111 ≈ 15",
			x:              &ast.BasicLit{Value: "1101", Kind: token.INT},
			y:              &ast.BasicLit{Value: "1011", Kind: token.INT},
			op:             token.OR,
			expectedResult: 15,
		},
		{
			name:           "1101 ^ 1011 = 0110 ≈ 6",
			x:              &ast.BasicLit{Value: "1101", Kind: token.INT},
			y:              &ast.BasicLit{Value: "1011", Kind: token.INT},
			op:             token.XOR,
			expectedResult: 6,
		},
		{
			name:           "1101 &^ 1011 can be written as 1101 & 0100 = 0100 ≈ 4",
			x:              &ast.BasicLit{Value: "1101", Kind: token.INT},
			y:              &ast.BasicLit{Value: "1011", Kind: token.INT},
			op:             token.AND_NOT,
			expectedResult: 4,
		},
		{
			name:           "0100 << 0110 = 1000000000000 ≈ 4 << 10 = 4096",
			x:              &ast.BasicLit{Value: "0100", Kind: token.INT},
			y:              &ast.BasicLit{Value: "1010", Kind: token.INT},
			op:             token.SHL,
			expectedResult: 4096,
		},
		{
			name:           "1111 >> 0010 = 0011 ≈ 15 >> 2 = 3",
			x:              &ast.BasicLit{Value: "1111", Kind: token.INT},
			y:              &ast.BasicLit{Value: "0010", Kind: token.INT},
			op:             token.SHR,
			expectedResult: 3,
		},
		{
			name: "!1111 & 0010 = ErrUnsupportedOperator on x",
			x: &ast.UnaryExpr{
				X:  &ast.BasicLit{Value: "1111", Kind: token.INT},
				Op: token.NOT,
			},
			y:           &ast.BasicLit{Value: "0010", Kind: token.INT},
			op:          token.AND,
			expectedErr: ErrUnsupportedOperator,
		},
		{
			name: "1111 & !0010 = ErrUnsupportedOperator on x",
			x:    &ast.BasicLit{Value: "1111", Kind: token.INT},
			y: &ast.UnaryExpr{
				X:  &ast.BasicLit{Value: "0010", Kind: token.INT},
				Op: token.NOT,
			},
			op:          token.AND,
			expectedErr: ErrUnsupportedOperator,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			binaryExpr := &ast.BinaryExpr{X: tc.x, Y: tc.y, Op: tc.op}
			v := &visitor{}
			v.bitwise(binaryExpr)
			if !errors.Is(v.err, tc.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, v.err)
			}
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %d, got: %d", tc.expectedResult, v.res)
			}
		})
	}
}

func TestVisitBinary(t *testing.T) {
	tt := []struct {
		name           string
		binaryExpr     *ast.BinaryExpr
		expectedResult int
		expectedErr    error
	}{
		{
			name: "1 + 4 = 5",
			binaryExpr: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Y:  &ast.BasicLit{Value: "4", Kind: token.INT},
				Op: token.ADD,
			},
			expectedResult: 5,
		},
		{
			name: "0001 & 1111 = 1",
			binaryExpr: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "0001", Kind: token.INT},
				Y:  &ast.BasicLit{Value: "1111", Kind: token.INT},
				Op: token.AND,
			},
			expectedResult: 1,
		},
		{
			name: "0001 = 1111 = ErrUnsupportedOperator",
			binaryExpr: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "0001", Kind: token.INT},
				Y:  &ast.BasicLit{Value: "1111", Kind: token.INT},
				Op: token.ASSIGN,
			},
			expectedErr: ErrUnsupportedOperator,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := &visitor{}
			_ = v.visitBinary(tc.binaryExpr)

			if !errors.Is(v.err, tc.expectedErr) {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, v.err)
			}
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %d, got: %d", tc.expectedResult, v.res)
			}
		})
	}

}

func TestVisit(t *testing.T) {
	tt := []struct {
		name           string
		node           ast.Node
		expectedResult int
	}{
		{
			name: "*ast.ParenExpr, (1) = 1",
			node: &ast.ParenExpr{
				X: &ast.BasicLit{Value: "1", Kind: token.INT},
			},
			expectedResult: 1,
		},
		{
			name: "*ast.BinaryExpr, 1 + 2 = 3",
			node: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
				Op: token.ADD,
			},
			expectedResult: 3,
		},
		{
			name:           "*ast.BasicLit, 1",
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

			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %d, got: %d", tc.expectedResult, v.res)
			}
		})
	}
}

func TestResult(t *testing.T) {
	tt := []struct {
		name           string
		v              *visitor
		expectedResult int
		expectedErr    error
	}{
		{name: "9.9", v: &visitor{res: 9}, expectedResult: 9},
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
