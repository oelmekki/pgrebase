package main

import (
	"os"
)

/*
 * Check if file exists and is a directory
 */
func IsDir( path string ) bool {
	info, err := os.Stat( path )
	if err != nil { return false }
	return info.IsDir()
}
