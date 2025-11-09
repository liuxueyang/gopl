package main

import (
	"flag"
	"fmt"
	"log"

	"gopl.io/ch5/links"
)

var (
	maxDepth *int = flag.Int("depth", 3, "maximimum crawl depth")
	tokens        = make(chan struct{}, 20)
)

type WorkList struct {
	links []string
	dep   int
}

func main() {
	flag.Parse()

	worklist := make(chan WorkList)
	var n int

	n++
	go func() {
		worklist <- WorkList{flag.Args(), 0}
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := <-worklist
		curDep := list.dep
		for _, link := range list.links {
			if !seen[link] {
				seen[link] = true
				ok := curDep+1 <= *maxDepth
				if ok {
					n++
				}
				go func(link string) {
					if ok {
						worklist <- WorkList{crawl(link), curDep + 1}
					}
				}(link)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens
	if err != nil {
		log.Println(err)
		return nil
	}
	return list
}
