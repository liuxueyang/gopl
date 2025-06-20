package main

import (
	"fmt"
	"gitee.com/liuxueyang/gopl/ch7/database"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// Exercise 7.11
	db := DB(database.Db)
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/list_html", db.listHTML)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type DB database.Database

func (db DB) list(w http.ResponseWriter, req *http.Request) {

	for name, price := range db {
		_, _ = fmt.Fprintf(w, "%s: %s\n", name, price)
	}
}

func (db DB) listHTML(w http.ResponseWriter, req *http.Request) {
	// Exercise 7.12
	var itemList = template.Must(
		template.New("itemList").Parse(
			`
<table>
<tr style='text-align: left'>
       <th>Item</th>
       <th>Price</th>
</tr>
{{range .Items}}
<tr>
<td>{{.Item}}</td>
<td>{{.Price}}</td>
</tr>
{{end}}
</table>
`))
	type Product struct {
		Items []struct {
			Item  string
			Price database.Dollars
		}
	}
	product := Product{}

	for name, price := range db {
		product.Items = append(
			product.Items, struct {
				Item  string
				Price database.Dollars
			}{Item: name, Price: price})
	}
	if err := itemList.Execute(w, product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "failed to render template")
		return
	}
}

func (db DB) price(w http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("item")
	if v, ok := db[name]; ok {
		_, _ = fmt.Fprintf(w, "%s\n", v)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "item not found: %s\n", name)
	}
}

func (db DB) create(w http.ResponseWriter, req *http.Request) {

	values := req.URL.Query()
	item, price := values.Get("item"), values.Get("price")

	priceF, err := strconv.ParseFloat(price, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid parameter: %v\n", err)
		return
	}

	if priceF < 0 || len(item) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid parameter: price:%f, item:%s\n",
			priceF, item)
		return
	}

	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "item already exists: %s", item)
		return
	}

	db[item] = database.Dollars(priceF)
}

func (db DB) update(w http.ResponseWriter, req *http.Request) {

	values := req.URL.Query()
	item, price := values.Get("item"), values.Get("price")

	priceF, err := strconv.ParseFloat(price, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid parameter: %v\n", err)
		return
	}

	if priceF < 0 || len(item) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid parameter: price:%f, item:%s\n",
			priceF, item)
		return
	}

	db[item] = database.Dollars(priceF)
}

func (db DB) delete(w http.ResponseWriter, req *http.Request) {

	values := req.URL.Query()
	item, price := values.Get("item"), values.Get("price")

	priceF, err := strconv.ParseFloat(price, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid parameter: %v\n", err)
		return
	}

	if priceF < 0 || len(item) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "invalid parameter: price:%s, item:%s\n",
			price, item)
		return
	}

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "item not exists: %s", item)
		return
	}

	delete(db, item)
}

func (db DB) read(w http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("item")
	if v, ok := db[name]; ok {
		_, _ = fmt.Fprintf(w, "%s: %s\n", name, v)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "item not found: %s\n", name)
	}
}
