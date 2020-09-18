package expr_test

import (
	"fmt"

	"github.com/muktihari/expr"
)

func ExampleInt() {
	values := []string{
		"1 + 2 + 3 + 4 + 5",
		"(2 + 2) * 4 / 4",
		"((2 + 2) * 4 / 4) * 10",
		"((2 + 2) * 4 / 4) * 10 + 2",
		"((2 + 2) * 4 / 4) * 10 + 2 - 2",
		"((2 + 2) * 4 / 4) * 10 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.5 + 4.234567",
		"((2 + 2) * 4 / 4) * 10.7 + 4.234567 * (50 + 50)",
	}

	for _, value := range values {
		v, _ := expr.Int(value)
		fmt.Println(v)
	}

	// Output:
	// 15
	// 4
	// 40
	// 42
	// 40
	// 44
	// 44
	// 440
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
	}

	for _, value := range values {
		v, _ := expr.Float64(value)
		fmt.Println(v)
	}

	// Output:
	// 15
	// 4
	// 40
	// 42
	// 40
	// 44.234567
	// 46.234567
	// 466.2567
}

func ExampleBool() {
	values := []string{
		"1 < 2 && 3 < 4 && ( 1==1 || 12 > 4)",
		"\"expr\" == \"expr\" && \"Expr\" != \"expr\"",
		"\"expr\" == \"expr\" && \"Expr\" == \"expr\"",
		"(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || 1 ==1 ",
		"(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || true == true ",
		"(\"expr\" == \"expr\" && \"Expr\" == \"expr\") || true == false ",
	}

	for _, value := range values {
		v, _ := expr.Bool(value)
		fmt.Println(v)
	}

	// Output:
	// true
	// true
	// false
	// true
	// true
	// false
}
