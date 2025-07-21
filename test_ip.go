package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	domain := "example.com"
	
	// Test IP lookup functionality
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("IP lookup failed: %v\n", err)
		return
	}
	
	var ipv4List, ipv6List []string
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4List = append(ipv4List, ip.String())
		} else {
			ipv6List = append(ipv6List, ip.String())
		}
	}
	
	fmt.Printf("Domain: %s\n", domain)
	if len(ipv4List) > 0 {
		fmt.Printf("IPv4: %s\n", strings.Join(ipv4List, ", "))
	} else {
		fmt.Printf("IPv4: -\n")
	}
	
	if len(ipv6List) > 0 {
		fmt.Printf("IPv6: %s\n", strings.Join(ipv6List, ", "))
	} else {
		fmt.Printf("IPv6: -\n")
	}
}
