package main

import (
	"fmt"
	"github"
	"log"
	"os"
	"time"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues:\n", result.TotalCount)

	var lessThanMonth, lessThanYear, moreThanYear []*github.Issue

	for _, item := range result.Items {
		if item.CreatedAt.Before(time.Now().AddDate(0, -1, 0)) {
			lessThanMonth = append(lessThanMonth, item)
		} else if item.CreatedAt.Before(time.Now().AddDate(-1, 0, 0)) {
			lessThanYear = append(lessThanYear, item)
		} else {
			moreThanYear = append(moreThanYear, item)
		}
	}

	fmt.Printf("==== Issues created in less than a month (%d) ====\n", len(lessThanMonth))
	for _, item := range lessThanMonth {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	fmt.Println()

	fmt.Printf("==== Issues created in less than a year (%d) ====\n", len(lessThanYear))
	for _, item := range lessThanYear {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	fmt.Println()

	fmt.Printf("==== Issues created more than a year ago (%d) ====\n", len(moreThanYear))
	for _, item := range moreThanYear {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}
