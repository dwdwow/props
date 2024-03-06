package props

import (
	"fmt"
	"sync"
	"testing"
)

type oneMuxTest struct {
	SafeRWData[int]
}

func TestOneMux(t *testing.T) {
	muxTest := new(oneMuxTest)
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i, "start to change", muxTest.Get())
			muxTest.Set(i)
			fmt.Println(i, "finish changing", muxTest.Get())
		}()
	}
	wg.Wait()
}

func TestSafeSliceLastOne(t *testing.T) {
	s := NewSafeRWSlice[int]()
	fmt.Println(s.LastOne())
	s.Append(1)
	s.Append(3)
	s.Append(4)
	fmt.Println(s.LastOne())
}
