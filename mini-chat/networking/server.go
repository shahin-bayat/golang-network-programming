package networking

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Server struct {
	Host     string
	Port     int
	listener net.Listener
}

func NewServer(host string, port int) *Server {
	return &Server{
		Host: host,
		Port: port,
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
	defer ln.Close()
	s.listener = ln
	s.acceptConnection() // should be blocking, otherwise this function will return
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("stopping server; listener closed")
				break
			}
			slog.Error("failed to accept", "error", err)
			continue
		}
		slog.Info("client connected", "remote", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	_, err := c.Write([]byte("Welcome to minichat application\n"))
	if err != nil {
		slog.Error("failed to send welcome message", "error", err)
		return
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("client disconnected", "remote", c.RemoteAddr(), "error", err)
			} else {
				slog.Error("read error from client", "error", err)
			}
			return // drop this client
		}
		fmt.Println(strings.TrimSpace(line))
	}
}
