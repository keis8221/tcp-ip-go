package transport

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type State int

const (
	StateListen = iota
	StateSynReceived
	StateEstablished
	StateCloseWait
	StateLastAck
	StateClosed
)

type Connection struct {
	SrcPort uint16
	DstPort uint16
	State   State
	Pkt     TcpPacket
	N       uintptr

	initialSeqNum   uint32
	incrementSeqNum uint32

	isAccept bool
}

type ConnectionManager struct {
	Connections           []Connection
	AcceptConnectionQueue chan Connection
	lock                  sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Connection:            make([]Connection, 0),
		AcceptConnectionQueue: make(chan Connection, QUEUE_SIZE),
	}
}

func (m *ConnectionManager) recv(queue *TcpPacketQueue, pkt TcpPacket) {
	conn, ok := m.find(pkt)
	if ok {
		conn.Pkt = pkt
	} else {
		conn = m.addConnection(pkt)
	}

	if pkt.TcpHeader.Flags.SYN && !ok {
		log.Printf("Received SYN Packet")

		queue.Write(conn, HeaderFlags{
			SYN: true,
			ACK: true,
		}, nil)

		m.update(pkt, StateSynReceived, false)
	}

	if ok && pkt.TcpHeader.Flags.ACK && conn.State == StateSynReceived {
		log.Printf("Received ACK Packet")
		m.update(pkt, StateEstablished, false)
	}

	if ok && pkt.TcpHeader.Flags.PSH && conn.State == StateEstablished {
		log.Printf("Received PSH Packet")

		queue.Write(conn, HeaderFlags{
			ACK: true,
		}, nil)
		m.update(pkt, StateEstablished, true)

		m.AcceptConnectionQueue <- conn
	}

	if ok && pkt.TcpHeader.Flags.FIN && conn.State == StateEstablished {
		log.Printf("Received FIN Packet")

		queue.Write(conn, HeaderFlags{
			ACK: true,
		}, nil)
		m.update(pkt, StateCloseWait, false)

		queue.Write(conn, HeaderFlags{
			FIN: true,
			ACK: true,
		}, nil)
		m.update(pkt, StateLastAck, false)
	}
}

func (m *ConnectionManager) addConnection(pkt TcpPacket) Connection {
	m.lock.Lock()
	defer m.lock.Unlock()
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	conn := Connection{
		SrcPort:         pkt.TcpHeader.SrcPort,
		DstPort:         pkt.TcpHeader.DstPort,
		State:           StateSynReceived,
		N:               pkt.Packet.N,
		Pkt:             pkt,
		initialSeqNum:   uint32(r.Int31()),
		incrementSeqNum: 0,
	}
	m.Connections = append(m.Connections, conn)

	return conn
}

func (m *ConnectionManager) remove(pkt TcpPacket) {
	m.lock.Lock()
	defer m.lock.Unlock()
}

func (m *ConnectionManager) update(pkt TcpPacket, state State, isAccept bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for i, conn := range m.Connections {
		if conn.SrcPort == pkt.TcpHeader.SrcPort && conn.DstPort == pkt.TcpHeader.DstPort {
			m.Connections[i].State = state
			m.Connections[i].isAccept = isAccept
			return
		}
	}
}
func find() {
	m.lock.Lock()
	defer m.lock.Unlock()
}
