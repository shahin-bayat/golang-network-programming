package main

import (
	"encoding/gob"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/shahin-bayat/mini-ftp/config"
	"github.com/shahin-bayat/mini-ftp/networking"
)

func init() {
	gob.Register(&networking.Message{})
	opts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
		NoColor:    false,
	}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	cfg, err := config.ParseAndValidateConfig()
	if err != nil {
		slog.Error("Error parsing config", "error", err)
		os.Exit(1)
	}
	mode, host, port, file, action := cfg.Mode,
		cfg.Host,
		cfg.Port,
		cfg.File,
		cfg.Action

	if mode == "server" {
		slog.Info("Starting server", "host", host, "port", port)
		server := networking.NewServer("", port)
		server.Run()
	} else {
		slog.Info("Starting client", "mode", mode, "remote host", host, "remote port", port, "file", file, "action", action)
		client := networking.NewClient(host, port)
		defer client.Close()
		if err := client.Connect(); err != nil {
			slog.Error("Error connecting to server", "error", err)
			os.Exit(1)
		}
		f, err := os.Open(file)
		if err != nil {
			slog.Error("Error opening file", "file", file, "error", err)
			os.Exit(1)
		}
		if err = client.Send(f); err != nil {
			slog.Error("Error sending file", "file", file, "error", err)
			os.Exit(1)
		}
		if err = client.Close(); err != nil {
			slog.Error("Error closing client connection", "error", err)
			os.Exit(1)
		}
		slog.Info("File transfer completed successfully", "file", file, "action", action)
	}
}
