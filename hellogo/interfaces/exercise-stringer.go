package main

import "fmt"

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.

func (ip IPAddr) String() string {
	finalStr := ""
	for i, number := range ip {
		if i == len(ip) - 1 {
			finalStr += fmt.Sprintf("%d", number)
		} else {
			finalStr += fmt.Sprintf("%d.", number)
		}
	}

	return finalStr
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
		"invalid": {},
		"invalid2": {0,0,1},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}
