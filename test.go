package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter int32
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		atomic.AddInt32(&counter, 1)
	}()

	go func() {
		defer wg.Done()
		atomic.AddInt32(&counter, 2)
	}()

	wg.Wait()

	fmt.Println("counter:", atomic.LoadInt32(&counter)) // counter: 3

	old := atomic.SwapInt32(&counter, 4)
	fmt.Printf("old: %v, new: %v\n", old, counter) // old: 3, new: 4

	comp := atomic.CompareAndSwapInt32(&counter, 4, 5)
	fmt.Printf("comp: %v, counter: %v\n", comp, counter) // comp: true, counter: 5

	atomic.StoreInt32(&counter, 6)
	fmt.Printf("counter: %v\n", counter) // counter: 6
}
