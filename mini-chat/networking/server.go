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
	host       string
	port       int
	welcomeMsg string
	listener   net.Listener
	rooms      map[string]map[net.Conn]struct{} // map[net.Conn]struct{} is idiomatic way of defining set in go
	users      map[net.Conn]string
	mu         sync.RWMutex
	handler    ServerHandler
}

func NewServer(host string, port int) *Server {
	rooms := make(map[string]map[net.Conn]struct{})
	users := make(map[net.Conn]string)
	return &Server{
		host:       host,
		port:       port,
		welcomeMsg: "Welcome to the minichat application! Version 1.0.0\n",
		rooms:      rooms,
		users:      users,
		handler:    NewServerHandler(),
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

	var h Handshake
	if err := h.Deserialize(c); err != nil {
		fmt.Fprintf(c, "ERR: %s\n", err)
		slog.Error("handshake failed", "error", err, "remote", c.RemoteAddr())
		return // drop this client
	}

	if err := s.handshake(c, &h); err != nil {
		fmt.Fprintf(c, "ERR: %s\n", err)
		slog.Error("handshake error", "error", err, "remote", c.RemoteAddr())
		return // drop this client
	}

	r := bufio.NewReader(c)
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
			s.broadcast(cmd.Room, c, cmd.Text)
		case "LEAVE":
			fmt.Fprintf(c, "OK LEAVE %s\n", cmd.Room)
			s.disconnect(c)
		case "QUIT":
			fmt.Fprintln(c, "Goodbye!")
			return // close the connection
		case "LIST":
			s.mu.RLock()
			var rooms []string
			for room := range s.rooms {
				rooms = append(rooms, room)
			}
			s.mu.RUnlock()
			if len(rooms) == 0 {
				fmt.Fprintln(c, "No rooms available")
			} else {
				fmt.Fprintf(c, "Available rooms: %s\n", strings.Join(rooms, ", "))
			}
		default:
			fmt.Fprintf(c, "ERR: unsupported command %s\n", cmd.Verb)
			slog.Warn("unsupported command", "command", cmd.Verb, "from", c.RemoteAddr())
			continue // continue to read next command
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
	s.broadcast(room, c, "has joined the room")
}

func (s *Server) broadcast(room string, from net.Conn, text string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peers, ok := s.rooms[room]
	if !ok {
		slog.Warn("broadcast to non-existing room", "room", room)
		return // no one in this room, nothing to broadcast
	}
	user := s.users[from]
	for peer := range peers {
		fmt.Fprintf(peer, "[%s] %s:%s\n", room, user, text)
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
	delete(s.users, c)
	s.mu.Unlock()

	for _, room := range leftRooms {
		s.broadcast(room, c, "has left the room")
	}
}

func (s *Server) handshake(c net.Conn, h *Handshake) error {
	s.mu.RLock()
	for _, user := range s.users {
		if user == h.User {
			s.mu.RUnlock()
			return fmt.Errorf("user %s is already connected", h.User)
		}
	}
	s.mu.RUnlock()

	_, err := c.Write([]byte(s.welcomeMsg))
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}
	s.mu.Lock()
	s.users[c] = h.User
	s.mu.Unlock()
	fmt.Fprintf(c, "OK USER %s\n", h.User)
	return nil
}
