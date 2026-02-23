package main

import (
	"fmt"
	"net/http"
	"sync"
)

func fetch(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(url, resp.Status)

}

func main() {

	var wg sync.WaitGroup

	sites := []string{
		"https://google.com",
		"https://ya.ru",
		"https://mail.ru",
	}

	for _, url := range sites {
		wg.Add(1)
		go fetch(url, &wg)
	}

	wg.Wait()

}
