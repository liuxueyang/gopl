package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	for _, arg := range os.Args[1:] {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			log.Fatalf("Usage: %s key=value ...", os.Args[0])
			continue
		}

		location, address := parts[0], parts[1]
		go func(location, address string) {
			conn, err := net.Dial("tcp", address)

			if err != nil {
				log.Printf("Error connecting to %s: %v\n", location, err)
				return
			}

			defer conn.Close()
			mustCopy(os.Stdout, conn)
		}(location, address)
	}
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
