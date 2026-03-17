package main

import (
	"fmt"
	"sort"
	"strings"
)

type indexedLine struct {
	Index int
	Line  string
}

// ProcessLogs normalizes lines using a simple pipeline and returns results in order.
func ProcessLogs(lines []string) []string {
	var out []string

	in := make(chan indexedLine)

	normal := normalize(in)
	filtered := filterEmpty(normal)

	go func() {
		for i, l := range lines {

			in <- indexedLine{
				Index: i,
				Line:  l,
			}
		}
		close(in)
	}()

	var results []indexedLine

	for v := range filtered {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	for _, v := range results {
		fmt.Println("Добавляем : ", v.Line)
		out = append(out, v.Line)
	}

	for i, l := range out {
		fmt.Println(i, " Значение:", l)
		fmt.Println(len(out))
	}

	return out
}

func normalize(in <-chan indexedLine) <-chan indexedLine {
	out := make(chan indexedLine)

	go func() {
		defer close(out)
		for value := range in {

			//fmt.Println("normalize in, val:", value)
			res := strings.ToLower(strings.TrimSpace(value.Line))
			out <- indexedLine{
				Index: value.Index,
				Line:  res,
			}
		}
	}()

	return out
}

func filterEmpty(in <-chan indexedLine) <-chan indexedLine {
	out := make(chan indexedLine)

	go func() {
		defer close(out)
		for value := range in {
			if value.Line != "" {
				fmt.Println("empty in, val:", value)
				out <- value
			}
		}

	}()

	return out
}

func main() {
	lines := []string{"INFO: Started  ", "", " WARN: Disk FULL ", "error: failed"}
	fmt.Println(ProcessLogs(lines))

}
