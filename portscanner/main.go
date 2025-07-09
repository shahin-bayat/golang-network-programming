package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"sync"
	"time"
)

type ScanResult struct {
	Port int
	Open bool
}

func main() {
	host := flag.String("host", "", "hostname (required)")
	start := flag.Int("start", 0, "start port (required)")
	end := flag.Int("end", 0, "end port (required)")
	flag.Parse()

	if *host == "" ||
		*start < 1 || *start > 65535 ||
		*end < 1 || *end > 65535 ||
		*start > *end {
		flag.Usage()
		log.Fatalf("invalid arguments: host=%q start=%d end=%d", *host, *start, *end)
	}
	fmt.Printf("Target: %s, Ports: %d–%d\n", *host, *start, *end)
	startTime := time.Now()
	results := make(chan ScanResult)
	semaphore := make(chan struct{}, 100) // Buffered channel as semaphore, 100 concurrent goroutines
	var wg sync.WaitGroup
	for p := *start; p <= *end; p++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}        // add to semaphore - blocks if it is full until one of goroutines finishes
			defer func() { <-semaphore }() // release semaphore

			open := scanPort(*host, p, time.Millisecond*500)
			results <- ScanResult{Port: p, Open: open}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	var all []ScanResult
	for r := range results {
		all = append(all, r)
	}
	elapsed := time.Since(startTime)
	printResults(all)
	fmt.Printf("%d ports scan completed in %v\n", *end-*start+1, elapsed)
}

func scanPort(host string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func printResults(res []ScanResult) {
	sort.Slice(res, func(i, j int) bool {
		return res[i].Port < res[j].Port
	})
	for _, r := range res {
		status := "⛔️ close"
		if r.Open {
			status = "✅ open"
			fmt.Printf("Port %d: %s\n", r.Port, status)
		}
	}
}
