package main

import (
	"fmt"
)

var prereqs = map[string][]string{
	"algorithms":            {"data structures"},
	"calculus":              {"linear algebra"},
	"compilers":             {"data structures", "formal languages", "computer organization"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	ans := topoSort(prereqs)
	for i, v := range ans {
		fmt.Printf("%d:\t%s\n", i+1, v)
	}
}

func topoSort(m map[string][]string) []string {
	var ans []string
	vis := make(map[string]bool)

	var dfs func(lst []string)

	dfs = func(lst []string) {
		for _, v := range lst {
			if vis[v] {
				continue
			}
			vis[v] = true
			dfs(m[v])
			ans = append(ans, v)
		}
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	dfs(keys)

	return ans
}
