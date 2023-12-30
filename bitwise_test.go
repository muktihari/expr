// Copyright 2023 The Expr Authors
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

package expr

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"testing"
)

func TestBitwise(t *testing.T) {
	tt := []struct {
		v, vx, vy     *Visitor
		op            token.Token
		expectedValue int64
		expectedErr   error
	}{
		{
			v:           &Visitor{options: options{numericType: NumericTypeFloat}},
			vx:          &Visitor{value: int64(0b1000), kind: KindInt},
			vy:          &Visitor{value: int64(0b1001), kind: KindInt},
			op:          token.AND, // "&"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: float64(2.0), kind: KindFloat},
			vy:          &Visitor{value: int64(0b1001), kind: KindInt},
			op:          token.AND, // "&"
			expectedErr: nil,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: float64(2.2), kind: KindFloat},
			vy:          &Visitor{value: int64(0b1001), kind: KindInt},
			op:          token.AND, // "&"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: int64(0b1001), kind: KindInt},
			vy:          &Visitor{value: float64(2.0), kind: KindFloat},
			op:          token.AND, // "&"
			expectedErr: nil,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: int64(0b1001), kind: KindInt},
			vy:          &Visitor{value: float64(2.2), kind: KindFloat},
			op:          token.AND, // "&"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: int64(0b1000), kind: KindInt},
			vy:            &Visitor{value: int64(0b1001), kind: KindInt},
			op:            token.AND, // "&"
			expectedValue: int64(0b1000),
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: int64(0b1000), kind: KindInt},
			vy:            &Visitor{value: int64(0b0001), kind: KindInt},
			op:            token.OR, // "|"
			expectedValue: int64(0b1001),
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: int64(0b1000), kind: KindInt},
			vy:            &Visitor{value: int64(0b1001), kind: KindInt},
			op:            token.XOR, // "^"
			expectedValue: int64(0b0001),
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: int64(0b1100), kind: KindInt},
			vy:            &Visitor{value: int64(0b0101), kind: KindInt},
			op:            token.AND_NOT, // "&^"
			expectedValue: int64(0b1000),
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: int64(0b1001), kind: KindInt},
			vy:            &Visitor{value: int64(0b0001), kind: KindInt},
			op:            token.SHL, // "<<"
			expectedValue: int64(0b10010),
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: int64(0b1000), kind: KindInt},
			vy:            &Visitor{value: int64(0b0001), kind: KindInt},
			op:            token.SHR, // ">>"
			expectedValue: int64(0b0100),
		},
		{
			v:             &Visitor{},
			vx:            &Visitor{value: float64(4.0), kind: KindFloat},
			vy:            &Visitor{value: int64(10), kind: KindInt},
			op:            token.SHL, // "<<"
			expectedValue: int64(0b1000000000000),
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: true, kind: KindBoolean},
			vy:          &Visitor{value: int64(10), kind: KindInt},
			op:          token.SHL, // "<<"
			expectedErr: ErrBitwiseOperation,
		},
		{
			v:           &Visitor{},
			vx:          &Visitor{value: int64(10), kind: KindInt},
			vy:          &Visitor{value: true, kind: KindBoolean},
			op:          token.SHL, // "<<"
			expectedErr: ErrBitwiseOperation,
		},
	}

	for _, tc := range tt {
		tc := tc
		name := fmt.Sprintf("%v %s %v", tc.vx.value, tc.op, tc.vy.value)
		t.Run(name, func(t *testing.T) {
			opPos := token.Pos(len(fmt.Sprintf("%v", tc.vx.value)) + 1)
			be := &ast.BinaryExpr{
				X:     &ast.BasicLit{Value: fmt.Sprintf("%v", tc.vx.value), ValuePos: 1},
				Op:    tc.op,
				OpPos: opPos,
				Y:     &ast.BasicLit{Value: fmt.Sprintf("%v", tc.vy.value), ValuePos: opPos + 1},
			}

			bitwise(tc.v, tc.vx, tc.vy, be)
			if !errors.Is(tc.v.err, tc.expectedErr) {
				t.Fatalf("expected err: %v, got: %v", tc.expectedErr, tc.v.err)
			}

			value, _ := tc.v.value.(int64)
			if value != tc.expectedValue {
				t.Fatalf("expected value: %v (%T), got: %v (%T)", tc.expectedValue, tc.expectedValue,
					tc.v.value, tc.v.value)
			}
		})
	}
}
