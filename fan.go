package props

import (
	"log/slog"
	"os"
	"slices"
	"sync"
	"time"
)

type Fanout[D any] struct {
	mux    sync.Mutex
	outers []chan D
	outDur time.Duration
	chCap  int
	logger *slog.Logger
}

type FanoutOption[D any] func(*Fanout[D])

func WithFanoutDur[D any](outDur time.Duration) FanoutOption[D] {
	return func(f *Fanout[D]) {
		f.outDur = outDur
	}
}

func WithFanoutLogger[D any](logger *slog.Logger) FanoutOption[D] {
	return func(f *Fanout[D]) {
		if logger != nil {
			f.logger = logger
		}
	}
}

func WithFanoutChCap[D any](chCap int) FanoutOption[D] {
	return func(f *Fanout[D]) {
		f.chCap = chCap
	}
}

func NewFanout[D any](opts ...FanoutOption[D]) *Fanout[D] {
	f := &Fanout[D]{
		outDur: time.Second,
		chCap:  1024,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *Fanout[D]) ListenerNum() int {
	f.mux.Lock()
	defer f.mux.Unlock()
	return len(f.outers)
}

func (f *Fanout[D]) Sub() <-chan D {
	f.mux.Lock()
	defer f.mux.Unlock()
	outer := make(chan D, f.chCap)
	f.outers = append(f.outers, outer)
	return outer
}

func (f *Fanout[D]) SubAt(ch chan D) {
	f.mux.Lock()
	defer f.mux.Unlock()
	for _, o := range f.outers {
		if o == ch {
			return
		}
	}
	f.outers = append(f.outers, ch)
}

func (f *Fanout[D]) Unsub(ch <-chan D) {
	f.mux.Lock()
	defer f.mux.Unlock()
	defer func() {
		recErr := recover()
		if recErr != nil {
			// should not be here
			f.logger.Error("Fanout: Unsub Recovered", "err", recErr)
		}
	}()
	for i, o := range f.outers {
		if o == ch {
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
					f.logger.Error("Fanout: Broadcast Recovered", "err", recErr)
				}
			}()
			t := time.NewTimer(f.outDur)
			defer t.Stop()
			select {
			case <-t.C:
				f.logger.Error("Fanout: One Channel Cannot Read Data", "duration", f.outDur.Milliseconds())
			case o <- d:
				// outer may be closed, should recover
			}
		}()
	}
}

type Radio[D any] struct {
	fans *SafeRWMap[string, *Fanout[D]]
	opts []FanoutOption[D]
}

func NewRadio[D any](opts ...FanoutOption[D]) *Radio[D] {
	return &Radio[D]{
		fans: NewSafeRWMap[string, *Fanout[D]](),
		opts: opts,
	}
}

func (r *Radio[D]) newFanout(channel string) *Fanout[D] {
	outer := NewFanout(r.opts...)
	outer.logger = outer.logger.With("radio_channel", channel)
	return outer
}

func (r *Radio[D]) Channels() []string {
	return r.fans.Keys()
}

func (r *Radio[D]) ListenerNum(channel string) int {
	fan, ok := r.fans.GetVWithOk(channel)
	if ok {
		return fan.ListenerNum()
	}
	return 0
}

func (r *Radio[D]) Sub(channel string) <-chan D {
	outer := r.fans.GetVWithNew(channel, func() *Fanout[D] {
		return r.newFanout(channel)
	})
	return outer.Sub()
}

func (r *Radio[D]) Unsub(channel string, ch <-chan D) {
	fan, ok := r.fans.GetVWithOk(channel)
	if ok {
		fan.Unsub(ch)
	}
}

func (r *Radio[D]) UnsubAll(ch <-chan D) {
	r.fans.Iterate(func(key string, fan *Fanout[D]) {
		fan.Unsub(ch)
	})
}

func (r *Radio[D]) SubWithCh(channel string, ch chan D) {
	r.fans.GetVWithNew(channel, func() *Fanout[D] {
		return r.newFanout(channel)
	}).SubAt(ch)
}

func (r *Radio[D]) Broadcast(channel string, d D) {
	fan, ok := r.fans.GetVWithOk(channel)
	if ok {
		fan.Broadcast(d)
	}
}
