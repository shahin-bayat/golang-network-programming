package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	addr string
	conn net.Conn
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Send(data []byte) error {
	n, err := c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("write %d bytes: %w", n, err)
	}
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func NewClient(ip string, port int) *Client {
	return &Client{
		addr: fmt.Sprintf("%s:%d", ip, port),
	}
}

func main() {
	ip := flag.String("ip", "", "Server IP address (required)")
	port := flag.Int("port", 0, "Server Port number (required)")
	flag.Parse()

	if *ip == "" || *port == 0 {
		flag.Usage()
		os.Exit(1)
	}
	client := NewClient(*ip, *port)
	if err := client.Connect(); err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer client.Close()

	msg := "Hello Server, how are you doing?"
	if err := client.Send([]byte(msg)); err != nil {
		log.Fatalf("send failed: %v", err)
	}

	buff := make([]byte, 2048)
	n, err := client.conn.Read(buff)
	if err != nil && err != io.EOF {
		log.Fatalf("Read failed: %v", err)
	}
	fmt.Printf("Received %d bytes: %s\n", n, string(buff[:n]))
}
