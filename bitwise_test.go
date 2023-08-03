package expr

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"testing"
)

func TestBitwise(t *testing.T) {
	tt := []struct {
		v, vx, vy     *Visitor
		op            token.Token
		expectedValue string
		expectedErr   error
	}{
		{
			v:           &Visitor{options: options{numericType: NumericTypeFloat}},
			vx:          &Visitor{value: "0b1000", kind: KindInt},
			vy:          &Visitor{value: "0b1001", kind: KindInt},
			op:          token.AND, // "&"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: "2.0", kind: KindFloat},
			vy:          &Visitor{value: "0b1001", kind: KindInt},
			op:          token.AND, // "&"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: "0b1001", kind: KindInt},
			vy:          &Visitor{value: "2.0", kind: KindFloat},
			op:          token.AND, // "&"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: "0b1000", kind: KindInt},
			vy:            &Visitor{value: "0b1001", kind: KindInt},
			op:            token.AND, // "&"
			expectedValue: "0b1000",
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: "0b1000", kind: KindInt},
			vy:            &Visitor{value: "0b0001", kind: KindInt},
			op:            token.OR, // "|"
			expectedValue: "0b1001",
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: "0b1000", kind: KindInt},
			vy:            &Visitor{value: "0b1001", kind: KindInt},
			op:            token.XOR, // "^"
			expectedValue: "0b0001",
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: "0b1100", kind: KindInt},
			vy:            &Visitor{value: "0b0101", kind: KindInt},
			op:            token.AND_NOT, // "&^"
			expectedValue: "0b1000",
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: "0b1001", kind: KindInt},
			vy:            &Visitor{value: "0b0001", kind: KindInt},
			op:            token.SHL, // "<<"
			expectedValue: "0b10010",
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: "0b1000", kind: KindInt},
			vy:            &Visitor{value: "0b0001", kind: KindInt},
			op:            token.SHR, // ">>"
			expectedValue: "0b0100",
		},
	}

	for _, tc := range tt {
		tc := tc
		name := fmt.Sprintf("%s %s %s", tc.vx.value, tc.op, tc.vy.value)
		t.Run(name, func(t *testing.T) {
			opPos := token.Pos(len(tc.vx.value) + 1)
			be := &ast.BinaryExpr{
				X:     &ast.BasicLit{Value: tc.vx.value, ValuePos: 1},
				Op:    tc.op,
				OpPos: opPos,
				Y:     &ast.BasicLit{Value: tc.vy.value, ValuePos: opPos + 1},
			}

			bitwise(tc.v, tc.vx, tc.vy, be)
			if !errors.Is(tc.v.err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, tc.v.err)
			}

			expectedValue, _ := strconv.ParseInt(tc.expectedValue, 0, 64)
			value, _ := strconv.ParseInt(tc.v.value, 0, 64)
			if value != expectedValue {
				t.Fatalf("expected value: %s, got: %s", tc.expectedValue, tc.v.value)
			}
		})
	}
}
