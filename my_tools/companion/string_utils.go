package main

import (
	"regexp"
	"strings"
)

func validateAtCoder(filename string) string {
	re := regexp.MustCompile(`\((AtCoder.*?)\)`)
	matches := re.FindAllStringSubmatch(filename, -1)
	if matches == nil {
		return filename
	}
	for _, match := range matches {
		if len(match) >= 2 {
			return match[1]
		}
	}
	return filename
}

// 'Educational Codeforces Round 182 (Rated for Div. 2)'
func validateCodeforces(filename string) string {
	re := regexp.MustCompile(`\(Rated for .*?\)`)
	filename = re.ReplaceAllString(filename, "")
	return strings.TrimSpace(filename)
}
