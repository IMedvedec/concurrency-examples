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
