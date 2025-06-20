package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	elementCount := make(map[string]int)
	countElements(elementCount, doc)

	for element, count := range elementCount {
		fmt.Printf("%s: %d\n", element, count)
	}
}

func countElements(mp map[string]int, n *html.Node) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode {
		mp[n.Data]++
	}
	countElements(mp, n.FirstChild)
	countElements(mp, n.NextSibling)
}
