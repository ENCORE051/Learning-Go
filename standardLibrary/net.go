package main

import (
	"fmt"
	"net"
)

func usingIPV4() {
	ipv4Addr, err := net.ResolveIPAddr("ip4", "192.168.2.1")
	if err != nil {
		panic(err)
	}
	fmt.Println(ipv4Addr)

	ipv6Addr, err := net.ResolveIPAddr("ip6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	if err != nil {
		panic(err)
	}
	fmt.Println(ipv6Addr)
}

func usingDNS() {
	addrs, err := net.LookupHost("github.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(addrs)
}

func main() {
	usingIPV4()
	usingDNS()
}
