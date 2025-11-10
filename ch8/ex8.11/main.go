package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var done = make(chan struct{})

func main() {
	flag.Parse()
	urls := flag.Args()
	if len(urls) == 0 {
		log.Fatalf("provide at least one mirror")
	}
	resp := mirroredQuery(urls)

	fmt.Println(resp[:100])
	// panic("test")
}

func mirroredQuery(mirrors []string) string {
	responses := make(chan string, 3)
	for _, url := range mirrors {
		go func(u string) {
			responses <- request(u)
		}(url)
	}

	ans := <-responses
	close(done)
	return ans
}

func request(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("new request error: %v", err)
		return ""
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("request error: %v", err)
		return ""
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("status is not ok: %d", resp.StatusCode)
		return ""
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v", err)
		return ""
	}

	return string(content)
}
