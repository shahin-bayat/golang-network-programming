package config

import "fmt"

type ServerConfig struct {
	Host string
	Port int
}

func (c *ServerConfig) Validate() error {
	if c.Port == 0 {
		return fmt.Errorf("port is required for server")
	}
	return nil
}

type ClientConfig struct {
	Host string
	Port int
	User string
}

func (c *ClientConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required for client")
	}
	if c.Port == 0 {
		return fmt.Errorf("port is required for client")
	}
	if c.User == "" {
		return fmt.Errorf("user is required for client")
	}
	return nil
}
