package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	quitch     chan struct{}
	bufferSize int
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.listener = ln
	go s.accept()

	<-s.quitch

	return nil
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		fmt.Println("new connection", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, s.bufferSize)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			switch {
			case errors.Is(err, io.EOF):
				fmt.Printf("client %s disconnected: %s", conn.RemoteAddr(), err)
				return
			default:
				fmt.Println("read error:", err)
				return // fatal error — stop
			}
		}
		msg := buff[:n]
		fmt.Println(string(msg))
		_, err = conn.Write([]byte("Hey Client, how are you dude?"))
		if err != nil {
			fmt.Println("write error:", err)
			return // fatal error — stop
		}
	}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		bufferSize: 2048,
	}
}

func main() {
	server := NewServer(":8081")
	server.Start()
}
