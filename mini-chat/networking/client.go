package networking

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	remoteHost string
	remotePort int
	user       string
	timeout    time.Duration
	done       chan struct{}
}

func NewClient(remoteHost string, remotePort int, user string) *Client {
	return &Client{
		remoteHost: remoteHost,
		remotePort: remotePort,
		user:       user,
		timeout:    time.Millisecond * 500,
		done:       make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.remoteHost, c.remotePort)
	conn, err := net.DialTimeout("tcp", addr, c.timeout)
	if err != nil {
		slog.Error("failed to connect remote host", "error", err)
		return err
	}
	defer conn.Close()

	h := Handshake{
		User: c.user,
	}
	if err := h.Serialize(conn); err != nil {
		slog.Error("handshake failed", "error", err)
		close(c.done)
	}

	go c.read(conn)
	go c.write(conn)

	<-c.done
	return nil
}

func (c *Client) read(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		slog.Error("read error", "error", err)
	}
	close(c.done)
}

func (c *Client) write(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		// check if the read goroutine has finished
		select {
		case <-c.done:
			return
		default:
		}
		// send the line to the server
		if _, err := fmt.Fprintln(conn, line); err != nil {
			slog.Error("write error", "error", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Error("stdin error", "error", err)
	}
}

// •	The read goroutine keeps scanning until the connection is closed by the peer (or an error occurs).
// •	As soon as scanner.Scan() returns false, you log any error and then call close(c.done).
// •	That unblocks the <-c.done in Connect(), so Connect returns.
// •	Meanwhile, the write goroutine is also watching c.done in its select; once c.done is closed it’ll exit its loop and stop trying to write.
