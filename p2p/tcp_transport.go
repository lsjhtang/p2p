package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{conn: conn, outbound: outbound}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddress string
	HandlerFunc   HandlerFunc
	Decoder       Coder
	OnPeer        OnPeer
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	RPCch    chan RPC

	mu    sync.RWMutex
	peers map[net.Addr]peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		RPCch:            make(chan RPC),
		mu:               sync.RWMutex{},
		peers:            map[net.Addr]peer{},
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.acceptLoop()

	return nil
}
func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Println("listener accept error:", err)
		}

		go t.handleConn(conn, true)
	}
}
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound)
	if t.HandlerFunc != nil {
		if err := t.HandlerFunc(peer); err != nil {
			if outbound {
				log.Printf("TCP handlerFunc error:%v\n", err)
				_ = conn.Close()
				return
			}
		}
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			log.Printf("TCP read err: %v\n", err)
			return
		}
		rpc.From = conn.RemoteAddr()
		t.RPCch <- *rpc
	}

}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.RPCch
}
