package boot

import (
	"com.lsjhtang.p2p/p2p"
	"com.lsjhtang.p2p/store"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Opts struct {
	StoreRoot            string
	PathTransportFromFun store.PathTransportFromFun
	Transport            p2p.Transport
}

type FileServer struct {
	Opts
	store *store.FileStore

	quitCh chan os.Signal
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
	}
}

func (s *FileServer) Start() error {
	err := s.Transport.ListenAndAccept()
	if err != nil {
		return err
	}

	return s.loop()
}

func (s *FileServer) loop() error {
	defer func() {

	}()

	for {
		select {
		case <-s.quitCh:
			log.Printf("file server is closed")
			return s.Transport.Close()
		case msg := <-s.Transport.Consume():
			log.Printf("%+v\n", msg)
		}
	}
}

func (s *FileServer) Stop() {
	close(s.quitCh)
}
