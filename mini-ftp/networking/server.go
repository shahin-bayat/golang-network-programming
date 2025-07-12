package networking

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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

	reader := bufio.NewReader(c)
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Printf("client %v closed the connection\n", c.RemoteAddr())
		} else {
			log.Printf("read error from %v: %v\n", c.RemoteAddr(), err)
		}
		return
	}
	cmd := strings.TrimSpace(line)
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		c.Write([]byte("ERR invalid command\n"))
		return
	}
	action := parts[0]
	filename := parts[1]

	fmt.Printf("received action=%q, filename=%q from %v\n", action, filename, c.RemoteAddr())

	switch action {
	case "put":
		if err := handlePut(reader, filename); err != nil {
			c.Write([]byte("ERR " + err.Error() + "\n"))
		} else {
			c.Write([]byte("OK\n"))
		}
	case "get":
		if err := handleGet(c, filename); err != nil {
			c.Write([]byte("ERR " + err.Error() + "\n"))
		}
	default:
		c.Write([]byte("ERR unknown action\n"))
	}
}

func handlePut(reader *bufio.Reader, filename string) error {
	if err := os.MkdirAll("files", 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fullpath := filepath.Join("files", filename)

	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("could not create file %q: %w", fullpath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write to file %q: %w", fullpath, err)
	}
	return nil
}

func handleGet(writer io.Writer, filename string) error {
	fullpath := filepath.Join("files", filename)
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
