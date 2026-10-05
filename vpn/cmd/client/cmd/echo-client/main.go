package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vpn/cmd/client/cmd/internal/frame"
)

func main() {
	certPEM, err := os.ReadFile("./server_vpn_cert.crt")
	if err != nil {
		log.Fatalf("client error: %v", err)
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(certPEM); !ok {
		log.Fatal("client error: server certificate is not valid PEM")
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:9000", &tls.Config{
		RootCAs:    roots,
		ServerName: "localhost",
		MinVersion: tls.VersionTLS13,
	})
	if err != nil {
		log.Fatalf("TLS connection failed: %v", err)
	}
	defer conn.Close()

	payload, err := frame.ReadFrame(conn)

	if err != nil {
		log.Fatalf("client error, reading from socket: %v", err)
	}
	log.Printf("greeting: %q", string(payload))

	state := conn.ConnectionState()
	log.Printf("TLS version: %s", tls.VersionName(state.Version))
	log.Printf("cipher suite: %s", tls.CipherSuiteName(state.CipherSuite))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("client cancelled: %v", ctx.Err())
			return

		case <-ticker.C:
			if err := frame.WriteFrame(conn, []byte("ping\n")); err != nil {
				log.Printf("writing ping: %v", err)
				return
			}

			echo, err := frame.ReadFrame(conn)
			if err != nil {
				log.Printf("reading echo: %v", err)
				return
			}

			log.Printf("echo: %q", string(echo))
		}
	}

}
