package networking

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
)

type Server struct {
	host  string
	port  int
	store map[string]string

	listener net.Listener
	mu       sync.RWMutex
}

type Command struct {
	op    string
	key   string
	value string
}

func NewServer(host string, port int) *Server {
	return &Server{
		host:  host,
		port:  port,
		store: make(map[string]string),
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("unable to start server", "error", err)
		os.Exit(1)
	}
	defer ln.Close()
	s.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("unable to accept connection", "error", err)
			continue
		}

		slog.Info("Connection accepted", "remote_addr", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		ln, err := reader.ReadString('\n')
		if err != nil {
			c.Write([]byte("ERR " + err.Error() + "\n"))
			break
		}
		cmd, err := parseCommand(ln)
		if err != nil {
			c.Write([]byte("ERR " + err.Error() + "\n"))
			continue
		}
		if err := s.handleCommand(cmd, c); err != nil {
			c.Write([]byte("ERR " + err.Error() + "\n"))
			continue
		}
	}
}

func (s *Server) handleCommand(cmd *Command, c net.Conn) error {
	switch cmd.op {
	case "SET":
		s.mu.Lock()
		s.store[cmd.key] = cmd.value
		s.mu.Unlock()
		c.Write([]byte("OK\n"))

	case "GET":
		s.mu.RLock()
		v, ok := s.store[cmd.key]
		s.mu.RUnlock()
		if !ok {
			c.Write([]byte("ERR no such key\n"))
		} else {
			c.Write([]byte(v + "\n"))
		}

	case "DEL":
		s.mu.Lock()
		_, ok := s.store[cmd.key]
		if ok {
			delete(s.store, cmd.key)
		}
		s.mu.Unlock()
		if ok {
			c.Write([]byte("OK\n"))
		} else {
			c.Write([]byte("ERR no such key\n"))
		}
	default:
		return fmt.Errorf("unhandled op %s", cmd.op)
	}
	return nil
}

func parseCommand(ln string) (*Command, error) {
	ln = strings.TrimSpace(ln)
	parts := strings.SplitN(ln, " ", 3)
	op := strings.ToLower(parts[0])
	switch op {
	case "set":
		if len(parts) != 3 {
			return nil, fmt.Errorf("SET requires key and value")
		}
		return &Command{op: "SET", key: parts[1], value: parts[2]}, nil

	case "get", "del":
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s requires exactly one key", strings.ToUpper(op))
		}
		return &Command{op: strings.ToUpper(op), key: parts[1]}, nil

	default:
		return nil, fmt.Errorf("unsupported operation %q", parts[0])
	}
}
