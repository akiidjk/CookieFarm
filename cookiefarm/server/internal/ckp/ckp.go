package ckp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"pool"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	soReusePort = 0x0F
	tcpFastOpen = 0x17
)

type Server struct {
	listenAddr           *net.TCPAddr
	listener             *net.TCPListener
	shutdown             bool
	shutdownDeadline     time.Time
	requestHandler       RequestHandlerFunc
	connectionCreator    ConnectionCreatorFunc
	ctx                  *context.Context
	activeConnections    int32
	maxAcceptConnections int32
	acceptedConnections  int32
	listenConfig         *ListenConfig
	connWaitGroup        sync.WaitGroup
	connStructPool       sync.Pool
	wp                   *pool.WorkerPool[*net.TCPConn]
	allowThreadLocking   bool
	ballast              []byte
}

type ListenConfig struct {
	lc                     net.ListenConfig
	SocketReusePort        bool
	SocketFastOpen         bool
	SocketFastOpenQueueLen int
	SocketDeferAccept      bool
}

type RequestHandlerFunc func(conn Connection)

type ConnectionCreatorFunc func() Connection

type controlFunc func(network, address string, c syscall.RawConn) error

func applyListenSocketOptions(lc *ListenConfig) controlFunc {
	return func(network, address string, c syscall.RawConn) error {
		var err error

		c.Control(func(fd uintptr) {
			if lc.SocketReusePort {
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, soReusePort, 1)
			}
			if lc.SocketFastOpen {
				qlen := lc.SocketFastOpenQueueLen
				if qlen <= 0 {
					qlen = 256
				}
				err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, tcpFastOpen, qlen)
			}
			if lc.SocketDeferAccept {
				err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, 1)
			}
		})
		return err
	}
}

var defaultListenConfig *ListenConfig = &ListenConfig{
	SocketReusePort: true,
}

func NewServer(listenAddr string) (*Server, error) {
	la, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	var s *Server

	s = &Server{
		listenAddr:   la,
		listenConfig: defaultListenConfig,
		connStructPool: sync.Pool{
			New: func() interface{} {
				conn := s.connectionCreator()
				conn.SetServer(s)
				return conn
			},
		},
	}

	s.connectionCreator = func() Connection {
		return &TCPConn{}
	}

	s.SetBallast(20)

	return s, nil
}

func (s *Server) SetListenConfig(config *ListenConfig) {
	s.listenConfig = config
}

func (s *Server) GetListenConfig() *ListenConfig {
	return s.listenConfig
}

func (s *Server) Listen() error {
	network := "tcp4"
	if isIPv6Addr(s.listenAddr) {
		network = "tcp6"
	}

	s.listenConfig.lc.Control = applyListenSocketOptions(s.listenConfig)
	l, err := s.listenConfig.lc.Listen(*s.GetContext(), network, s.listenAddr.String())
	if err != nil {
		return err
	}

	tcpl, ok := l.(*net.TCPListener)
	if !ok {
		return errors.New("listener must be of type net.TCPListener")
	}

	s.listener = tcpl
	return nil
}

func (s *Server) SetMaxAcceptConnections(limit int32) {
	atomic.StoreInt32(&s.maxAcceptConnections, limit)
}

func (s *Server) GetActiveConnections() int32 {
	return s.activeConnections
}

func (s *Server) GetAcceptedConnections() int32 {
	return s.acceptedConnections
}

func (s *Server) GetListenAddr() *net.TCPAddr {
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr().(*net.TCPAddr)
}

func (s *Server) Shutdown(d time.Duration) (err error) {
	s.shutdownDeadline = time.Time{}
	if d > 0 {
		s.shutdownDeadline = time.Now().Add(d)
	}
	s.shutdown = true
	err = s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Halt() (err error) {
	return s.Shutdown(-1 * time.Second)
}

func (s *Server) Serve() error {
	maxProcs := runtime.GOMAXPROCS(0)

	s.wp = pool.NewWorkerPool(s.serveConn)
	s.wp.SetNumShards(maxProcs * 2)
	s.wp.SetIdleWorkerLifetime(5 * time.Second)
	s.wp.Start()
	defer s.wp.Stop()

	if s.activeConnections == 0 {
		return nil
	}

	if s.shutdownDeadline.IsZero() {
		s.connWaitGroup.Wait()
	} else {
		diff := s.shutdownDeadline.Sub(time.Now())
		if diff > 0 {
			time.Sleep(diff)
		}
	}

	return nil
}

func (s *Server) SetRequestHandler(requestHandler RequestHandlerFunc) {
	s.requestHandler = requestHandler
}

func (s *Server) SetConnectionCreator(f ConnectionCreatorFunc) {
	s.connectionCreator = f
}

func (s *Server) SetContext(ctx *context.Context) {
	s.ctx = ctx
}

func (s *Server) GetContext() *context.Context {
	if s.ctx == nil {
		ctx := context.Background()
		s.ctx = &ctx
	}
	return s.ctx
}

func (s *Server) SetAllowThreadLocking(allow bool) {
	s.allowThreadLocking = allow
}

func (s *Server) SetBallast(sizeInMiB int) {
	s.ballast = make([]byte, sizeInMiB*1024*1024)
}

func (s *Server) acceptLoop(id int) error {
	for {
		if s.maxAcceptConnections > 0 && s.acceptedConnections >= s.maxAcceptConnections {
			s.Shutdown(0)
		}

		if s.shutdown {
			_ = s.listener.Close()
			break
		}

		tcpConn, err := s.listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if !(opErr.Temporary() && opErr.Timeout()) && s.shutdown {
					break
				}
			}

			s.listener.Close()
			return err
		}

		newAcceptedConns := atomic.AddInt32(&s.acceptedConnections, 1)
		if s.maxAcceptConnections > 0 && newAcceptedConns > s.maxAcceptConnections {
			tcpConn.Close()
			continue
		}

		s.wp.AddTask(tcpConn)
		go s.serveConn(tcpConn)
		tcpConn = nil
	}

	return nil
}

func (s *Server) serveConn(netConn *net.TCPConn) {
	conn := s.connStructPool.Get().(*TCPConn)

	atomic.AddInt32(&s.activeConnections, 1)

	conn.Reset(netConn)
	conn.Start()
	s.requestHandler(conn)
	conn.Close()
	atomic.AddInt32(&s.activeConnections, -1)

	s.connStructPool.Put(conn)
}

func isIPv6Addr(addr *net.TCPAddr) bool {
	return addr.IP.To4() == nil && len(addr.IP) == net.IPv6len
}

func StartServer(port uint16) error {
	s, err := NewServer(fmt.Sprintf(":%u", port))
	if err != nil {
		return err
	}

	s.SetListenConfig(&ListenConfig{
		SocketReusePort:   true,
		SocketFastOpen:    false,
		SocketDeferAccept: false,
	})

	s.SetRequestHandler(handler)
	s.SetAllowThreadLocking(true)

	err = s.Listen()
	if err != nil {
		return err
	}

	return s.Serve()
}
