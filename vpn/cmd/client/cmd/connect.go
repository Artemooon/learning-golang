/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "net"
    "log"
    "time"
	"github.com/spf13/cobra"
	"vpn/utils"
	"github.com/songgao/water"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the vpn server",
	RunE: func(cmd *cobra.Command, args []string) error {

	    server, err := cmd.Flags().GetString("server")
	    if err != nil {
	        return err
	    }
        if server == "" {
            fmt.Println("--server flag is required", server)
            return nil
        }
        host, port, err := net.SplitHostPort(server)
        if err != nil {
            return fmt.Errorf(
               "--server must be in the host:port format, example: 127.0.0.1:8000: %w",
               err,
            )
        }
        if host == "" || port == "" {
            return fmt.Errorf(
                "--server must include host and port",
            )
        }

        tun := utils.CreateTUN("10.0.5.2", "10.0.6.2", "utun6")
        log.Println("Interface Name:", tun.Name())

        clientSocket(tun, "127.0.0.1:9000")
		return nil
	},
}

func clientSocket(incomingTun *water.Interface, serverAddress string) {
    dialer := net.Dialer{
        Timeout: 3 * time.Second,
    }
    conn, err := dialer.Dial("udp", serverAddress)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    _, err = conn.Write([]byte("hello from client"))
    if err != nil {
        panic(err)
    }

    go readFromUdpToTun(incomingTun, conn)
    go readFromTunToUdp(incomingTun, conn)
    select {}

//     buffer := make([]byte, 1024)
//
//     n, err := conn.Read(buffer)
//     if err != nil {
//         panic(err)
//     }
//
//     fmt.Printf("server response: %s", buffer[:n])
//
//     n, err = conn.Read(buffer)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Printf("server2 response: %s", buffer[:n])
}

func readFromUdpToTun(tun *water.Interface, conn net.Conn) {
    buffer := make([]byte, 1024)

    for {
        n, err := conn.Read(buffer)
        if err != nil {
            log.Printf(
                "Error reading from UDP socket: %w", err,
            )
            continue
        }
        packet := utils.ParseIPv4(buffer[:n])
        log.Printf("Will write to %s: %+v", tun.Name(), packet)

        _, err = tun.Write(buffer[:n])

        if err != nil {
            log.Println("error writing to TUN:", err)
        }
        log.Printf("%s: Wrote: %d bytes: %+v", tun.Name(), n, buffer[:n])
    }
}

func readFromTunToUdp(tun *water.Interface, conn net.Conn) {
    buffer := make([]byte, 1024)

    for {
        n, err := tun.Read(buffer)
        if err != nil {
            log.Printf(
                "Error reading from tun: %w", err,
            )
            continue
        }

        packet := utils.ParseIPv4(buffer[:n])
        log.Printf("Read from TUN: %s: %+v", tun.Name(), packet)

        _, err = conn.Write(buffer[:n])
        if err != nil {
            log.Println("error writing to UDP socket:", err)
        }
        log.Printf("%s: Udp socket: %d bytes: %+v", conn.LocalAddr(), n, buffer[:n])
    }
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	connectCmd.Flags().StringP("server", "s", "", "VPN server in host:port form")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
