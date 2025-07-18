package networking

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

type Server struct {
	host             string
	port             int
	backends         []string
	nextBackendIndex int

	listener net.Listener
	mu       sync.Mutex
}

func NewServer(host string, port int, backends []string) *Server {
	return &Server{
		host:     host,
		port:     port,
		backends: backends,
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		return
	}
	defer ln.Close()
	slog.Info("server running", "host", s.host, "port", s.port)

	s.listener = ln
	s.accept()
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				slog.Warn("temporary accept error", "error", err)
				// retry
				time.Sleep(1 * time.Second)
				continue
			}
			slog.Error("failed to accept connection", "error", err)
			break
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()
	slog.Info("accepted connection", "from", c.RemoteAddr())

	nextBackend := s.nextBackend()
	backendConn, err := net.DialTimeout("tcp", nextBackend, 500*time.Millisecond)
	if err != nil {
		slog.Error("backend server is down", "error", err)
		return
	}
	defer backendConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go s.forwardTraffic(c, backendConn, &wg)
	go s.forwardTraffic(backendConn, c, &wg)

	wg.Wait()
}

func (s *Server) nextBackend() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	backend := s.backends[s.nextBackendIndex]
	// round robin
	s.nextBackendIndex = (s.nextBackendIndex + 1) % len(s.backends)

	return backend
}

func (s *Server) forwardTraffic(src io.Reader, dest io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(dest, src)
	if err != nil {
		slog.Error("failed to transfer traffic", "from", src, "to", dest)
		return
	}
}
