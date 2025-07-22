# Network Programming in Go: Learning Roadmap

This roadmap is designed to take you from the fundamentals of network programming in Go to advanced topics like building a VPN.

## Basic Roadmap

This section covers the foundational concepts and common protocols.

- [x] 1. **Simple Port Scanner:** Probe a range of TCP ports to report “open” vs. “closed.”
- [x] 2. **File-Transfer Utility (“mini-FTP”):** Implement a simple GET/PUT protocol over TCP.
- [x] 3. **In-Memory Key-Value Store (“mini-Redis”):** Build a server for SET/GET/DEL commands, using a mutex for concurrent access.
- [x] 4. **Group Chat Server & Client:** Broadcast messages from each client to all connected peers in real-time.
- [x] 5. **Bare-Bones HTTP Server:** Parse HTTP requests and serve static files or dynamic content.
  - [x] **Step 1: Project Setup:** Created the `http-server` directory and initialized the Go module.
  - [x] **Step 2: Basic TCP Server:** Implemented the initial TCP listener and connection handling loop.
  - [x] **Step 3: Parse Request Line:** Read from the connection and parsed the HTTP method, path, and protocol.
  - [x] **Step 4: Send Hardcoded Response:** Sent a valid `HTTP/1.1 200 OK` response with a "Hello, World!" body.
  - [x] **Step 5: Serve Specific Static File:** Implemented logic to serve `www/index.html` for `/` and `/index.html` paths.
  - [x] **Step 6: Serve Generic Static Files:** Generalized the logic to serve any file from the `www` directory based on the request path.
- [x] 6. **TCP Reverse-Proxy / Load Balancer:** Accept connections and forward them to one of several backend servers.
  - [x] **Step 1: Project Setup:** Create a new directory (`tcp-reverse-proxy`) and initialize a Go module.
  - [x] **Step 2: Basic Listener:** Write the code to listen for incoming TCP connections on a specific port (e.g., `:8080`).
  - [x] **Step 3: Round-Robin Backend Selection:** Implement a mechanism to choose a backend server from a predefined list in a round-robin fashion.
  - [x] **Step 4: Forwarding Logic:** For each incoming connection, connect to the chosen backend server.
  - [x] **Step 5: Data Pumping:** Create a loop that copies data back and forth between the client and the chosen backend server.
- [x] 7. **Simple DNS Server (UDP):** Respond to basic A-record queries using the UDP protocol.
  - [x] **Step 1: Project Setup:** Create a new directory (`simple-dns`) and initialize a Go module.
  - [x] **Step 2: Basic UDP Listener:** Write the code to listen for incoming UDP packets on port 53.
  - [x] **Step 3: Parsing DNS Queries:** Use the `golang.org/x/net/dns/dnsmessage` library to parse the incoming byte buffer into a DNS message.
  - [x] **Step 4: Crafting a Response:** Create a DNS response message with a hardcoded IP address for a specific domain query (e.g., `test.local.`).
  - [x] **Step 5: Sending the Response:** Marshal the response message back into bytes and send it back to the client's address.
  - [x] **Step 6: (Optional) Dynamic Responses:** Create a map to store domain-to-IP mappings and respond dynamically based on the query.
  - [x] **Step 7: (Future) Iterative Resolution:**
    - [x] **Step 7.1: Root Hints:** Load a pre-defined list of root DNS server IP addresses.
    - [x] **Step 7.2: Iterative Query Function:** Implement a function to perform iterative queries, following referrals from root to TLD to authoritative servers.
    - [x] **Step 7.3: Caching:** Implement an in-memory cache for resolved DNS records to improve performance.
    - [x] **Step 7.4: Handling Other Record Types:** Extend parsing and response building to support AAAA, MX, NS, CNAME, etc.
    - [x] **Step 7.5: Integrate Resolver with Server:** Modify the UDP server handler to use the iterative resolver for any queries not found in its local records.
- [x] 8. **Port-Forwarding Tunnel (like `ssh -L`):** Forward a local port through a TCP connection to a remote host/port.
  - [x] **Step 8.1: Project Setup:** Create a new directory (`port-forwarder`) and initialize a Go module.
  - [x] **Step 8.2: Parse Command-Line Arguments:** Write code to accept the local and remote addresses from the command line.
  - [x] **Step 8.3: Implement TCP Listener:** Write the code to listen for incoming TCP connections on the specified local address.
  - [x] **Step 8.4: Implement Forwarding Logic:** For each incoming connection, connect to the specified remote address.
  - [x] **Step 8.5: Data Pumping:** Create two goroutines to copy data back and forth between the local and remote connections using `io.Copy`.
- [ ] 9. **Simple TLS Terminating Proxy:** Accept TLS-encrypted TCP, decrypt it, and proxy the clear-text to a backend.
  - [ ] **Step 9.1: Project Setup:** Create a new directory (`tls-proxy`) and initialize a Go module.
  - [ ] **Step 9.2: Generate Self-Signed Certificate:** Create a private key and a self-signed certificate for the server.
  - [ ] **Step 9.3: Implement TLS Listener:** Write the code to listen for incoming TLS connections using `tls.Listen`.
  - [ ] **Step 9.4: Implement Forwarding Logic:** For each incoming connection, open a plain TCP connection to a backend service.
  - [ ] **Step 9.5: Data Pumping:** Use `io.Copy` to move the decrypted data between the client and the backend.
- [ ] 10. **Custom RPC Framework:** Design a remote procedure call framework using `gob` or Protocol Buffers.
- [ ] **(Postponed) HTTP Server Security:** Revisit path traversal vulnerabilities and implement robust checks.

## Advanced Roadmap: Building a VPN

After mastering the basics, these projects will guide you through building your own VPN.

- [ ] 11. **TUN Device "Hello, World!":** Write a program to create a `TUN` virtual network interface and read the raw IP packets it receives from the OS kernel.
- [ ] 12. **Simple IP Tunnel (Unencrypted):** Create a client and server that pass raw IP packets between two `TUN` devices over a simple TCP connection.
- [ ] 13. **Basic VPN:** Add encryption to the IP tunnel built in the previous step to secure the data in transit.

---

## Assistant Instructions

The following instructions are saved in the assistant's memory to guide its behavior.

- **Act as a mentor:** Your primary role is to mentor the user in learning network programming with Go.
- **Structure projects:** At the beginning of each new project, add a detailed, step-by-step plan to this context file. Mark steps as complete (`[x]`) as the user finishes them.
- **Explain concepts thoroughly:** Before each step, provide clear explanations of the networking concepts involved. Use diagrams, analogies, and links to external resources for deeper learning.
- **Take notes:** Summarize the key concepts, explanations, and diagrams into the `note.md` file for the user's reference.
- **Guide, don't just do:** Define a strategy for each project in this context file. Ask the user to implement each step themselves. Review their code, and coach them to fix or improve it.
- **Provide code on request:** Only provide complete, working code snippets when the user explicitly asks for them after getting stuck.
- **Maintain the roadmap:** Keep the learning roadmap in this file updated to reflect the user's progress.
- **Focus on Architecture & Design Patterns:** For all future projects, prioritize good architectural design, maintainability, and scalability. Actively seek opportunities to apply suitable design patterns (e.g., Dependency Injection, Strategy, Builder) where they enhance code quality and understanding.

### Commit Guidelines

Use the [Conventional Commits](https://www.conventionalcommits.org/) specification.

**Format:** `type(scope): description`

*   **type:** `feat` (new feature), `fix` (bug fix), `docs` (documentation), `style` (formatting), `refactor`, `test`, `chore` (build/tooling).
*   **scope:** The project folder name (e.g., `mini-ftp`, `http-server`).
*   **description:** A concise, imperative-mood summary.

**Example:** `feat(mini-reverse-proxy): implement round-robin load balancing`

### Pre-defined Directives

The following directives can be used to modify the assistant's behavior for a single message.

1. **+++ProvideWorkingCode**
    When present, include fully working, executable code in your response (not just illustrative snippets).
    Scope: Message-scoped.
2. **+++Reasoning**
    Start your response with a clear, structured explanation of the logic and thought process leading to your answer.
    Scope: Message-scoped.
3. **+++StepByStep**
    Break your answer into explicitly numbered steps (e.g. [Step 1] → [Step 2] → … → [Final Step]).
    Scope: Message-scoped.
4. **+++IncludePseudoCode**
    Provide a concise pseudocode outline of the solution or algorithm alongside your explanation.
    Scope: Message-scoped.
