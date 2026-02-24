package main

import (
	"fmt"
	"time"
)

func CountMultiples(nums []int, k int) int {
	if k == 0 {
		return 0
	}

	ch := make(chan int)

	for _, num := range nums {

		go func(n int, k2 int) {

			if n%k2 == 0 {
				ch <- 1
				fmt.Println(n)
				time.Sleep(11100 * time.Millisecond)
			} else {
				ch <- 0
			}
		}(num, k)

	}
	//close(ch)
	sum := 0

	for i := 0; i < len(nums); i++ {
		sum = sum + <-ch
	}

	return sum

}

func main() {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 10, 12, 15, 18, 21}
	fmt.Println(CountMultiples(nums, 3))
}
