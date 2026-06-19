//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		resp, err := http.Get(url)
		check(err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		check(err)
		fmt.Println(string(body))
	}
}
