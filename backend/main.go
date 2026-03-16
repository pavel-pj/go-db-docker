package main

import "fmt"

func get(nums ...int) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()

	return out

}

func getSquares(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for val := range in {
			out <- val * val
		}
	}()

	return out
}

func double(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range in {
			out <- v * 2
		}
	}()

	return out
}

func main() {

	nums := get(1, 2, 3, 4, 5)
	squares := getSquares(nums)
	dbl := double(squares)

	for v := range dbl {
		fmt.Println(v)
	}

}
