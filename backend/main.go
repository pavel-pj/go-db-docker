package main

import (
	"fmt"
	"time"
)

func WaitFirst(left <-chan string, right <-chan string, timeout time.Duration) (string, bool) {

	for {
		select {
		case data, ok := <-left:
			if ok {
				return data, true
			}
			return "", false
		case data, ok := <-right:
			if ok {
				return data, true
			}
			return "", false
		case <-time.After(timeout):
			return "", false
		}

	}

}

func main() {
	left := make(chan string)
	right := make(chan string)

	_, ok := WaitFirst(left, right, 5*time.Millisecond)
	fmt.Println(ok)

}
