package pool

import (
	"net"
	"sync"
)

type Pool struct {
	mu sync.Mutex
    conns chan net.Conn
	handle func(net.Conn)
}

func NewPool(size int, handler func(net.Conn)) *Pool {
    pool := &Pool{
		conns: make(chan net.Conn, 100),
		handle: handler,
	}
	
    return pool
}

func (p *Pool) Submit(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case p.conns <- conn:
	default:
		conn.Close()
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
}
