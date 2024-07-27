package main

import (
	"com.lsjhtang.p2p/p2p"
	"log"
)

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandlerFunc:   p2p.NOPHandlerFunc,
		Decoder:       &p2p.DefaultCode{},
		OnPeer:        p2p.DefaultPeerFunc,
	}
	tcpTransport := p2p.NewTCPTransport(opts)
	err := tcpTransport.ListenAndAccept()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg := <-tcpTransport.RPCch
			log.Printf("%+v\n", msg)
		}
	}()

	select {}
}
