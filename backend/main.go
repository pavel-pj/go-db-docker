package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {

	time.Sleep(100 * time.Millisecond)
	for job := range jobs {
		result := job * 2
		results <- result
		fmt.Println("worker :", id, "job:", job, "result:", result)
	}

}

func main() {

	var wg sync.WaitGroup
	limit := 10
	jobs := make(chan int)
	results := make(chan int)

	for i := 0; i < limit; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			worker(id, jobs, results)
		}(i)

	}

	go func() {

		for x := 0; x < 100; x++ {
			jobs <- x
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println("RESULT:", result)
	}

}
