package main

import (
	"encoding/json"
	"fmt"
	"os"
)

//	{
//		"guid": "lE7Dg4CP1k5T",
//		"title": "Get Help",
//		"index": 0,
//		"dateAdded": 1685075377811000,
//		"lastModified": 1748596282000000,
//		"id": 8,
//		"typeCode": 1,
//		"iconUri": "fake-favicon-uri:https://support.mozilla.org/products/firefox",
//		"type": "text/x-moz-place",
//		"uri": "https://support.mozilla.org/products/firefox"
//	}
type Bookmark struct {
	// The `GUID` field is a globally unique identifier for the bookmark.
	GUID string `json:"guid"`
	// The `Title` field is the name of the bookmark.
	Title string `json:"title"`
	// The `Index` field is the position of the bookmark in its parent folder.
	Index int `json:"index"`
	// The `DateAdded` and `LastModified` fields are timestamps in microseconds.
	DateAdded    int64 `json:"dateAdded"`
	LastModified int64 `json:"lastModified"`
	// The `ID` field is a unique identifier for the bookmark within the browser's database.
	ID int `json:"id"`
	// The `TypeCode` field is an integer that represents the type of the bookmark.
	TypeCode int `json:"typeCode"`
	// The `IconURI` field is a string that represents the URI of the bookmark's icon.
	IconURI string `json:"iconUri"`
	// Type is a reserved keyword in Go, so we use Kind instead
	// The `Type` field is a string that represents the type of the bookmark.
	Kind string `json:"type"`
	// The `URI` field is the URL of the bookmark.
	URI string `json:"uri"`
	// The `root` field is used to indicate the root folder of the bookmark.
	Root string `json:"root,omitempty"` // Optional field, may not be present in all bookmarks
	// Note: The `children` field is used for folders that can contain other bookmarks.
	Children []Bookmark `json:"children,omitempty"` // Optional field for nested bookmarks
}

func (b Bookmark) String() string {
	if len(b.URI) > 0 {
		if _, ok := URImp[b.URI]; !ok {
			fmt.Printf("%d. [%s](%s)\n", idx+1, b.Title, b.URI)
			idx++
			URImp[b.URI] = true // Mark this URI as printed
		}
	}
	if len(b.Children) > 0 {
		for _, child := range b.Children {
			fmt.Printf("%s", child)
		}
	}
	return ""
}

var URImp map[string]bool
var idx int

// URImp is a map to keep track of unique URIs
// It is used to ensure that each URI is only printed once, avoiding duplicates.
// It is initialized to an empty map.

func init() {
	URImp = make(map[string]bool)
	idx = 0
}

func main() {
	var bookmark Bookmark

	// read file bookmarks-2025-06-18.json
	fh, err := os.Open("bookmarks-2025-06-18.json")
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	// Use a JSON decoder to read the file and decode the JSON into the bookmarks slice
	decoder := json.NewDecoder(fh)
	if err := decoder.Decode(&bookmark); err != nil {
		panic(err)
	}

	fmt.Printf("Read %s bookmarks\n", bookmark)
}
