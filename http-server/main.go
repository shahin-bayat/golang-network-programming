package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	host     string
	port     int
	listener net.Listener
}

type Config struct {
	host string
	port int
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		return
	}
	defer ln.Close()
	s.listener = ln
	s.accept()
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			slog.Error("failed to accept", "error", err)
			continue
		}
		slog.Info("connection accepted", "remote address", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		slog.Error("failed to read message", "error", err)
		c.Write([]byte("ERR: failed to read message\n"))
		return
	}
	msg := string(buf[:n])
	lines := strings.Split(msg, "\n")

	httpReq := strings.TrimSpace(lines[0])
	parts := strings.SplitN(httpReq, " ", 3)
	if len(parts) < 3 {
		c.Write([]byte("ERR: invalid request format\n"))
		slog.Error("invalid request format", "line", httpReq)
		return
	}
	method, path, protocol := parts[0], parts[1], parts[2]
	slog.Info("received request", "method", method, "path", path, "protocol", protocol)
	if method != "GET" {
		c.Write([]byte("ERR: unsupported method\n"))
		slog.Error("unsupported method", "method", method)
		return
	}
	path = strings.TrimPrefix(path, "/")
	if strings.Contains(path, "..") {
		c.Write([]byte("ERR: invalid path\n"))
		slog.Error("invalid path", "path", path)
		return
	}
	if path == "" {
		path = "index.html"
	}
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	header := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\n\r\n", mimeType)
	file, err := os.Open(filepath.Join("www", path))
	if err != nil {
		slog.Error("failed to open file", "path", path, "error", err)
		c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\nNot Found"))
		return
	}
	defer file.Close()
	_, err = c.Write([]byte(header))
	if err != nil {
		slog.Error("failed to write response header", "error", err)
		return
	}
	_, err = io.Copy(c, file)
	if err != nil {
		slog.Error("failed to write response body", "error", err)
		c.Write([]byte("ERR: failed to write response body\n"))
		return
	}
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.host, "host", "", "Host to bind the server to")
	flag.IntVar(&cfg.port, "port", 0, "Port to bind the server to")
	flag.Parse()

	if cfg.port == 0 {
		flag.Usage()
		os.Exit(1)
	}
	server := &Server{
		host: cfg.host,
		port: cfg.port,
	}
	server.Run()
}
