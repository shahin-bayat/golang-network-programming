package server

import (
	"fmt"
	"log/slog"
	"net"
	"simple-dns/resolver"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Server struct {
	host     string
	port     int
	resolver *resolver.Resolver
}

func NewServer(host string, port int) *Server {
	resolver := resolver.NewResolver()
	return &Server{
		host:     host,
		port:     port,
		resolver: resolver,
	}
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		slog.Error("failed to resolve address", "address", addr, "error", err)
		return
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		slog.Error("failed to start server", "error", err)
		return
	}
	defer udpConn.Close()
	slog.Info("server listening", "address", laddr)
	s.read(udpConn)
}

func (s *Server) read(c *net.UDPConn) {
	buf := make([]byte, 512)
	for {
		n, remoteAddr, err := c.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				slog.Error("timeout error", "error", err)
				time.Sleep(1 * time.Second)
				continue
			}
			slog.Error("fatal read error", "error", err)
			break
		}

		query := buf[:n]
		question, header, err := s.parseQuery(query)
		if err != nil {
			slog.Error("failed to parse DNS query", "error", err)
			continue
		}
		if question == nil {
			slog.Info("no question found in the DNS query, skipping")
			continue
		}

		if question.Type != dnsmessage.TypeA {
			errResponse, err := s.buildErrorResponse(header, question, dnsmessage.RCodeNotImplemented)
			if err != nil {
				slog.Error("failed to build DNS response", "error", err)
				continue
			}
			_, err = c.WriteToUDP(errResponse, remoteAddr)
			if err != nil {
				slog.Error("failed to send DNS response", "error", err)
				continue
			}
		} else {
			answers, err := s.resolver.Resolve(question.Name, question.Type)
			if err != nil {
				errResponse, err := s.buildErrorResponse(header, question, dnsmessage.RCodeNameError)
				if err != nil {
					slog.Error("failed to build DNS response", "error", err)
					continue
				}
				_, err = c.WriteToUDP(errResponse, remoteAddr)
				if err != nil {
					slog.Error("failed to send DNS response", "error", err)
					continue
				}
				continue
			}

			response, err := s.buildResponse(header, question, answers)
			if err != nil {
				slog.Error("failed to build DNS response", "error", err)
				continue
			}
			_, err = c.WriteToUDP(response, remoteAddr)
			if err != nil {
				slog.Error("failed to send DNS response", "error", err)
				continue
			}
		}
	}
}
