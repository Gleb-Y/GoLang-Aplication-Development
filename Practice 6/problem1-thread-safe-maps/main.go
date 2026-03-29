package main

import (
	"fmt"
	"sync"
)

// Method 1: Using sync.Map
func demonstrateSyncMap() {
	fmt.Println("=== Method 1: sync.Map ===")
	var m sync.Map

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(key int) {
			defer wg.Done()
			m.Store(fmt.Sprintf("key%d", key), key)
		}(i)
	}

	wg.Wait()

	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	fmt.Printf("sync.Map - Stored %d items\n\n", count)
}

// Method 2: Using sync.RWMutex with regular map
func demonstrateRWMutex() {
	fmt.Println("=== Method 2: sync.RWMutex with regular map ===")
	var mu sync.RWMutex
	m := make(map[string]int)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(key int) {
			defer wg.Done()
			mu.Lock()
			m[fmt.Sprintf("key%d", key)] = key
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.RLock()
	count := len(m)
	mu.RUnlock()

	fmt.Printf("sync.RWMutex - Stored %d items\n\n", count)
}

func main() {
	fmt.Println("Thread-Safe Maps Demo\n")
	fmt.Println("Problem: Multiple goroutines writing to the same map\n")

	for i := 1; i <= 3; i++ {
		fmt.Printf("--- Run %d ---\n", i)
		demonstrateSyncMap()
		demonstrateRWMutex()
	}
}
