package ckp

import (
	"bufio"
	"errors"
	"logger"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn   *net.TCPConn
	reader *bufio.Reader
	writer *bufio.Writer
	mu     sync.Mutex // serialize writes if multiple goroutines write
}

func NewClient(addr string) (*Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	// Performance tuning
	conn.SetNoDelay(true)   // disable Nagle → no batching delay
	conn.SetKeepAlive(true) // detect dead connections
	conn.SetKeepAlivePeriod(30 * time.Second)
	conn.SetReadBuffer(65536)
	conn.SetWriteBuffer(65536)

	conn.SetKeepAliveConfig(net.KeepAliveConfig{
		Enable:   true,
		Idle:     15 * time.Second, // send first probe after 15s idle
		Interval: 5 * time.Second,  // probe every 5s after that
		Count:    3,                // drop after 3 unanswered probes
	})

	return &Client{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, 65536),
		writer: bufio.NewWriterSize(conn, 65536),
	}, nil
}

func (c *Client) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err := c.writer.Write(data); err != nil {
		return err
	}
	return c.writer.Flush() // flush bufio buffer to the socket
}

func (c *Client) Receive() ([]byte, error) {
	c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	return c.reader.ReadBytes('\n') // adjust delimiter to your protocol
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) reconnect(addr string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.Close()
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	conn.SetNoDelay(true)
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(30 * time.Second)
	conn.SetReadBuffer(65536)
	conn.SetWriteBuffer(65536)

	c.conn = conn
	c.reader = bufio.NewReaderSize(conn, 65536)
	c.writer = bufio.NewWriterSize(conn, 65536)
	return nil
}

func (c *Client) SendWithRetry(addr string, data []byte, maxRetries int) error {
	for i := range maxRetries {
		if err := c.Send(data); err != nil {
			logger.Log.Error().Err(err).Msgf("Error sending data to CKP server, attempt %d/%d", i+1, maxRetries)
			if rerr := c.reconnect(addr); rerr != nil {
				logger.Log.Error().Err(rerr).Msg("Reconnect failed, retrying with backoff")
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond) // backoff
				continue
			}
			continue
		}
		return nil
	}
	return errors.New("max retries reached")
}

func (c *Client) ReadPump() {
	reader := bufio.NewReaderSize(c.conn, 65536)
	for {
		response, err := reader.ReadBytes('\n')
		if err != nil {
			if rerr := c.reconnect(ADDR); rerr != nil {
				logger.Log.Error().Err(rerr).Msg("Reconnect failed, retrying in 1s")
				time.Sleep(1 * time.Second)
				continue
			}
			reader = bufio.NewReaderSize(c.conn, 65536)
			continue
		}
		logger.Log.Debug().Str("response", string(response)).Msg("Received response from CKP server")
		handleConfig(response)
	}
}
