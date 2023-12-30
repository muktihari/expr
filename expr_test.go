// Copyright 2020-2023 The Expr Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expr_test

import (
	"errors"
	"math"
	"testing"

	"github.com/muktihari/expr"
)

func TestAny(t *testing.T) {
	tt := []struct {
		In  string
		Eq  interface{}
		Err error
	}{
		{In: "2", Eq: int64(2)},
		{In: "\"2\"", Eq: "2"},
		{In: "2.5", Eq: float64(2.5)},
		{In: "4 == 2", Eq: false},
		{In: "1 + 1 + (4 == 2)", Err: expr.ErrArithmeticOperation},
		{In: "2 && 2", Err: expr.ErrLogicalOperation},
		{In: "4 - 2", Eq: int64(2)},
		{In: "4 * 2", Eq: int64(8)},
		{In: "4 / 2", Eq: int64(2)},
		{In: "(2 + 2) * 10", Eq: int64(40)},
		{In: "(2 + 2) / 10", Eq: float64(0.4)},
		{In: "(2.0 + 2) / 10", Eq: float64(0.4)},
		{In: "(2 * 2) * (8 + 2)", Eq: int64(40)},
		{In: "(2 * 2) * (8 + 2) * 2", Eq: int64(80)},
		{In: "((2 * 2) * (8 + 2) * 2) + 1", Eq: int64(81)},
		{In: "((2 * 2) * (8 + 2) * 2) + 1.5", Eq: float64(81.5)},
		{In: "((2 * 2) * (8 + 2) * 2) + 2.56789", Eq: float64(82.56789)},
		{In: "1 + 2 + 3 + 4 + 5", Eq: int64(15)},
		{In: "(2 + 2) * 4 / 4", Eq: int64(4)},
		{In: "((2 + 2) * 4 / 4) * 10", Eq: int64(40)},
		{In: "((2 + 2) * 4 / 4) * 10 + 2", Eq: int64(42)},
		{In: "((2 + 2) * 4 / 4) * 10 + 2 - 2", Eq: int64(40)},
		{In: "((2 + 2) * 4 / 4) * 10 + 4.234567", Eq: float64(44.234567)},
		{In: "((2 + 2) * 4 / 4) * 10.5 + 4.234567", Eq: float64(46.234567)},
		{In: "((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)", Eq: float64(466.2567)},
		{In: "10 * -5", Eq: int64(-50)},
		{In: "10 * (-5-5)", Eq: int64(-100)},
		{In: "10 * -5 + (-5.5)", Eq: float64(-55.5)},
		{In: "10 + (10 * -10)", Eq: int64(-90)},
		{In: "10 + ((-5 * -10) * 10)", Eq: int64(510)},
		{In: "10 + ((-5 * -10) / -10) - 2", Eq: int64(3)},
		{In: "0 / 10", Eq: int64(0)},
		{In: "12.5 | 4.3", Err: expr.ErrBitwiseOperation},
		{In: "12 | 4", Eq: int64(12)},
		{In: "4 << 10", Eq: int64(4096)},
		{In: "(10+5i) + (10+7i)", Eq: complex(20, 12)},
		{In: "(2+3i) - (2+2i)", Eq: complex(0, 1)},
		{In: "(2+2i) * (2+2i)", Eq: complex(0, 8)},
		{In: "(2+2i) / (2+2i)", Eq: complex(1, 0)},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Any(tc.In)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected error: %v, got: %v", tc.Err, err)
			}

			if v != tc.Eq {
				t.Fatalf("expected value: %f, got: %f", tc.Eq, v)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tt := []struct {
		In  string
		Eq  bool
		Err error
	}{
		{In: "1", Err: expr.ErrValueTypeMismatch},
		{In: "!7", Err: expr.ErrUnaryOperation},
		{In: "1 < 2", Eq: true},
		{In: "2 < 1", Eq: false},
		{In: "2 < 1 && (1 + 1) > 1", Eq: false},
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
		{In: "true", Eq: true},
		{In: "!false", Eq: true},
		{In: "!false || false", Eq: true},
		{In: "(-10 < -2) && -1 > -2", Eq: true},
		{In: "-(-1) > -1", Eq: true},
		{In: "-(-1.5) > +1.3", Eq: true},
		{In: "-4 * -2 > -1", Eq: true},
		{In: "10 % 2 > -2", Eq: true},
		{In: "10 % 2 < 1", Eq: true},
		{In: "10.2 % 2 > 2", Eq: false}, // 0.2 > 2
		{In: "10.2 % 2 < 2", Eq: true},  // 0.2 < 2
		{In: `"a" > "a"`, Eq: false},
		{In: `"a" >= "b"`, Eq: false},
		{In: `"b" > "a"`, Eq: true},
		{In: `"b" >= "a"`, Eq: true},
		{In: `"a" < "b"`, Eq: true},
		{In: `"a" <= "b"`, Eq: true},
		{In: `"b" < "a"`, Eq: false},
		{In: `"b" <= "a"`, Eq: false},
		{In: `"a" <= "a"`, Eq: true},
		{In: `"a" >= "a"`, Eq: true},
		{In: `0x4 << 0xA > 1024`, Eq: true},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Bool(tc.In)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected %v, got: %v", tc.Err, err)
			}
			if v != tc.Eq {
				t.Fatalf("expected %v, got: %v", tc.Eq, v)
			}
		})
	}
}

func TestComplex128(t *testing.T) {
	tt := []struct {
		In  string
		Eq  complex128
		Err error
	}{
		{In: "2", Eq: complex(2, 0)},
		{In: "2.5", Eq: complex(2.5, 0)},
		{In: "4 == 2", Err: expr.ErrValueTypeMismatch},
		{In: "1 + 1 + (4 == 2)", Err: expr.ErrArithmeticOperation},
		{In: "2 && 2", Eq: 0, Err: expr.ErrLogicalOperation},
		{In: "4 - 2", Eq: complex(2, 0)},
		{In: "4 * 2", Eq: complex(8, 0)},
		{In: "4 / 2", Eq: complex(2, 0)},
		{In: "(2 + 2) * 10", Eq: complex(40, 0)},
		{In: "(2 + 2) / 10", Eq: complex(0.4, 0)},
		{In: "(2 * 2) * (8 + 2)", Eq: complex(40, 0)},
		{In: "(2 * 2) * (8 + 2) * 2", Eq: complex(80, 0)},
		{In: "((2 * 2) * (8 + 2) * 2) + 1", Eq: complex(81, 0)},
		{In: "((2 * 2) * (8 + 2) * 2) + 1.5", Eq: complex(81.5, 0)},
		{In: "((2 * 2) * (8 + 2) * 2) + 2.56789", Eq: complex(82.56789, 0)},
		{In: "1 + 2 + 3 + 4 + 5", Eq: complex(15, 0)},
		{In: "(2 + 2) * 4 / 4", Eq: complex(4, 0)},
		{In: "((2 + 2) * 4 / 4) * 10", Eq: complex(40, 0)},
		{In: "((2 + 2) * 4 / 4) * 10 + 2", Eq: complex(42, 0)},
		{In: "((2 + 2) * 4 / 4) * 10 + 2 - 2", Eq: complex(40, 0)},
		{In: "((2 + 2) * 4 / 4) * 10 + 4.234567", Eq: complex(44.234567, 0)},
		{In: "((2 + 2) * 4 / 4) * 10.5 + 4.234567", Eq: complex(46.234567, 0)},
		{In: "((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)", Eq: complex(466.2567, 0)},
		{In: "10 * -5", Eq: complex(-50, 0)},
		{In: "10 * (-5-5)", Eq: complex(-100, 0)},
		{In: "10 * -5 + (-5.5)", Eq: complex(-55.5, 0)},
		{In: "10 + (10 * -10)", Eq: complex(-90, 0)},
		{In: "10 + ((-5 * -10) * 10)", Eq: complex(510, 0)},
		{In: "10 + ((-5 * -10) / -10) - 2", Eq: complex(3, 0)},
		{In: "10 / 0", Eq: complex(math.Inf(+1), math.NaN())}, // IEEE 754 says that only NaNs satisfy f != f.
		{In: "0 / 10", Eq: 0},
		{In: "12.5 | 4.3", Err: expr.ErrBitwiseOperation},
		{In: "12 | 4", Err: expr.ErrBitwiseOperation},
		{In: "(10+5i) + (10+7i)", Eq: complex(20, 12)},
		{In: "(2+3i) - (2+2i)", Eq: complex(0, 1)},
		{In: "(2+2i) * (2+2i)", Eq: complex(0, 8)},
		{In: "(2+2i) / (2+2i)", Eq: complex(1, 0)},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Complex128(tc.In)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected error: %v, got: %v", tc.Err, err)
			}

			if vIsNaN, eqIsNaN := isComplexNaN(v), isComplexNaN(tc.Eq); vIsNaN == true || eqIsNaN == true {
				// If result is a NaN, the expected value should also be a NaN.
				if vIsNaN != eqIsNaN {
					t.Fatalf("expected value: %v (isNaN: %v), got: %v (isNaN: %v)",
						tc.Eq, eqIsNaN, v, vIsNaN)
				}
				return
			}

			if v != tc.Eq {
				t.Fatalf("expected value: %f, got: %f", tc.Eq, v)
			}
		})
	}
}

func isComplexNaN(v complex128) bool { return math.IsNaN(real(v)) || math.IsNaN(imag(v)) }

func TestFloat64(t *testing.T) {
	tt := []struct {
		In  string
		Eq  float64
		Err error
	}{
		{In: "2", Eq: 2},
		{In: "4 == 2", Err: expr.ErrValueTypeMismatch},
		{In: "1 + 1 + (4 == 2)", Err: expr.ErrArithmeticOperation},
		{In: "2 && 2", Eq: 0, Err: expr.ErrLogicalOperation},
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
		{In: "10 * -5", Eq: -50},
		{In: "10 * (-5-5)", Eq: -100},
		{In: "10 * -5 + (-5.5)", Eq: -55.5},
		{In: "10 + (10 * -10)", Eq: -90},
		{In: "10 + ((-5 * -10) * 10)", Eq: 510},
		{In: "10 + ((-5 * -10) / -10) - 2", Eq: 3},
		{In: "10 / 0", Eq: math.Inf(+1)},
		{In: "0 / 10", Eq: 0},
		{In: "10.0 % 2.6", Eq: 2.2},
		{In: "12.5 | 4.3", Err: expr.ErrBitwiseOperation},
		{In: "12 | 4", Err: expr.ErrBitwiseOperation},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Float64(tc.In)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected error: %v, got: %v", tc.Err, err)
			}

			v = math.Round(v*1000000) / 1000000 // round up to 6 decimal
			if v != tc.Eq {
				t.Fatalf("expected value: %f, got: %f", tc.Eq, v)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tt := []struct {
		In  string
		Eq  int
		Err error
	}{
		{In: "4.23", Eq: 4},
		{In: "4 == 2", Err: expr.ErrValueTypeMismatch},
		{In: "1 + 1 + (4 == 2)", Err: expr.ErrArithmeticOperation},
		{In: "2 && 2", Eq: 0, Err: expr.ErrLogicalOperation},
		{In: "2", Eq: 2},
		{In: "2 + 2", Eq: 4},
		{In: "4 - 2", Eq: 2},
		{In: "4 * 2", Eq: 8},
		{In: "4 / 2", Eq: 2},
		{In: "4 || 2", Eq: 0, Err: expr.ErrLogicalOperation},
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
		{In: "10 * -5", Eq: -50},
		{In: "10 * (-5-5)", Eq: -100},
		{In: "10 * -5 + (-5)", Eq: -55},
		{In: "10 + (10 * -10)", Eq: -90},
		{In: "10 + ((-5 * -10) * 10)", Eq: 510},
		{In: "10 + ((-5 * -10) / -10) - 2", Eq: 3},
		{In: "10 / 0", Eq: 0},
		{In: "0b1100 | 0b0100", Eq: 12}, // = 1111
		{In: "0b1100 ^ 0b0100", Eq: 8},  // = 1011
		{In: "0b1100 & 0b0100", Eq: 4},  // = 0100
		{In: "0b1100 &^ 0b0100", Eq: 8}, // = 1011
		{In: "12 | 4", Eq: 12},
		{In: "12 ^ 4", Eq: 8},
		{In: "12 & 4", Eq: 4},
		{In: "12 &^ 4", Eq: 8},
		{In: "4 << 10", Eq: 4096},
		{In: "10 >> 1", Eq: 5},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Int(tc.In)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected error: %v, got: %v", tc.Err, err)
			}
			if v != tc.Eq {
				t.Fatalf("expected %d, got: %d", tc.Eq, v)
			}
		})
	}
}

func TestInt64Strict(t *testing.T) {
	tt := []struct {
		In  string
		Eq  int64
		Err error
	}{
		{In: "4.23", Eq: 4},
		{In: "4/0", Err: expr.ErrIntegerDividedByZero},
		{In: "11 + 7", Eq: 18},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.In, func(t *testing.T) {
			v, err := expr.Int64Strict(tc.In)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected error: %v, got: %v", tc.Err, err)
			}
			if v != tc.Eq {
				t.Fatalf("expected %d, got: %d", tc.Eq, v)
			}
		})
	}
}
