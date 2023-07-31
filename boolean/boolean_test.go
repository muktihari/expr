package boolean

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
		expectedResult string
		expectedErr    error
	}{
		{
			name: "!true",
			unaryExpr: &ast.UnaryExpr{
				Op: token.NOT,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "true",
				},
			},
			expectedResult: "false",
		},
		{
			name: "+2",
			unaryExpr: &ast.UnaryExpr{
				Op: token.ADD,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2",
				},
			},
			expectedResult: "2",
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
			expectedResult: "-2",
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
			expectedResult: "2",
		},
		{
			name: "*2",
			unaryExpr: &ast.UnaryExpr{
				Op: token.MUL,
				X: &ast.BasicLit{
					Kind:  token.INT,
					Value: "2",
				},
			},
			expectedErr: ErrUnsupportedOperator,
		},
		{
			name: "-(*2)",
			unaryExpr: &ast.UnaryExpr{
				Op: token.SUB,
				X: &ast.ParenExpr{
					X: &ast.UnaryExpr{
						Op: token.MUL,
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
				t.Fatalf("expected result: %s, got: %s", tc.expectedResult, v.res)
			}
		})
	}
}

func TestArithmetic(t *testing.T) {
	tt := []struct {
		name           string
		x              *visitor
		y              *visitor
		op             token.Token
		expectedResult string
		expectedErr    error
	}{
		{
			name:           "float, 8.5 + 1.5 = 10",
			x:              &visitor{res: "8.5", kind: token.FLOAT},
			y:              &visitor{res: "1.5", kind: token.FLOAT},
			op:             token.ADD,
			expectedResult: "10.000000",
		},
		{
			name:           "float, 8.5 - 1.5 = 7",
			x:              &visitor{res: "8.5", kind: token.FLOAT},
			y:              &visitor{res: "1.5", kind: token.FLOAT},
			op:             token.SUB,
			expectedResult: "7.000000",
		},
		{
			name:           "float, 8.5 * 1.5 = 12.75",
			x:              &visitor{res: "8.5", kind: token.FLOAT},
			y:              &visitor{res: "1.5", kind: token.FLOAT},
			op:             token.MUL,
			expectedResult: "12.750000",
		},
		{
			name:           "float, 8.5 / 1.5 = 5.666667",
			x:              &visitor{res: "8.5", kind: token.FLOAT},
			y:              &visitor{res: "1.5", kind: token.FLOAT},
			op:             token.QUO,
			expectedResult: "5.666667",
		},
		{
			name:        "float, 8.5 % 1.5 = ErrInvalidOperationOnFloat",
			x:           &visitor{res: "8.5", kind: token.FLOAT},
			y:           &visitor{res: "1.5", kind: token.FLOAT},
			op:          token.REM,
			expectedErr: ErrInvalidOperationOnFloat,
		},
		{
			name:           "int, 8 + 1 = 9",
			x:              &visitor{res: "8", kind: token.INT},
			y:              &visitor{res: "1", kind: token.INT},
			op:             token.ADD,
			expectedResult: "9",
		},
		{
			name:           "int, 8 - 1 = 7",
			x:              &visitor{res: "8", kind: token.INT},
			y:              &visitor{res: "1", kind: token.INT},
			op:             token.SUB,
			expectedResult: "7",
		},
		{
			name:           "int, 8 * 2 = 16",
			x:              &visitor{res: "8", kind: token.INT},
			y:              &visitor{res: "2", kind: token.INT},
			op:             token.MUL,
			expectedResult: "16",
		},
		{
			name:           "int, 8 / 2 = 4",
			x:              &visitor{res: "8", kind: token.INT},
			y:              &visitor{res: "2", kind: token.INT},
			op:             token.QUO,
			expectedResult: "4",
		},
		{
			name:        "int, 8 / 0 = ErrIntegerDividedByZero",
			x:           &visitor{res: "8", kind: token.INT},
			y:           &visitor{res: "0", kind: token.INT},
			op:          token.QUO,
			expectedErr: ErrIntegerDividedByZero,
		},
		{
			name:           "int, 8 % 3 = 2",
			x:              &visitor{res: "8", kind: token.INT},
			y:              &visitor{res: "3", kind: token.INT},
			op:             token.REM,
			expectedResult: "2",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := &visitor{}
			v.arithmetic(tc.x, tc.y, tc.op)
			if !errors.Is(v.err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, v.err)
			}
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %s, got: %s", tc.expectedResult, v.res)
			}
		})
	}
}

func TestComparison(t *testing.T) {
	tt := []struct {
		name           string
		x              *visitor
		y              *visitor
		op             token.Token
		expectedResult string
	}{
		{
			name:           "aaaa == aaaa = true",
			x:              &visitor{res: "a", kind: token.STRING},
			y:              &visitor{res: "a", kind: token.STRING},
			op:             token.EQL,
			expectedResult: "true",
		},
		{
			name:           "aaaa != aaaa = false",
			x:              &visitor{res: "a", kind: token.STRING},
			y:              &visitor{res: "a", kind: token.STRING},
			op:             token.NEQ,
			expectedResult: "false",
		},
		{
			name:           "b > a = true",
			x:              &visitor{res: "b", kind: token.STRING},
			y:              &visitor{res: "a", kind: token.STRING},
			op:             token.GTR,
			expectedResult: "true",
		},
		{
			name:           "b >= a = true",
			x:              &visitor{res: "b", kind: token.STRING},
			y:              &visitor{res: "a", kind: token.STRING},
			op:             token.GEQ,
			expectedResult: "true",
		},
		{
			name:           "a < b = true",
			x:              &visitor{res: "a", kind: token.STRING},
			y:              &visitor{res: "b", kind: token.STRING},
			op:             token.LSS,
			expectedResult: "true",
		},
		{
			name:           "a <= b = true",
			x:              &visitor{res: "a", kind: token.STRING},
			y:              &visitor{res: "b", kind: token.STRING},
			op:             token.LEQ,
			expectedResult: "true",
		},
		{
			name:           "10.567 > 10.234 = true",
			x:              &visitor{res: "10.567", kind: token.FLOAT},
			y:              &visitor{res: "10.234", kind: token.FLOAT},
			op:             token.GTR,
			expectedResult: "true",
		},
		{
			name:           "10.567 >= 10.234 = true",
			x:              &visitor{res: "10.567", kind: token.FLOAT},
			y:              &visitor{res: "10.234", kind: token.FLOAT},
			op:             token.GEQ,
			expectedResult: "true",
		},
		{
			name:           "10 >= 10 = true",
			x:              &visitor{res: "10", kind: token.FLOAT},
			y:              &visitor{res: "10", kind: token.FLOAT},
			op:             token.GEQ,
			expectedResult: "true",
		},
		{
			name:           "10.234 < 10.567 = true",
			x:              &visitor{res: "10.234", kind: token.FLOAT},
			y:              &visitor{res: "10.567", kind: token.FLOAT},
			op:             token.LSS,
			expectedResult: "true",
		},
		{
			name:           "10.234 <= 10.567 = true",
			x:              &visitor{res: "10.234", kind: token.FLOAT},
			y:              &visitor{res: "10.567", kind: token.FLOAT},
			op:             token.LEQ,
			expectedResult: "true",
		},
		{
			name:           "10 <= 10 = true",
			x:              &visitor{res: "10", kind: token.FLOAT},
			y:              &visitor{res: "10", kind: token.FLOAT},
			op:             token.LEQ,
			expectedResult: "true",
		},
		{
			name:           "10 > 9 = true",
			x:              &visitor{res: "10", kind: token.INT},
			y:              &visitor{res: "9", kind: token.INT},
			op:             token.GTR,
			expectedResult: "true",
		},
		{
			name:           "10 >= 9 = true",
			x:              &visitor{res: "10", kind: token.INT},
			y:              &visitor{res: "9", kind: token.INT},
			op:             token.GEQ,
			expectedResult: "true",
		},
		{
			name:           "10 >= 10 = true",
			x:              &visitor{res: "10", kind: token.INT},
			y:              &visitor{res: "10", kind: token.INT},
			op:             token.GEQ,
			expectedResult: "true",
		},
		{
			name:           "9 < 10 = true",
			x:              &visitor{res: "9", kind: token.INT},
			y:              &visitor{res: "10", kind: token.INT},
			op:             token.LSS,
			expectedResult: "true",
		},
		{
			name:           "9 <= 10 = true",
			x:              &visitor{res: "9", kind: token.INT},
			y:              &visitor{res: "10", kind: token.INT},
			op:             token.LEQ,
			expectedResult: "true",
		},
		{
			name:           "10 <= 10 = true",
			x:              &visitor{res: "10", kind: token.INT},
			y:              &visitor{res: "10", kind: token.INT},
			op:             token.LEQ,
			expectedResult: "true",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := &visitor{}
			v.comparison(tc.x, tc.y, tc.op)
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %s, got: %s", tc.expectedResult, v.res)
			}
		})
	}
}

func TestLogical(t *testing.T) {
	tt := []struct {
		name           string
		x              *visitor
		y              *visitor
		op             token.Token
		expectedResult string
	}{
		{
			name:           "true && true = true",
			x:              &visitor{res: "true"},
			y:              &visitor{res: "true"},
			op:             token.LAND,
			expectedResult: "true",
		},
		{
			name:           "true && false = false",
			x:              &visitor{res: "true"},
			y:              &visitor{res: "false"},
			op:             token.LAND,
			expectedResult: "false",
		},
		{
			name:           "false && false = false",
			x:              &visitor{res: "false"},
			y:              &visitor{res: "false"},
			op:             token.LAND,
			expectedResult: "false",
		},
		{
			name:           "true || true = true",
			x:              &visitor{res: "true"},
			y:              &visitor{res: "true"},
			op:             token.LOR,
			expectedResult: "true",
		},
		{
			name:           "true || false = true",
			x:              &visitor{res: "true"},
			y:              &visitor{res: "false"},
			op:             token.LOR,
			expectedResult: "true",
		},
		{
			name:           "false || false = false",
			x:              &visitor{res: "false"},
			y:              &visitor{res: "false"},
			op:             token.LOR,
			expectedResult: "false",
		},
		{
			name:           "true && 10 = ?",
			x:              &visitor{res: "true"},
			y:              &visitor{res: "10"},
			op:             token.LAND,
			expectedResult: "",
		},
		{
			name:           "10 || true = ?",
			x:              &visitor{res: "10"},
			y:              &visitor{res: "true"},
			op:             token.LOR,
			expectedResult: "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := &visitor{}
			v.logical(tc.x, tc.y, tc.op)
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %s, got: %s", tc.expectedResult, v.res)
			}
		})
	}
}

func TestVisitBinary(t *testing.T) {
	tt := []struct {
		name            string
		x               ast.Expr
		y               ast.Expr
		ops             []token.Token
		expectedResults []string
		expectedErrs    []error
	}{
		{
			name:            "comparison",
			x:               &ast.BasicLit{Value: "7", Kind: token.INT},
			y:               &ast.BasicLit{Value: "3", Kind: token.INT},
			ops:             []token.Token{token.EQL, token.NEQ, token.GTR, token.GEQ, token.LSS, token.LEQ},
			expectedResults: []string{"false", "true", "true", "true", "false", "false"},
		},
		{
			name:            "arithmetic",
			x:               &ast.BasicLit{Value: "45", Kind: token.INT},
			y:               &ast.BasicLit{Value: "5", Kind: token.INT},
			ops:             []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedResults: []string{"50", "40", "225", "9", "0"},
		},
		{
			name:            "logical",
			x:               &ast.BasicLit{Value: "true"},
			y:               &ast.BasicLit{Value: "false"},
			ops:             []token.Token{token.LAND, token.LOR},
			expectedResults: []string{"false", "true"},
		},
		{
			name:            "ErrUnsupportedOperator",
			x:               &ast.BasicLit{Value: "true"},
			y:               &ast.BasicLit{Value: "false"},
			ops:             []token.Token{token.ASSIGN},
			expectedResults: []string{""},
			expectedErrs:    []error{ErrUnsupportedOperator},
		},
		{
			name: "err on x: ErrUnsupportedOperator",
			x: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "true"},
				Y:  &ast.BasicLit{Value: "false"},
				Op: token.ASSIGN,
			},
			y:               &ast.BasicLit{Value: "false"},
			ops:             []token.Token{token.ASSIGN},
			expectedResults: []string{""},
			expectedErrs:    []error{ErrUnsupportedOperator},
		},
		{
			name: "err on Y: ErrUnsupportedOperator",
			x:    &ast.BasicLit{Value: "false"},
			y: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "true"},
				Y:  &ast.BasicLit{Value: "false"},
				Op: token.ASSIGN,
			},
			ops:             []token.Token{token.ASSIGN},
			expectedResults: []string{""},
			expectedErrs:    []error{ErrUnsupportedOperator},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := &visitor{}
			for i, op := range tc.ops {
				var (
					op             = op
					expectedResult = tc.expectedResults[i]
					expectedErr    error
				)

				if tc.expectedErrs != nil {
					expectedErr = tc.expectedErrs[i]
				}

				t.Run(op.String(), func(t *testing.T) {
					binaryExpr := &ast.BinaryExpr{X: tc.x, Y: tc.y, Op: op}
					_ = v.visitBinary(binaryExpr)
					if !errors.Is(v.err, expectedErr) {
						t.Fatalf("expected error: %v, got: %v", expectedErr, v.err)
					}
					if v.res != expectedResult {
						t.Fatalf("expected result: %s, got: %s", expectedResult, v.res)
					}
				})
			}
		})
	}
}

func TestVisit(t *testing.T) {
	tt := []struct {
		name           string
		node           ast.Node
		expectedResult string
	}{
		{
			name: "*ast.ParenExpr, (1) = 1",
			node: &ast.ParenExpr{
				X: &ast.BasicLit{Value: "1", Kind: token.INT},
			},
			expectedResult: "1",
		},
		{
			name: "*ast.UnaryExpr, !true = false",
			node: &ast.UnaryExpr{
				X:  &ast.BasicLit{Value: "true", Kind: token.STRING},
				Op: token.NOT,
			},
			expectedResult: "false",
		},
		{
			name: "*ast.BinaryExpr, 1 + 2 = 3",
			node: &ast.BinaryExpr{
				X:  &ast.BasicLit{Value: "1", Kind: token.INT},
				Y:  &ast.BasicLit{Value: "2", Kind: token.INT},
				Op: token.ADD,
			},
			expectedResult: "3",
		},
		{
			name:           "*ast.BasicLit, 1",
			node:           &ast.BasicLit{Value: "1", Kind: token.INT},
			expectedResult: "1",
		},
		{
			name:           "*ast.Ident, true",
			node:           &ast.Ident{Name: "true"},
			expectedResult: "true",
		},
		{
			name:           "not supported, e.g. *ast.ArrayType",
			node:           &ast.ArrayType{},
			expectedResult: "",
		},
		{
			name:           "nil node",
			node:           nil,
			expectedResult: "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := NewVisitor().(*visitor)
			_ = v.Visit(tc.node)
			if v.res != tc.expectedResult {
				t.Fatalf("expected result: %s, got: %s", tc.expectedResult, v.res)
			}
		})
	}
}

func TestResult(t *testing.T) {
	tt := []struct {
		name           string
		v              *visitor
		expectedResult bool
		expectedErr    error
	}{
		{name: "true", v: &visitor{res: "true"}, expectedResult: true},
		{name: "false", v: &visitor{res: "false"}, expectedResult: false},
		{name: "false", v: &visitor{err: ErrIntegerDividedByZero}, expectedErr: ErrIntegerDividedByZero},
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
