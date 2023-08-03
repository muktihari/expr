package expr

import (
	"errors"
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
		expectedValue string
		expectedErr   error
	}{
		{
			name:          "arithmetic non numeric error: x numeric y boolean",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "1", kind: KindInt},
			vy:            &Visitor{value: "true", kind: KindBoolean},
			op:            token.ADD,
			expectedValue: "",
			expectedErr:   ErrArithmeticOperation,
		},
		{
			name:          "arithmetic non numeric error: x boolean y numeric",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "false", kind: KindBoolean},
			vy:            &Visitor{value: "10", kind: KindInt},
			op:            token.ADD,
			expectedValue: "",
			expectedErr:   ErrArithmeticOperation,
		},
		{
			name:          "arithmetic numeric complex",
			v:             newVisitor(NumericTypeComplex),
			vx:            &Visitor{value: "(1+0i)", kind: KindImag},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "(3+0i)",
		},
		{
			name:          "arithmetic numeric float",
			v:             newVisitor(NumericTypeFloat),
			vx:            &Visitor{value: "1.5", kind: KindFloat},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "3.5",
		},
		{
			name:          "arithmetic numeric int",
			v:             newVisitor(NumericTypeInt),
			vx:            &Visitor{value: "1.5", kind: KindFloat},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "3",
		},
		{
			name:          "arithmetic numeric auto: x imag y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "(1.5+1i)", kind: KindImag},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "(3.5+1i)",
		},
		{
			name:          "arithmetic numeric auto: x int y imag",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "2", kind: KindInt},
			vy:            &Visitor{value: "(1.5+1i)", kind: KindImag},
			op:            token.ADD,
			expectedValue: "(3.5+1i)",
		},
		{
			name:          "arithmetic numeric auto: x float y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "1.5", kind: KindFloat},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "3.5",
		},
		{
			name:          "arithmetic numeric auto: x int y float",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "1.5", kind: KindFloat},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "3.5",
		},
		{
			name:          "arithmetic numeric auto: x int y int",
			v:             newVisitor(NumericTypeAuto),
			vx:            &Visitor{value: "1", kind: KindInt},
			vy:            &Visitor{value: "2", kind: KindInt},
			op:            token.ADD,
			expectedValue: "3",
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
				t.Fatalf("expected value: %s, got: %s", tc.expectedValue, tc.v.value)
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
		expectedValues []string
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: "1", kind: KindInt},
			vy:             &Visitor{value: "2", kind: KindInt},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []string{"(3+0i)", "(-1+0i)", "(2+0i)", "(0.5+0i)"},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "calculate floats",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: "1.0", kind: KindFloat},
			vy:             &Visitor{value: "2.0", kind: KindFloat},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []string{"(3+0i)", "(-1+0i)", "(2+0i)", "(0.5+0i)"},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: "(1+1i)", kind: KindImag},
			vy:             &Visitor{value: "(2+1i)", kind: KindImag},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO},
			expectedValues: []string{"(3+2i)", "(-1+0i)", "(1+3i)", "(0.6+0.2i)"},
			expectedErrs:   []error{nil, nil, nil, nil},
		},
		{
			name:           "unsupported complex operation",
			v:              newComplexVisitor(),
			vx:             &Visitor{value: "(1+1i)", kind: KindImag},
			vy:             &Visitor{value: "(2+1i)", kind: KindImag},
			ops:            []token.Token{token.REM},
			expectedValues: []string{""},
			expectedErrs:   []error{ErrArithmeticOperation},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := tc.vx.value + op.String() + tc.vy.value
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateComplex(tc.v, tc.vx, tc.vy, be)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %s, got: %s", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value != tc.expectedValues[i] {
						t.Fatalf("expected value: %s, got: %s", tc.expectedValues[i], tc.v.value)
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
		expectedValues []string
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: "10", kind: KindInt},
			vy:             &Visitor{value: "3", kind: KindInt},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []string{"13", "7", "30", "3.3333333333333335", "1"},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate floats",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: "10.0", kind: KindFloat},
			vy:             &Visitor{value: "3.0", kind: KindFloat},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []string{"13", "7", "30", "3.3333333333333335", "1"},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newFloatVisitor(),
			vx:             &Visitor{value: "(10+1i)", kind: KindImag},
			vy:             &Visitor{value: "(3+1i)", kind: KindImag},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []string{"13", "7", "30", "3.3333333333333335", "1"},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := tc.vx.value + op.String() + tc.vy.value
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateFloat(tc.v, tc.vx, tc.vy, be)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %s, got: %s", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value != tc.expectedValues[i] {
						t.Fatalf("expected value: %s, got: %s", tc.expectedValues[i], tc.v.value)
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
		expectedValues []string
		expectedErrs   []error
	}{
		{
			name:           "calculate integers",
			v:              newIntVisitor(),
			vx:             &Visitor{value: "10", kind: KindInt},
			vy:             &Visitor{value: "3", kind: KindInt},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []string{"13", "7", "30", "3", "1"},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate integers allowIntegerDividedByZero == true",
			v:              &Visitor{options: options{numericType: NumericTypeInt, allowIntegerDividedByZero: true}},
			vx:             &Visitor{value: "10", kind: KindInt},
			vy:             &Visitor{value: "0", kind: KindInt},
			ops:            []token.Token{token.QUO},
			expectedValues: []string{"0"},
			expectedErrs:   []error{nil},
		},
		{
			name:           "calculate integers allowIntegerDividedByZero == false",
			v:              &Visitor{options: options{numericType: NumericTypeInt, allowIntegerDividedByZero: false}},
			vx:             &Visitor{value: "10", kind: KindInt},
			vy:             &Visitor{value: "0", kind: KindInt},
			ops:            []token.Token{token.QUO},
			expectedValues: []string{""},
			expectedErrs:   []error{ErrIntegerDividedByZero},
		},
		{
			name:           "calculate floats",
			v:              newIntVisitor(),
			vx:             &Visitor{value: "10.0", kind: KindFloat},
			vy:             &Visitor{value: "3.0", kind: KindFloat},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []string{"13", "7", "30", "3", "1"},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
		{
			name:           "calculate complex numbers",
			v:              newIntVisitor(),
			vx:             &Visitor{value: "(10+1i)", kind: KindImag},
			vy:             &Visitor{value: "(3+1i)", kind: KindImag},
			ops:            []token.Token{token.ADD, token.SUB, token.MUL, token.QUO, token.REM},
			expectedValues: []string{"13", "7", "30", "3", "1"},
			expectedErrs:   []error{nil, nil, nil, nil, nil},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range tc.ops {
				i, op := i, op
				name := tc.vx.value + op.String() + tc.vy.value
				t.Run(name, func(t *testing.T) {
					be := &ast.BinaryExpr{Op: op}
					calculateInt(tc.v, tc.vx, tc.vy, be)
					if !errors.Is(tc.v.err, tc.expectedErrs[i]) {
						t.Fatalf("expected err: %s, got: %s", tc.expectedErrs[i], tc.v.err)
					}
					if tc.v.value != tc.expectedValues[i] {
						t.Fatalf("expected value: %s, got: %s", tc.expectedValues[i], tc.v.value)
					}
				})
			}
		})
	}
}
