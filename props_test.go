package props

import (
	"reflect"
	"testing"
)

type myErr struct {
}

func (e *myErr) Error() string {
	return "myErr"
}

func TestIsNil(t *testing.T) {
	err := new(myErr)
	err = nil
	if error(err) == nil {
		panic(1)
	}
	if !IsNil(err) {
		panic(2)
	}
	var s string
	if IsNil(s) {
		panic(3)
	}
	var ps *string
	if !IsNil(ps) {
		panic(4)
	}
	sErr := myErr{}
	if IsNil(sErr) {
		panic(5)
	}
	pSErr := &sErr
	if IsNil(pSErr) {
		panic(6)
	}
	pSErr = nil
	if !IsNil(pSErr) {
		panic(7)
	}
}

func TestPanicIfNotNil(t *testing.T) {
	err := new(myErr)
	err = nil
	if error(err) == nil {
		panic(1)
	}
	PanicIfNotNil(err)
}

func TestPrintMarshalIndent(t *testing.T) {
	type S struct {
		A string `json:"a,omitempty"`
		B int    `json:"b,omitempty"`
	}
	PrintlnIndent(S{"hhhh", 111})
}

func TestDivideIntoGroups(t *testing.T) {
	tests := map[string]struct {
		input       []int
		groupSize   int
		expected    [][]int
		expectPanic bool
	}{
		"regular division": {
			input:     []int{1, 2, 3, 4, 5, 6},
			groupSize: 2,
			expected:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		"with leftover": {
			input:     []int{1, 2, 3, 4, 5},
			groupSize: 2,
			expected:  [][]int{{1, 2}, {3, 4}, {5}},
		},
		"single group": {
			input:     []int{1, 2, 3, 4, 5},
			groupSize: 5,
			expected:  [][]int{{1, 2, 3, 4, 5}},
		},
		"empty slice": {
			input:     []int{},
			groupSize: 2,
			expected:  [][]int{},
		},
		"single element": {
			input:     []int{1},
			groupSize: 2,
			expected:  [][]int{{1}},
		},
		"invalid group size": {
			input:       []int{1, 2, 3, 4, 5},
			groupSize:   0,
			expectPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tc.expectPanic {
					t.Errorf("unexpected panic")
				} else if r == nil && tc.expectPanic {
					t.Errorf("expected panic, got none")
				}
			}()
			result := DivideIntoGroups(tc.input, tc.groupSize)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
