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

package bind

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrKeyValsLengthIsOdd occurs when given keyvals is not matched [key, val] structure.
	ErrKeyValsLengthIsOdd = errors.New("keyvals's length is odd")
	// ErrKeyvalsIsEmptyOrNil occurs when keyvals is empty or nil
	ErrKeyvalsIsEmptyOrNil = errors.New("keyvals is empty or nil")
	// ErrKeyIsNotAString occurs when given key in keyvals contains non-string type.
	ErrKeyIsNotAString = errors.New("key in keyvals is not a string")
	// ErrMalformedVariablePattern occurs when s is malformed or is not valid
	ErrMalformedVariablePattern = errors.New("malformed variable pattern")
	// ErrEmptyPrefix occurs when prefix is empty "" while it's a mandatory to bind the variables.
	ErrEmptyPrefix = errors.New("empty prefix")
)

var std = &Binder{Ident: DefaultIdent(), Formatter: DefaultFormater()}

// Bind binds given keyvals values into the given s. Key in keyvals should be a string that consist of alphanumeric [a-z, A-Z, 0-9] and symbol ['-', '_'] only.
//
// - e.g. price after discount calculation expression:
//
//   - s: "{price} - ({price} * {discount-percentage})"
//
//   - keyvals: ["price", 100, "discount-percentage", 0.1]
//
//   - resulting value: "100 - (100 * 0.1)"
//
// Note: If s is a really big string (len(s) > 60k for example) consider creating your own binder using strings.Replacer, see [bind_benchmark_test.go] file.
//
// Otherwise, use this for faster process with low memory footprint and low memory alloc.
func Bind(s string, keyvals ...interface{}) (string, error) {
	return std.Bind(s, keyvals...)
}

// SetIdent sets custom variable identifier to std. See bind.Ident{} for details.
func SetIdent(ident *Ident) {
	if ident != nil {
		std.Ident = ident
	}
}

// SetIdent sets custom keyvals formatter to std. See bind.Formatter for details.
func SetFormatter(formatter Formatter) {
	if formatter != nil {
		std.Formatter = formatter
	}
}

// Ident is variable name identifier. Prefix is mandatory when Suffix is optional.
//
// e.g.
//   - "{price}" : the "{" is the prefix identifier and "}" is the suffix identifier of variable named price.
//   - ":price:" : the ":" is the prefix identifier and ":" is the suffix identifier of variable named price.
//   - ":price" : the ":" is the prefix identifier and "" is the suffix identifier of variable named price.
type Ident struct {
	Prefix string // Prefix is mandatory
	Suffix string // Suffix is optional
}

func DefaultIdent() *Ident {
	return &Ident{
		Prefix: "{",
		Suffix: "}",
	}
}

// Formatter formats keyvals values into string values. Key will never be quoted, only the Value will be quoted.
//
// e.g.
//   - "price" -> "price"
//   - 100 -> "100"
//   - 2.1 -> "2.1"
//   - struct{}{} -> "{}"
//   - nil -> "<nil>"
type Formatter func(v interface{}) string

// DefaultFormater returns format
func DefaultFormater() Formatter { return Format }

// Binder binds variable values into string expression, it finds the variable name using specified identifier bind.Ident{}.
type Binder struct {
	Ident     *Ident    // variable identifier on string expression
	Formatter Formatter // keyvals values formatter.
}

type SyntaxError struct {
	Msg   string
	Begin int
	End   int
	Value string
	Err   error
}

func (s *SyntaxError) Error() string {
	return fmt.Sprintf("%s [value:\"%s\",beg:%d,end:%d]: %v", s.Msg, s.Value, s.Begin, s.End, s.Err)
}

func (s *SyntaxError) Unwrap() error { return s.Err }

// Bind binds keyvals values into s, key should be a string and val can be any. If keyvals is nil, s will be returned.
func (b *Binder) Bind(s string, keyvals ...interface{}) (string, error) {
	if len(keyvals) == 0 {
		return "", ErrKeyvalsIsEmptyOrNil
	}

	if len(keyvals)%2 != 0 {
		return "", ErrKeyValsLengthIsOdd
	}

	if b.Ident == nil {
		b.Ident = DefaultIdent()
	}

	if b.Ident.Prefix == "" {
		return "", ErrEmptyPrefix
	}

	if b.Formatter == nil {
		b.Formatter = DefaultFormater()
	}

	prefix, suffix := b.Ident.Prefix, b.Ident.Suffix
	lenPrefix, lenSuffix := len(prefix), len(suffix)

	m := make(map[string]string)
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return "", fmt.Errorf("key '%v' is not a string, err: %w", key, ErrKeyIsNotAString)
		}
		m[key] = b.Formatter(keyvals[i+1])
	}

	var isPrefixBegin, isBreakBySuffix bool
	var begin, end int

	var strbuf strings.Builder
	var cur int
	for i := 0; i < len(s); i++ {
		if !isPrefixBegin {
			if i+lenPrefix < len(s) && s[i:i+lenPrefix] == prefix { // find beginning of a prefix
				isPrefixBegin = true
				isBreakBySuffix = false
				begin = i
				i += lenPrefix - 1
			}
			continue
		}

		if lenSuffix != 0 && i+lenSuffix <= len(s) { // check breaking point by a proper suffix if specified
			if s[i:i+lenSuffix] == suffix {
				end = i + lenSuffix
				i += lenSuffix - 1

				strbuf.WriteString(s[cur:begin])
				strbuf.WriteString(m[s[begin+lenPrefix:end-lenSuffix]])
				cur = end

				isPrefixBegin = false
				isBreakBySuffix = true
				continue
			}
		}

		// check breaking point
		r := rune(s[i])
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-') {
			end = i
			strbuf.WriteString(s[cur:begin])
			strbuf.WriteString(m[s[begin+lenPrefix:end-lenSuffix]])
			cur = end

			isPrefixBegin = false
			isBreakBySuffix = false

			if lenSuffix != 0 { // not broken by suffix when it should
				return "", &SyntaxError{
					Msg:   "suffix is specified but it is broken by '" + string(r) + "' before reaching suffix",
					Begin: begin,
					End:   end,
					Value: s[begin:end],
					Err:   ErrMalformedVariablePattern,
				}
			}
		}
	}

	if isPrefixBegin {
		if lenSuffix != 0 && !isBreakBySuffix {
			return "", &SyntaxError{
				Msg:   "suffix is specified but missing suffix at the end of s when it should be ended by a proper suffix",
				Begin: begin,
				End:   len(s),
				Value: s[begin:],
				Err:   ErrMalformedVariablePattern,
			}
		}
	}

	strbuf.WriteString(s[cur:])

	return strbuf.String(), nil
}

// Format formats given v type into string.
func Format(v interface{}) string {
	// declared common used types for faster conversion
	switch val := v.(type) {
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(val, 'f', -1, 128)
	case string:
		return strconv.Quote(val)
	case bool:
		return strconv.FormatBool(val)
	case error:
		return strconv.Quote(val.Error())
	case fmt.Stringer:
		return strconv.Quote(val.String())
	default: // slower but it can handle "{}" "[1, 2]" "<nil>", etc.
		s := fmt.Sprintf("%v", v) // e.g. int32(2) -> 2
		if idx := strings.IndexFunc(s, func(r rune) bool {
			return r == '[' || r == ']' || r == '{' || r == '}' || r == '<' || r == '>'
		}); idx != -1 {
			return strconv.Quote(s)
		}
		return s
	}
}
