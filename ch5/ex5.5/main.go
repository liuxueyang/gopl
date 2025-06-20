package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := os.Args[1]
	words, images, err := CountWordsAndImages(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting words and images: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Words: %d, Images: %d\n", words, images)
}

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}

	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n == nil {
		return
	}

	if n.Type == html.TextNode {
		words += countWordsInString(n.Data)
	} else if n.Type == html.ElementNode {
		if n.Data == "img" {
			images++
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		wordsInChildren, imagesInChildren := countWordsAndImages(c)
		words += wordsInChildren
		images += imagesInChildren
	}

	return
}

func countWordsInString(s string) (count int) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		count++
	}
	return
}
