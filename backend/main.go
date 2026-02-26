package main

import (
	"fmt"
	"strings"
	"time"
)

// Item represents an input unit.
type Item struct {
	ID    int
	Value string
}

// Result represents a processed item.
type Result struct {
	ID    int
	Value string
}

type indexedResult struct {
	Index  int
	Result Result
}

// ProcessItems processes items concurrently and returns results in input order.
func ProcessItems(items []Item) []Result {
	result := make([]Result, len(items))
	ch := make(chan indexedResult, 3)

	for idx, item := range items {

		go func(index int, i Item) {

			val := process(i.Value)
			ch <- indexedResult{
				Index:  index,
				Result: Result{ID: i.ID, Value: val},
			}
		}(idx, item)
	}

	for i := 0; i < len(items); i++ {
		data := <-ch
		result[data.Index] = data.Result

	}
	close(ch)
	return result
}

func process(value string) string {
	_ = value
	// Simulate work.
	time.Sleep(5 * time.Millisecond)
	return strings.ToUpper(strings.TrimSpace(value))
}

func main() {

	items := []Item{
		{ID: 1, Value: "first"},
		{ID: 2, Value: "second"},
		{ID: 3, Value: "third"},
	}

	fmt.Println(ProcessItems(items))

}
