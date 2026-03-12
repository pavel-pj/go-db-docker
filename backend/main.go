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
			fmt.Println("worker", id, "остановлен:", ctx.Err())
			return
		default:
			fmt.Println("worker:", id, "работает")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	for x := 1; x < 3; x++ {
		go worker(ctx, x)
	}

	time.Sleep(6 * time.Second)

}

/*
func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done(): // сигнал отмены или дедлайна
			fmt.Println("worker", id, "остановлен:", ctx.Err())
			return // корректное завершение горутины
		default:
			fmt.Println("worker", id, "работает")
			time.Sleep(100 * time.Millisecond) // имитация полезной работы
		}
	}
}

func main() {
	// Контекст сам отменится через 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // ВСЕ РАВНО НУЖНО ВЫЗВАТЬ! (освобождает ресурсы)

	go worker(ctx, 1)

	// Можно даже не вызывать cancel() вручную - таймаут сработает сам
	time.Sleep(3 * time.Second)
}
*/
