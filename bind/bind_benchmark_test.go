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

package bind_test

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/muktihari/expr/bind"
)

var (
	benchDefaultIdent  = bind.DefaultIdent()
	benchDefaultFormat = bind.DefaultFormater()
)

// benchBindWithStringsReplacer is the fastest when len(s) > 6k chars, but not that significant in <=600. This implementation
// only replace the exact pattern, when there is small typo like added " " space between prefix and suffix, it will not being replaced.
// example: should be {price}, but written {price }. while default Bind can return an error.
func benchBindWithStringsReplacer(s string, keyvals ...interface{}) (string, error) {
	if len(keyvals) == 0 {
		return s, nil
	}

	if len(keyvals)%2 != 0 {
		return "", bind.ErrKeyValsLengthIsOdd
	}

	strkeyvals := make([]string, 0, len(keyvals))
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return "", fmt.Errorf("key '%v' is not a string, err: %w", key, bind.ErrKeyIsNotAString)
		}
		val := benchDefaultFormat(keyvals[i+1])
		key = benchDefaultIdent.Prefix + key + benchDefaultIdent.Suffix
		strkeyvals = append(strkeyvals, key, val)
	}

	return strings.NewReplacer(strkeyvals...).Replace(s), nil
}

// benchBindWithStringsReplaceAll is the fastest for extra small len(s) (tested: len = 76), but getting slower when len(s) is increasing.
func benchBindWithStringsReplaceAll(s string, keyvals ...interface{}) (string, error) {
	if len(keyvals) == 0 {
		return s, nil
	}

	if len(keyvals)%2 != 0 {
		return "", bind.ErrKeyValsLengthIsOdd
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return "", fmt.Errorf("key '%v' is not a string, err: %w", key, bind.ErrKeyIsNotAString)
		}
		key, val := keyvals[i].(string), benchDefaultFormat(keyvals[i+1])
		key = benchDefaultIdent.Prefix + key + benchDefaultIdent.Suffix
		s = strings.ReplaceAll(s, key, val)
	}

	return s, nil
}

var reg = regexp.MustCompile(`{(.*?)}`) // compiling is expensive, so compile at build time

// benchBindWithRegexp is the slowest most of the time, but when s is big, it's the second slowest after strings.ReplaceAll. Used for benchmarking.
func benchBindWithRegexp(s string, keyvals ...interface{}) (string, error) {
	if len(keyvals) == 0 {
		return s, nil
	}

	if len(keyvals)%2 != 0 {
		return "", bind.ErrKeyValsLengthIsOdd
	}

	m := make(map[string]string, len(keyvals)/2)

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return "", fmt.Errorf("key '%v' is not a string, err: %w", key, bind.ErrKeyIsNotAString)
		}
		val := benchDefaultFormat(keyvals[i+1])
		m[benchDefaultIdent.Prefix+key+benchDefaultIdent.Suffix] = val
	}

	s = reg.ReplaceAllStringFunc(s, func(s string) string { return m[s] })

	return s, nil
}

func TestBindAllBenchmarkBindAlternatives(t *testing.T) {
	t.SkipNow() // only comment when you need to test the bind alternatives.

	multiplier := 1

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
		func() struct {
			in      string
			keyvals []interface{}
			out     string
			err     error
		} {
			s := "{price} - ({price} * {discount-percentage})"
			for i := 0; i < multiplier; i++ {
				s += fmt.Sprintf(" + {Variable%d}", i)
			}
			r := "100 - (100 * 0.1)"
			for i := 0; i < multiplier; i++ {
				r += " + 100.213"
			}

			keyvals := make([]interface{}, 0, multiplier)
			keyvals = append(keyvals,
				"price", 100,
				"discount-percentage", 0.1,
			)
			for i := 0; i < multiplier; i++ {
				keyvals = append(keyvals,
					fmt.Sprintf("Variable%d", i), 100.213,
				)
			}

			return struct {
				in      string
				keyvals []interface{}
				out     string
				err     error
			}{
				in: s, keyvals: keyvals, out: r, err: nil,
			}
		}(),
	}

	for _, tc := range tt {
		tc := tc
		t.Run("last 10 chars: "+string(tc.in[len(tc.in)-10:]), func(t *testing.T) {
			t.Run("benchBindWithStringsReplaceAll", func(t *testing.T) {
				out, err := benchBindWithStringsReplaceAll(tc.in, tc.keyvals...)
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				if out != tc.out {
					t.Fatalf("expected out: %v, got: %v", tc.out, out)
				}
			})
			t.Run("benchBindWithRegexp", func(t *testing.T) {
				out, err := benchBindWithRegexp(tc.in, tc.keyvals...)
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				if out != tc.out {
					t.Fatalf("expected out: %v, got: %v", tc.out, out)
				}
			})
			t.Run("benchBindWithStringsReplacer", func(t *testing.T) {
				out, err := benchBindWithStringsReplacer(tc.in, tc.keyvals...)
				if !errors.Is(err, tc.err) {
					t.Fatalf("expected error: %v, got: %v", tc.err, err)
				}
				if out != tc.out {
					t.Fatalf("expected out: %v, got: %v", tc.out, out)
				}
			})
		})
	}
}

func multiplyStringTemplate(multiplier int) (str string, strResult string, varInStrCount int, keyvals []interface{}) {
	templateStr := "{price} - ({price} * {discount-percentage})"
	templateResult := "100 - (100 * 0.1)"

	keyvals = append(keyvals,
		"price", 100,
		"discount-percentage", 0.1,
	)

	str, varInStrCount = templateStr, 3
	for i := 0; i < multiplier; i++ {
		str += fmt.Sprintf(" + %s + {Variable%d}", templateStr, i)
		keyvals = append(keyvals,
			fmt.Sprintf("Variable%d", i), 100.213,
		)
		varInStrCount += 4
	}

	strResult = templateResult
	for i := 0; i < multiplier; i++ {
		strResult += fmt.Sprintf(" + %s + 100.213", templateResult)
	}

	return str, strResult, varInStrCount, keyvals
}

func BenchmarkBind(b *testing.B) {
	type strct struct {
		scenario      string
		str           string
		strResult     string
		varInStrCount int
		keyvals       []interface{}
	}

	tt := []strct{
		{
			scenario:  "extra small",
			str:       "(({price} * (1 - {discount-percentage})) - {another-discount}) * (1 + {tax})",
			strResult: "((100 * (1 - 0.1)) - 1) * (1 + 0.2)",
			keyvals: []interface{}{
				"price", 100,
				"discount-percentage", 0.1,
				"another-discount", 1,
				"tax", 0.2,
			},
			varInStrCount: 4,
		},
		func() strct {
			multiplier := 10

			str, strResult, varInStrCount, keyvals := multiplyStringTemplate(multiplier)

			return strct{
				scenario:      "small",
				str:           str,
				strResult:     strResult,
				varInStrCount: varInStrCount,
				keyvals:       keyvals,
			}
		}(),
		func() strct {
			multiplier := 100

			str, strResult, varInStrCount, keyvals := multiplyStringTemplate(multiplier)

			return strct{
				scenario:      "medium",
				str:           str,
				strResult:     strResult,
				varInStrCount: varInStrCount,
				keyvals:       keyvals,
			}
		}(),
		func() strct {
			multiplier := 1000

			str, strResult, varInStrCount, keyvals := multiplyStringTemplate(multiplier)

			return strct{
				scenario:      "large",
				str:           str,
				strResult:     strResult,
				varInStrCount: varInStrCount,
				keyvals:       keyvals,
			}
		}(),
	}

	for _, tc := range tt {
		tc := tc
		b.Run(tc.scenario+fmt.Sprintf(" len(s):%d, len(varInStrCount): %d, vars:%d", len(tc.str), tc.varInStrCount, len(tc.keyvals)/2), func(b *testing.B) {
			b.Run("Bind", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v, err := bind.Bind(tc.str, tc.keyvals...)
					if err != nil {
						b.Fatalf("expected nil, got: %v", err)
					}
					if v != tc.strResult {
						b.Fatalf("expected value: %s, got: %s", tc.strResult, v)
					}
				}
			})
			b.Run("bindWithStringsReplaceAll", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v, err := benchBindWithStringsReplaceAll(tc.str, tc.keyvals...)
					if err != nil {
						b.Fatalf("expected nil, got: %v", err)
					}
					if v != tc.strResult {
						b.Fatalf("expected value: %s, got: %s", tc.strResult, v)
					}
				}
			})
			b.Run("bindWithRegex", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v, err := benchBindWithRegexp(tc.str, tc.keyvals...)
					if err != nil {
						b.Fatalf("expected nil, got: %v", err)
					}
					if v != tc.strResult {
						b.Fatalf("expected value: %s, got: %s", tc.strResult, v)
					}
				}
			})
			b.Run("benchBindWithStringsReplacer", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v, err := benchBindWithStringsReplacer(tc.str, tc.keyvals...)
					if err != nil {
						b.Fatalf("expected nil, got: %v", err)
					}
					if v != tc.strResult {
						b.Fatalf("expected value: %s, got: %s", tc.strResult, v)
					}
				}
			})
		})
	}
}
