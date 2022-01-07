package basic

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func PrintToN(n int) {
	log.Printf("Starting printing to %d.\n", n)

	var wg sync.WaitGroup

	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("Printing #%d\n", i)
		}(i)
	}

	wg.Wait()
	log.Printf("Printing to %d ended.\n", n)
}

func CountToN(n int) {
	log.Printf("Starting counting to %d.\n", n)

	var wg sync.WaitGroup
	var mtx sync.Mutex
	var counter int

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mtx.Lock()
			defer mtx.Unlock()

			counter++
		}()
	}

	wg.Wait()
	log.Printf("Counting to %d ended. Counter = %d.\n", n, counter)
}

func RWCountToN(n int) {
	log.Printf("Starting reading and counting to %d.\n", n)

	var (
		wg      sync.WaitGroup
		rwmtx   sync.RWMutex
		counter int

		counterFunc = func() {
			defer wg.Done()
			fmt.Printf("CounterFunc waiting for write lock!\n")
			rwmtx.Lock()
			defer rwmtx.Unlock()
			fmt.Printf("CounterFunc took the write lock! Counter = #%d\n", counter)

			counter++
			fmt.Printf("CounterFunc releasing write lock! Counter = #%d\n", counter)
		}

		printFunc = func() {
			defer wg.Done()
			fmt.Printf("\tPrintFunc waiting for read lock!\n")
			rwmtx.RLock()
			defer rwmtx.RUnlock()
			fmt.Printf("\tPrintFunc took the read lock! Counter = #%d\n", counter)

			time.Sleep(time.Millisecond)
			fmt.Printf("\tPrintFunc releasing read lock! Counter = #%d\n", counter)
		}
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go counterFunc()

		for j := 0; j < 3; j++ {
			wg.Add(1)
			go printFunc()
		}
	}

	wg.Done()
	time.Sleep(time.Second)
	log.Printf("Reading and counting to %d ended.\n", n)
}
