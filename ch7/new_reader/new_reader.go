package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
)

type ReadString struct {
	s string
	i int64
}

func (r *ReadString) Read(p []byte) (n int, err error) {
	n = copy(p, r.s[r.i:])
	r.i += int64(n)
	if r.i >= int64(len(r.s)) {
		err = io.EOF
	}
	return
}

func NewReader(s string) io.Reader {
	return &ReadString{s: s, i:0}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, v := range n.Attr {
			if v.Key == "href" {
				links = append(links, v.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}

func main() {

	r := NewReader(`<html> <a href="a.out"> </html>`)
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findLinks1: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}
