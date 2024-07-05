package props

import (
	"testing"
	"time"
)

func TestFanout(t *testing.T) {
	fo := NewFanout[int](time.Second)
	go func() {
		for i := 0; i < 10; i++ {
			i := i
			go func() {
				o := fo.Sub()
				go func() {
					time.Sleep(time.Second * 10)
					fo.Unsub(o)
				}()
				for m := range o {
					t.Log("index", i, "message", m)
				}
				t.Log("index", i, "channel closed")
			}()
			time.Sleep(time.Second)
		}
	}()
	i := 0
	for {
		i++
		time.Sleep(time.Second)
		fo.Broadcast(123)
		if i == 22 {
			break
		}
	}
}
