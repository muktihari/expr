package bind

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestStdBinder(t *testing.T) {
	tt := []struct {
		in      string
		keyvals []interface{}
		out     string
		err     error
	}{
		{
			in: "{price} - ({price} * {discount-percentage})",
			keyvals: []interface{}{
				"price", 100,
				"discount-percentage", 0.1,
			},
			out: "100 - (100 * 0.1)",
		},
		{
			in: "{is-eligible} && {age} >= 60",
			keyvals: []interface{}{
				"is-eligible", true,
				"age", 28,
			},
			out: "true && 28 >= 60",
		},
		{
			in: "{is-eligible} && {age} >= 60",
			keyvals: []interface{}{
				"is-eligible", true,
				28, 28,
			},
			err: ErrKeyIsNotAString,
		},
		{
			in: "{is-eligible} && {age} >= 60",
			keyvals: []interface{}{
				"is-eligible", true,
				"age",
			},
			err: ErrKeyValsLengthIsOdd,
		},
		{
			in:      "{is-eligible} && {age} >= 60",
			keyvals: nil,
			err:     ErrKeyvalsIsEmptyOrNil,
		},
		{
			in: "{is-eligible} == \"<nil>\"",
			keyvals: []interface{}{
				"is-eligible", nil,
			},
			out: "\"<nil>\" == \"<nil>\"",
		},
		{
			in: "{is-eligible} && {age } >= 60",
			keyvals: []interface{}{
				"is-eligible", true,
				"age", 28,
			},
			err: ErrMalformedVariablePattern,
		},
		{
			in: "{is-eligible} && {age} >= {old",
			keyvals: []interface{}{
				"is-eligible", true,
				"age", 28,
				"old", 28,
			},
			err: ErrMalformedVariablePattern,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			out, err := Bind(tc.in, tc.keyvals...)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected error: %v, got: %v", tc.err, err)
			}
			if out != tc.out {
				t.Fatalf("expected out: %v, got: %v", tc.out, out)
			}
		})
	}
}

func TestCustomBinder(t *testing.T) {
	tt := []struct {
		in      string
		keyvals []interface{}
		binder  *Binder
		out     string
		err     error
	}{
		{
			in: ":price: - (:price: * :discount-percentage:)",
			keyvals: []interface{}{
				"price", 100,
				"discount-percentage", 0.1,
			},
			binder: &Binder{Ident: &Ident{
				Prefix: ":", Suffix: ":",
			}},
			out: "100 - (100 * 0.1)",
		},
		{
			in: ":price - (:price * :discount-percentage)",
			keyvals: []interface{}{
				"price", 100,
				"discount-percentage", 0.1,
			},
			binder: &Binder{Ident: &Ident{
				Prefix: ":", Suffix: "",
			}},
			out: "100 - (100 * 0.1)",
		},
		{
			in: "{price} - ({price} * {discount-percentage})",
			keyvals: []interface{}{
				"price", 100,
				"discount-percentage", 0.1,
			},
			binder: &Binder{},
			out:    "100 - (100 * 0.1)",
		},
		{
			in: "{price} - ({price} * {discount-percentage})",
			keyvals: []interface{}{
				"price", 100,
				"discount-percentage", 0.1,
			},
			binder: &Binder{Ident: &Ident{Prefix: ""}},
			err:    ErrEmptyPrefix,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			out, err := tc.binder.Bind(tc.in, tc.keyvals...)
			if !errors.Is(err, tc.err) {
				t.Fatalf("expected error: %v, got: %v", tc.err, err)
			}
			if out != tc.out {
				t.Fatalf("expected out: %v, got: %v", tc.out, out)
			}
		})
	}
}

func TestSetIdent(t *testing.T) {
	ident := &Ident{Prefix: ":", Suffix: ""}
	stdIdent := std.Ident
	SetIdent(ident)
	if std.Ident != ident {
		t.Fatalf("expected: %v, got: %v", ident, std.Ident)
	}
	std.Ident = stdIdent
}

func TestSetFormatter(t *testing.T) {
	var formatter Formatter = func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}

	stdFormatter := std.Formatter
	SetFormatter(formatter)
	if fmt.Sprintf("%p", std.Formatter) != fmt.Sprintf("%p", formatter) {
		t.Fatalf("expected: %v, got: %v", formatter, std.Formatter)
	}
	std.Formatter = stdFormatter
}

func TestFormat(t *testing.T) {
	type TestStruct struct{ Field string }
	var emptyErr error

	tt := []struct {
		in  interface{}
		out string
	}{
		{in: 12, out: "12"},
		{in: int64(12), out: "12"},
		{in: int32(12), out: "12"},
		{in: float64(12.3), out: "12.3"},
		{in: complex(2, -9), out: "(2-9i)"},
		{in: true, out: "true"},
		{in: "expr", out: "\"expr\""},
		{in: time.Time{}, out: "\"0001-01-01 00:00:00 +0000 UTC\""}, // test fmt.Stringer
		{in: "expr", out: "\"expr\""},
		{in: struct{}{}, out: "\"{}\""},
		{in: emptyErr, out: "\"<nil>\""},
		{in: fmt.Errorf("error something"), out: "\"error something\""},
		{in: []byte("c"), out: "\"[99]\""},
		{in: []string{"abc", "def"}, out: "\"[abc def]\""},
		{in: TestStruct{}, out: "\"{}\""},
		{in: TestStruct{Field: "value"}, out: "\"{value}\""},
		{in: &TestStruct{}, out: "\"&{}\""},
		{in: &TestStruct{Field: "value"}, out: "\"&{value}\""},
	}

	defaultFormat := DefaultFormater()
	for _, tc := range tt {
		tc := tc
		t.Run(fmt.Sprintf("%T: %v", tc.in, tc.in), func(t *testing.T) {
			s := defaultFormat(tc.in)
			if s != tc.out {
				t.Fatalf("expected string value: %s, got: %s", tc.out, s)
			}
		})
	}
}

func TestSyntaxError(t *testing.T) {
	err := &SyntaxError{
		Msg:   "test",
		Begin: 0,
		End:   5,
		Value: "value",
		Err:   ErrEmptyPrefix,
	}
	expectedErrorString := "test [value:\"value\",beg:0,end:5]: " + ErrEmptyPrefix.Error()
	if err.Error() != expectedErrorString {
		t.Fatalf("expected err string: %s, got: %s", expectedErrorString, err.Error())
	}
}
