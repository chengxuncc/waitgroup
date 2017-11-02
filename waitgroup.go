/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package waitgroup

import (
	"sync"
	"sync/atomic"
)

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
//
// A WaitGroup must not be copied after first use.
type WaitGroup struct {
	m       sync.Mutex
	counter int32
	ch      chan struct{}
	sync.WaitGroup
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
		if wg.ch != nil {
			close(wg.ch)
			wg.ch = nil

		}
		wg.m.Unlock()

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
