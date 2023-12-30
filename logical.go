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
	"fmt"
	"go/ast"
	"go/token"

	"github.com/muktihari/expr/internal/conv"
)

func logical(v, vx, vy *Visitor, binaryExpr *ast.BinaryExpr) {
	if vx.kind != KindBoolean {
		v.err = newLogicalNonBooleanError(vx, binaryExpr.X)
		return
	}

	if vy.kind != KindBoolean {
		v.err = newLogicalNonBooleanError(vy, binaryExpr.Y)
		return
	}

	x := vx.value.(bool)
	y := vy.value.(bool)

	v.kind = KindBoolean
	if binaryExpr.Op == token.LAND {
		v.value = x && y
		return
	}

	v.value = x || y // token.LOR
}

func newLogicalNonBooleanError(v *Visitor, e ast.Expr) error {
	s := conv.FormatExpr(e)
	return &SyntaxError{
		Msg: "result of \"" + s + "\" is \"" + fmt.Sprintf("%v", v.value) + "\" which is not a boolean",
		Pos: v.pos,
		Err: ErrLogicalOperation,
	}
}
