package main

import (
	"log"

	"github.com/keis8221/tcp-ip-go/application"
)

func main() {
	s := application.NewServer()
	defer s.Close()
	s.ListenAndServe()

	for {
		conn, err := s.Accept()
		if err != nil {
			log.Printf("parse error: %s", err)
			continue
		}

		reqRaw := string(conn.Pkt.Packet.Buf[conn.Pkt.IpHeader.IHL*4+conn.Pkt.TcpHeader.DataOff*4:])
		req, err := application.ParseHttpRequest(reqRaw)
		if err != nil {
			log.Printf("parse error: %s", err)
			continue
		}

		log.Printf("request: %v", req)
		if req.Method == "GET" && req.URI == "/" {
			resp := application.NewHttpResponse(application.HttpStatusOK, "Hello, world!\r\n")
			s.Write(conn, resp)
		}
	}
}
