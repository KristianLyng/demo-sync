/*
This is a demonstration of naive synchronization versus mutex-based synchronization.

The goal is to have some piece of code that only runs once on a given data
structure (e.g.: once per Skogul module), even if the function is called
multiple times.

There are three implementations: A naive approach which will not work in a parallel
processing environment (like Skogul), a manually implemented function, and one using
the convenience function of sync.Once.

Build it and try to run it a few times to see what the difference is.

This is a common issue in Skogul since Skogul frequently uses multiple
threads/go processes. Encoders and parses, for example must be able to run
in parallel.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)


// NaiveSync has a variable Synced. By default Synced is 0.
type NaiveSync struct {
	Value int
	Synced int
}

// SyncManual does the exact same thing as NaiveSync, but with an added
// Mutex.
type SyncManual struct {
	Value int
	Synced int
	Lock sync.Mutex
}

// SyncOnce uses the sync-packages convenience function sync.Once() which
// is (probably) slightly faster than the manual sync.
type SyncOnce struct {
	Value int
	Synced sync.Once
}

// Change this to change how many times to loop
const ITERATIONS = 10

// The basic idea is: Increase s.Value, but only do it once. In this
// example/demo, s.Value++ is a substitute for a "real" operation, e.g.:
// opening and parsing a file.
func (s *NaiveSync) AddMaxOne() {
	if s.Synced == 0 {
		s.Value++
		fmt.Println("Set naive value to ", s.Value)
		s.Synced = 1
	}
}

// SyncManual works the same way as NaiveSync, but before it does anything,
// it acquires a lock which means the code between Lock and Unlock will run
// exclusively on a single CPU core, never in parallel. This is the core
// mechanism of a mutex.
func (s *SyncManual) AddMaxOne() {
	s.Lock.Lock()
	if s.Synced == 0 {
		s.Value++
		fmt.Println("Set manual value to ", s.Value)
		s.Synced = 1
	}
	s.Lock.Unlock()
}

// A common use case for mutexes is to run initilization code (e.g.:
// reading a schema from disk) just once. So the sync-library provides a
// convenience function for us through sync.Once which does what we did
// "by hand" above.
func (s *SyncOnce) AddMaxOne() {
	s.Synced.Do(func() {
		s.Value++
		fmt.Println("Set synced value to ", s.Value)
	})
}

func main() {
	n := NaiveSync{}	
	s := SyncOnce{}	
	m := SyncManual{}

	// Try to remove the "go"-keyword. Instead of starting all of these
	// functions in separate go functions in parallel, they will run
	// sequentially.
	for i := 0; i < ITERATIONS; i++ {
		go n.AddMaxOne()
		go s.AddMaxOne()
		go m.AddMaxOne()
	}

	// We need to sleep before exiting, otherwise we will exit before
	// any code is run. Try to remove time.Sleep and see for yourself.
	time.Sleep(time.Second)
	fmt.Printf("Naive value: %d\n", n.Value)
	fmt.Printf("Synced value: %d\n", s.Value)
	fmt.Printf("Manually synced value: %d\n", m.Value)
}
