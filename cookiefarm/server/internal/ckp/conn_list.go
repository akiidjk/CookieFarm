package ckp

import (
	"net"
	"sync"
)

type Connections struct {
	mu    sync.RWMutex
	conns []Connection
}

func (connections *Connections) Add(conn Connection) {
	connections.mu.Lock()
	defer connections.mu.Unlock()
	connections.conns = append(connections.conns, conn)
}

func (connections *Connections) GetAll() []Connection {
	connections.mu.RLock()
	defer connections.mu.RUnlock()
	res := make([]Connection, len(connections.conns))
	copy(res, connections.conns)
	return res
}

func (connections *Connections) Clear() {
	connections.mu.Lock()
	defer connections.mu.Unlock()
	connections.conns = nil
}

func (connections *Connections) Count() int {
	connections.mu.RLock()
	defer connections.mu.RUnlock()
	return len(connections.conns)
}

func cmpAddrs(conn1, conn2 *net.TCPAddr) bool {
	if conn1 == nil || conn2 == nil {
		return false
	}

	if conn1.IP.Equal(conn2.IP) && conn1.Port == conn2.Port {
		return true
	}

	return false
}

func (connections *Connections) Remove(conn Connection) {
	connections.mu.Lock()
	defer connections.mu.Unlock()
	for i, c := range connections.conns {
		if cmpAddrs(c.GetClientAddr(), conn.GetClientAddr()) {
			connections.conns = append(connections.conns[:i], connections.conns[i+1:]...)
			return
		}
	}
}
