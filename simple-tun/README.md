# Simple TUN "Hello World" - Testing Guide

This guide explains the shell commands used to test the `simple-tun` application inside its Docker container.

## Testing Commands Explained

After getting a shell into the running container (`docker exec -it <container_id> /bin/sh`), you need to configure the virtual `tun0` interface and send a packet to it.

### 1. Assign an IP Address

```sh
ip addr add 10.0.0.1/24 dev tun0
```

- **`ip addr add`**: This is the command used to add a new IP address to a network interface.
- **`10.0.0.1/24`**: This is the IP address (10.0.0.1) and the subnet mask in CIDR notation (`/24`, which is equivalent to `255.255.255.0`). This defines a private network for our TUN device.
- **`dev tun0`**: This specifies that the command applies to the network device named `tun0`. The `dev` keyword is used to indicate the target device.

### 2. Bring the Interface Up

```sh
ip link set tun0 up
```

- **`ip link set`**: This command is used to change the attributes of a network device.
- **`tun0`**: The name of the target network interface.
- **`up`**: This action changes the state of the interface to "up," enabling it to send and receive packets.

```sh
ip addr add 10.0.0.1/24 dev tun-server
ip link set tun-server up

ip addr add 10.0.0.2/24 dev tun-client
ip link set tun-client up
```

```

```

### 3. Send a Test Packet

```sh
ping -c 1 10.0.0.2
```

- **`ping`**: A standard utility to test network connectivity. It works by sending an ICMP "ECHO_REQUEST" packet to a destination.
- **`-c 1`**: This option tells `ping` to send only one packet and then stop. Without this, `ping` would run continuously.
- **`10.0.0.2`**: The destination IP address. Since our `tun0` interface is on the `10.0.0.0/24` network, the operating system knows to route this packet through the `tun0` device. Your Go application, which is listening on this device, will then receive this packet.

### 4. Bringing down

```sh
ip addr del 10.0.0.1/24 dev tun0
ip addr del 10.0.0.2/24 dev tun1
```

### 5. Checking Interface Status

```sh
 ip addr show dev tun0
 ip route show
 tcpdump -i tun0

```

### 5. How to run

```sh
docker compose up -d
docker compose exec vpn-server sh
docker compose exec vpn-client sh

go run . -server
go run . -remote vpn-server:9090

# on the client
ip addr add 10.0.0.1/24 dev tun-client
ip link set tun-client up

# on the server
ip addr add 10.0.0.2/24 dev tun-server
ip link set tun-server up

```
