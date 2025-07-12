package networking

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Message struct {
	Action   string
	Filename string
	Size     int64
}

func (m *Message) String() string {
	return fmt.Sprintf("Action: %s, Filename: %s, Size: %d", m.Action, m.Filename, m.Size)
}

func (m *Message) Send(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(m)
}

func (m *Message) Receive(r io.Reader) error {
	dec := gob.NewDecoder(r)
	err := dec.Decode(m)
	if err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}
	return nil
}
