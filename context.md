# Network Programming in Go: Learning Roadmap

This roadmap is designed to take you from the fundamentals of network programming in Go to advanced topics like building a VPN.

## Basic Roadmap

This section covers the foundational concepts and common protocols.

*   [x] 1. **Simple Port Scanner:** Probe a range of TCP ports to report “open” vs. “closed.”
*   [x] 2. **File-Transfer Utility (“mini-FTP”):** Implement a simple GET/PUT protocol over TCP.
*   [x] 3. **In-Memory Key-Value Store (“mini-Redis”):** Build a server for SET/GET/DEL commands, using a mutex for concurrent access.
*   [x] 4. **Group Chat Server & Client:** Broadcast messages from each client to all connected peers in real-time.
*   [x] 5. **Bare-Bones HTTP Server:** Parse HTTP requests and serve static files or dynamic content.
    *   [x] **Step 1: Project Setup:** Created the `http-server` directory and initialized the Go module.
    *   [x] **Step 2: Basic TCP Server:** Implemented the initial TCP listener and connection handling loop.
    *   [x] **Step 3: Parse Request Line:** Read from the connection and parsed the HTTP method, path, and protocol.
    *   [x] **Step 4: Send Hardcoded Response:** Sent a valid `HTTP/1.1 200 OK` response with a "Hello, World!" body.
    *   [x] **Step 5: Serve Specific Static File:** Implemented logic to serve `www/index.html` for `/` and `/index.html` paths.
    *   [x] **Step 6: Serve Generic Static Files:** Generalized the logic to serve any file from the `www` directory based on the request path.
    *   [ ] **Step 7: Improve Security:** Implement a more robust path traversal check using absolute paths to prevent access to files outside the `www` directory.
*   [ ] 6. **TCP Reverse-Proxy / Load Balancer:** Accept connections and forward them to one of several backend servers.
*   [ ] 7. **Simple DNS Server (UDP):** Respond to basic A-record queries using the UDP protocol.
*   [ ] 8. **Port-Forwarding Tunnel (like `ssh -L`):** Forward a local port through a TCP connection to a remote host/port.
*   [ ] 9. **Simple TLS Terminating Proxy:** Accept TLS-encrypted TCP, decrypt it, and proxy the clear-text to a backend.
*   [ ] 10. **Custom RPC Framework:** Design a remote procedure call framework using `gob` or Protocol Buffers.

## Advanced Roadmap: Building a VPN

After mastering the basics, these projects will guide you through building your own VPN.

*   [ ] 11. **TUN Device "Hello, World!":** Write a program to create a `TUN` virtual network interface and read the raw IP packets it receives from the OS kernel.
*   [ ] 12. **Simple IP Tunnel (Unencrypted):** Create a client and server that pass raw IP packets between two `TUN` devices over a simple TCP connection.
*   [ ] 13. **Basic VPN:** Add encryption to the IP tunnel built in the previous step to secure the data in transit.

---

## Assistant Instructions

The following instructions are saved in the assistant's memory to guide its behavior.

1.  **+++ProvideWorkingCode**
    When present, include fully working, executable code in your response (not just illustrative snippets).
    Scope: Message-scoped.
2.  **+++Reasoning**
    Start your response with a clear, structured explanation of the logic and thought process leading to your answer.
    Scope: Message-scoped.
3.  **+++StepByStep**
    Break your answer into explicitly numbered steps (e.g. [Step 1] → [Step 2] → … → [Final Step]).
    Scope: Message-scoped.
4.  **+++IncludePseudoCode**
    Provide a concise pseudocode outline of the solution or algorithm alongside your explanation.
    Scope: Message-scoped.