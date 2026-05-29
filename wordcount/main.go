//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	files := os.Args[1:]
	wordCount := make(map[string]int)
	for _, file := range files {
		dat, err := os.ReadFile(file)
		check(err)
		words := strings.Split(string(dat), "\n")
		for _, word := range words {
			wordCount[word]++
		}
	}
	for word, cnt := range wordCount {
		if cnt > 1 {
			fmt.Printf("%d\t%s\n", cnt, word)
		}
	}
}
