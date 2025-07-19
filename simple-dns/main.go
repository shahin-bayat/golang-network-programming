package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Server struct {
	host string
	port int
}

func NewServer(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
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
		question, header, err := s.parseQuery(query)
		if err != nil {
			slog.Error("failed to parse DNS query", "error", err)
			continue
		}
		if question == nil {
			slog.Info("no question found in the DNS query, skipping")
			continue
		}
		slog.Info("Parsed question", "name", question.Name, "type", question.Type)

		targetDomain, err := dnsmessage.NewName("test.local.com.")
		if err != nil {
			slog.Error("failed to create target dns name", "error", err)
		}
		hardcodedIP := [4]byte{192, 168, 1, 1}

		if question.Type == dnsmessage.TypeA && question.Name.String() == targetDomain.String() {
			slog.Info("matched query for test.local. A record. Building response...")
			response, err := s.buildResponse(header.ID, question, hardcodedIP)
			if err != nil {
				slog.Error("failed to build DNS response", "error", err)
			}
			_, err = c.WriteToUDP(response, remoteAddr)
			if err != nil {
				slog.Error("failed to send DNS response", "error", err)
			}
		} else {
			slog.Info("query does not match test.local. A record, ignoring.", "name", question.Name, "type", question.Type, "target domain", targetDomain)
		}
	}
}

func (s *Server) parseQuery(query []byte) (*dnsmessage.Question, *dnsmessage.Header, error) {
	var p dnsmessage.Parser
	header, err := p.Start(query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start DNS parser: %w", err)
	}
	// for { // if you want to get all questions
	question, err := p.Question()
	if err != nil {
		if err != dnsmessage.ErrSectionDone {
			return nil, nil, fmt.Errorf("failed to get DNS question: %w", err)
		}
		slog.Info("no more questions in the DNS message")
		return nil, nil, nil
	}
	return &question, &header, nil
}

func (s *Server) buildResponse(queryID uint16, question *dnsmessage.Question, ip [4]byte) ([]byte, error) {
	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:       queryID,
			Response: true,
			OpCode:   dnsmessage.OpCode(0),
			RCode:    dnsmessage.RCodeSuccess,
		},
		Questions: []dnsmessage.Question{*question},
		Answers: []dnsmessage.Resource{
			{
				Header: dnsmessage.ResourceHeader{
					Name:  question.Name,
					Type:  dnsmessage.TypeA,
					Class: dnsmessage.ClassINET,
				},
				Body: &dnsmessage.AResource{A: ip},
			},
		},
	}

	// Pack the message into a byte slice
	packed, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DNS message: %w", err)
	}
	slog.Info("Built DNS response", "queryID", queryID, "questionName", question.Name, "ip", ip)
	return packed, nil
}

func main() {
	server := NewServer("", 53)
	server.Run()
}
