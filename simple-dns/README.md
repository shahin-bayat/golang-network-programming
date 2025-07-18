# Simple DNS Server

This project implements a basic DNS server in Go that listens for UDP queries on port 53. Currently, it's set up to log incoming queries. In later steps, it will be extended to respond to specific queries.

## How to Run

1.  Navigate to the project directory:
    ```bash
    cd /Users/sbayat/Projects/personal/go/tutorials/networking/simple-dns
    ```
2.  Run the server:
    ```bash
    sudo go run main.go
    ```
    *(Note: `sudo` is required because DNS typically runs on port 53, which is a privileged port.)*

## How to Test

You can test this DNS server using standard DNS client tools like `nslookup`, `dig`, and `host`. You'll need to explicitly tell these tools to query your local server (usually `127.0.0.1` or `localhost`).

### Using `nslookup`

`nslookup` is a command-line tool for querying DNS servers.

```bash
nslookup example.com 127.0.0.1
```
Replace `example.com` with any domain you want to query. The `127.0.0.1` at the end tells `nslookup` to use your local server.

### Using `dig`

`dig` (Domain Information Groper) is a more flexible and powerful tool for interrogating DNS name servers.

```bash
dig @127.0.0.1 example.com
```
The `@127.0.0.1` specifies your local DNS server.

### Using `host`

`host` is a simple utility for performing DNS lookups.

```bash
host example.com 127.0.0.1
```
Again, `127.0.0.1` directs the query to your running server.

---

When you run these commands, you should see output in your `simple-dns` server's console indicating that it received a query from your client tool.
