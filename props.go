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

func DivideIntoGroups[D any](d []D, perGroupNum int) [][]D {
	n := len(d) / perGroupNum
	groups := make([][]D, n)
	for i := 0; i < n; i++ {
		groups[i] = d[i*perGroupNum : (i+1)*perGroupNum]
	}
	if len(d)%perGroupNum != 0 {
		groups = append(groups, d[n*perGroupNum:])
	}
	return groups
}

func DivideMapIntoGroups[K comparable, D any](d map[K]D, perGroupNum int) []map[K]D {
	n := len(d) / perGroupNum
	if len(d)%perGroupNum != 0 {
		n++
	}
	groups := make([]map[K]D, n)
	m := make(map[K]D)
	var i int
	for k, v := range d {
		i++
		if i%perGroupNum == 0 {
			groups = append(groups, m)
			m = make(map[K]D)
		}
		m[k] = v
		if i == len(d) {
			groups = append(groups, m)
		}
	}
	return groups
}
