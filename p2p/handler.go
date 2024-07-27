package p2p

type Handler interface {
	handle(peer) error
}

type HandlerFunc func(peer) error

// 处理connet之前的逻辑
func NOPHandlerFunc(p peer) error {

	return nil
}

type OnPeer func(peer) error

// 处理connet之后的逻辑
func DefaultPeerFunc(p peer) error {
	return nil
}
