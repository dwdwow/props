package props

import (
	"log/slog"
	"slices"
	"sync"
	"time"
)

type Fanout[D any] struct {
	mux    sync.RWMutex
	in     chan D
	outers []chan D
	outDur time.Duration
}

func (f *Fanout[D]) NewOuter() <-chan D {
	f.mux.Lock()
	defer f.mux.Unlock()
	outer := make(chan D)
	f.outers = append(f.outers, outer)
	return outer
}

func (f *Fanout[D]) RemoveOuter(outer <-chan D) {
	f.mux.Lock()
	defer f.mux.Unlock()
	for i, o := range f.outers {
		if o == outer {
			f.outers = slices.Delete(f.outers, i, i)
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
			t := time.NewTimer(f.outDur)
			defer t.Stop()
			select {
			case <-t.C:
				slog.Error("Fanout: One Channel Cannot Read Data", "duration", f.outDur.Milliseconds())
			case o <- d:
			}
		}()
	}
}
