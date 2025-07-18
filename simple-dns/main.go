package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	host string
	port int
}

func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		slog.Error("failed to resolve address", "address", addr, "error", err)
		return
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		return
	}
	defer udpConn.Close()
	slog.Info("server listening", "address", laddr)
	s.read(udpConn)
}

func (s *Server) read(c *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := c.ReadFromUDP(buf)
		if err != nil {
			slog.Error("error read from udp connection", "error", err)
			break
		}

		slog.Info("Received packet", "from", remoteAddr.String(), "data", string(buf[:n]))
	}
}

func main() {
	server := NewServer("", 53)
	server.Run()
}
