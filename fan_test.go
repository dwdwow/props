package props

import (
	"testing"
	"time"
)

func TestFanout(t *testing.T) {
	fo := NewFanout(WithFanoutDur[int](time.Second))
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

func testRadio(t *testing.T, radio *Radio[int]) {
	// Test subscribing and broadcasting
	ch1 := radio.Sub("channel1")
	ch2 := radio.Sub("channel2")

	// Test ListenerNum
	if n := radio.ListenerNum("channel1"); n != 1 {
		t.Errorf("Expected 1 listener for channel1, got %d", n)
	}

	// Test Channels
	channels := radio.Channels()
	if len(channels) != 2 {
		t.Errorf("Expected 2 channels, got %d", len(channels))
	}

	// Test broadcasting
	go func() {
		radio.Broadcast("channel1", 100)
		radio.Broadcast("channel2", 200)
	}()

	// Test receiving on ch1
	select {
	case msg := <-ch1:
		if msg != 100 {
			t.Errorf("Expected 100 on channel1, got %d", msg)
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for message on channel1")
	}

	// Test SubWithCh
	customCh := make(chan int)
	radio.SubWithCh("channel3", customCh)
	if n := radio.ListenerNum("channel3"); n != 1 {
		t.Errorf("Expected 1 listener for channel3, got %d", n)
	}

	// Test Unsub
	radio.Unsub("channel1", ch1)
	if n := radio.ListenerNum("channel1"); n != 0 {
		t.Errorf("Expected 0 listeners after unsub, got %d", n)
	}

	// Test UnsubAll
	radio.UnsubAll(ch2)
	if n := radio.ListenerNum("channel2"); n != 0 {
		t.Errorf("Expected 0 listeners after UnsubAll, got %d", n)
	}

	ch3 := radio.Sub("channel3")
	ch4 := radio.Sub("channel4")

	// Test BroadcastAll
	radio.BroadcastAll(300)
	msg := <-ch3
	if msg != 300 {
		t.Errorf("Expected 300 on channel2, got %d", msg)
	}
	msg = <-ch4
	if msg != 300 {
		t.Errorf("Expected 300 on channel2, got %d", msg)
	}
}

func TestRadio(t *testing.T) {
	radio := NewRadio(WithFanoutDur[int](time.Second))
	testRadio(t, radio)
}

func TestRadioNoDur(t *testing.T) {
	radio := NewRadio[int]()
	testRadio(t, radio)
}
