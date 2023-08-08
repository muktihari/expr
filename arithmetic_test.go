package expr

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"testing"
)

func TestArithmetic(t *testing.T) {
	newVisitor := func(numericType NumericType) *Visitor {
		return &Visitor{options: options{numericType: numericType}}
	}
	tt := []struct {
		name          string
		v, vx, vy     *Visitor
		op            token.Token
		expectedValue interface{}
		expectedErr   error
	}{
		{
			name:          "arithmetic non numeric error: x numeric y boolean",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: int64(1), kind: KindInt},
			vy:            &Visitor{value: true, kind: KindBoolean},
			op:            token.ADD,
			expectedValue: nil,
			expectedErr:   ErrArithmeticOperation,
		},
		{
			name:          "arithmetic non numeric error: x boolean y numeric",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: false, kind: KindBoolean},
			vy:            &Visitor{value: int64(10), kind: KindInt},
			op:            token.ADD,
			expectedValue: nil,
			expectedErr:   ErrArithmeticOperation,
		},
		{
			name:          "arithmetic numeric complex",
			v:             newVisitor(NumericTypeComplex),
			vx:            &Visitor{value: (1 + 0i), kind: KindImag},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: (3 + 0i),
		},
		{
			name:          "arithmetic numeric float",
			v:             newVisitor(NumericTypeFloat),
			vx:            &Visitor{value: float64(1.5), kind: KindFloat},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: 3.5,
		},
		{
			name:          "arithmetic numeric int",
			v:             newVisitor(NumericTypeInt),
			vx:            &Visitor{value: float64(1.5), kind: KindFloat},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: int64(3),
		},
		{
			name:          "arithmetic numeric auto: x imag y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: (1.5 + 1i), kind: KindImag},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: (3.5 + 1i),
		},
		{
			name:          "arithmetic numeric auto: x int y imag",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: int64(2), kind: KindInt},
			vy:            &Visitor{value: (1.5 + 1i), kind: KindImag},
			op:            token.ADD,
			expectedValue: (3.5 + 1i),
		},
		{
			name:          "arithmetic numeric auto: x float y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: float64(1.5), kind: KindFloat},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: float64(3.5),
		},
		{
			name:          "arithmetic numeric auto: x int y float",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: float64(1.5), kind: KindFloat},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: float64(3.5),
		},
		{
			name:          "arithmetic numeric auto: x int y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: int64(1), kind: KindInt},
			vy:            &Visitor{value: int64(2), kind: KindInt},
			op:            token.ADD,
			expectedValue: float64(3),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			be := &ast.BinaryExpr{Op: tc.op}
			arithmetic(tc.v, tc.vx, tc.vy, be)
			if !errors.Is(tc.v.err, tc.expectedErr) {
				t.Fatalf("expected err: %s, got: %s", tc.expectedErr, tc.v.err)
			}
			if tc.v.value != tc.expectedValue {
				t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValue, tc.expectedValue,
					tc.v.value, tc.v.value)
			}
		})
	}
}

func TestCalculateComplex(t *testing.T) {
	newComplexVisitor := func() *Visitor {
		return &Visitor{options: options{numericType: NumericTypeComplex}}
	}
	tt := []struct {
		name           string
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []interface{}
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: int64(1), kind: KindInt},
			vy:             &Visitor{value: int64(2), kind: KindInt},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []interface{}{(3 + 0i), (-1 + 0i), (2 + 0i), (0.5 + 0i)},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "calculate floats",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: float64(1.0), kind: KindFloat},
			vy:             &Visitor{value: float64(2.0), kind: KindFloat},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []interface{}{(3 + 0i), (-1 + 0i), (2 + 0i), (0.5 + 0i)},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: (1 + 1i), kind: KindImag},
			vy:             &Visitor{value: (2 + 1i), kind: KindImag},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []interface{}{(3 + 2i), (-1 + 0i), (1 + 3i), (0.6 + 0.2i)},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "unsupported complex operation",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: (1 + 1i), kind: KindImag},
			vy:             &Visitor{value: (2 + 1i), kind: KindImag},
			ops:            []token.Token{token.REM},
			expectedValues: []interface{}{nil},
			expectedErrs:   []error{ErrArithmeticOperation},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := fmt.Sprintf("%v%s%v", tc.vx.value, op, tc.vy.value)
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateComplex(tc.v, tc.vx, tc.vy, be)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %s, got: %s", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value != tc.expectedValues[i] {
						t.Fatalf("expected value: %v (%T), got: % (%T)", tc.expectedValues[i], tc.expectedValues[i],
							tc.v.value, tc.v.value)
					}
				})
			}
		})
	}
}

func TestCalculateFloat(t *testing.T) {
	newFloatVisitor := func() *Visitor {
		return &Visitor{options: options{numericType: NumericTypeFloat}}
	}

	tt := []struct {
		name           string
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []interface{}
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: int64(10), kind: KindInt},
			vy:             &Visitor{value: int64(3), kind: KindInt},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []interface{}{float64(13), float64(7), float64(30), float64(3.3333333333333335), float64(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate floats",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: float64(10.0), kind: KindFloat},
			vy:             &Visitor{value: float64(3.0), kind: KindFloat},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []interface{}{float64(13), float64(7), float64(30), float64(3.3333333333333335), float64(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: (10 + 1i), kind: KindImag},
			vy:             &Visitor{value: (3 + 1i), kind: KindImag},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []interface{}{float64(13), float64(7), float64(30), float64(3.3333333333333335), float64(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := fmt.Sprintf("%v%s%v", tc.vx.value, op, tc.vy.value)
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateFloat(tc.v, tc.vx, tc.vy, be)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %s, got: %s", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value != tc.expectedValues[i] {
						t.Fatalf("expected value: %v (%T), got: %s (%T)", tc.expectedValues[i], tc.expectedValues[i],
							tc.v.value, tc.v.value)
					}
				})
			}
		})
	}
}

func TestCalculateInt(t *testing.T) {
	newIntVisitor := func() *Visitor {
		return &Visitor{options: options{
			numericType:               NumericTypeInt,
			allowIntegerDividedByZero: true,
		}}
	}

	tt := []struct {
		name           string
		v, vx, vy      *Visitor
		ops            []token.Token
		expectedValues []interface{}
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newIntVisitor(),
			vx:             &Visitor{value: int64(10), kind: KindInt},
			vy:             &Visitor{value: int64(3), kind: KindInt},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []interface{}{int64(13), int64(7), int64(30), int64(3), int64(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate integers allowIntegerDividedByZero == true",
			v:              &Visitor{options: options{numericType: NumericTypeInt, allowIntegerDividedByZero: true}},
			vx:             &Visitor{value: int64(10), kind: KindInt},
			vy:             &Visitor{value: int64(0), kind: KindInt},
			ops:            []token.Token{token.QUO},
			expectedValues: []interface{}{int64(0)},
			expectedErrs:   []error{nil},
		},
		{
			name:           "calculate integers allowIntegerDividedByZero == false",
			v:              &Visitor{options: options{numericType: NumericTypeInt, allowIntegerDividedByZero: false}},
			vx:             &Visitor{value: int64(10), kind: KindInt},
			vy:             &Visitor{value: int64(0), kind: KindInt},
			ops:            []token.Token{token.QUO},
			expectedValues: []interface{}{nil},
			expectedErrs:   []error{ErrIntegerDividedByZero},
		},
		{
			name:           "calculate floats",
			v:              newIntVisitor(),
			vx:             &Visitor{value: float64(10.0), kind: KindFloat},
			vy:             &Visitor{value: float64(3.0), kind: KindFloat},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []interface{}{int64(13), int64(7), int64(30), int64(3), int64(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newIntVisitor(),
			vx:             &Visitor{value: (10 + 1i), kind: KindImag},
			vy:             &Visitor{value: (3 + 1i), kind: KindImag},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []interface{}{int64(13), int64(7), int64(30), int64(3), int64(1)},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := fmt.Sprintf("%v%s%v", tc.vx.value, op, tc.vy.value)
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateInt(tc.v, tc.vx, tc.vy, be)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %s, got: %s", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value != tc.expectedValues[i] {
						t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValues[i], tc.expectedValues[i],
							tc.v.value, tc.v.value)
					}
				})
			}
		})
	}
}
