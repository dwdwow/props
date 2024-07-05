package props

import (
	"log/slog"
	"slices"
	"sync"
	"time"
)

type Fanout[D any] struct {
	mux    sync.Mutex
	outers []chan D
	outDur time.Duration
}

func NewFanout[D any](outDur time.Duration) *Fanout[D] {
	return &Fanout[D]{
		outDur: outDur,
	}
}

func (f *Fanout[D]) Sub() <-chan D {
	f.mux.Lock()
	defer f.mux.Unlock()
	outer := make(chan D)
	f.outers = append(f.outers, outer)
	return outer
}

func (f *Fanout[D]) Unsub(outer <-chan D) {
	f.mux.Lock()
	defer f.mux.Unlock()
	defer func() {
		recErr := recover()
		if recErr != nil {
			// should not be here
			slog.Error("Fanout: Unsub Recovered", "err", recErr)
		}
	}()
	for i, o := range f.outers {
		if o == outer {
			f.outers = slices.Delete(f.outers, i, i+1)
			close(o)
			break
		}
	}
}

func (f *Fanout[D]) Broadcast(d D) {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, o := range f.outers {
		o := o
		go func() {
			defer func() {
				recErr := recover()
				if recErr != nil {
					slog.Error("Fanout: Broadcast Recovered", "err", recErr)
				}
			}()
			t := time.NewTimer(f.outDur)
			defer t.Stop()
			select {
			case <-t.C:
				slog.Error("Fanout: One Channel Cannot Read Data", "duration", f.outDur.Milliseconds())
			case o <- d:
				// outer may be closed, should recover
			}
		}()
	}
}
