package main

import (
	"fmt"

	"github.com/keis8221/tcp-ip-go/internet"
	"github.com/keis8221/tcp-ip-go/network"
	"github.com/keis8221/tcp-ip-go/transport"
)

func main() {
	network, _ := network.NewTun()
	network.Bind()
	ip := internet.NewIpPacketQueue()
	ip.ManageQueues(network)
	tcp := transport.NewTcpPacketQueue()
	tcp.ManageQueues(ip)

	for {
		pkt, _ := tcp.ReadAcceptConnection()
		fmt.Printf("TCP Header: %+v\n", pkt.Pkt.TcpHeader)
	}
}
