package expr_test

import (
	"fmt"
	"testing"

	"github.com/muktihari/expr"
)

func BenchmarkInt(b *testing.B) {
	intExpr := "((2 + 2) * 4 / 4) * 10 + 4 * (50 + 50)"
	multiply := func(s string, n int) (string, int) {
		eq := 440
		for i := 1; i < n; i++ {
			s = s + " + (" + s + ")"
			eq += eq
		}
		return s, eq
	}

	for n := 1; n <= 16; n *= 2 {
		s, eq := multiply(intExpr, n)
		b.Run(fmt.Sprintf("%d [^%d]", len(s), n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, err := expr.Int(s)
				if err != nil {
					b.Fatalf("expected nil, got: %v", err)
				}
				if v != eq {
					b.Fatalf("expected %d, got: %d", eq, v)
				}
			}
		})
	}
}

func BenchmarkFloat64(b *testing.B) {
	floatExpr := "((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)"
	multiply := func(s string, n int) (string, float64) {
		eq := 466.2567
		for i := 1; i < n; i++ {
			s = s + " + (" + s + ")"
			eq += eq
		}
		return s, eq
	}

	for n := 1; n <= 16; n *= 2 {
		s, eq := multiply(floatExpr, n)
		b.Run(fmt.Sprintf("%d [^%d]", len(s), n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, err := expr.Float64(s)
				if err != nil {
					b.Fatalf("expected nil, got: %v", err)
				}
				if v != eq {
					b.Fatalf("expected %f, got: %f", eq, v)
				}
			}
		})
	}

	b.Fail()
}

func BenchmarkBoolean(b *testing.B) {
	booleanExpr := "(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || !false && 1 < 2 && (1 > 2 || -1 > -2)"
	multiply := func(s string, n int) string {
		for i := 1; i < n; i++ {
			s += " && (" + s + ")"
		}
		return s
	}

	for n := 1; n <= 16; n *= 2 {
		s := multiply(booleanExpr, n)
		b.Run(fmt.Sprintf("%d [^%d]", len(s), n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, err := expr.Bool(s)
				if err != nil {
					b.Fatalf("expected nil, got: %v", err)
				}

				if v != true {
					b.Fatalf("expected nil, got: %v", v)
				}
			}
		})
	}
}
