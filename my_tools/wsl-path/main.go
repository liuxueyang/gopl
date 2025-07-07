package main

import (
	"flag"
	"strings"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		println("Usage: wsl-path <path>")
		return
	}

	path := flag.Arg(0)
	if path == "" {
		println("Error: No path provided")
		return
	}

	// Convert the Windows path to WSL path
	wslPath := convertToWSLPath(path)
	println(wslPath)
}

func convertToWSLPath(path string) string {
	// Check if the path is a Windows path
	if len(path) > 2 && path[1] == ':' {
		// Convert Windows path to WSL path
		return "/mnt/" + strings.ToLower(string(path[0])) + strings.ReplaceAll(path[2:], "\\", "/")
	}
	// If it's already a WSL path, return it as is
	return path
}
