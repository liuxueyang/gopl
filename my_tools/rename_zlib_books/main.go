package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var dir_path = flag.String("dir", "", "Directory path to rename zlib books")

func main() {
	flag.Parse()

	if *dir_path == "" {
		fmt.Fprintf(os.Stderr, "Error: Directory path is required.\n")
		return
	}

	os.Chdir(*dir_path)
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		oldName := file.Name()

		// get extension and base name
		ext := filepath.Ext(oldName)
		base := oldName[:len(oldName)-len(ext)]

		base = strings.TrimSuffix(base, "(Z-Library)")
		base = strings.TrimSuffix(base, " ")
		base = strings.TrimPrefix(base, "[Kmoe]")
		base = strings.TrimPrefix(base, " ")
		newName := base + ext

		if newName != oldName {
			err := os.Rename(oldName, newName)
			if err != nil {
				panic(err)
			}
			println("Renamed:", oldName, "to", newName)
		}
	}
	println("Renaming completed.")
}
