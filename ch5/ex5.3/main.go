package main

import (
	"fmt"
	"os"

	"strings"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	printInsideScript(false, doc)
}

func printInsideScript(insideScript bool, n *html.Node) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode {
		if n.Data == "script" || n.Data == "style" {
			insideScript = true
		}
	}

	if insideScript && n.Type == html.TextNode {
		data := strings.TrimSpace(n.Data)
		if len(data) > 0 {
			fmt.Printf("%s\n-------\n", data)
		}
	}
	printInsideScript(insideScript, n.FirstChild)
	printInsideScript(insideScript, n.NextSibling)
}
