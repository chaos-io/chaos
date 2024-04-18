package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

func main() {
	// original functionality
	sometimesFlakes()

	fmt.Println("-----------------------------------------------------------")

	// functionality where context wraps the entire operation: "do request + read body"
	shouldNotFlake()
}

func sometimesFlakes() {
	req, err := http.NewRequest("GET", "https://swapi.co/api/people/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := fetch(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(readBody(resp))
}

func shouldNotFlake() {
	req, err := http.NewRequest("GET", "https://swapi.co/api/people/1", nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := fetchWithContext(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(readBody(resp))
}

func fetch(req *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return ctxhttp.Do(ctx, http.DefaultClient, req)
}

func fetchWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	return ctxhttp.Do(ctx, http.DefaultClient, req)
}

func readBody(resp *http.Response) (string, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), err
}
