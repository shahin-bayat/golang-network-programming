# Simple TLS Terminating Proxy

This project is a simple TLS terminating proxy written in Go. It accepts encrypted TLS connections, decrypts the traffic, and forwards it to a backend server as plain TCP.

## Setup

### Generate Self-Signed Certificate

To run this proxy, you first need a TLS certificate and a private key. For development, you can generate a self-signed certificate using `openssl`.

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost"
```

**What this command does:**

*   `openssl req`: A command to create and process certificate signing requests (CSRs).
*   `-x509`: Specifies that we want to create a self-signed certificate instead of a certificate request.
*   `-newkey rsa:2048`: Creates a new 2048-bit RSA private key.
*   `-nodes`: (No DES) tells OpenSSL not to encrypt the private key with a passphrase.
*   `-keyout key.pem`: Specifies the filename for the private key.
*   `-out cert.pem`: Specifies the filename for the certificate.
*   `-days 365`: Sets the certificate's validity period to 365 days.
*   `-subj "/CN=localhost"`: Sets the "Common Name" for the certificate to `localhost`, avoiding interactive prompts.
