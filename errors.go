package expr

import (
	"errors"
	"fmt"
)

var (
	// ErrUnsupportedOperator is error unsupported operator
	ErrUnsupportedOperator = errors.New("unsupported operator")
	// ErrUnaryOperation occurs when unary operation failed
	ErrUnaryOperation = errors.New("unary operation")
	// ErrArithmeticOperation occurs when either x or y is not int or float
	ErrArithmeticOperation = errors.New("arithmetic operation")
	// ErrIntegerDividedByZero occurs when x/y and y equals to 0 and AllowIntDivByZero == false (default).
	// Go does not allow integer to be divided by zero by default.
	ErrIntegerDividedByZero = errors.New("integer divided by zero")
	// ErrInvalidBitwiseOperation occurs when neither x nor y is an int
	ErrBitwiseOperation = errors.New("bitwise operation")
	// ErrBitwiseOperation occurs when either x or y is boolean and given operator is neither '==' nor '!='
	ErrComparisonOperation = errors.New("comparison operation")
	// ErrLogicalOperation occurs when either x or y is not boolean
	ErrLogicalOperation = errors.New("logical operation")
	// ErrValueTypeMismatch occurs when the result of expr evaluation is not match with desired type
	ErrValueTypeMismatch = errors.New("returned value's type is not match with desired type")
)

// SyntaxError is syntax error
type SyntaxError struct {
	Msg string
	Pos int
	Err error
}

func (e SyntaxError) Error() string { return fmt.Sprintf("%s [pos: %d]: %v", e.Msg, e.Pos, e.Err) }

func (e SyntaxError) Unwrap() error { return e.Err }
