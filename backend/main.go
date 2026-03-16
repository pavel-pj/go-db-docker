package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	ID  int
	Val string
}

type Result struct {
	ID     int
	Output string
	Err    error
}

func worker(id int, jobs <-chan Job, result chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range jobs {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Worker:", id, "Запись ID:", val.ID)
		result <- Result{
			ID:     val.ID,
			Output: val.Val,
			Err:    nil,
		}
	}

}

func main() {
	var wg sync.WaitGroup
	workers := 4

	jobs := make(chan Job, 16)
	results := make(chan Result, 16)

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(i, jobs, results, &wg)
	}

	go func() {
		for x := 0; x < 45; x++ {
			jobs <- Job{
				ID:  x,
				Val: fmt.Sprintf("Значение Задачи : %s", strconv.Itoa(x)),
			}
		}
		close(jobs)
	}()

	var resultWg sync.WaitGroup
	resultInfo := []Result{}

	resultWg.Add(1)

	go func() {
		defer resultWg.Done()
		for res := range results {
			resultInfo = append(resultInfo, res)
		}
	}()

	wg.Wait()
	close(results)
	resultWg.Wait()

	fmt.Println("\n=== ВСЕ РЕЗУЛЬТАТЫ ===")
	for _, res := range resultInfo {
		fmt.Printf("РЕЗУЛЬТАТ: id: %02d, Значение : %s\n", res.ID, res.Output)
	}

}
