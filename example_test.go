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
		"15 % 4",
		"1 + 1 + (4 == 2)",
		"(1 * 2))",
	}

	for _, value := range values {
		v, err := expr.Int(value)
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
	// 0 unsupported operator
	// 0 1:8: expected 'EOF', found ')'

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
	// 0 unsupported operator
	// 0 unsupported operator
	// 0 1:7: expected ')', found newline
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
	// false unsupported operator
	// false unsupported operator
	// false 1:8: expected operand, found 'EOF'
}
