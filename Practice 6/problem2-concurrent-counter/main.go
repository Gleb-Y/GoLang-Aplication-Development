package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Solution 1: Using sync.Mutex
func solutionWithMutex() {
	fmt.Println("=== Solution 1: Using sync.Mutex ===")
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d\n\n", counter)
}

// Solution 2: Using atomic operations
func solutionWithAtomic() {
	fmt.Println("=== Solution 2: Using atomic.AddInt64 ===")
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d\n\n", counter)
}

// Buggy version to demonstrate the problem
func buggyVersion() {
	fmt.Println("=== BUGGY VERSION (without synchronization) ===")
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d (Expected: 1000)\n\n", counter)
}

func main() {
	fmt.Println("Concurrent Counter Problem\n")
	fmt.Println("Explanation: Without synchronization, multiple goroutines")
	fmt.Println("read-modify-write the counter concurrently, causing lost updates.\n")

	buggyVersion()

	solutionWithMutex()
	solutionWithAtomic()
}
