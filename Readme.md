# Wireguard Connector

rewrite the wgsimple script in go

Goals:
    - simple
    - boring technology (json)

## Router

- init
- get endpoint
- add endpoint (out)
- rm endpoint
- switch endpoint
- health endpoint
- enable/disable endpoint
- get wg config
- get config
- add nameserver

## Network Modes

- create groups of nodes
- group mode internal only
- group mode forward only
- group mode forward to endpoint and internal
- chain groups together

## Nodes

- get nodes
- add node
- rm node
- update node (refresh keys)
  - sign requets with keys/keybase
- auto ip (ipv4/ipv6?)

## Client

- add endpoint (router)
- rm endpoint (router)
- get endpoint (router)
- switch endpoint (router)
- connect endpoint (router)
- get wg config
- get config
- add nameserver

## Architecture

- Central Router (traffic in/out)
- Signaling Router (only signaling nodes and endpoints)

## API

All internal and external calls go through the go api

## Untracked

- simple json config for provisioning
- switch node/router mode, use a node as a endpoint
- add stun tun signaling
- add Automatic key rotation
- add WsTunnel (websocket) support
- add a function that check if the endpoint ovpn is alive.
- self healing, download fresh list from nordvpn,expressvpn,cyberghost
- set on each restart different server
- check if connection is correctly etablished
- output iptables rules
- add and delete rules from iptables
- select different countries
- make it testable with docker
- build cli chaining tool
- save ip and port of the client (nat, multiple nodes one nat ip)
- simple webui
- use net namespaces
- use fmark tables
- exec cmd in net namespace

## Question

### how verify the person/device?

- gpg, keybase, OAuth with github/gitlab?
- simple password on init?

### how to get the initial connection?

Maybe join with keybase or other gpg tool together and encrypt/decrypt the request or key?

### use iptables or nftables?

?

### Internal DNS?

with node names?

### Simple go ui?

- windows
- linux
- darwin
- arm

### Who will be use this?

- companies
- Gamers no techs
- techs
- self hoster

## docs

- [tailscale](https://tailscale.com/blog/)
  - use maschine and node keys curve
  - node key: human identity, OAuth2 identity provider
