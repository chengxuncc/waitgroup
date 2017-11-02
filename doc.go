// Package waitgroup implements Go sync.WaitGroup(https://golang.org/src/sync/waitgroup.go)
// semantics with channel notification and does not have concurrent race issue.
// Similarities with official sync.WaitGroup:
// 1. Same interface [Wait, Done, Wait].
// 2. Same Function: The main goroutine calls Add to set the number of
//    goroutines to wait for. Then each of the goroutines runs and
//    calls Done when finished. At the same time, Wait can be used
//    to block until all goroutines have finished.
// Differences between official sync.WaitGroup and WaitGroup here:
// 1. WaitGroup here utilizes channel notification, and does not record
//    waiters number, only close ch to wake up all blocking waiters.
// 2. Unlike sync.WaitGroup, WaitGroup.Add here can be called at any time
//    without race issue. Even call with a positive delta that occur when the
//    counter is zero, which is prevented from sync.WaitGroup.
package waitgroup //import "github.com/hzxuzhonghu/waitgroup"
