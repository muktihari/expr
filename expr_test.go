package expr_test

import (
	"testing"

	"github.com/muktihari/expr"
)

func TestInt(t *testing.T) {
	tt := []struct {
		In string
		Eq int
	}{
		{In: "2 + 2", Eq: 4},
		{In: "4 - 2", Eq: 2},
		{In: "4 * 2", Eq: 8},
		{In: "4 / 2", Eq: 2},
		{In: "(2 + 2) * 10", Eq: 40},
		{In: "(2 + 2) / 10", Eq: 0},
		{In: "(2 * 2) * (8 + 2)", Eq: 40},
		{In: "(2 * 2) * (8 + 2) * 2", Eq: 80},
		{In: "((2 * 2) * (8 + 2) * 2) + 1", Eq: 81},
		{In: "((2 * 2) * (8 + 2) * 2) + 1.5", Eq: 81},
		{In: "((2 * 2) * (8 + 2) * 2) + 2.56789", Eq: 82},
		{In: "1 + 2 + 3 + 4 + 5", Eq: 15},
		{In: "(2 + 2) * 4 / 4", Eq: 4},
		{In: "((2 + 2) * 4 / 4) * 10", Eq: 40},
		{In: "((2 + 2) * 4 / 4) * 10 + 2", Eq: 42},
		{In: "((2 + 2) * 4 / 4) * 10 + 2 - 2", Eq: 40},
		{In: "((2 + 2) * 4 / 4) * 10 + 4.234567", Eq: 44},
		{In: "((2 + 2) * 4 / 4) * 10.5 + 4.234567", Eq: 44},
		{In: "((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)", Eq: 440},
	}

	for _, tc := range tt {
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Int(tc.In)
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			if v != tc.Eq {
				t.Fatalf("expected %d, got: %d", tc.Eq, v)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tt := []struct {
		In string
		Eq float64
	}{
		{In: "2 + 2", Eq: 4},
		{In: "4 - 2", Eq: 2},
		{In: "4 * 2", Eq: 8},
		{In: "4 / 2", Eq: 2},
		{In: "(2 + 2) * 10", Eq: 40},
		{In: "(2 + 2) / 10", Eq: 0.4},
		{In: "(2 * 2) * (8 + 2)", Eq: 40},
		{In: "(2 * 2) * (8 + 2) * 2", Eq: 80},
		{In: "((2 * 2) * (8 + 2) * 2) + 1", Eq: 81},
		{In: "((2 * 2) * (8 + 2) * 2) + 1.5", Eq: 81.5},
		{In: "((2 * 2) * (8 + 2) * 2) + 2.56789", Eq: 82.56789},
		{In: "1 + 2 + 3 + 4 + 5", Eq: 15},
		{In: "(2 + 2) * 4 / 4", Eq: 4},
		{In: "((2 + 2) * 4 / 4) * 10", Eq: 40},
		{In: "((2 + 2) * 4 / 4) * 10 + 2", Eq: 42},
		{In: "((2 + 2) * 4 / 4) * 10 + 2 - 2", Eq: 40},
		{In: "((2 + 2) * 4 / 4) * 10 + 4.234567", Eq: 44.234567},
		{In: "((2 + 2) * 4 / 4) * 10.5 + 4.234567", Eq: 46.234567},
		{In: "((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)", Eq: 466.2567},
	}

	for _, tc := range tt {
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Float64(tc.In)
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}

			if v != tc.Eq {
				t.Fatalf("expected: %f, got: %f", tc.Eq, v)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tt := []struct {
		In string
		Eq bool
	}{
		{In: "1 < 2", Eq: true},
		{In: "2 < 1", Eq: false},
		{In: "(1 < 2 && 3 > 4) || 1 == 1", Eq: true},
		{In: "((1 < 2 && 3 > 4) || 1 == 1) && 4 > 5", Eq: false},
		{In: "false && false", Eq: false},
		{In: "true && false", Eq: false},
		{In: "true && true", Eq: true},
		{In: "true && false || true", Eq: true},
		{In: "true && (false || true)", Eq: true},
		{In: "1 < 2 && 3 < 4 && ( 1==1 || 12 > 4)", Eq: true},
		{In: "\"expr\" == \"expr\" && \"Expr\" != \"expr\"", Eq: true},
		{In: "\"expr\" == \"expr\" && \"Expr\" == \"expr\"", Eq: false},
		{In: "(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || 1 == 1 ", Eq: true},
		{In: "(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || true == true ", Eq: true},
		{In: "(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || true == false ", Eq: false},
	}

	for _, tc := range tt {
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Bool(tc.In)
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
			if v != tc.Eq {
				t.Fatalf("expected %v, got: %v", tc.Eq, v)
			}
		})
	}
}
