/*
Copyright © 2023 Blair Gillam <ns1h@airmada.net>
*/
package cmd

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var listenPort int

// Check if a certain ip in a cidr range.
func cidrRangeContains(cidrRange string, checkIP string) (bool, error) {
	_, ipnet, err := net.ParseCIDR(cidrRange)
	if err != nil {
		return false, err
	}
	secondIP := net.ParseIP(checkIP)
	return ipnet.Contains(secondIP), err
}

// handlePacket processes UDP packet payload in a goroutine
func handlePacket(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	//responseStr := fmt.Sprintf("%v", string(buf))

	//fmt.Println(addr)
	//fmt.Println(responseStr)
	//fmt.Println(udpServer.LocalAddr())

	conn, _ := net.Dial("udp", "255.255.255.255:4992")
	defer conn.Close()
	conn.Write(buf)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flextool",
	Short: "A flexradio broadcast packet receiver for VPN clients",
	Long:  `A flexradio broadcast packet receiver for VPN clients`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listening for FlexRadio Discovery packets... Use Ctrl+C to stop listening.")

		if math.Signbit(float64(listenPort)) || listenPort >= 65536 {
			fmt.Println("Port number must be a valid port between 0 and 65535")
			os.Exit(1)
		}

		// Listen on all interfaces for broadcast traffic to UDP listenPort (default: 14992/udp)
		addr := net.UDPAddr{
			Port: listenPort,
			IP:   net.ParseIP("0.0.0.0"),
		}
		udpServer, err := net.ListenUDP("udp", &addr)
		if err != nil {
			panic(err)
		}
		defer udpServer.Close()

		// Retransmit payload on local interfaces to 255.255.255.255 4992/udp
		for {
			buf := make([]byte, 1024)
			_, addr, err := udpServer.ReadFrom(buf)
			if err != nil {
				continue
			}

			// DEBUG
			// fmt.Println("addr.String: ", addr.(*net.UDPAddr).IP.String())

			// Ignore packets from 127.0.0.1 or packets from the VPN client subnet. (Currently hardcoded for a test network)
			firstCheck, err := cidrRangeContains("10.5.0.0/24", addr.(*net.UDPAddr).IP.String())
			if err != nil {
				log.Println(err)
			}

			if !firstCheck && addr.(*net.UDPAddr).IP.String() != "127.0.0.1" {
				go handlePacket(udpServer, addr, buf)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("debug", "d", false, "NOT IMPLEMENTED YET. Turns on verbose debug output to the console")
	rootCmd.Flags().IntVarP(&listenPort, "port", "p", 14992, "UDP port to listen for FlexRadio discovery packets")
}
