package main

import (
	"simple-dns/records"
	"simple-dns/server"
)

func main() {
	recordStore := records.InMemoryRecordStore{
		Records: map[string][4]byte{
			"example.com.": {192, 0, 2, 1},
			"another.net.": {10, 0, 0, 5},
		},
	}
	server := server.NewServer("", 53, &recordStore)
	server.Run()
}
