package expr

import (
	"math"
	"strconv"
)

// Kind of value (value's type)
type Kind byte

const (
	KindIllegal Kind = iota
	KindBoolean      // true false

	// Identifiers of numeric type
	numeric_beg
	KindInt   // 12345
	KindFloat // 123.45
	KindImag  // 123.45i
	numeric_end

	KindString // "abc" 'abc' `abc`
)

var kinds = [...]string{
	KindIllegal: "KindIllegal",
	KindBoolean: "KindBoolean",
	KindInt:     "KindInt",
	KindFloat:   "KindFloat",
	KindImag:    "KindImag",
	KindString:  "KindString",
}

func (k Kind) String() string {
	if k < Kind(len(kinds)) {
		return kinds[k]
	}
	return "kind(" + strconv.Itoa(int(k)) + ")"
}

// value is a custom value to reduce memory allocation,
// so we don't allocate if the value is bool, int64 or float64.
type value struct {
	_   [0]func()   // disallow ==
	num uint64      // storage for bool, int64 or float64 value.
	any interface{} // storage for Kind (only if bool, int64 or float64), complex128 value or string value.
}

// Kind returns value's kind.
func (v *value) Kind() Kind {
	switch k := v.any.(type) {
	case Kind:
		return k
	case complex128:
		return KindImag
	case string:
		return KindString
	}
	return KindIllegal
}

// SetKind sets value's kind.
func (v *value) SetKind(k Kind) { v.any = k }

// Bool returns value as bool.
func (v *value) Bool() bool { return v.num == 1 }

// Int64 returns value as int64.
func (v *value) Int64() int64 { return int64(v.num) }

// Float64 returns value as float64.
func (v *value) Float64() float64 { return math.Float64frombits(v.num) }

// Complex128 returns value as complex128.
func (v *value) Complex128() complex128 {
	val, _ := v.any.(complex128)
	return val
}

// String returns value as string.
func (v *value) String() string {
	s, _ := v.any.(string)
	return s
}

// Any returns underlying value as interface{}.
func (v *value) Any() interface{} {
	switch v.Kind() {
	case KindBoolean:
		return v.num == 1
	case KindInt:
		return int64(v.num)
	case KindFloat:
		return math.Float64frombits(v.num)
	case KindImag:
		v, _ := v.any.(complex128)
		return v
	case KindString:
		s, _ := v.any.(string)
		return s
	}
	return nil
}

// boolValue creates boolean value.
func boolValue(v bool) value {
	var num uint64
	if v {
		num = 1
	}
	return value{num: num, any: KindBoolean}
}

// int64Value creates int64 value.
func int64Value(v int64) value { return value{num: uint64(v), any: KindInt} }

// float64Value creates float64 value.
func float64Value(v float64) value { return value{num: math.Float64bits(v), any: KindFloat} }

// complex128Value creates complex128 value.
func complex128Value(v complex128) value { return value{any: v} }

// stringValue creates string value.
func stringValue(v string) value { return value{any: v} }
