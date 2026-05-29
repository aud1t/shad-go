//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type fetchLog struct {
	url     string
	status  bool
	size    int
	elapsed time.Duration
	err     error
}

func main() {
	urls := os.Args[1:]
	start := time.Now()
	fetchLogs := make(chan fetchLog)
	for _, url := range urls {
		go func(url string, ch chan<- fetchLog) {
			f := fetchLog{url: url}
			st := time.Now()
			resp, err := http.Get(url)
			if err != nil {
				f.err = err
				ch <- f
				return
			}
			defer resp.Body.Close()
			dat, err := io.ReadAll(resp.Body)
			if err != nil {
				f.err = err
				ch <- f
				return
			}
			f.status = true
			f.size = len(dat)
			f.elapsed = time.Since(st)
			ch <- f
		}(url, fetchLogs)
	}
	for range urls {
		log := <-fetchLogs
		if log.status {
			fmt.Printf("%v\t %d %s\n", log.elapsed, log.size, log.url)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("%v elapsed\n", elapsed)
}
