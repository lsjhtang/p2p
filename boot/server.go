package boot

import (
	"bytes"
	"com.lsjhtang.p2p/p2p"
	"com.lsjhtang.p2p/store"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Opts struct {
	StoreRoot            string
	PathTransportFromFun store.PathTransportFromFun
	Transport            p2p.Transport
	BootStrapNodes       []string
}

type FileServer struct {
	Opts
	store  *store.FileStore
	quitCh chan os.Signal

	peerLock sync.RWMutex
	peers    map[string]p2p.Peer
}

func NewFileServer(opts Opts) *FileServer {
	if len(opts.StoreRoot) == 0 {
		opts.StoreRoot = "root"
	}
	if opts.PathTransportFromFun == nil {
		opts.PathTransportFromFun = store.CASPathTransportFromFun
	}

	if opts.Transport == nil {
		opts.Transport = p2p.NewTCPTransport(p2p.DefaultTCPTransportOpts())
	}

	if opts.BootStrapNodes == nil {
		opts.BootStrapNodes = make([]string, 0)
	}

	storeOpts := store.Opts{
		Root:                 opts.StoreRoot,
		PathTransportFromFun: opts.PathTransportFromFun,
	}
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return &FileServer{
		Opts:   opts,
		store:  store.NewFileStore(storeOpts),
		quitCh: quitCh,
		peers:  make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Start() error {
	err := s.Transport.ListenAndAccept()
	if err != nil {
		return err
	}

	err = s.bootStrapNetwork()
	if err != nil {
		_ = s.Transport.Close()
		return err
	}

	return s.loop()
}

func (s *FileServer) loop() error {
	for {
		select {
		case <-s.quitCh:
			log.Printf("file server is closed")
			return s.Transport.Close()
		case msg := <-s.Transport.Consume():
			log.Printf("local:%s, receive: %+v %s\n", s.Transport.GetListenAddress(), msg, msg.Payload)
		}
	}
}

func (s *FileServer) Stop() {
	close(s.quitCh)
}

func (s *FileServer) bootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}

		if err := s.Transport.Dial(addr); err != nil {
			return err
		}
		log.Printf("%s attemp to connet wiht remote addr: %s\n", s.Transport.GetListenAddress(), addr)
	}

	return nil
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[peer.RemoteAddr().String()] = peer
	return nil
}

func (s *FileServer) broadcast(msg *p2p.RPC) error {
	//todo encode
	code := &p2p.GobCode{}
	for _, peer := range s.peers {
		err := code.Encode(peer, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	err := s.store.Write(key, r)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)
	if err != nil {
		return err
	}
	msg := &p2p.RPC{
		From:    s.Transport.GetListenAddress(),
		Payload: buf.Bytes(),
	}
	return s.broadcast(msg)
}
