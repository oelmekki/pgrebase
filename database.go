package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

/*
 * Wrapper for sql.Query, meant to be the main query interface
 */
func Query( query string, parameters ...interface{} ) ( rows *sql.Rows, err error ) {
	var co *sql.DB

	co, err = sql.Open( "postgres", Cfg.DatabaseUrl )
	if err != nil {
		fmt.Printf( "can't connect to database : %v\n", err )
		return rows, err
	}
	defer co.Close()

	rows, err = co.Query( query, parameters... )
	if err != nil {
		fmt.Printf( "Can't execute query : %v\n%v\n", query, err )
		return rows, err
	}

	return
}

/*
 * Find a sensible default for max connections.
 *
 * We don't want to just use pg max conn, because some other
 * clients may be using it, so we'll take half that number.
 *
 * This will be useful for concurrently loading sql files.
 */
func FindMaxConnection() ( max int, err error ) {
	rows, err := Query( `SELECT setting FROM pg_settings WHERE name = 'max_connections'` )
	if err != nil { return max, err }
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan( &max ) ; err != nil { return max, err }
	}

	max = max / 2

	return
}
