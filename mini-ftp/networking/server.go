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

	switch msg.Action {
	case "put":
		if err := handlePut(c, msg); err != nil {
			c.Write([]byte("ERR " + err.Error()))
		} else {
			c.Write([]byte("OK"))
		}
	case "get":
		if err := handleGet(c, msg); err != nil {
			slog.Error("failed to send file", "error", err.Error())
			c.Write([]byte("ERR " + err.Error() + "\n"))
		}
	default:
		c.Write([]byte("ERR unknown action\n"))
	}
}

func handlePut(r io.Reader, message Message) error {
	// 1. get Message from client
	// 2. create 'files' directory if not exist
	// 3. save file to 'files' directory
	if err := os.MkdirAll("files", 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	fullpath := filepath.Join("files", message.Filename)

	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %s", fullpath, err)
	}
	defer file.Close()

	if _, err := io.CopyN(file, r, message.Size); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", fullpath, err)
	}
	return nil
}

func handleGet(w io.Writer, message Message) error {
	// 1. get Message from client
	// 2. check if filename inside Message exist inside 'files' directory
	// 3. create a Message with file info
	// 4. send Message to client
	// 5. send file to client
	fullpath := filepath.Join("files", message.Filename)

	fileInfo, err := os.Stat(fullpath)
	if err != nil {
		return fmt.Errorf("failed to locate file %s: %s", fullpath, err)
	}
	file, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", fullpath, err)
	}
	defer file.Close()

	msg := &Message{
		Action:   "get",
		Filename: filepath.Base(file.Name()),
		Size:     fileInfo.Size(),
	}
	if err := msg.Send(w); err != nil {
		return fmt.Errorf("unable to send message: %s", err)
	}
	if _, err := io.CopyN(w, file, fileInfo.Size()); err != nil {
		return fmt.Errorf("failed to write from %s: %s", fullpath, err)
	}
	return nil
}
