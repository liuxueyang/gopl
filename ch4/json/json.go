package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

func main() {
	movies := []Movie{
		{Title: "Casa", Year: 1942},
		{Title: "Cool Hand"},
		{Title: "Bullitt", Year: 1968, Color: true},
	}

	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("marshal failed: %s", err)
	}
	fmt.Printf("%s\n", data)

	fmt.Printf("ts=%s", time.Now())
}
