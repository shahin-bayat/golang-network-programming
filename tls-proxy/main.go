package main

import (
	"crypto/tls"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		slog.Error("failed to load tls certificate files", "error", err)
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	listener, err := tls.Listen("tcp", ":8443", tlsConfig)
	if err != nil {
		slog.Error("failed to start tls listener", "error", err)
	}
	defer listener.Close()

	accept(listener)
}

func accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()
	slog.Info("connection accepted:", "client", clientConn.RemoteAddr())
	backendConn, err := net.Dial("tcp", ":8080") // the server which receives decrypted data
	if err != nil {
		slog.Error("failed to dial server", "error", err)
		return
	}
	defer backendConn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done) // client has finished or disconnected
		io.Copy(backendConn, clientConn)
	}()

	io.Copy(clientConn, backendConn)
	<-done
}
