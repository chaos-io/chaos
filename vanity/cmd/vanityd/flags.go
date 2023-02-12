package main

import (
	"flag"
)

var flagRunAddr string

func init() {
	flag.StringVar(&flagRunAddr, "a", "[::]:8081", "Address to run server on")
	flag.Parse()
}
