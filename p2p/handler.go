package p2p

type Handler interface {
	Handle(Peer) error
}

type HandlerFunc func(Peer) error

// 处理connet之前的逻辑
func NOPHandlerFunc(p Peer) error {

	return nil
}

type OnPeer func(Peer) error

// 处理connet之后的逻辑
func DefaultPeerFunc(p Peer) error {
	return nil
}
