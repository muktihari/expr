package expr_test

import (
	"math"
	"testing"

	"github.com/muktihari/expr"
)

var exprs = [...]string{
	expr.KindBoolean: "((3 + 2 == 5) && (7 * 2 <= 10)) || (!(4 + 1 == 5) && (2 / 1 < 5)) || ((8 >= 10) && (!(6 + 3 > 9))) || (9 % 4 == 1)", // true
	expr.KindInt:     "((10 + 2 * 5) / (7 - 3)) + ((15 - 2 * 3) * 2)",                                                                      // int64: 23
	expr.KindFloat:   "(((12 + 5) * (7 - 3)) / (2 + 1) + (6 * (4 - 2)) - 5) * 2 + (8 / 2.2)",                                               // float64: 62.9696969697
	expr.KindImag:    "((3 + 2i) * (2 - 4i)) / ((1 + 3i) + (2i - 1))",                                                                      // complex128: (-1.6-2.8i)
}

func BenchmarkAny(b *testing.B) {
	results := [...]interface{}{
		expr.KindBoolean: true,
		expr.KindInt:     int64(23),
		expr.KindFloat:   float64(62.9696969697),
		expr.KindImag:    complex(-1.6, -2.8),
	}

	for i, e := range exprs {
		i, e := i, e
		kind := expr.Kind(i)
		if e == "" {
			continue
		}

		b.Run(kind.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v, err := expr.Any(e)
				if err != nil {
					b.Fatalf("expected nil, got: %v", err)
				}

				// handle float prec
				vf, vok := v.(float64)
				rf, rok := results[kind].(float64)
				if vok && rok {
					digit := 1e10
					v := math.Round(vf*digit) / digit
					r := math.Round(rf*digit) / digit
					if v != r {
						b.Fatalf("expected value %v, got: %v", r, v)
					}
					continue
				}

				if v != results[kind] {
					b.Fatalf("expected value %v, got: %v", results[kind], v)
				}
			}
		})
	}
}

func BenchmarkBoolean(b *testing.B) {
	var (
		e string = exprs[expr.KindBoolean]
		r bool   = true
	)

	b.Run(expr.KindBoolean.String(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := expr.Bool(e)
			if err != nil {
				b.Fatalf("expected nil, got: %v", err)
			}
			if v != r {
				b.Fatalf("expected %t, got: %t", true, v)
			}
		}
	})

}

func BenchmarkComplex128(b *testing.B) {
	var (
		e string     = exprs[expr.KindImag]
		r complex128 = complex(-1.6, -2.8)
	)

	b.Run(expr.KindImag.String(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := expr.Complex128(e)
			if err != nil {
				b.Fatalf("expected nil, got: %v", err)
			}
			if v != r {
				b.Fatalf("expected value: %f, got: %f", r, v)
			}
		}
	})
}

func BenchmarkFloat64(b *testing.B) {
	var (
		e string  = exprs[expr.KindFloat]
		r float64 = 62.9696969697
	)

	b.Run(expr.KindFloat.String(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := expr.Float64(e)
			if err != nil {
				b.Fatalf("expected nil, got: %v", err)
			}
			// handle float prec
			digit := 1e10
			v = math.Round(v*digit) / digit
			r = math.Round(r*digit) / digit
			if v != r {
				b.Fatalf("expected value: %f, got: %f", r, v)
			}
		}
	})
}

func BenchmarkInt64(b *testing.B) {
	var (
		e string = exprs[expr.KindInt]
		r int64  = 23
	)

	b.Run(expr.KindInt.String(), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := expr.Int64(e)
			if err != nil {
				b.Fatalf("expected nil, got: %v", err)
			}
			if v != r {
				b.Fatalf("expected value: %d, got: %d", r, v)
			}
		}
	})
}
