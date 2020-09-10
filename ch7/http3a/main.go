package main

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch7/database"
	"log"
	"net/http"
)

type DB database.Database

func main() {
	db := DB(database.Db)
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func (db DB) list(w http.ResponseWriter, r *http.Request) {

	for name, price := range db {
		_, _ = fmt.Fprintf(w, "%s:%s\n", name, price)
	}
}

func (db DB) price(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("item")
	if v, ok := db[name]; !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "item not found: %s", name)
		return
	} else {
		_, _ = fmt.Fprintf(w, "%s: %s\n", name, v)
	}
}
