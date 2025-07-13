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

func (c *Client) Send(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %w", filename, err)
	}
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

	serverResponse := make([]byte, 1024)
	n, err := c.conn.Read(serverResponse)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading server response: %w", err)
	}
	slog.Info("Server response", "response", string(serverResponse[:n]))
	return nil
}

func (c *Client) Receive(filename string) error {
	// 1. create a get request message
	// 2. send message to server
	// 3. receive get response message from server
	// 4. create 'received' directory if not exist
	// 5. create a file with filename received from get response message
	// 6. save file
	reqMsg := &Message{
		Action:   "get",
		Filename: filepath.Base(filename),
	}
	if err := reqMsg.Send(c.conn); err != nil {
		return fmt.Errorf("failed to receive message: %s", err)
	}
	var resMsg Message
	if err := resMsg.Receive(c.conn); err != nil {
		return fmt.Errorf("failed to receive response message from server %s", err)
	}
	if err := os.MkdirAll("received", 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s", err)
	}

	fullpath := filepath.Join("received", filepath.Base(resMsg.Filename))
	file, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer file.Close()

	if _, err := io.CopyN(file, c.conn, resMsg.Size); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", fullpath, err)
	}
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
