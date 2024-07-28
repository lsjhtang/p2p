package p2p

import "net"

type Peer interface {
	net.Conn
	Send([]byte) error
}

type Transport interface {
	Dial(address string) error
	GetListenAddress() string
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
