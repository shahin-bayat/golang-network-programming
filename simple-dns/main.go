package main

import (
	"fmt"
	"log/slog"
	"os"
	"simple-dns/resolver"

	"golang.org/x/net/dns/dnsmessage"
)

func main() {
	dnsNama, _ := dnsmessage.NewName("shahinbayat.ir.")
	answers, err := resolver.Resolve(dnsNama, dnsmessage.TypeA)
	if err != nil {
		slog.Error("IP not found", "error", err)
		os.Exit(1)
	}
	for _, a := range answers {
		fmt.Println("found IP:", a.Body)
	}
	// server := server.NewServer("", 530, &recordStore)
	// server.Run()
}
