package server

import (
	"fmt"

	"golang.org/x/net/dns/dnsmessage"
)

func (s *Server) parseQuery(query []byte) (*dnsmessage.Question, *dnsmessage.Header, error) {
	var p dnsmessage.Parser
	header, err := p.Start(query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start DNS parser: %w", err)
	}
	question, err := p.Question()
	if err != nil {
		if err != dnsmessage.ErrSectionDone {
			return nil, nil, fmt.Errorf("failed to get DNS question: %w", err)
		}
		return nil, nil, nil
	}
	return &question, &header, nil
}

func (s *Server) buildResponse(header *dnsmessage.Header, question *dnsmessage.Question, answers []dnsmessage.Resource) ([]byte, error) {
	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:       header.ID,
			Response: true,
			OpCode:   header.OpCode,
			RCode:    dnsmessage.RCodeSuccess,
		},
		Questions: []dnsmessage.Question{*question},
		Answers:   answers,
	}

	packed, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DNS message: %w", err)
	}
	return packed, nil
}

func (s *Server) buildErrorResponse(header *dnsmessage.Header, question *dnsmessage.Question, rcode dnsmessage.RCode) ([]byte, error) {
	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:       header.ID,
			Response: true,
			OpCode:   header.OpCode,
			RCode:    rcode,
		},
		Questions: []dnsmessage.Question{*question},
	}

	packed, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DNS message: %w", err)
	}
	return packed, nil
}
