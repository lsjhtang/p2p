package main

import (
	"bytes"
	"com.lsjhtang.p2p/boot"
	"com.lsjhtang.p2p/p2p"
	"com.lsjhtang.p2p/store"
	"log"
	"path"
	"strings"
	"time"
)

func main() {
	go func() {
		time.Sleep(2 * time.Second)
		bs := makeFileServer(":4000", ":3000")
		go func() {
			err := bs.Start()
			if err != nil {
				log.Fatal(err)
			}
			//time.Sleep(3 * time.Second)
			//bs.Stop()
		}()

		time.Sleep(2 * time.Second)
		buf := bytes.NewBuffer([]byte("send msg 4000"))
		bs.StoreData("lsjhtang", buf)
		select {}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		bs := makeFileServer(":5000", ":3000", ":4000")
		go func() {
			err := bs.Start()
			if err != nil {
				log.Fatal(err)
			}
			//time.Sleep(3 * time.Second)
			//bs.Stop()
		}()

		time.Sleep(2 * time.Second)
		buf := bytes.NewBuffer([]byte("send msg 5000"))
		bs.StoreData("lsjhtang", buf)
		select {}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		bs := makeFileServer(":6000", ":3000", ":5000", ":4000")
		go func() {
			err := bs.Start()
			if err != nil {
				log.Fatal(err)
			}
			//time.Sleep(3 * time.Second)
			//bs.Stop()
		}()

		time.Sleep(2 * time.Second)
		buf := bytes.NewBuffer([]byte("send msg 6000"))
		bs.StoreData("lsjhtang", buf)
		select {}
	}()

	bs := makeFileServer(":3000")
	//go func() {
	//	time.Sleep(3 * time.Second)
	//	bs.Stop()
	//}()
	err := bs.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func makeFileServer(addr string, node ...string) *boot.FileServer {
	TCPTransportOpts := p2p.TCPTransportOpts{
		ListenAddress: addr,
		HandlerFunc:   p2p.NOPHandlerFunc,
		Decoder:       &p2p.GobCode{},
		OnPeer:        p2p.DefaultPeerFunc,
	}
	tcpTransport := p2p.NewTCPTransport(TCPTransportOpts)
	opts := boot.Opts{
		StoreRoot:            path.Join("root", strings.ReplaceAll(addr, ":", "_")),
		PathTransportFromFun: store.CASPathTransportFromFun,
		Transport:            tcpTransport,
		BootStrapNodes:       node,
	}
	bs := boot.NewFileServer(opts)
	tcpTransport.OnPeer = bs.OnPeer

	return bs
}
