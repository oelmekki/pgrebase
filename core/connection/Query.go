package connection

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Query is a wrapper for sql.Query, meant to be the main query interface.
func Query(query string, parameters ...interface{}) (rows *sql.Rows, err error) {
	var co *sql.DB

	co, err = sql.Open("postgres", conf.DatabaseUrl)
	if err != nil {
		fmt.Printf("can't connect to database : %v\n", err)
		return rows, err
	}
	defer co.Close()

	rows, err = co.Query(query, parameters...)
	if err != nil {
		return rows, err
	}

	return
}
