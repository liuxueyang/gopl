package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	snippetsDir := path.Join(homeDir, ".config/helix/snippets")

	// list *.snippet files under snippetsDir
	// write the file content to the corresponding field in the cpp.json file

	// read file
	jsonFile := path.Join(snippetsDir, "cpp.json")
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Error reading json file: %v", err)
	}

	var snippets map[string]*Snippet
	err = json.Unmarshal(data, &snippets)
	if err != nil {
		log.Fatalf("error unmarshalling file: %v", err)
	}

	entries, err := os.ReadDir(snippetsDir)
	if err != nil {
		log.Fatalf("error list directory <%s>: %v", snippetsDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".snippet") {
			baseName := strings.TrimSuffix(name, ".snippet")
			fullPath := path.Join(snippetsDir, name)
			fileContent, err := os.ReadFile(fullPath)
			if err != nil {
				log.Fatalf("error reading file <%s>: %v", fullPath, err)
			}

			if v, ok := snippets[baseName]; ok {
				v.Body = string(fileContent)
			}
		}
	}

	for key := range snippets {
		fmt.Printf("%s\n", key)
	}

	result, err := json.MarshalIndent(snippets, "", "  ")
	if err != nil {
		log.Fatalf("error marshalling snippets: %v", err)
	}

	err = os.WriteFile(jsonFile, result, 0o755)
	if err != nil {
		log.Fatalf("error write to json file: %v", err)
	}
}

type Snippet struct {
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
	Body        string `json:"body"`
}
