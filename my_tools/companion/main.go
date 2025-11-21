package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func main() {
	http.HandleFunc("/", companionHandler)

	http.ListenAndServe(":12347", nil)
}

func sanitizeFilename(filename string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := re.ReplaceAllString(filename, "_")
	sanitized, _ = strings.CutPrefix(sanitized, "Codeforces - ")
	sanitized, _ = strings.CutPrefix(sanitized, "AtCoder - ")
	sanitized = validateAtCoder(sanitized)
	sanitized = validateCodeforces(sanitized)

	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

func companionHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	var problem Problem
	if err := json.Unmarshal(body, &problem); err != nil {
		log.Fatalf("JSON unmarshaling failed: %v", err)
	}

	// log.Printf("problem = %#v\n\nbody=%#v\n", problem, string(body))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}
	dirPath := path.Join(homeDir, "competitive", sanitizeFilename(problem.Group),
		sanitizeFilename(problem.Name))

	err = os.MkdirAll(dirPath, 0o755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	fmt.Printf("Directory '%s' ensured to exist.\n", dirPath)
}

type Problem struct {
	Name        string
	Group       string
	URL         string
	Interactive bool
	MemoryLimit int
	Tests       []Case
	TestType    string
}

type Case struct {
	Input  string
	Output string
}
