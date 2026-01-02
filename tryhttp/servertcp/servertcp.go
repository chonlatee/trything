package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	const name = "tcpupperecho"
	log.SetPrefix(name + "\t")

	port := flag.Int("p", 8080, "port to listen on")
	flag.Parse()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Printf("listening at localhost: %s", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go echoUpper(conn, conn)
	}
}

func echoUpper(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		log.Printf("recieve: %s\n", line)

		fmt.Fprintf(w, "%s\n", strings.ToUpper(line))

		if err := scanner.Err(); err != nil {
			log.Printf("error: %s", err)
		}
	}
}
