package p2p

type peer interface {
	Close() error
}

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
