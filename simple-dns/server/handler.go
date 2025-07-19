package server

import (
	"fmt"
	"log/slog"

	"golang.org/x/net/dns/dnsmessage"
)

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
