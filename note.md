### Project 7: Simple DNS Server (UDP)

**Concept: UDP (User Datagram Protocol)**

| Feature             | TCP (Transmission Control Protocol) | UDP (User Datagram Protocol) |
|---------------------|-------------------------------------|------------------------------|
| **Analogy**         | Package with tracking & confirmation| Postcard                     |
| **Connection**      | Connection-Oriented (Handshake)     | Connectionless               |
| **Reliability**     | Reliable (guaranteed order, re-sends) | Unreliable (best-effort)     |
| **Data Model**      | Stream-based                        | Datagram-based (packets)     |
| **Use Cases**       | Web (HTTP), File Transfer (FTP), Email | DNS, Online Gaming, VoIP     |

**Why DNS uses UDP?**
DNS queries are typically small and require quick responses. If a query is lost, it's often faster to just re-send it than to establish a reliable TCP connection. UDP's speed and low overhead make it ideal for this kind of request-response pattern.

---

**Concept: How DNS Works**

DNS translates domain names (e.g., `google.com`) to IP addresses (e.g., `1.2.3.4`). A client sends a UDP packet containing a **query** to a DNS server on port 53. The server sends back a UDP packet with an **answer**.

**DNS Packet Structure (Simplified):**
1.  **Header:** Contains flags and counters (e.g., query ID, number of questions, answers).
2.  **Question Section:** The domain name being queried and the type of record requested (e.g., A record for IPv4 address).
3.  **Answer Section:** Contains the resource records (RRs) that answer the question (e.g., the IP address for an A record).

We will use the `golang.org/x/net/dns/dnsmessage` library to handle the binary format of these packets.