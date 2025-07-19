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
		s.parseQuery(query, remoteAddr)
	}
}

func (s *Server) parseQuery(query []byte, remoteAddr *net.UDPAddr) {
	var p dnsmessage.Parser
	h, err := p.Start(query)
	if err != nil {
		slog.Error("error read dns message", "error", err)
		return
	}
	// for {
	question, err := p.Question()
	if err != nil {

		if err != dnsmessage.ErrSectionDone {
			slog.Error("failed to get dns question", "error", err)
			break
		}
		return
	}
	slog.Info("question:", "type", question.Type, "name", question.Name)

	targetDomain, err := dnsmessage.NewName("test.local.com.")
	if err != nil {
		slog.Error("failed to create target dns name", "error", err)
	}
	hardcodedIP := [4]byte{192, 168, 1, 1}

	if question.Type == dnsmessage.TypeA && question.Name.String() == targetDomain.String() {
		slog.Info("matched query for test.local. A record. Building response...")
		_, err = s.buildResponse(h.ID, question, hardcodedIP)
		if err != nil {
			slog.Error("failed to build DNS response", "error", err)
		}
		// TODO: send response to client
	} else {
		slog.Info("query does not match test.local. A record, ignoring.", "name", question.Name, "type", question.Type, "target domain", targetDomain)
	}
	// }
}

func (s *Server) buildResponse(queryID uint16, question dnsmessage.Question, ip [4]byte) ([]byte, error) {
	header := dnsmessage.Header{
		ID:       queryID,
		Response: true,
		OpCode:   dnsmessage.OpCode(0),
		RCode:    dnsmessage.RCodeSuccess,
	}

	resourceHeader := dnsmessage.ResourceHeader{
		Name:  question.Name,
		Type:  question.Type,
		Class: question.Class,
	}

	res := make([]byte, 1024)

	b := dnsmessage.NewBuilder(res, header)
	b.EnableCompression()
	b.AResource(resourceHeader, dnsmessage.AResource{A: ip})
	return b.Finish()
}

func main() {
	server := NewServer("", 53)
	server.Run()
}
