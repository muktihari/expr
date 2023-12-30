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
	"fmt"

	"github.com/muktihari/expr"
)

func ExampleAny() {
	values := []string{
		"1 + 2 + 3 + 4 + 5",
		"(2 + 2) * 4 / 4",
		"((2 + 2) * 4 / 4) * 10.5 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)",
		"(10+5i) + (10+7i)",
		"(2+3i) - (2+2i)",
		"(2+2i) * (2+2i)",
		"(2+2i) / (2+2i)",
		"4 % 2",
		"1 + 1 + (4 == 2)",
		"(1 + 1",
	}

	for _, value := range values {
		v, err := expr.Any(value)
		fmt.Println(v, err)
	}

	// Output:
	// 15 <nil>
	// 4 <nil>
	// 46.234567 <nil>
	// 466.2567 <nil>
	// (20+12i) <nil>
	// (0+1i) <nil>
	// (0+8i) <nil>
	// (1+0i) <nil>
	// 0 <nil>
	// <nil> result of "(4 == 2)" is "false" which is not a number [pos: 10]: arithmetic operation
	// <nil> 1:7: expected ')', found newline
}

func ExampleBool() {
	values := []string{
		"1 < 2 && 3 < 4 && ( 1==1 || 12 > 4)",
		"\"expr\" == \"expr\" && \"Expr\" != \"expr\"",
		"\"expr\" == \"expr\" && \"Expr\" == \"expr\"",
		"(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || 1 ==1 ",
		"(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || true == true ",
		"(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || true == false ",
		"1 + 2 > 1 * 2",
		"1 + 2 < (2 + 2) * 10",
		"1 < 1 <",
		"-1 > -10",
		"true",
		"!false",
		"!false || false",
	}

	for _, value := range values {
		v, err := expr.Bool(value)
		fmt.Println(v, err)
	}

	// Output:
	// true <nil>
	// true <nil>
	// false <nil>
	// true <nil>
	// true <nil>
	// false <nil>
	// true <nil>
	// true <nil>
	// false 1:8: expected operand, found 'EOF'
	// true <nil>
	// true <nil>
	// true <nil>
	// true <nil>
}

func ExampleComplex128() {
	values := []string{
		"1 + 2 + 3 + 4 + 5",
		"(2 + 2) * 4 / 4",
		"((2 + 2) * 4 / 4) * 10.5 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)",
		"(10+5i) + (10+7i)",
		"(2+3i) - (2+2i)",
		"(2+2i) * (2+2i)",
		"(2+2i) / (2+2i)",
		"4 % 2",
		"1 + 1 + (4 == 2)",
		"(1 + 1",
	}

	for _, value := range values {
		v, err := expr.Complex128(value)
		fmt.Println(v, err)
	}

	// Output:
	// (15+0i) <nil>
	// (4+0i) <nil>
	// (46.234567+0i) <nil>
	// (466.2567+0i) <nil>
	// (20+12i) <nil>
	// (0+1i) <nil>
	// (0+8i) <nil>
	// (1+0i) <nil>
	// (0+0i) operator "%" is not supported to do arithmetic on complex number [pos: 3]: arithmetic operation
	// (0+0i) result of "(4 == 2)" is "false" which is not a number [pos: 10]: arithmetic operation
	// (0+0i) 1:7: expected ')', found newline
}

func ExampleFloat64() {
	values := []string{
		"1 + 2 + 3 + 4 + 5",
		"(2 + 2) * 4 / 4",
		"((2 + 2) * 4 / 4) * 10",
		"((2 + 2) * 4 / 4) * 10 + 2",
		"((2 + 2) * 4 / 4) * 10 + 2 - 2",
		"((2 + 2) * 4 / 4) * 10 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.5 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)",
		"4 % 2",
		"1 + 1 + (4 == 2)",
		"(1 + 1",
	}

	for _, value := range values {
		v, err := expr.Float64(value)
		fmt.Println(v, err)
	}

	// Output:
	// 15 <nil>
	// 4 <nil>
	// 40 <nil>
	// 42 <nil>
	// 40 <nil>
	// 44.234567 <nil>
	// 46.234567 <nil>
	// 466.2567 <nil>
	// 0 <nil>
	// 0 result of "(4 == 2)" is "false" which is not a number [pos: 10]: arithmetic operation
	// 0 1:7: expected ')', found newline
}

func ExampleInt() {
	values := []string{
		"1 + 2 + 3 + 4 + 5",
		"(2 + 2) * 4 / 4",
		"(2 + 2) - 2 * 2",
	}

	for _, value := range values {
		v, err := expr.Int64(value)
		fmt.Println(v, err)
	}

	// Output:
	// 15 <nil>
	// 4 <nil>
	// 0 <nil>
}

func ExampleInt64() {
	values := []string{
		"1 + 2 + 3 + 4 + 5",
		"(2 + 2) * 4 / 4",
		"((2 + 2) * 4 / 4) * 10",
		"((2 + 2) * 4 / 4) * 10 + 2",
		"((2 + 2) * 4 / 4) * 10 + 2 - 2",
		"((2 + 2) * 4 / 4) * 10 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.5 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)",
		"15 % 4",
		"1 + 1 + (4 == 2)",
		"(1 * 2))",
	}

	for _, value := range values {
		v, err := expr.Int64(value)
		fmt.Println(v, err)
	}

	// Output:
	// 15 <nil>
	// 4 <nil>
	// 40 <nil>
	// 42 <nil>
	// 40 <nil>
	// 44 <nil>
	// 44 <nil>
	// 440 <nil>
	// 3 <nil>
	// 0 result of "(4 == 2)" is "false" which is not a number [pos: 10]: arithmetic operation
	// 0 1:8: expected 'EOF', found ')'
}

func ExampleInt64Strict() {
	values := []string{
		"(2 + 2) * 4 / 4",
		"(2 + 2) - 2 / 0",
	}

	for _, value := range values {
		v, err := expr.Int64Strict(value)
		fmt.Println(v, err)
	}

	// Output:
	// 4 <nil>
	// 0 could not divide x with zero y, allowIntegerDividedByZero == false [pos: 15]: integer divided by zero
}
