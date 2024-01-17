package main

import (
	"fmt"
	"log"
	"net/http"
)

var addr = "8092"

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("addr: %v hello handle started\n", addr)
	fmt.Fprintf(w, fmt.Sprintf("hello %v\n", addr))
}

func main() {
	http.HandleFunc("/"+addr, hello)
	err := http.ListenAndServe(":"+addr, nil)
	if err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
