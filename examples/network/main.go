package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/kawa1214/tcp-ip-go/network"
)

func main() {
	network, err := network.NewTun()
	if err != nil {
		log.Fatalf("Failed to create TUN device: %v", err)
	}
	network.Bind()

	for {
		pkt, _ := network.Read()
		fmt.Print(hex.Dump(pkt.Buf[:pkt.N]))
		network.Write(pkt)
	}
}
