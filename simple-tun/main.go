package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/songgao/water"
)

// LoggingWriter is a simple wrapper to log the number of bytes written.
type LoggingWriter struct {
	writer io.Writer
	name   string
}

// Write implements the io.Writer interface.
func (lw *LoggingWriter) Write(p []byte) (n int, err error) {
	n, err = lw.writer.Write(p)
	if err != nil {
		log.Printf("LOG > %s: write error: %v", lw.name, err)
	} else {
		log.Printf("LOG > %s: wrote %d bytes", lw.name, n)
	}
	return
}

type Config struct {
	Server bool
	Remote string
}

func main() {
	var cfg Config
	flag.BoolVar(&cfg.Server, "server", false, "Server (boolean)")
	flag.StringVar(&cfg.Remote, "remote", "", "remote server host and ip")
	flag.Parse()

	if !cfg.Server && cfg.Remote == "" {
		flag.Usage()
	}

	if cfg.Server {
		runServer()
	} else {
		runClient(cfg.Remote)
	}
}

func runServer() {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen for tcp connections: %v", err)
	}
	defer ln.Close()
	log.Println("Server is listening on port 9090")

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("Failed to accept  connections: %v", err)
	}

	log.Printf("Accepted connection from %s ", conn.RemoteAddr())

	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: "tun-server",
		},
	})
	if err != nil {
		log.Fatalf("Error creating TUN device: %v", err)
	}

	log.Println("TUN device created: ", iface.Name())

	go func() {
		_, err := io.Copy(&LoggingWriter{writer: iface, name: "server-iface"}, conn)
		if err != nil {
			log.Printf("Error copying from connection to interface: %v", err)
		}
	}()

	_, err = io.Copy(&LoggingWriter{writer: conn, name: "server-conn"}, iface)
	if err != nil {
		log.Printf("Error copying from interface to connection: %v", err)
	}
}

func runClient(remoteAddr string) {
	log.Printf("Connecting to server at %s ", remoteAddr)

	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatalf("Failed to dial server: %v ", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: "tun-client",
		},
	})
	if err != nil {
		log.Fatalf("Error creating TUN device: %v", err)
	}

	log.Println("TUN device created: ", iface.Name())

	go func() {
		_, err := io.Copy(&LoggingWriter{writer: iface, name: "client-iface"}, conn)
		if err != nil {
			log.Printf("Error copying from connection to interface: %v", err)
		}
	}()

	_, err = io.Copy(&LoggingWriter{writer: conn, name: "client-conn"}, iface)
	if err != nil {
		log.Printf("Error copying from interface to connection: %v", err)
	}
}
