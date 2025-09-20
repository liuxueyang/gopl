package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var dir_path *string = flag.String("dir", ".", "Directory to search")

var skipFileExtensions map[string]bool = map[string]bool{
	".suo": true,
}
var removeDirs map[string]bool = map[string]bool{
	".vs": true,
	"x64": true,
}

func path_contain_dir(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		if removeDirs[part] {
			return true
		}
	}
	return false
}

func my_walk_dir(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}
		if path_contain_dir(path) && !skipFileExtensions[filepath.Ext(d.Name())] {
			fmt.Println("removing file:", path)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", path, err)
			}
		}
		return nil
	})
}

func main() {
	flag.Parse()

	if err := my_walk_dir(*dir_path); err != nil {
		fmt.Printf("Failed to walk directory %s: %v\n", *dir_path, err)
		os.Exit(1)
	}
}
