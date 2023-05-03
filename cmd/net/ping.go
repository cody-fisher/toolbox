/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package net

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var (
	urlPath string
)

func ping(domain string) (int, error) {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, err
	}
	defer c.Close()

	// Resolve the domain to an IP address
	ip, err := net.ResolveIPAddr("ip4", domain)
	if err != nil {
		return 0, err
	}

	// Create the ICMP message
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   0,
			Seq:  0,
			Data: []byte(""),
		},
	}

	// Marshal the ICMP message
	b, err := m.Marshal(nil)
	if err != nil {
		return 0, err
	}

	// Send the ICMP message
	start := time.Now()
	n, err := c.WriteTo(b, ip)
	if err != nil {
		return 0, err
	}

	// Read the response
	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		return 0, err
	}
	duration := time.Since(start)
	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", n, peer, 0, duration)
	return n, nil
}

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "This pings a remote URL and returns the response",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := ping(urlPath)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	pingCmd.Flags().StringVarP(&urlPath, "url", "u", "", "URL to ping")

	if err := pingCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
