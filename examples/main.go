package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/docker/libcontainer/netlink"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func show(dev string) {

	c, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	defer c.Close()

	d, err := c.Device(dev)
	if err != nil {
		log.Fatalf("failed to get device %q: %v", dev, err)
	}
	printWg(d)

	//devices, err = c.Devices()

}

func printWg(wg *wgtypes.Device) {
	printDevice(wg)
	for _, p := range wg.Peers {
		printPeer(p)
	}
}

func config(dev string) {

	peerpub, err := wgtypes.ParseKey("YkWRYOOFndLypbSTiEiN22hHzaIvOKYAmUSfSpmbLz8=")
	if err != nil {
		fmt.Println(err)
		return
	}
	endpoint, _, err := net.ParseCIDR("192.168.399.11/32")
	if err != nil {
		fmt.Println(err)
	}
	_, n, err := net.ParseCIDR("192.168.399.0/24")
	var allowedIPs []net.IPNet
	allowedIPs = append(allowedIPs, *n)

	if err != nil {
		fmt.Println(err)
	}
	peers := []wgtypes.PeerConfig{
		{PublicKey: peerpub,
			Endpoint: &net.UDPAddr{
				IP:   endpoint,
				Port: 35700,
			},
			AllowedIPs: allowedIPs,
		},
	}

	priv, err := wgtypes.GeneratePrivateKey()

	cfg := wgtypes.Config{
		PrivateKey: &priv,
		Peers:      peers,
	}

	c, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	defer c.Close()

	err = c.ConfigureDevice(dev, cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	d, _ := c.Device(dev)
	printWg(d)

}

func create(name string) {
	// ip link add dev wg0 type wireguard
	// ip link del dev wg0 type wireguard

	err := netlink.NetworkLinkAdd(name, "wireguard")
	if err != nil {
		fmt.Println(err)
		return
	}

}

func delete(name string) {
	err := netlink.NetworkLinkDel(name)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func main() {
	dev := "wg1"
	delete(dev)
	return

	create(dev)
	wgd, err := net.InterfaceByName(dev)
	if err != nil {
		fmt.Println(err)
		return
	}
	ip, netip, err := net.ParseCIDR("192.168.399.2/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = netlink.NetworkLinkAddIp(wgd, ip, netip)
	if err != nil {
		fmt.Println(err)
	}

	err = netlink.NetworkLinkUp(wgd)
	if err != nil {
		fmt.Println(err)
		return
	}

	config(dev)

	//show(dev)

}

func printDevice(d *wgtypes.Device) {
	const f = `interface: %s (%s)
  public key: %s
  private key: (hidden)
  listening port: %d
`

	fmt.Printf(
		f,
		d.Name,
		d.Type.String(),
		d.PublicKey.String(),
		d.ListenPort)
}

func printPeer(p wgtypes.Peer) {
	const f = `peer: %s
  endpoint: %s
  allowed ips: %s
  latest handshake: %s
  transfer: %d B received, %d B sent
`

	fmt.Printf(
		f,
		p.PublicKey.String(),
		// TODO(mdlayher): get right endpoint with getnameinfo.
		p.Endpoint.String(),
		ipsString(p.AllowedIPs),
		p.LastHandshakeTime.String(),
		p.ReceiveBytes,
		p.TransmitBytes,
	)
}

func ipsString(ipns []net.IPNet) string {
	ss := make([]string, 0, len(ipns))
	for _, ipn := range ipns {
		ss = append(ss, ipn.String())
	}

	return strings.Join(ss, ", ")
}
