package basic

import (
	"fmt"
	"log"
	"sync"
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
