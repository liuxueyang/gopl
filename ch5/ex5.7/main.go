package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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

	forEachNode(doc, startElement, endElement)
}

var depth int
var hasChild bool

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		var attrs []string
		for _, att := range n.Attr {
			attrs = append(attrs, fmt.Sprintf("%s=%q", att.Key, att.Val))
		}

		if n.Data == "img" && !hasChild {
			if len(attrs) > 0 {
				fmt.Printf("%*s<%s %s/>\n", depth*2, "", n.Data, strings.Join(attrs, " "))
			} else {
				fmt.Printf("%*s<%s/>\n", depth*2, "", n.Data)
			}
		} else {
			if len(attrs) > 0 {
				fmt.Printf("%*s<%s %s>\n", depth*2, "", n.Data, strings.Join(attrs, " "))
			} else {
				fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
			}
			depth++
		}
	} else if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Printf("%*s%s\n", depth*2, "", text)
		}
	} else if n.Type == html.CommentNode {
		fmt.Printf("%*s<!--%s-->\n", depth*2, "", n.Data)
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		if n.Data == "img" && !hasChild {
			return
		}

		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

func nodeHasChild(n *html.Node) bool {
	if n == nil {
		return false
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode || c.Type == html.TextNode || c.Type == html.CommentNode {
			return true
		}
	}
	return false
}

func forEachNode(n *html.Node, pre, post func(*html.Node)) {
	if n == nil {
		return
	}

	if pre != nil {
		if nodeHasChild(n) {
			hasChild = true
		}
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
		hasChild = false
	}
}
