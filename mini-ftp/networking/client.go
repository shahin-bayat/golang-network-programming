package networking

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	remoteHost string
	remotePort int
	conn       net.Conn
}

func NewClient(remoteHost string, remotePort int) *Client {
	return &Client{
		remoteHost: remoteHost,
		remotePort: remotePort,
	}
}

func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.remoteHost, c.remotePort)
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Send(file *os.File) error {
	fileInfo, err := file.Stat()
	if err != nil {
		slog.Error("unable to get file size", "file", file.Name(), "error", err)
	}
	msg := &Message{
		Action:   "put",
		Filename: filepath.Base(file.Name()),
		Size:     fileInfo.Size(),
	}
	if err := msg.Send(c.conn); err != nil {
		slog.Error("unable to send message", "error", err)
		os.Exit(1)
	}
	_, err = io.Copy(c.conn, file)
	if err != nil {
		log.Printf("unable to send file: %s", err)
		return err
	}
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
