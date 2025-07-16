package main

import (
	"encoding/gob"

	cmd "github.com/shahin-bayat/mini-chat/cmd/minichat"
	"github.com/shahin-bayat/mini-chat/networking"
)

func main() {
	cmd.Execute()
}

func init() {
	gob.Register(networking.Handshake{})
}
