package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// Message represents a notification request.
type Message struct {
	ID      int
	Email   string
	Subject string
	Body    string
}

// SendResult represents a delivery result.
type SendResult struct {
	ID        int
	Delivered bool
}

func SendNotifications(messages []Message, limit int) []SendResult {

	if limit <= 0 {
		return []SendResult{}
	}

	semaphore := make(chan struct{}, limit)
	var result []SendResult

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, msg := range messages {
		wg.Add(1)
		go func(m Message) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			res := send(m)
			mu.Lock()
			result = append(result, SendResult{
				ID:        m.ID,
				Delivered: res,
			})
			mu.Unlock()

		}(msg)
	}

	wg.Wait()

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result

}

func send(msg Message) bool {
	_ = msg
	// Simulate I/O latency.
	time.Sleep(10 * time.Millisecond)
	return strings.Contains(msg.Email, "@")
}

func main() {
	messages := []Message{
		{ID: 1, Email: "ada@example.com", Subject: "hi", Body: "hello"},
		{ID: 2, Email: "invalid", Subject: "hi", Body: "hello"},
	}

	results := SendNotifications(messages, 1)
	fmt.Println(results)
}

/*
package notifier

import (
	"strings"
	"sync"
	"time"
)

// Message represents a notification request.
type Message struct {
	ID      int
	Email   string
	Subject string
	Body    string
}

// SendResult represents a delivery result.
type SendResult struct {
	ID        int
	Delivered bool
}

// BEGIN
// SendNotifications sends messages with limited parallelism and returns results in input order.
func SendNotifications(messages []Message, limit int) []SendResult {
	results := make([]SendResult, len(messages))
	if limit <= 0 {
		for i, msg := range messages {
			results[i] = SendResult{ID: msg.ID, Delivered: false}
		}
		return results
	}

	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	wg.Add(len(messages))

	for i, msg := range messages {
		go func(idx int, m Message) {
			defer wg.Done()
			sem <- struct{}{}
			delivered := send(m)
			<-sem

			results[idx] = SendResult{ID: m.ID, Delivered: delivered}
		}(i, msg)
	}

	wg.Wait()
	return results
}
// END

func send(msg Message) bool {
	_ = msg
	// Simulate I/O latency.
	time.Sleep(10 * time.Millisecond)
	return strings.Contains(msg.Email, "@")
}
*/
