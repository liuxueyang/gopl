package main

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch7/database"
	"log"
	"net/http"
)

func main() {

	db := DB(database.Db)
	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/price", http.HandlerFunc(db.price))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type DB database.Database

func (db DB) list(
	w http.ResponseWriter,
	req *http.Request,
) {

	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db DB) price(
	w http.ResponseWriter,
	req *http.Request,
) {

	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}
