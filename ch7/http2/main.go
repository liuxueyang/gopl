package main

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch7/database"
	"log"
	"net/http"
)

func main() {

	db := DB(database.Db)
	log.Fatal(http.ListenAndServe(
		"localhost:8000",
		db))
}

type DB database.Database

func (db DB) ServeHTTP(
	w http.ResponseWriter,
	req *http.Request) {

	switch req.URL.Path {
	case "/list":
		for item, price := range db {
			fmt.Fprintf(w, "%s: %s\n", item, price)
		}
	case "/price":
		item := req.URL.Query().Get("item")
		price, ok := db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no such item: %q\n", item)
			return
		}
		fmt.Fprintf(w, "%s\n", price)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}
}
