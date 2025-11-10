package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

var done = make(chan struct{})

func main() {
	worklist := make(chan []string)
	var n int

	n++
	go func() {
		worklist <- os.Args[1:]
	}()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		close(done)
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		if cancelled() {
			break
		}
		list := <-worklist
		for _, link := range list {
			if cancelled() {
				break
			}
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.ExtractV1(url, done)
	<-tokens
	if err != nil {
		log.Println(err)
		return nil
	}
	return list
}
