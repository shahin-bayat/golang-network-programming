package networking

import (
	"fmt"
	"io"
	"log"
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
	filename := filepath.Base(file.Name())
	cmd := fmt.Sprintf("%s %s\n", "put", filename)
	_, err := c.conn.Write([]byte(cmd))
	if err != nil {
		log.Printf("unable to send command: %s", err)
		return err
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
