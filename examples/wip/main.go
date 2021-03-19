package main

import (
	"fmt"
	"syscall"

	"github.com/mdlayher/genetlink"
	"golang.org/x/sys/unix"
)

type Client struct {
	c      *genetlink.Conn
	family genetlink.Family

	interfaces func() ([]string, error)
}

func main() {
	c, err := genetlink.Dial(nil)
	if err != nil {
		return
	}
	f, err := c.GetFamily("wireguard")

	tab, err := syscall.NetlinkRIB(unix.RTM_GETLINK, unix.AF_UNSPEC)
	if err != nil {
		return
	}

	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return
	}

	//var ifis []string
	for _, m := range msgs {
		// Only deal with link messages, and they must have an ifinfomsg
		// structure appear before the attributes.
		fmt.Println(m.Header.Type)
		fmt.Println(unix.RTM_NEWLINK)
	}

	fmt.Println(c)
	fmt.Println(f)
	//fmt.Println(msgs)
}
