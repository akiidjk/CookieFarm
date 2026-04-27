package ckp

import (
	"context"
	"fmt"
	"net"
	"sync"

	"pool"
)

const SIZE int = 100

type Server struct {
	listener net.Listener
	maxConns int
	connSem  chan struct{}
	wg       sync.WaitGroup
}

func NewServer(port uint16, maxConns int) (*Server, error) {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%u", port))
    if err != nil {
        return nil, err
    }

    return &Server{
        listener:    listener,
        maxConns:    maxConns,
        connSem:     make(chan struct{}, maxConns),
    }, nil
}

func (s *Server) Serve(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        conn, err := s.listener.Accept()
        if err != nil {
            if ne, ok := err.(net.Error); ok && ne.Timeout() {
                continue
            }
            return err
        }

        select {
        case s.connSem <- struct{}{}:
            s.wg.Add(1)
            go s.handleConnection(conn)
        default:
            conn.Close()
        }
    }
}

func readMessage(conn net.Conn) ([]byte, error) {
	return nil, nil
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
}

func StartServer(port uint16) error {
	s, err := NewServer(port, SIZE)
	if err != nil {
		return err
	}
	
	return s.Serve()
}
