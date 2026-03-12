package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker:", id, "завершил работу")
			return
		default:
			fmt.Println("worker:", id, "Работает")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i < 5; i++ {
		go worker(ctx, i)
	}

	time.Sleep(2 * time.Second)
	cancel()
	time.Sleep(200 * time.Millisecond)

}
