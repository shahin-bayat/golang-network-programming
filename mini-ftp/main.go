package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/shahin-bayat/mini-ftp/networking"
)

func main() {
	mode := flag.String("mode", "", "mode: server or client (required)")
	host := flag.String("host", "", "host (required)")
	port := flag.Int("port", 0, "port (required)")
	file := flag.String("file", "", "file location")
	action := flag.String("action", "", "action: put or get (required)")
	flag.Parse()

	if (*mode != "server" && *mode != "client") ||
		*host == "" ||
		*port == 0 ||
		(*mode == "client" && (*action != "put" && *action != "get")) {

		flag.Usage()
		log.Fatalf("invalid arguments: mode=%q host=%q port=%d action=%q ", *mode, *host, *port, *action)
	}
	if *mode == "client" && *file == "" {
		flag.Usage()
		log.Fatalf("file is required for %q mode", *mode)
	}
	if *mode == "server" {
		fmt.Printf("Mode: %q, Listening on %q:%d\n", *mode, *host, *port)
		server := networking.NewServer("127.0.0.1", *port)
		server.Run()
	} else {
		fmt.Printf("Mode: %q, Connecting to %q:%d, File: %q, Action: %q\n", *mode, *host, *port, *file, *action)
	}
}
