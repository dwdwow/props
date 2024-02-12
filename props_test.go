package props

import "testing"

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
