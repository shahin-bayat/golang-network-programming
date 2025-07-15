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
	"sync"
)

type Server struct {
	host     string
	port     int
	listener net.Listener
	rooms    map[string]map[net.Conn]struct{} // map[net.Conn]struct{} is idiomatic way of defining set in go
	mu       sync.RWMutex
	handler  ServerHandler
}

func NewServer(host string, port int) *Server {
	rooms := make(map[string]map[net.Conn]struct{})
	return &Server{
		host:    host,
		port:    port,
		rooms:   rooms,
		handler: NewServerHandler(),
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
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
	defer func() {
		c.Close()
		slog.Info("client disconnected", "remote", c.RemoteAddr())
		s.disconnect(c)
	}()

	r := bufio.NewReader(c)
	_, err := c.Write([]byte("Welcome to minichat application\n"))
	if err != nil {
		slog.Error("failed to send welcome message", "error", err)
		return // when retuns, connection will be closed to this client because of defer on top of this function
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

		line = strings.TrimSpace(line)
		if line == "" {
			continue // skip empty lines
		}

		cmd, err := s.handler.Parse(line)
		if err != nil {
			fmt.Fprintf(c, "ERR: %s\n", err)
			continue // continue to read next command
		}

		switch cmd.Verb {
		case "JOIN":
			s.joinRoom(cmd.Room, c)
			fmt.Fprintf(c, "OK JOIN %s\n", cmd.Room)
		case "MSG":
			s.broadcast(cmd.Room, c.RemoteAddr(), cmd.Text)
		}
	}
}

func (s *Server) joinRoom(room string, c net.Conn) {
	s.mu.Lock()
	if s.rooms[room] == nil {
		s.rooms[room] = make(map[net.Conn]struct{})
	}
	s.rooms[room][c] = struct{}{}
	s.mu.Unlock()
	s.broadcast(room, c.RemoteAddr(), "has joined the room")
}

func (s *Server) broadcast(room string, from net.Addr, text string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peers, ok := s.rooms[room]
	if !ok {
		slog.Warn("broadcast to non-existing room", "room", room)
		return // no one in this room, nothing to broadcast
	}

	for peer := range peers {
		fmt.Fprintf(peer, "[%s] %s:%s\n", room, from, text)
	}
}

func (s *Server) disconnect(c net.Conn) {
	var leftRooms []string

	s.mu.Lock()
	for room, peers := range s.rooms {
		if _, ok := peers[c]; !ok {
			continue // this client is not in this room
		}
		leftRooms = append(leftRooms, room)
		delete(peers, c)
		if len(peers) == 0 {
			delete(s.rooms, room) // remove empty rooms
		}
	}
	s.mu.Unlock()

	for _, room := range leftRooms {
		s.broadcast(room, c.RemoteAddr(), "has left the room")
	}
}
