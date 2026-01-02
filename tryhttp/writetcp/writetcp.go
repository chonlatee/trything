package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const name = "writetcp"
	log.SetPrefix(name + "\t")

	port := flag.Int("p", 8080, "port to connect to")
	flag.Parse()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: *port})
	if err != nil {
		log.Fatalf("error connecting to localhost:%d %v", *port, err)
	}

	log.Printf("connected to %s: will foward stdin", conn.RemoteAddr())
	defer conn.Close()

	go func() {
		for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
			fmt.Printf("%s\n", connScanner.Text())

			if err := connScanner.Err(); err != nil {
				log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}()

	for stdinScanner := bufio.NewScanner(os.Stdin); stdinScanner.Scan(); {
		log.Printf("sent: %s\n", stdinScanner.Text())

		if _, err := conn.Write(stdinScanner.Bytes()); err != nil {
			log.Fatalf("error writing to %s: %v", conn.RemoteAddr(), err)
		}

		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Fatalf("error writting to %s: %v", conn.RemoteAddr(), err)
		}

		if stdinScanner.Err() != nil {
			log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
		}
	}

}
