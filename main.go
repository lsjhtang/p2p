package main

import (
	"com.lsjhtang.p2p/boot"
	"com.lsjhtang.p2p/p2p"
	"com.lsjhtang.p2p/store"
	"log"
	"time"
)

func main() {
	TCPTransportOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandlerFunc:   p2p.NOPHandlerFunc,
		Decoder:       &p2p.DefaultCode{},
		OnPeer:        p2p.DefaultPeerFunc,
	}
	tcpTransport := p2p.NewTCPTransport(TCPTransportOpts)
	opts := boot.Opts{
		StoreRoot:            "root",
		PathTransportFromFun: store.CASPathTransportFromFun,
		Transport:            tcpTransport,
	}
	bs := boot.NewFileServer(opts)
	go func() {
		time.Sleep(3 * time.Second)
		bs.Stop()
	}()
	err := bs.Start()
	if err != nil {
		log.Fatal(err)
	}
}
