package server

import (
	"fmt"
	"log/slog"
	"net"
	"simple-dns/records"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Server struct {
	host        string
	port        int
	recordStore records.RecordStore
}

func NewServer(host string, port int, recordStore records.RecordStore) *Server {
	return &Server{
		host:        host,
		port:        port,
		recordStore: recordStore,
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
	buf := make([]byte, 1024)
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

		slog.Info("Received packet", "from", remoteAddr.String())

		query := buf[:n]
		question, header, err := s.parseQuery(query) // Call parseQuery method
		if err != nil {
			slog.Error("failed to parse DNS query", "error", err)
			continue
		}
		if question == nil {
			slog.Info("no question found in the DNS query, skipping")
			continue
		}
		slog.Info("Parsed question", "name", question.Name, "type", question.Type)

		if question.Type != dnsmessage.TypeA {
			slog.Info("dns record not supported", "type", question.Type)
			// In a real server, you might send a NOTIMP (Not Implemented) response
			continue
		}

		// Look up the IP in the recordStore
		ip, err := s.recordStore.Get(question.Name.String())
		if err != nil {
			slog.Info("record not found", "domain", question.Name)
			// In a real server, you might send an NXDOMAIN (Non-Existent Domain) response
			continue
		}

		response, err := s.buildResponse(header.ID, question, ip) // Call buildResponse method
		if err != nil {
			slog.Error("failed to build DNS response", "error", err)
			// Consider if you want to break or continue here
			continue
		}
		_, err = c.WriteToUDP(response, remoteAddr)
		if err != nil {
			slog.Error("failed to send DNS response", "error", err)
			// Consider if you want to break or continue here
			continue
		}
	}
}
