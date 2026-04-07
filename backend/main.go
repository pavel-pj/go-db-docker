package main

import (
	"fmt"
	"time"
)

func WaitFirst(left <-chan string, right <-chan string, timeout time.Duration) (string, bool) {
	select {
	case a := <-left:
		return a, true
	case b := <-right:
		return b, true
	case <-time.After(timeout):
		return "", false
	}

}

func main() {
	// BRANCH DEV
	left := make(chan string)
	right := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		left <- "left"
	}()

	value, ok := WaitFirst(left, right, 50*time.Millisecond)
	fmt.Println(value, ok) // left true
}
