package main

import (
	"fmt"
	"sort"
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
	"linear algebra":        {"calculus"},
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
	onPath := make(map[string]bool)
	var hasCycle bool

	var dfs func(lst []string)

	dfs = func(lst []string) {
		if hasCycle {
			return
		}

		for _, v := range lst {
			if hasCycle {
				return
			}

			if onPath[v] {
				fmt.Printf("cycle detected: %s\n", v)
				hasCycle = true
				return
			}
			if vis[v] {
				continue
			}

			onPath[v] = true
			vis[v] = true
			dfs(m[v])
			onPath[v] = false

			ans = append(ans, v)
		}
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	dfs(keys)

	if hasCycle {
		fmt.Println("topological sort failed due to cycle")
		return nil
	}

	return ans
}
