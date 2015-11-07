package main

import (
	"sync"
	"time"

	"github.com/gosuri/multibar"
)

func main() {
	bars := multibar.New()
	count1, count2 := 2000, 2500
	bar1 := bars.MakeBar(count1, "bar1")
	bar2 := bars.MakeBar(count2, "bar2")

	// Start listening for updates
	bars.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i <= count1; i++ {
			bars.Increment(bar1, i)
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i <= count2; i++ {
			bars.Increment(bar2, i)
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Wait()
}
