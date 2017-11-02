package waitgroup

import (
	"sync"
	"sync/atomic"
)

// WaitGroup must not be copied after first use.
type WaitGroup struct {
	m       sync.Mutex
	counter int32
	ch      chan struct{}
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
func (wg *WaitGroup) Add(delta int) {
	v := atomic.AddInt32(&wg.counter, int32(delta))
	// Panic if negative
	if v < 0 {
		panic("sync: negative WaitGroup counter")

	} else if v == 0 {
		wg.m.Lock()
		defer wg.m.Unlock()
		if wg.ch != nil {
			close(wg.ch)
			wg.ch = nil
		}
	}
}

// Done decrements the WaitGroup counter.
func (wg *WaitGroup) Done() {
	wg.Add(-1)

}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	if atomic.LoadInt32(&wg.counter) == 0 {
		return

	}
	wg.m.Lock()
	if wg.ch == nil {
		wg.ch = make(chan struct{})

	}
	ch := wg.ch
	wg.m.Unlock()
	// In case when counter=0, but wg.ch has not been initiated.
	if atomic.LoadInt32(&wg.counter) == 0 {
		return

	}
	<-ch
}
