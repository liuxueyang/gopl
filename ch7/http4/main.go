package main

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch7/database"
	"log"
	"net/http"
)

func main() {
	db := DB(database.Db)
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type DB database.Database

func (db DB) price(w http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("item")
	if v, ok := db[name]; ok {
		_, _ = fmt.Fprintf(w, "%s: %s\n", name, v)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "item not found: %s\n", name)
		return
	}
}

func (db DB) list(w http.ResponseWriter, req *http.Request) {

	for name, price := range db {
		_, _ = fmt.Fprintf(w, "%s: %s\n", name, price)
	}
}
