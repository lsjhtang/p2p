package p2p

import (
	"errors"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddress: ":3000",
		HandlerFunc:   NOPHandlerFunc,
		Decoder:       &DefaultCode{},
	}
	TCPTransport := NewTCPTransport(opts)
	err := TCPTransport.ListenAndAccept()
	if err != nil {
		t.Error(err)
	}
	if TCPTransport.TCPTransportOpts.ListenAddress != opts.ListenAddress {
		t.Error(errors.New("listenAndAccept err"))
	}
}
