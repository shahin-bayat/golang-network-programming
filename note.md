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

---

**Concept: Subdomains and Authoritative Zones**

Our current simple DNS server acts as an **authoritative server** for the domains explicitly listed in its `records` map. This means it's the "source of truth" for those specific domain names.

*   **Exact Match:** If you query `example.com.`, the server looks for `example.com.` in its map.
*   **Subdomains:** If you query `test.example.com.`, the server looks for `test.example.com.` in its map. It does *not* automatically infer that `test.example.com.` should resolve if `example.com.` is present. For a subdomain to resolve, it must be explicitly added to the `records` map.

In a real DNS system, a server is authoritative for a **zone** (e.g., `example.com.`). This zone includes the main domain and all its subdomains (`www.example.com.`, `mail.example.com.`, `test.example.com.`, etc.) unless a subdomain is explicitly delegated to another DNS server.

---

**Concept: Recursive vs. Iterative DNS Resolution**

When a client (like your web browser) needs to resolve a domain name, it typically sends a **recursive query** to a local DNS resolver (e.g., your ISP's DNS server, or public ones like 8.8.8.8).

*   **Recursive Resolver:**
    *   **Role:** Takes on the full responsibility of resolving the query for the client.
    *   **Process:** If it doesn't have the answer in its cache, it will perform a series of **iterative queries** on behalf of the client to find the answer.
    *   **Analogy:** A concierge who finds all the information you need and gives you the final answer.

*   **Iterative Resolver:**
    *   **Role:** Responds to queries by providing the "next step" in the resolution process, rather than the final answer.
    *   **Process:**
        1.  Queries a **Root Server** (e.g., "Where is `.com`?").
        2.  Root Server responds with **referral** to TLD servers (e.g., "Go ask these `.com` servers.").
        3.  Queries a **TLD Server** (e.g., "Where is `example.com`?").
        4.  TLD Server responds with **referral** to Authoritative Name Servers for `example.com`.
        5.  Queries an **Authoritative Name Server** for `example.com`.
        6.  Authoritative Server provides the final IP address.
    *   **Analogy:** A librarian who tells you which section or other library to go to next to find your book.

Our current server is a very simple authoritative server. If we extend it to perform lookups for domains it doesn't know about, we could implement either a recursive approach (forwarding the query to an upstream recursive resolver) or an iterative approach (performing the full lookup chain ourselves).