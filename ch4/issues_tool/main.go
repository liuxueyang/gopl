package main

import (
	"encoding/json"
	"fmt"
	"gitee.com/liuxueyang/gopl/ch4/github"
	"log"
	"net/http"
)

var githubURL = "https://api.github.com"

type Issue struct {
	URL         string
	Title       string
	CommentsURL string `json:"comments_url"`
	HTMLURL     string `json:"html_url"`
	Body        string
	Number      int
	User *github.User
}

func GetIssue(
	owner, repo string,
	issueNumber int,
) {
	url := fmt.
		Sprintf("%s/repos/%s/%s/issues/%d",
			githubURL, owner, repo, issueNumber)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("http GET failed: %#v", resp)
		return
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("issue=%#v", issue)
}

func main() {
	// Exercise 4.11
	// 触发了频控，之后再做

	GetIssue("octocat", "hello-world", 42)
}
