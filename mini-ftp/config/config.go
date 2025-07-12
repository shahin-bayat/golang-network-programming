package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Mode   string
	Host   string
	Port   int
	File   string
	Action string
}

func ParseAndValidateConfig() (*Config, error) {
	mode := flag.String("mode", "", "mode: server or client")
	host := flag.String("host", "", "host (required)")
	port := flag.Int("port", 0, "port (required)")
	file := flag.String("file", "", "file location")
	action := flag.String("action", "", "action: put or get (required)")
	flag.Parse()

	if *mode != "server" && *mode != "client" {
		return nil, fmt.Errorf("mode must be server or client")
	}
	if *port == 0 {
		return nil, fmt.Errorf("port is required")
	}
	if *mode == "client" {
		if *host == "" {
			return nil, fmt.Errorf("host is required in client mode")
		}
		if *action != "put" && *action != "get" {
			return nil, fmt.Errorf("action must be put or get in client mode")
		}
		if *file == "" {
			return nil, fmt.Errorf("file is required in client mode")
		}
	}
	if *mode == "server" && *host == "" {
		*host = "127.0.0.1"
	}
	config := &Config{
		Mode:   *mode,
		Host:   *host,
		Port:   *port,
		File:   *file,
		Action: *action,
	}
	return config, nil
}
