package networking

import (
	"fmt"
	"strings"
)

// Command represents a client request, either joining a room or sending a message.
type Command struct {
	Verb string // "JOIN" or "MSG"
	Room string
	Text string // only set for MSG
}

type ServerHandler interface {
	Parse(line string) (*Command, error)
}

type serverHandler struct{}

func NewServerHandler() ServerHandler {
	return &serverHandler{}
}

func (h *serverHandler) Parse(line string) (*Command, error) {
	parts := strings.SplitN(line, " ", 3)

	verb := strings.ToUpper(parts[0])
	switch verb {
	case "JOIN":
		// Expect: JOIN <room>
		if len(parts) != 2 {
			return nil, fmt.Errorf("JOIN requires exactly one argument: room name")
		}
		room := strings.TrimSpace(strings.ToLower(parts[1]))
		return &Command{Verb: verb, Room: room}, nil

	case "MSG":
		// Expect: MSG <room> <message>
		if len(parts) != 3 {
			return nil, fmt.Errorf("MSG requires a room and a message")
		}
		room := strings.TrimSpace(strings.ToLower(parts[1]))
		return &Command{Verb: verb, Room: room, Text: parts[2]}, nil
	case "QUIT":
		if len(parts) != 1 {
			return nil, fmt.Errorf("QUIT does not require any arguments")
		}
		return &Command{Verb: verb}, nil
	case "LEAVE":
		if len(parts) != 2 {
			return nil, fmt.Errorf("LEAVE requires exactly one argument: room name")
		}
		room := strings.TrimSpace(strings.ToLower(parts[1]))
		return &Command{Verb: verb, Room: room}, nil
	case "LIST":
		if len(parts) != 1 {
			return nil, fmt.Errorf("LIST does not require any arguments")
		}
		return &Command{Verb: verb}, nil
	default:
		return nil, fmt.Errorf("unsupported command: %s", verb)
	}
}
