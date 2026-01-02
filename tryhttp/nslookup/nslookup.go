package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("%s: usage: <host>", os.Args[0])
		log.Fatalf("expected exactly one argument; got %d", len(os.Args)-1)
	}

	host := os.Args[1]
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatalf("loopup ip: %s: %v", host, err)
	}

	if len(ips) == 0 {
		log.Fatalf("no ips found for %s", host)
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			fmt.Printf("ipV4: %s\n", ip)
			goto IPV6
		}
	}

	fmt.Print("none\n")

IPV6:
	for _, ip := range ips {
		if ip.To4() == nil {
			fmt.Printf("ipV6: %s\n", ip)
			return
		}
	}
	fmt.Print("none\n")
}
