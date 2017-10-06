package main

import (
	"fmt"
	"os"
)

func errexit(format string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", params...)
	os.Exit(2)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
