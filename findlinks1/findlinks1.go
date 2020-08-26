package main

import (
	"fmt"

	"golang.org/x/net/html"
)

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, v := range n.Attr {
			if v.Key == "href" {
				links = append(links, v.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = n.NextSibling {
		links = visit(links, c)
	}

	return links
}

func main() {
	fmt.Println("Hello")
}
