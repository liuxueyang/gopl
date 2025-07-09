package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
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

		base = strings.TrimSpace(base)
		base = strings.TrimSuffix(base, "(Z-Library)")
		base = strings.TrimPrefix(base, "[Kmoe]")
		base = removeLastBracket(base)
		base = strings.TrimSpace(base)

		if len(base) == 0 {
			// file name is enclosed in brackets
			continue
		}

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
	idxl1 := strings.LastIndex(s, "(")
	idxl2 := strings.LastIndex(s, "（")

	if idxl1 == -1 && idxl2 == -1 {
		return s
	} else if idxl1 > idxl2 {
		idxr1 := strings.Index(s[idxl1:], ")")
		if idxr1 != -1 {
			idxr1 += idxl1
			return s[:idxl1] + s[idxr1+1:]
		} else {
			return s
		}
	} else {
		idxr2 := strings.Index(s[idxl2:], "）")
		if idxr2 != -1 {
			idxr2 += idxl2
			_, sz := utf8.DecodeRuneInString(s[idxr2:])
			return s[:idxl2] + s[idxr2+sz:]
		} else {
			return s
		}
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
