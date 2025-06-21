package links

import (
	"fmt"
	"testing"
)

// Write test for Extract function
func TestExtract(t *testing.T) {
	url := "https://go.dev/tos"
	links, err := Extract(url)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(links) == 0 {
		t.Errorf("Expected links, got none")
	}

	for _, link := range links {
		if link == "" {
			t.Error("Found an empty link")
		}
		fmt.Printf("%s\n", link)
	}
}
