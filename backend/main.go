package main

import (
	"fmt"
	"sync"
)

func producer(ch chan<- int) {
	for x := 1; x < 11; x++ {
		ch <- x
	}
	close(ch)
}

func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ch {

		fmt.Println("consumer:", id, "; сообещние :", msg)
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 1)

	go producer(ch)
	for i := 1; i < 4; i++ {
		wg.Add(1)
		go consumer(i, ch, &wg)
	}

	wg.Wait()

}
