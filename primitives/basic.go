package primitives

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	wg.Wait()
	log.Printf("Reading and counting to %d ended.\n", n)
}

func PrintOnceToN(n int) {
	printerToN := func() {
		for i := 1; i <= n; i++ {
			fmt.Printf("Printing #%d\n", i)
		}
	}

	var once sync.Once
	for i := 0; i < 100; i++ {
		once.Do(printerToN)
	}
}

func CondJobWhenNotified() {
	cond := sync.NewCond(&sync.Mutex{})

	job := func() {
		cond.L.Lock()
		defer cond.L.Unlock()

		fmt.Println("Waiting for event!")
		time.Sleep(time.Second)
		cond.Wait()
		time.Sleep(time.Second)
		fmt.Println("\tDoing my job!")
	}

	go job()
	go job()
	go job()

	<-time.NewTimer(5 * time.Second).C
	fmt.Println("Broadcasting!")
	cond.Broadcast()

	end := make(chan os.Signal, 1)
	signal.Notify(end, os.Interrupt, syscall.SIGTERM)
	<-end
}
