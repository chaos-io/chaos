package main

import (
	"fmt"
	"log"
	"net/http"
)

var addr = "10018"

func main() {
	http.HandleFunc("/"+addr, func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Printf("path /%v hello handle started\n", addr)
		fmt.Fprintf(w, fmt.Sprintf("hello %v\n", addr))
	})

	fmt.Printf("http listen %s start\n", addr)
	if err := http.ListenAndServe(":"+addr, nil); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
