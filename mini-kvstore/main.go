package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/shahin-bayat/mini-kvstore/networking"
)

type Config struct {
	host string
	port int
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.host, "host", "", "host")
	flag.IntVar(&cfg.port, "port", 0, "port (required)")
	flag.Parse()

	if cfg.port == 0 {
		slog.Error("port is required")
		os.Exit(1)
	}
	slog.Info("Starting server", "host", cfg.host, "port", cfg.port)
	server := networking.NewServer(cfg.host, cfg.port)
	server.Run()
}
