package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"vpn/cmd/client/cmd/internal/frame"
)

func main() {
	cert, err := tls.LoadX509KeyPair("./server_vpn_cert.crt", "./server_ec_private.key")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:9000", &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Println("echo server is listening on 127.0.0.1:9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept err:", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			log.Printf("client connected: %s", conn.RemoteAddr())

			err := frame.WriteFrame(conn, []byte("Just an echo\n"))
			if err != nil {
				log.Println("greeting write error:", err)
				return
			}

			for {
				payload, err := frame.ReadFrame(conn)
				if err == io.EOF {
					log.Printf("client disconnected: %s", conn.RemoteAddr())
					return
				}
				if err != nil {
					log.Printf("read failed: %v", err)
					return
				}

				log.Printf("received: %q", string(payload))

				if err := frame.WriteFrame(conn, payload); err != nil {
					log.Printf("echo write failed: %v", err)
					return
				}
			}

		}(conn)
	}
}
