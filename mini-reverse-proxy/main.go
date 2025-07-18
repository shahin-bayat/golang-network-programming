package main

import (
	"github.com/shahin-bayat/mini-reverse-proxy/networking"
)

func main() {
	host := "127.0.0.1"
	port := 8081
	backends := []string{"127.0.0.1:8082", "127.0.0.1:8083"}

	server := networking.NewServer(host, port, backends)
	server.Run()
}
