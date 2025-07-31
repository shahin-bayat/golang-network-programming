# Network Programming in Go: A Learning Journey

This repository documents my progress through a network programming learning roadmap using Go. Each directory contains a separate project that builds on fundamental concepts, starting from basic TCP/UDP clients and servers to more advanced applications like a custom RPC framework and a simple VPN.

## Basic Roadmap

This section covers the foundational concepts and common protocols.

- **1. Simple Port Scanner:** A tool to probe a range of TCP ports and report their status as "open" or "closed."
- **2. File-Transfer Utility ("mini-FTP"):** A simple implementation of a GET/PUT protocol over TCP for file transfers.
- **3. In-Memory Key-Value Store ("mini-Redis"):** A server that handles SET, GET, and DEL commands with mutex-based concurrency control.
- **4. Group Chat Server & Client:** A real-time chat application that broadcasts messages from each client to all connected peers.
- **5. Bare-Bones HTTP Server:** A simple HTTP server capable of parsing requests and serving static files.
- **6. TCP Reverse-Proxy / Load Balancer:** A reverse proxy that accepts TCP connections and forwards them to one of several backend servers in a round-robin fashion.
- **7. Simple DNS Server (UDP):** A basic DNS server that responds to A-record queries over UDP, with support for iterative resolution and caching.
- **8. Port-Forwarding Tunnel:** A tool that forwards a local port to a remote host/port, similar to `ssh -L`.
- **9. Simple TLS Terminating Proxy:** A proxy that accepts TLS-encrypted TCP connections, decrypts the traffic, and forwards it to a backend service.
- **10. TUN Device "Hello, World!":** A basic implementation to create a TUN device and read a packet.
- **11. Basic IP Tunnel (Unencrypted):** A simple IP tunnel that passes raw IP packets between two TUN devices over a TCP connection.
