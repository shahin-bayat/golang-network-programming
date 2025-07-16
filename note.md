# HTTP Request Format: A Quick Reference

This note explains the structure of HTTP requests.

## The HTTP/1.x Request-Line

For any request using HTTP/1.0 or HTTP/1.1, the communication **must** start with a specific line of text called the **Request-Line**. This line always follows the format:

`METHOD /path HTTP/version`

**Example:** `GET /index.html HTTP/1.1`

This format is strictly defined by the HTTP specification and contains three distinct parts separated by single spaces:

1.  **Method:** A verb that tells the server *what action to perform*.
    *   `GET`: Retrieve a resource. (This is what you handle).
    *   `POST`: Submit data to be processed (e.g., a web form).
    *   `PUT`: Replace a resource.
    *   `DELETE`: Remove a resource.
    *   `HEAD`: Same as `GET`, but only asks for the headers, not the response body.
    *   Others include `OPTIONS`, `PATCH`, `CONNECT`.

2.  **Path (Request Target):** This tells the server *which resource* the action applies to. It's the part of the URL that comes after the domain (e.g., `/about-us.html` or `/search?q=networking`).

3.  **Protocol Version:** This tells the server which version of the HTTP rules the client (browser) is following (e.g., `HTTP/1.1`). This is essential for ensuring compatibility.

## The Shift to HTTP/2 and HTTP/3

The text-based `Request-Line` format is specific to **HTTP/1.x**. Modern web communication has evolved.

*   **HTTP/2 and HTTP/3 are binary protocols.**
*   They do not send human-readable text. Instead, they use a system of binary "frames" to communicate the same information (method, path, headers, etc.).
*   This binary format is much more efficient, faster, and less error-prone to parse.

### Why Does a Simple Text-Based Server Still Work?

Your server works because of **protocol negotiation**.

1.  A modern browser first tries to connect using HTTP/2 or HTTP/3.
2.  The server (yours) doesn't understand this new protocol.
3.  The browser then automatically **falls back to the universally supported HTTP/1.1**.
4.  Your server receives a classic, text-based HTTP/1.1 request that your code can parse.

Building a server that parses the text-based format is a fundamental exercise in network programming that teaches you the core concepts of the web's most important protocol.
