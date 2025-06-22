package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopl.io/ch5/links"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)

	for len(worklist) > 0 {
		items := worklist
		worklist = nil

		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)

	// save the page to a file, create directories as necessary
	newURL := strings.TrimPrefix(url, "http://")
	newURL = strings.TrimPrefix(newURL, "https://")
	parts := strings.Split(newURL, "/")
	baseFileName := parts[len(parts)-1]

	if len(baseFileName) > 0 {
		direName := strings.Join(parts[:len(parts)-1], "/")
		if err := os.MkdirAll(direName, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", direName, err)
			return nil
		}

		if !strings.HasSuffix(baseFileName, ".html") {
			baseFileName += ".html"
		}

		fullFileName := filepath.Join(direName, baseFileName)

		file, err := os.Create(fullFileName)
		if err != nil {
			log.Printf("Failed to create file %s: %v", fullFileName, err)
			return nil
		}
		defer file.Close()

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch %s: %v", url, err)
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error fetching %s: %s", url, resp.Status)
			return nil
		}
		if _, err := io.Copy(file, resp.Body); err != nil {
			log.Printf("Failed to write to file %s: %v", fullFileName, err)
			return nil
		}
	}

	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
		return nil
	}

	res := make([]string, 0, len(list))
	for _, link := range list {
		if strings.HasPrefix(link, domain) {
			res = append(res, link)
		}
	}
	return res
}

var domain string

func main() {
	domain = os.Args[1]
	breadthFirst(crawl, []string{domain})
}
