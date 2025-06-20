package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetching URL %s: %v\n", url, err)
		os.Exit(1)
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing HTML: %v\n", err)
		os.Exit(1)
	}

	res := ElementByID(doc, os.Args[2])
	if res != nil {
		fmt.Printf("Element with ID '%s' found:\n", os.Args[2])
		html.Render(os.Stdout, res)
	} else {
		fmt.Printf("Element with ID '%s' not found.", os.Args[2])
	}
}

var targetID string

func startElement(n *html.Node) bool {
	if n.Type == html.ElementNode {
		for _, att := range n.Attr {
			if att.Key == "id" && att.Val == targetID {
				fmt.Printf("%s\n", n.Data)
				return false
			}
		}
	}
	return true
}

func endElement(n *html.Node) bool {
	return true
}

func forEachNode(n *html.Node, pre, post func(*html.Node) bool) *html.Node {
	if n == nil {
		return n
	}

	if pre != nil {
		if !pre(n) {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cur := forEachNode(c, pre, post)
		if cur != nil {
			return cur
		}
	}

	if post != nil {
		post(n)
	}

	return nil
}

func ElementByID(doc *html.Node, id string) *html.Node {
	targetID = id
	return forEachNode(doc, startElement, endElement)
}
