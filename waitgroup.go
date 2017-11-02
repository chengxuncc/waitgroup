package waitgroup

import (
	"sync"
	"sync/atomic"
)

type WaitGroup struct {
	m       sync.Mutex
	counter int32
	ch      chan struct{}
}

func (wg *WaitGroup) Add(delta int) {
	v := atomic.AddInt32(&wg.counter, int32(delta))
	// Panic if negative
	if v < 0 {
		panic("sync: negative WaitGroup counter")

	} else if v == 0 {
		wg.m.Lock()
		if wg.ch != nil {
			close(wg.ch)
			wg.ch = nil

		}
		wg.m.Unlock()

	}

}

func (wg *WaitGroup) Done() {
	wg.Add(-1)

}

func (wg *WaitGroup) Wait() {
	if atomic.LoadInt32(&wg.counter) == 0 {
		return

	}
	wg.m.Lock()
	if wg.ch == nil {
		wg.ch = make(chan struct{})

	}
	wg.m.Unlock()
	// In case when counter=0, but wg.ch has not been initiated.
	if atomic.LoadInt32(&wg.counter) == 0 {
		return

	}
	<-wg.ch

}
