package ckp

import (
	"context"
	"net"
	"time"
)

type Connection interface {
	net.Conn
	GetNetConn() net.Conn
	GetServer() *Server
	GetClientAddr() *net.TCPAddr
	GetServerAddr() *net.TCPAddr
	GetStartTime() time.Time
	SetContext(ctx *context.Context)
	GetContext() *context.Context

	Start()
	Reset(netConn net.Conn)
	SetServer(server *Server)
}

type TCPConn struct {
	net.Conn
	server            *Server
	ctx               *context.Context
	ts                int64
	_cacheLinePadding [24]byte
}

func (conn *TCPConn) GetClientAddr() *net.TCPAddr {
	return conn.RemoteAddr().(*net.TCPAddr)
}

func (conn *TCPConn) GetServerAddr() *net.TCPAddr {
	return conn.LocalAddr().(*net.TCPAddr)
}

func (conn *TCPConn) GetStartTime() time.Time {
	return time.Unix(conn.ts/1e9, conn.ts%1e9)
}

func (conn *TCPConn) SetContext(ctx *context.Context) {
	conn.ctx = ctx
}

func (conn *TCPConn) GetContext() *context.Context {
	if conn.ctx == nil {
		ctx := context.Background()
		conn.ctx = &ctx
	}
	return conn.ctx
}

func (conn *TCPConn) GetNetConn() net.Conn {
	return conn.Conn
}

func (conn *TCPConn) GetNetTCPConn() (c *net.TCPConn) {
	c, _ = conn.Conn.(*net.TCPConn)
	return
}

func (conn *TCPConn) SetServer(s *Server) {
	conn.server = s
}

func (conn *TCPConn) GetServer() *Server {
	return conn.server
}

func (conn *TCPConn) Reset(netConn net.Conn) {
	conn.Conn = netConn
	conn.ctx = nil
}

func (conn *TCPConn) Start() {
	conn.ts = time.Now().UnixNano()
}
