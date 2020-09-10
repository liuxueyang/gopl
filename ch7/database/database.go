package database

import "fmt"

type Dollars float64

func (d Dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type Database map[string]Dollars

var Db  = Database{
	"socks": 5.0,
	"shoes": 50.0,
}
