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

	namemp := make(map[string]string)
	stay := make(map[string]struct{})
	vis := make(map[string]struct{})

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		oldName := file.Name()

		// get extension and base name
		ext := filepath.Ext(oldName)
		if !isEBookFile(ext) {
			println("Not ebook file:", oldName, ext)
			continue
		}

		base := oldName[:len(oldName)-len(ext)]

		base = strings.TrimSuffix(base, "(Z-Library)")
		base = strings.TrimSuffix(base, " ")
		base = strings.TrimPrefix(base, "[Kmoe]")
		base = strings.TrimPrefix(base, " ")
		base = removeLastBracket(base)
		base = strings.TrimSpace(base)

		newName := base + ext

		if newName != oldName {
			if _, err := os.Stat(newName); err == nil {
				println("File already exists:", newName)
				continue
			}

			if _, ok := vis[newName]; ok {
				stay[oldName] = struct{}{}
				stay[newName] = struct{}{}
				fmt.Println("two files rename to the same name:", newName)
				continue
			}

			namemp[oldName] = newName
			vis[newName] = struct{}{}
		}
	}

	for oldName, newName := range namemp {
		if _, ok := stay[oldName]; ok {
			continue
		}
		if _, ok := stay[newName]; ok {
			continue
		}

		err := os.Rename(oldName, newName)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%s] => [%s]\n", oldName, newName)
	}

	fmt.Println("Renaming completed.")
}

func removeLastBracket(s string) string {
	idxr1 := strings.LastIndex(s, "(")
	idxr2 := strings.LastIndex(s, "（")

	if idxr1 == -1 && idxr2 == -1 {
		return s
	} else if idxr1 > idxr2 {
		return s[:idxr1]
	} else {
		return s[:idxr2]
	}
}

// check normal book file extensions
func isEBookFile(ext string) bool {
	_, ok := validExtensions[ext]
	return ok
}

var validExtensions = make(map[string]struct{})

func init() {
	extensions := []string{".pdf", ".epub", ".mobi", ".djvu", ".azw3", ".azw"}
	for _, ex := range extensions {
		validExtensions[ex] = struct{}{}
	}
}
