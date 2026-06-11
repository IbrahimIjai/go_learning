package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// concurrency: not same as parallelism. It meaans multiple tasks in progress at same time. to make the tasks occur in parallel(parallel execution) we use multiple cores of the CPU. but in concurrency we can have multiple tasks in progress at the same time but not necessarily executing at the same time. it can be executed one after another but they are in progress at the same time.
var wg = sync.WaitGroup{} // wait group to wait for all goroutines to finish before exiting the main function
var mu = sync.RWMutex{}   //sync.Mutex{}    mutex to protect shared resources from concurrent access
var dbData = []string{"data1", "data2", "data3", "data4", "data5"}
var results []string // shared resource to store the results of the db calls
func main() {
	t0 := time.Now()
	for i := 0; i < len(dbData); i++ {
		wg.Add(1) // increment the wait group counter for each goroutine we are going to start
		go dbcall(i)
	}
	wg.Wait() // wait for all goroutines to finish
	fmt.Printf("Total time taken: %v\n", time.Since(t0))
}

func dbcall(i int) {
	var delay float32 = rand.Float32() * 2000
	time.Sleep(time.Duration(delay) * time.Millisecond)
	save(dbData[i])
	log()
	wg.Done()
}

func save(result string) {
	mu.Lock()
	results = append(results, result)
	mu.Unlock()
}

func log() {
	mu.RLock()
	fmt.Printf("results: %v\n", results)
	mu.RUnlock()
}
