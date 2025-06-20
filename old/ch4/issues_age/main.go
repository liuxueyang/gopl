package main

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch4/github"
	"log"
	"os"
	"time"
)

func main() {
	// Exercise 4.10
	result, err := github.SearchIssues(
		os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	issues := make(map[int][]*github.Issue, 3)
	for _, issue := range result.Items {
		d := time.Since(issue.CreatedAt).Hours() / 24
		if d < 30 {
			issues[0] = append(issues[1], issue)
		} else if d < 365 {
			issues[1] = append(issues[2], issue)
		} else {
			issues[2] = append(issues[3], issue)
		}
	}

	fmt.Printf("Less than a month old:\n")
	printIssues(issues[0])

	fmt.Printf("Less than a year old:\n")
	printIssues(issues[1])

	fmt.Printf("More than a year old:\n")
	printIssues(issues[2])
}

func printIssues(items []*github.Issue) {
	for _, item := range items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}