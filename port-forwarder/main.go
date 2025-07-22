package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type config struct {
	localAddress  string
	remoteAddress string
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.localAddress, "local", "", "local address")
	flag.StringVar(&cfg.remoteAddress, "remote", "", "remote address")
	flag.Parse()

	if cfg.localAddress == "" || cfg.remoteAddress == "" {
		flag.Usage()
	}

	fmt.Println("local address", cfg.localAddress, "remote address", cfg.remoteAddress)

	listener, err := net.Listen("tcp", cfg.localAddress)
	if err != nil {
		log.Fatalf("failed to listen on %s", cfg.localAddress)
	}
	defer listener.Close()

	acceptConnection(listener, cfg.remoteAddress)
}

func acceptConnection(l net.Listener, remoteAddress string) {
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, remoteAddress)
	}
}

func handleConnection(conn net.Conn, remoteAddress string) {
	defer conn.Close()
	var wg sync.WaitGroup

	remoteConn, err := net.Dial("tcp", remoteAddress)
	if err != nil {
		return
	}
	defer remoteConn.Close()

	wg.Add(2)
	go forwardTraffic(remoteConn, conn, &wg)
	go forwardTraffic(conn, remoteConn, &wg)

	wg.Wait()
}

func forwardTraffic(dst, src net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	io.Copy(dst, src)
	// When io.Copy returns, we're done reading from src.
	// We should tell the destination that we're done writing.
	dst.(*net.TCPConn).CloseWrite()
}
