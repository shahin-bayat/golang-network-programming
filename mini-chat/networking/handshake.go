package networking

import (
	"encoding/gob"
	"fmt"
	"net"
)

type Handshake struct {
	User string
}

func (h *Handshake) Serialize(c net.Conn) error {
	return gob.NewEncoder(c).Encode(h)
}

func (h *Handshake) Deserialize(c net.Conn) error {
	if err := gob.NewDecoder(c).Decode(h); err != nil {
		return err
	}
	if h.User == "" {
		return fmt.Errorf("handshake user cannot be empty")
	}
	return nil
}
