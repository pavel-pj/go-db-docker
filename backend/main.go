package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // сигнал от системы или отмены
			fmt.Println("завершение воркера:", ctx.Err())
			return
		default:
			fmt.Println("фонова работа")
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func main() {

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	go worker(ctx)
	time.Sleep(10 * time.Second)

}
