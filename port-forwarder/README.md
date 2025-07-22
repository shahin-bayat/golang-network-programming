# Port Forwarding Tunnel

This project is a simple TCP port forwarding tunnel, similar in function to `ssh -L`. It listens on a local port and forwards all traffic to a specified remote address.

## How to Build

From the root of the `networking` directory, run:

```bash
go build -o port-forwarder/port-forwarder ./port-forwarder
```

## How to Test

To test the tunnel, you need three separate terminal windows.

### 1. Terminal 1: Start a Destination Server

We need a server to forward traffic *to*. A simple echo server is perfect for this. The `ncat` (`nc`) utility can create one easily. This server will listen on port 9000.

```bash
ncat -l 9000 -k --exec "/bin/cat"
```
*   `-l 9000`: Listen on port 9000.
*   `-k`: Keep the server running after a client disconnects.
*   `--exec "/bin/cat"`: Creates the echo behavior.

### 2. Terminal 2: Start the Port Forwarder

Run the `port-forwarder` application. We will tell it to listen on local port `8080` and forward all traffic to our echo server at `localhost:9000`.

```bash
./port-forwarder -local :8080 -remote localhost:9000
```

You should see the message: `local address :8080 remote address localhost:9000`

### 3. Terminal 3: Connect as a Client

Now, connect to the tunnel's listener on port `8080`.

```bash
ncat localhost 8080
```

Anything you type into this terminal will now be sent through the port forwarder to the echo server (Terminal 1), echoed back through the port forwarder, and finally displayed on your screen in Terminal 3.
