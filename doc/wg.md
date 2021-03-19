
# Wireguard

- save private keys encrypted
- generate qrcode: qrencode -t ansiutf8 -r client.conf

examples: /usr/share/wireguard-tools/examples
    - json: convert current wg conf to json format
    - external-tests/go: go test example
    - extract-keys: extract the symmetric ChaCha20Poly1305 session keys from kernel mem (can not compile)
    - reresolve-dns/reresolve-dns.sh: recover from an endpoint that has changed its IP

wg:

```txt
Usage: wg set <interface> [listen-port <port>] [fwmark <mark>] [private-key <file path>] [peer <base64 public key> [remove] [preshared-key <file path>] [endpoint <ip>:<port>] [persistent-keepalive <interval seconds>] [allowed-ips <ip1>/<cidr1>[,<ip2>/<cidr2>]...] ]...
```

## Manual Setup

### Node-A

ip link add dev wg0 type wireguard
ip addr add 10.0.0.1/24 dev wg0
ip addr add fdc9:281f:04d7:9ee9::1/64 dev wg0
wg set wg0 listen-port 51871 private-key /path/to/peer_A.key
wg set wg0 peer PEER_B_PUBLIC_KEY preshared-key /path/to/peer_A-peer_B.psk endpoint peer-b.example:51902 allowed-ips 10.0.0.2/32,fdc9:281f:04d7:9ee9::2/128
wg set wg0 peer PEER_C_PUBLIC_KEY preshared-key /path/to/peer_A-peer_C.psk allowed-ips 10.0.0.3/32,fdc9:281f:04d7:9ee9::3/128
ip link set wg0 up

### Node-B

ip link add dev wg0 type wireguard
ip addr add 10.0.0.2/24 dev wg0
ip addr add fdc9:281f:04d7:9ee9::2/64 dev wg0
wg set wg0 listen-port 51902 private-key /path/to/peer_B.key
wg set wg0 peer PEER_A_PUBLIC_KEY preshared-key /path/to/peer_A-peer_B.psk endpoint 198.51.100.101:51871 allowed-ips 10.0.0.1/32,fdc9:281f:04d7:9ee9::1/128
wg set wg0 peer PEER_C_PUBLIC_KEY preshared-key /path/to/peer_B-peer_C.psk allowed-ips 10.0.0.3/32,fdc9:281f:04d7:9ee9::3/128
ip link set wg0 up

### Node-C

ip link add dev wg0 type wireguard
ip addr add 10.0.0.3/24 dev wg0
ip addr add fdc9:281f:04d7:9ee9::3/64 dev wg0
wg set wg0 listen-port 51993 private-key /path/to/peer_C.key
wg set wg0 peer PEER_A_PUBLIC_KEY preshared-key /path/to/peer_A-peer_C.psk endpoint 198.51.100.101:51871 allowed-ips 10.0.0.1/32,fdc9:281f:04d7:9ee9::1/128
wg set wg0 peer PEER_B_PUBLIC_KEY preshared-key /path/to/peer_B-peer_C.psk endpoint peer-b.example:51902 allowed-ips 10.0.0.2/32,fdc9:281f:04d7:9ee9::2/128
ip link set wg0 up

## Point-to-Site

allowed-ips 10.0.0.2/32,fdc9:281f:04d7:9ee9::2/128,192.168.35.0/24,fd7b:d0bd:7a6e::/64
ip route add 192.168.35.0/24 dev wg0
ip route add fd7b:d0bd:7a6e::/64 dev wg0

## Site-to-point

sysctl -w net.ipv4.ip_forward=1
sysctl -w net.ipv6.conf.all.forwarding=1

## Testing

```bash
# site on tunnel listen, port must be open
nc -vvlnp 2222
# other side sender
dd if=/dev/zero bs=1024K count=1024 | nc -v 10.0.0.203 2222
# status with
wg
```

## Debug Mode

modprobe wireguard
echo module wireguard +p > /sys/kernel/debug/dynamic_debug/control
dmesg -wH

## Reload Config

wg syncconf ${WGNET} <(wg-quick strip ${WGNET})

## docs

- [wiki](https://docs.sweeting.me/s/wireguard)
