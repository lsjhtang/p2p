package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPPeer struct {
	net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{Conn: conn, outbound: outbound}
}

func (p *TCPPeer) Send(msg []byte) error {
	_, err := p.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

type TCPTransportOpts struct {
	ListenAddress string
	HandlerFunc   HandlerFunc
	Decoder       Coder
	OnPeer        OnPeer
}

func DefaultTCPTransportOpts() TCPTransportOpts {
	return TCPTransportOpts{
		ListenAddress: ":3000",
		HandlerFunc:   NOPHandlerFunc,
		Decoder:       &DefaultCode{},
		OnPeer:        DefaultPeerFunc,
	}
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	RPCch    chan RPC

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		RPCch:            make(chan RPC),
		mu:               sync.RWMutex{},
		peers:            make(map[net.Addr]Peer),
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

func (t *TCPTransport) GetListenAddress() string {
	return t.ListenAddress
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
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

	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			log.Printf("TCP read err: %v\n", err)
			return
		}
		t.RPCch <- rpc
	}

}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.RPCch
}

func (t *TCPTransport) Close() error {
	close(t.RPCch)
	return t.listener.Close()
}
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}
