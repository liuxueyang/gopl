package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

var search_bytes = flag.String("search", "", "Search bytes")
var replace_bytes = flag.String("replace", "", "Replace bytes")
var file_path = flag.String("file", "", "File path to modify")

func main() {
	flag.Parse()

	if *file_path == "" {
		panic("File path is required")
	}
	if *search_bytes == "" {
		panic("Search bytes are required")
	}
	if *replace_bytes == "" {
		panic("Replace bytes are required")
	}
	if _, err := os.Stat(*file_path); os.IsNotExist(err) {
		panic("File does not exist: " + *file_path)
	}

	search_bytes1 := hexToBytes(removeSpace(*search_bytes))
	replace_bytes1 := hexToBytes(removeSpace(*replace_bytes))

	if len(search_bytes1) != len(replace_bytes1) {
		panic("Search bytes and replace bytes must have the same length")
	}

	// read the file content into a byte slice
	file_bytes, err := os.ReadFile(*file_path)
	if err != nil {
		panic(err)
	}

	// search for the search_bytes in the bytes slice
	cnt := 0
	for i := range file_bytes {
		if i+len(search_bytes1) <= len(file_bytes) && bytes.Equal(file_bytes[i:i+len(search_bytes1)], search_bytes1) {
			i += len(replace_bytes1) - 1 // skip over the replaced bytes
			cnt++
		}
	}
	fmt.Printf("Found %d occurrences of the search bytes.\n", cnt)

	file_bytes = bytes.ReplaceAll(file_bytes, search_bytes1, replace_bytes1)
	// write the modified bytes back to the file
	err = os.WriteFile(*file_path, file_bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func removeSpace(s string) string {
	return strings.ToLower(string(bytes.ReplaceAll([]byte(s), []byte(" "), []byte{})))
}

func hexToBytes(hex string) []byte {
	hex = removeSpace(hex)
	if len(hex)%2 != 0 {
		hex = "0" + hex // pad with leading zero if odd length
	}

	bytes := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		var b byte
		fmt.Sscanf(hex[i:i+2], "%x", &b)
		bytes[i/2] = b
	}
	return bytes
}
