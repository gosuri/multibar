package multibar_test

import (
	"github.com/gosuri/multibar"
	"sync"
	"time"
)

func Example() {
	bars := multibar.New()
	count1, count2 := 2000, 2500
	bar1 := bars.MakeBar(count1, "bar1")
	bar2 := bars.MakeBar(count2, "bar2")

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
