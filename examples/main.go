package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/docker/libcontainer/netlink"
	"go.uber.org/zap"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func show(dev string) {
	c, err := wgctrl.New()
	if err != nil {
		fmt.Printf("failed to open wgctrl: %v", err)
		return
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

type Config map[string]wgtypes.Config

// type Unpacker struct {
// 	Data interface{}
// }

// func (u *Unpacker) UnmarshalJSON(b []byte) error {
// 	device := &wgtypes.Device{}
// 	err := json.Unmarshal(b, device)
// 	// abort if we have an error other than the wrong type
// 	if _, ok := err.(*json.UnmarshalTypeError); err != nil && !ok {
// 		return err
// 	}
// 	fmt.Println(device)
// 	return nil
// }

func readConfig() error {
	b, err := ioutil.ReadFile("../wgc.json")
	if err != nil {
		return err
	}
	c := make(Config)
	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	fmt.Println(c)

	return nil

}

func config(dev string) {

	peerpub, err := wgtypes.ParseKey("YkWRYOOFndLypbSTiEiN22hHzaIvOKYAmUSfSpmbLz8=")
	if err != nil {
		fmt.Println(err)
		return
	}
	endpoint, _, err := net.ParseCIDR("192.168.200.11/32")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, n, err := net.ParseCIDR("192.168.200.0/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	var allowedIPs []net.IPNet
	allowedIPs = append(allowedIPs, *n)

	if err != nil {
		fmt.Println(err)
		return
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
	// d, err := c.Device(dev)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// printWg(d)

}

func create(name string) error {
	// ip link add dev wg0 type wireguard
	return netlink.NetworkLinkAdd(name, "wireguard")
}

func delete(name string) error {
	// ip link del dev wg0 type wireguard
	return netlink.NetworkLinkDel(name)
}

// func logFatal(log sugar ,err error){
// 	log.
// }

func main() {
	dev := "wg1"

	err := readConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	return
	//delete(dev)
	//return

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	// log.Infow("failed to fetch URL",
	// 	// Structured context as loosely typed key-value pairs.
	// 	"url", url,
	// 	"attempt", 3,
	// 	"backoff", time.Second,
	// )

	err = create(dev)
	if err != nil {
		log.Infof("create failed: %s", err)
		return
	}
	wgd, err := net.InterfaceByName(dev)
	if err != nil {
		fmt.Println(err)
		return
	}
	ip, netip, err := net.ParseCIDR("192.168.200.2/24")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = netlink.NetworkLinkAddIp(wgd, ip, netip)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = netlink.NetworkLinkUp(wgd)
	if err != nil {
		fmt.Println(err)
		return
	}

	config(dev)

	show(dev)

	delete(dev)

}

func printDevice(d *wgtypes.Device) {
	const f = `interface: %s (%s)
  public key: %s
  private key: %s
  listening port: %d
`

	fmt.Printf(
		f,
		d.Name,
		d.Type.String(),
		d.PublicKey.String(),
		d.PrivateKey.String(),
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
