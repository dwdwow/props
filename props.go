package props

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// IsNil judge that a value is nil or not.
// If v is value type, return false.
// If v is pointer unless interface, check with reflecting.
// If v is interface, return if v's underlying value is nil or not.
func IsNil(v any) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Interface,
		reflect.Pointer,
		reflect.Slice,
		reflect.Map,
		reflect.Func,
		reflect.Chan,
		reflect.UnsafePointer:
		return val.IsNil()
	case reflect.Invalid:
		return true
	default:
		return false
	}
}

// PanicIfNotNil use IsNil as nil checker,
// in case that the input is a value
// whose underlying value is nil, but underlying
// type is not nil.
func PanicIfNotNil(err error) {
	if !IsNil(err) {
		panic(err)
	}
}

func PrintlnIndent(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	var s string
	if err != nil {
		s = err.Error()
	} else {
		s = string(data)
	}
	fmt.Println(s)
}
