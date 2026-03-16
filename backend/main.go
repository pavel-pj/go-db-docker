package main

import (
	"context"
	"fmt"
	"time"
)

func get(ctx context.Context, nums ...int) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {

			select {
			case <-ctx.Done():
				fmt.Println("Отмена контекста в nums")
				return
			default:
				out <- n
			}

		}
	}()

	return out

}

func getSquares(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for val := range in {
			select {
			case <-ctx.Done():
				fmt.Println("Отмена контекста в squares")
				return
			default:
				out <- val * val
			}
		}
	}()

	return out
}

func double(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-ctx.Done():
				fmt.Println("Отмена контекста в double")
				return
			default:
				out <- v * 2
			}
		}
	}()

	return out
}

func main() {

	ctx, _ := context.WithTimeout(context.Background(), 100000*time.Nanosecond)
	nums := get(ctx, 1, 2, 3, 4, 5)
	squares := getSquares(ctx, nums)
	//cancel()
	dbl := double(ctx, squares)

	for v := range dbl {
		fmt.Println(v)
	}

}
