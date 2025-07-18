### Project 6: TCP Reverse-Proxy / Load Balancer

**Concept: What is a Reverse Proxy?**

A reverse proxy is a server that sits in front of one or more "backend" servers, forwarding client requests to them. To the outside world, it looks like a single server.

*   **Analogy: A Receptionist**
    Imagine a large office building with hundreds of employees. You don't know each person's direct phone number. Instead, you call the main reception desk. The receptionist (the reverse proxy) takes your message and forwards it to the correct employee (the backend server). You, the caller (the client), only ever interact with the receptionist.

**Diagram:**

```
                 +------------------+
                 |                  |
    +--------+   |   Reverse Proxy  |   +----------------+
    | Client |----->| (e.g., :8080)  |-->| Backend Server 1 |
    +--------+   |                  |   +----------------+
                 |                  |
                 |                  |   +----------------+
                 |                  |-->| Backend Server 2 |
                 +------------------+   +----------------+
```

**Key Functions:**
*   **Load Balancing:** Distribute incoming traffic among several backend servers.
*   **Security:** Hide the identity and characteristics of backend servers.
*   **SSL Termination:** Handle incoming HTTPS connections, decrypting them and passing unencrypted requests to the backends.
*   **Caching:** Store copies of responses to speed up future requests.

---

**Concept: Load Balancing**

Load balancing is the process of efficiently distributing network traffic across multiple servers. The goal is to prevent any single server from becoming a bottleneck and to ensure high availability and reliability.

#### Common Load Balancing Algorithms

1.  **Round-Robin**
    *   **How it works:** Requests are distributed to servers in a rotating sequence. The first request goes to Server 1, the second to Server 2, and so on. When the end of the list is reached, it starts over from the beginning.
    *   **Analogy: A Fair Restaurant Host:** A host seats arriving guests at each table one by one. Once every table has a group, the next group goes to the first table again. It's simple and ensures every server gets an equal number of connections over time.

2.  **Least Connections**
    *   **How it works:** The next incoming request is sent to the server that currently has the fewest active connections.
    *   **Analogy: A Smart Supermarket Cashier:** You look for the shortest checkout line. This algorithm sends new customers (requests) to the least busy cashier (server), which is often more efficient than just rotating.

---

**Concept: Thread Safety & Mutexes**

When a server handles multiple connections at once (concurrently), different parts of the code can run at the same time in separate **goroutines**.

If multiple goroutines try to read and write a shared piece of data—like the counter for our round-robin algorithm—they can interfere with each other, leading to a **race condition**.

*   **Analogy: The Key to a Shared Resource Room**
    Imagine a small supply room that only one person can enter at a time. To control access, there's a single key. To enter, you must take the key. If another person arrives, they have to wait until you come out and return the key.

A **Mutex** (Mutual Exclusion) is that key. Before a goroutine can access a shared variable, it must "lock" the mutex. When it's done, it "unlocks" it, allowing other goroutines to take their turn. This prevents race conditions and ensures data integrity.
