package networking

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	host     string
	port     int
	listener net.Listener
}

func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("unable to start server:", err)
	}
	defer ln.Close()
	s.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("unable to accept connection", err)
			continue
		}
		fmt.Printf("connection accepted from %v\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("client %v closed the connection\n", c.RemoteAddr())
			} else {
				log.Printf("read error from %v: %v\n", c.RemoteAddr(), err)
			}
			return
		}
		fmt.Printf("message received from %v, %v \n", c.RemoteAddr(), string(buf[:n]))
	}
}
