package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var counter int64
var wg sync.WaitGroup

func main() {

	for x := 0; x < 100; x++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				atomic.AddInt64(&counter, 1)
				fmt.Println(atomic.LoadInt64(&counter))
			}
		}(x)
	}

	wg.Wait()
	fmt.Println("total:", atomic.LoadInt64(&counter))

}
