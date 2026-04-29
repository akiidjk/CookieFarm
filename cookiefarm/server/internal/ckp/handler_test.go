package ckp

/*
Category Partition Methodology - server/internal/ckp
====================================================

Step 1 - Identify Testable Utility Units
  (a) readUntilDelimiter(r, delim, maxSize)
      Splits a byte stream into one frame without using networking.
  (b) findString(data)
      Extracts a null-terminated string and its terminator index.
  (c) parse(data)
      Converts a CKP binary frame payload into database.Flag.
  (d) HandlerConfig(conn, cfg)
      Writes newline-delimited config bytes to a Connection.
  (e) Connections
      Tracks active CKP connections and removes them by client address.

Step 2 - Define Categories
  A) Delimiter:
     A1 delimiter present
     A2 delimiter absent with partial payload
     A3 empty delimiter
     A4 payload exceeds maxSize before delimiter

  B) Null-terminated strings:
     B1 non-empty string followed by NUL
     B2 empty string followed by NUL
     B3 missing NUL

  C) CKP payload shape:
     C1 valid header + flag string + exploit string
     C2 payload shorter than 8-byte fixed header
     C3 valid header but missing flag terminator
     C4 valid flag but missing exploit terminator

  D) Connection registry:
     D1 add and enumerate
     D2 remove matching client address
     D3 remove non-matching client address
     D4 clear

  E) Config write:
     E1 writable connection receives exact bytes

Step 3 - Constraints
  - parse depends on the server config singleton to map port -> service name,
    so tests seed the shared service map directly.
  - HandlerConfig is tested with an in-memory net.Conn implementation, not a
    real socket.
  - Connections.Remove compares TCP client IP and port only.

Step 4 - Test Frames
  - TestReadUntilDelimiter_CategoryPartition covers A1-A4.
  - TestFindString_CategoryPartition covers B1-B3.
  - TestParse_CategoryPartition covers C1-C4.
  - TestHandlerConfig_WritesConfigBytes covers E1.
  - TestConnections_CategoryPartition covers D1-D4.
*/

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"sharedconfig"
	"testing"
	"time"

	"server/config"
)

func TestReadUntilDelimiter_CategoryPartition(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		delim       []byte
		maxSize     int
		want        []byte
		wantErrIs   error
		wantErrText string
	}{
		{
			name:    "delimiter_present",
			input:   []byte{'a', 'b', 0xBB, 'T', 0xCC, 'x'},
			delim:   DelimiterBytes,
			maxSize: 1024,
			want:    []byte("ab"),
		},
		{
			name:      "delimiter_absent_with_partial_payload",
			input:     []byte("partial"),
			delim:     DelimiterBytes,
			maxSize:   1024,
			wantErrIs: io.ErrUnexpectedEOF,
		},
		{
			name:        "empty_delimiter",
			input:       []byte("anything"),
			delim:       nil,
			maxSize:     1024,
			wantErrText: "empty delimiter",
		},
		{
			name:        "message_too_large",
			input:       []byte("abcd"),
			delim:       DelimiterBytes,
			maxSize:     3,
			wantErrText: "message too large",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := readUntilDelimiter(bufio.NewReader(bytes.NewReader(tc.input)), tc.delim, tc.maxSize)

			if tc.wantErrIs != nil {
				if !errors.Is(err, tc.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tc.wantErrIs, err)
				}
				return
			}
			if tc.wantErrText != "" {
				if err == nil || err.Error() != tc.wantErrText {
					t.Fatalf("expected error %q, got %v", tc.wantErrText, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(got, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestFindString_CategoryPartition(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantString  string
		wantIndex   int
		wantErrText string
	}{
		{
			name:       "non_empty_string",
			input:      []byte{'f', 'l', 'a', 'g', 0, 'x'},
			wantString: "flag",
			wantIndex:  4,
		},
		{
			name:       "empty_string",
			input:      []byte{0, 'x'},
			wantString: "",
			wantIndex:  0,
		},
		{
			name:        "missing_terminator",
			input:       []byte("unterminated"),
			wantIndex:   -1,
			wantErrText: "unterminted flag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, idx, err := findString(tc.input)
			if tc.wantErrText != "" {
				if err == nil || err.Error() != tc.wantErrText {
					t.Fatalf("expected error %q, got %v", tc.wantErrText, err)
				}
				if idx != tc.wantIndex {
					t.Fatalf("expected index %d, got %d", tc.wantIndex, idx)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantString || idx != tc.wantIndex {
				t.Fatalf("expected (%q, %d), got (%q, %d)", tc.wantString, tc.wantIndex, got, idx)
			}
		})
	}
}

func TestParse_CategoryPartition(t *testing.T) {
	config.GetInstance().SetShared(sharedconfig.Shared{
		Services: map[string]uint16{"cookie": 31337},
	})

	tests := []struct {
		name        string
		payload     []byte
		wantErrText string
	}{
		{
			name:    "valid_payload",
			payload: buildServerPayload(123456789, 31337, 42, "FLAG{cookie}", "/tmp/exploit.py"),
		},
		{
			name:        "short_header",
			payload:     []byte{1, 2, 3},
			wantErrText: "invalid length",
		},
		{
			name:        "missing_flag_terminator",
			payload:     buildServerPayloadRaw(1, 31337, 42, []byte("FLAG{cookie}")),
			wantErrText: "unterminted flag",
		},
		{
			name:        "missing_exploit_terminator",
			payload:     buildServerPayloadRaw(1, 31337, 42, append([]byte("FLAG{cookie}\x00"), []byte("exploit.py")...)),
			wantErrText: "unterminted flag",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse(tc.payload)
			if tc.wantErrText != "" {
				if err == nil || err.Error() != tc.wantErrText {
					t.Fatalf("expected error %q, got %v", tc.wantErrText, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.SubmitTime != 123456789 {
				t.Errorf("SubmitTime expected 123456789, got %d", got.SubmitTime)
			}
			if got.PortService != 31337 {
				t.Errorf("PortService expected 31337, got %d", got.PortService)
			}
			if got.TeamID != 42 {
				t.Errorf("TeamID expected 42, got %d", got.TeamID)
			}
			if got.FlagCode != "FLAG{cookie}" {
				t.Errorf("FlagCode expected %q, got %q", "FLAG{cookie}", got.FlagCode)
			}
			if got.ExploitName != "/tmp/exploit.py" {
				t.Errorf("ExploitName expected %q, got %q", "/tmp/exploit.py", got.ExploitName)
			}
			if got.ServiceName != "cookie" {
				t.Errorf("ServiceName expected %q, got %q", "cookie", got.ServiceName)
			}
			if got.Msg != "Flag found for team: 42" {
				t.Errorf("Msg expected team message, got %q", got.Msg)
			}
		})
	}
}

func TestHandlerConfig_WritesConfigBytes(t *testing.T) {
	conn := &memoryConnection{
		clientAddr: &net.TCPAddr{IP: net.IPv4(10, 0, 0, 2), Port: 1111},
		serverAddr: &net.TCPAddr{IP: net.IPv4(10, 0, 0, 1), Port: 7777},
	}

	HandlerConfig(conn, []byte(`{"configured":true}`+"\n"))

	if got := conn.String(); got != "{\"configured\":true}\n" {
		t.Fatalf("expected exact config bytes, got %q", got)
	}
}

func TestConnections_CategoryPartition(t *testing.T) {
	c1 := newFakeConnection("10.0.0.2", 1000)
	c2 := newFakeConnection("10.0.0.3", 1001)
	missing := newFakeConnection("10.0.0.4", 1002)

	conns := &Connections{}
	conns.Add(c1)
	conns.Add(c2)

	if conns.Count() != 2 {
		t.Fatalf("expected 2 connections after add, got %d", conns.Count())
	}
	if len(conns.GetAll()) != 2 {
		t.Fatalf("expected GetAll length 2, got %d", len(conns.GetAll()))
	}

	conns.Remove(missing)
	if conns.Count() != 2 {
		t.Fatalf("removing non-matching connection changed count to %d", conns.Count())
	}

	conns.Remove(c1)
	if conns.Count() != 1 {
		t.Fatalf("expected 1 connection after remove, got %d", conns.Count())
	}
	if !cmpAddrs(conns.GetAll()[0].GetClientAddr(), c2.GetClientAddr()) {
		t.Fatalf("remaining connection is not c2")
	}

	conns.Clear()
	if conns.Count() != 0 || len(conns.GetAll()) != 0 {
		t.Fatalf("expected no connections after clear, got count=%d len=%d", conns.Count(), len(conns.GetAll()))
	}
}

func buildServerPayload(submitTime uint32, port, teamID uint16, flagCode, exploitName string) []byte {
	body := make([]byte, 0, 8+len(flagCode)+1+len(exploitName)+1)
	body = appendUint32LE(body, submitTime)
	body = appendUint16LE(body, port)
	body = appendUint16LE(body, teamID)
	body = append(body, flagCode...)
	body = append(body, 0)
	body = append(body, exploitName...)
	body = append(body, 0)
	return body
}

func buildServerPayloadRaw(submitTime uint32, port, teamID uint16, rest []byte) []byte {
	body := make([]byte, 0, 8+len(rest))
	body = appendUint32LE(body, submitTime)
	body = appendUint16LE(body, port)
	body = appendUint16LE(body, teamID)
	body = append(body, rest...)
	return body
}

func appendUint32LE(dst []byte, v uint32) []byte {
	return append(dst, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
}

func appendUint16LE(dst []byte, v uint16) []byte {
	return append(dst, byte(v), byte(v>>8))
}

type fakeConnection struct {
	*memoryConnection
}

func newFakeConnection(ip string, port int) *fakeConnection {
	return &fakeConnection{
		memoryConnection: &memoryConnection{
			clientAddr: &net.TCPAddr{IP: net.ParseIP(ip), Port: port},
			serverAddr: &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7777},
		},
	}
}

type memoryConnection struct {
	bytes.Buffer
	server     *Server
	ctx        *context.Context
	clientAddr *net.TCPAddr
	serverAddr *net.TCPAddr
}

func (*memoryConnection) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

func (c *memoryConnection) Write(p []byte) (int, error) {
	return c.Buffer.Write(p)
}

func (*memoryConnection) Close() error {
	return nil
}

func (c *memoryConnection) LocalAddr() net.Addr {
	return c.serverAddr
}

func (c *memoryConnection) RemoteAddr() net.Addr {
	return c.clientAddr
}

func (*memoryConnection) SetDeadline(_ time.Time) error {
	return nil
}

func (*memoryConnection) SetReadDeadline(_ time.Time) error {
	return nil
}

func (*memoryConnection) SetWriteDeadline(_ time.Time) error {
	return nil
}

func (c *memoryConnection) GetNetConn() net.Conn {
	return c
}

func (c *memoryConnection) GetServer() *Server {
	return c.server
}

func (c *memoryConnection) GetClientAddr() *net.TCPAddr {
	return c.clientAddr
}

func (c *memoryConnection) GetServerAddr() *net.TCPAddr {
	return c.serverAddr
}

func (c *memoryConnection) SetContext(ctx *context.Context) {
	c.ctx = ctx
}

func (c *memoryConnection) GetContext() *context.Context {
	if c.ctx == nil {
		ctx := context.Background()
		c.ctx = &ctx
	}
	return c.ctx
}

func (*memoryConnection) Reset(_ net.Conn) {}

func (c *memoryConnection) SetServer(server *Server) {
	c.server = server
}
