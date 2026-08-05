package main

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	var n pgtype.Numeric
	err := n.Scan("10.5")
	fmt.Printf("Scan string: err=%v, valid=%v\n", err, n.Valid)
	
	var n2 pgtype.Numeric
	err2 := n2.Scan(10.5)
	fmt.Printf("Scan float64: err=%v, valid=%v\n", err2, n2.Valid)
}
