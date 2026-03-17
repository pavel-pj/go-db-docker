package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Task represents input data.
type Task struct {
	ID   int
	Text string
}

// Result represents processed data.
type Result struct {
	ID     int
	Output string
}

var ErrEmptyText = errors.New("empty text")

// ProcessTasks processes tasks concurrently and returns an error if any task fails.
func ProcessTasks(tasks []Task) ([]Result, error) {
	results := make([]Result, len(tasks))
	errs := make(chan error, len(tasks))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for idx, t := range tasks {
		wg.Add(1)
		go func(val Task) {
			defer wg.Done()

			res, err := process(t.Text)
			if err != nil {
				errs <- err
			}
			mu.Lock()
			results[idx] = Result{
				ID:     t.ID,
				Output: res,
			}
			mu.Unlock()

		}(t)

	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	// Ждем все ошибки
	var firstErr error
	for err := range errs {
		if firstErr == nil {
			firstErr = err
		}
	}

	// ТЕПЕРЬ сортируем (после того как все данные собраны)
	if firstErr != nil {
		return nil, firstErr
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results, nil
}

func process(text string) (string, error) {
	value := strings.TrimSpace(text)
	if value == "" {
		return "", ErrEmptyText
	}
	return strings.ToUpper(value), nil
}

func main() {
	tasks := []Task{
		{ID: 1, Text: "first"},
		{ID: 2, Text: "second"},
		{ID: 3, Text: "third"},
	}

	results, err := ProcessTasks(tasks)

	fmt.Println(results, err)

}
