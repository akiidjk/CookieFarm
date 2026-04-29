package ckp

import "net"

type Connections struct {
	conns []Connection
}

func (connections *Connections) Add(conn Connection) {
	connections.conns = append(connections.conns, conn)
}

func (connections *Connections) GetAll() []Connection {
	return connections.conns
}

func (connections *Connections) Clear() {
	connections.conns = nil
}

func (connections *Connections) Count() int {
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
	for i, c := range connections.conns {
		if cmpAddrs(c.GetClientAddr(), conn.GetClientAddr()) {
			connections.conns = append(connections.conns[:i], connections.conns[i+1:]...)
			return
		}
	}
}
