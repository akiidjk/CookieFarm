package ckp

import (
	"context"
	"net"
)

type Connection interface {
	net.Conn
	GetNetConn() net.Conn
	GetServer() *Server
	GetClientAddr() *net.TCPAddr
	GetServerAddr() *net.TCPAddr
	SetContext(ctx *context.Context)
	GetContext() *context.Context

	Reset(netConn net.Conn)
	SetServer(server *Server)
}

type TCPConn struct {
	net.Conn
	server            *Server
	ctx               *context.Context
	ts                int64    //nolint
	_cacheLinePadding [24]byte //nolint
}

func (conn *TCPConn) GetClientAddr() *net.TCPAddr {
	return conn.RemoteAddr().(*net.TCPAddr)
}

func (conn *TCPConn) GetServerAddr() *net.TCPAddr {
	return conn.LocalAddr().(*net.TCPAddr)
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
	return c
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
