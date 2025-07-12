package networking

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
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
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	var msg Message
	if err := msg.Receive(c); err != nil {
		slog.Error("error receiving message", "remote_addr", c.RemoteAddr(), "error", err)
	}
	action := msg.Action
	filename := msg.Filename
	fileSize := msg.Size

	switch action {
	case "put":
		if err := handlePut(c, msg); err != nil {
			slog.Error("error handling put command", "filename", filename, "file size", fileSize, "error", err)
			c.Write([]byte("ERR " + err.Error() + "\n"))
		} else {
			slog.Info("file received", "filename", filename, "remote_addr", c.RemoteAddr())
			c.Write([]byte("OK\n"))
		}
	case "get":
		if err := handleGet(c, msg); err != nil {
			c.Write([]byte("ERR " + err.Error() + "\n"))
		}
	default:
		c.Write([]byte("ERR unknown action\n"))
	}
}

func handlePut(reader io.Reader, message Message) error {
	if err := os.MkdirAll("files", 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fullpath := filepath.Join("files", message.Filename)

	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("could not create file %q: %w", fullpath, err)
	}
	defer file.Close()

	if _, err := io.CopyN(file, reader, message.Size); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", fullpath, err)
	}
	return nil
}

func handleGet(writer io.Writer, message Message) error {
	fullpath := filepath.Join("files", message.Filename)
	fmt.Println("fullpath", fullpath)
	if _, err := os.Stat(fullpath); err != nil {
		return fmt.Errorf("failed to locate file %q: %w", fullpath, err)
	}
	file, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", fullpath, err)
	}
	defer file.Close()
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("failed to write from %q: %w", fullpath, err)
	}
	return nil
}
