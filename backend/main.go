package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Job represents a unit of work.
type Job struct {
	ID   int
	Text string
}

// Result represents a processed job.
type Result struct {
	ID     int
	Output string
}

// ProcessJobs processes jobs concurrently and stops when the context is done.
func ProcessJobs(ctx context.Context, jobs []Job) []Result {

	ch := make(chan Result, len(jobs))

	for _, v := range jobs {

		go func(ctx context.Context, value Job) {

			select {
			case <-ctx.Done():
				return
			default:
				result := process(value.Text)
				ch <- Result{
					ID:     value.ID,
					Output: result,
				}
			}

		}(ctx, v)

	}

	results := make([]Result, len(jobs))

	for i := 0; i < len(jobs); i++ {
		select {
		case <-ctx.Done():
			return results[:i]
		case res := <-ch:
			results[i] = res
		}
	}

	return results

}

func process(text string) string {
	// Simulate work.
	time.Sleep(10 * time.Millisecond)
	return strings.ToUpper(strings.TrimSpace(text))
}

func main() {

	jobs := []Job{
		{ID: 1, Text: "first"},
		{ID: 2, Text: "second"},
		{ID: 3, Text: "third"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

	defer cancel()

	results := ProcessJobs(ctx, jobs)
	fmt.Println(results)

}
