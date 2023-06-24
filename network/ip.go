package network

import (
	"log"

	"github.com/kawa1214/tcp-ip-go/link"
)

type IpPacket struct {
	IpHeader *Header
	Packet   link.Packet
}

type IpPacketQueue struct {
	packetChan chan IpPacket
}

func NewIpPacketQueue() *IpPacketQueue {
	return &IpPacketQueue{}
}

func (q *IpPacketQueue) QueuePacket(d link.NetDevice) {
	packets := make(chan IpPacket, 10)
	q.packetChan = packets
	go func() {
		for {
			select {
			case pkt := <-d.PacketChan():
				log.Printf("network pkt: %d", pkt.N)
				ipHeader, err := Parse(pkt.Buf[:pkt.N])
				if err != nil {
					log.Printf("parse error: %s", err)
					continue
				}
				ipPacket := IpPacket{
					IpHeader: ipHeader,
					Packet:   pkt,
				}
				packets <- ipPacket
			}
		}
	}()
}

func (q *IpPacketQueue) PacketChan() chan IpPacket {
	return q.packetChan
}
