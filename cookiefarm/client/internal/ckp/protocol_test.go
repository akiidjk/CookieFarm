package ckp

/*
Category Partition Methodology - client/internal/ckp/protocol.go
================================================================

Goal
----
Test only CKP client utility behavior. These tests do not dial TCP, do not
start a CKP server, and do not emulate end-to-end connectivity.

Testable Units
--------------
1. buildPayload(flag database.Flag) []byte
   Encodes one flag into the client-to-server CKP binary frame.

2. handleConfig(payload []byte) error
   Decodes one server-to-client shared config payload and updates the client
   config singleton.

Categories
----------
A) Fixed-width numeric fields in buildPayload:
   A1 submit_time is encoded as uint32 little-endian.
   A2 port_service is encoded as uint16 little-endian.
   A3 team_id is encoded as uint16 little-endian.
   A4 values wider than their CKP field are truncated by the current encoder.

B) Variable-width string fields in buildPayload:
   B1 flag_code is copied as bytes and followed by one NUL byte.
   B2 exploit_name is filepath.Base(flag.ExploitName), not the full path.
   B3 exploit_name is followed by one NUL byte.
   B4 an empty flag_code still emits its NUL terminator.
   B5 an empty exploit_name follows filepath.Base("") and is encoded as ".".

C) Frame termination:
   C1 every built payload ends with DelimiterBytes.
   C2 delimiter bytes are emitted after both NUL-terminated strings.

D) Config payload validity in handleConfig:
   D1 valid JSON replaces Shared config.
   D2 valid JSON invokes OnNewConfig when callback is present.
   D3 valid JSON succeeds when OnNewConfig is nil.
   D4 invalid JSON returns an error.
   D5 invalid JSON does not invoke OnNewConfig.
   D6 invalid JSON leaves the previous Shared config unchanged.

Constraints
-----------
- buildPayload is deterministic for a given database.Flag.
- handleConfig mutates global package state: OnNewConfig and the client config
  singleton are reset for every subtest.
- OnNewConfig is invoked asynchronously by handleConfig, so callback assertions
  use a bounded channel wait.
*/

import (
	"bytes"
	"encoding/binary"
	"path/filepath"
	"server/database"
	"sharedconfig"
	"testing"
	"time"

	"client/config"
)

func TestBuildPayload_EncodesCanonicalFrame(t *testing.T) {
	flag := database.Flag{
		SubmitTime:  0x01020304,
		PortService: 0x1122,
		TeamID:      0x3344,
		FlagCode:    "FLAG{cookie}",
		ExploitName: filepath.Join("tmp", "nested", "exploit.py"),
	}

	payload := buildPayload(flag)

	assertPayloadHeader(t, payload, 0x01020304, 0x1122, 0x3344)
	assertPayloadTail(t, payload, []byte("FLAG{cookie}\x00exploit.py\x00"))
}

func TestBuildPayload_UsesExploitBasename(t *testing.T) {
	payload := buildPayload(database.Flag{
		FlagCode:    "FLAG{x}",
		ExploitName: filepath.Join("var", "tmp", "owned.py"),
	})

	assertPayloadTail(t, payload, []byte("FLAG{x}\x00owned.py\x00"))
}

func TestBuildPayload_EncodesEmptyStringsAccordingToCurrentImplementation(t *testing.T) {
	payload := buildPayload(database.Flag{})

	assertPayloadHeader(t, payload, 0, 0, 0)
	assertPayloadTail(t, payload, []byte("\x00.\x00"))
}

func TestBuildPayload_TruncatesFieldsWiderThanProtocol(t *testing.T) {
	payload := buildPayload(database.Flag{
		SubmitTime:  uint64(1)<<32 + 0xAABBCCDD,
		PortService: 0x7788,
		TeamID:      int64(1)<<16 + 0x1234,
		FlagCode:    "F",
		ExploitName: "E",
	})

	assertPayloadHeader(t, payload, 0xAABBCCDD, 0x7788, 0x1234)
	assertPayloadTail(t, payload, []byte("F\x00E\x00"))
}

func TestBuildPayload_SizeFormula(t *testing.T) {
	flag := database.Flag{
		SubmitTime:  1,
		PortService: 2,
		TeamID:      3,
		FlagCode:    "abc",
		ExploitName: filepath.Join("dir", "tool.py"),
	}

	payload := buildPayload(flag)
	wantLen := SizeTimestamp + SizePort + SizeTeamID +
		len("abc") + 1 +
		len("tool.py") + 1 +
		DelimiterSize

	if len(payload) != wantLen {
		t.Fatalf("payload length: want %d, got %d", wantLen, len(payload))
	}
}

func TestHandleConfig_ValidJSONUpdatesSharedConfigAndInvokesCallback(t *testing.T) {
	resetProtocolTestState(t)
	called := make(chan struct{}, 1)
	OnNewConfig = func() { called <- struct{}{} }

	err := handleConfig([]byte(`{"services":{"web":8080,"db":5432},"regex_flag":"FLAG\\{[A-Z0-9]+\\}","format_ip_teams":"10.10.{}.1","my_team_id":7,"url_flag_ids":"http://ids.local","nop_team":1,"range_ip_teams":20,"configured":true}`)) //nolint
	if err != nil {
		t.Fatalf("handleConfig returned unexpected error: %v", err)
	}

	got := config.GetInstance().Get().Shared
	if got.Services["web"] != 8080 || got.Services["db"] != 5432 {
		t.Fatalf("services were not updated correctly: %+v", got.Services)
	}
	if got.RegexFlag != `FLAG\{[A-Z0-9]+\}` {
		t.Fatalf("regex_flag: want %q, got %q", `FLAG\{[A-Z0-9]+\}`, got.RegexFlag)
	}
	if got.FormatIPTeams != "10.10.{}.1" {
		t.Fatalf("format_ip_teams: want %q, got %q", "10.10.{}.1", got.FormatIPTeams)
	}
	if got.MyTeamID != 7 || got.URLFlagIds != "http://ids.local" || got.NOPTeam != 1 || got.RangeIPTeams != 20 || !got.Configured {
		t.Fatalf("shared config fields were not updated correctly: %+v", got)
	}

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("OnNewConfig was not invoked")
	}
}

func TestHandleConfig_ValidJSONWithNilCallbackStillSucceeds(t *testing.T) {
	resetProtocolTestState(t)
	OnNewConfig = nil

	err := handleConfig([]byte(`{"services":{"pwn":31337},"configured":true}`))
	if err != nil {
		t.Fatalf("handleConfig returned unexpected error: %v", err)
	}

	got := config.GetInstance().Get().Shared
	if got.Services["pwn"] != 31337 {
		t.Fatalf("expected pwn service to be updated, got %+v", got.Services)
	}
	if !got.Configured {
		t.Fatal("expected configured=true")
	}
}

func TestHandleConfig_InvalidJSONReturnsErrorPreservesConfigAndSkipsCallback(t *testing.T) {
	resetProtocolTestState(t)
	called := make(chan struct{}, 1)
	OnNewConfig = func() { called <- struct{}{} }

	err := handleConfig([]byte(`{"services":`))
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}

	got := config.GetInstance().Get().Shared
	if got.Services["original"] != 1 || got.RegexFlag != "ORIGINAL" || got.Configured {
		t.Fatalf("invalid JSON should preserve previous config, got %+v", got)
	}

	select {
	case <-called:
		t.Fatal("OnNewConfig must not be invoked after invalid JSON")
	case <-time.After(25 * time.Millisecond):
	}
}

func assertPayloadHeader(t *testing.T, payload []byte, wantSubmit uint32, wantPort uint16, wantTeam uint16) {
	t.Helper()

	if len(payload) < SizeTimestamp+SizePort+SizeTeamID+DelimiterSize {
		t.Fatalf("payload too short: %d bytes", len(payload))
	}
	if got := binary.LittleEndian.Uint32(payload[0:4]); got != wantSubmit {
		t.Fatalf("submit_time: want %#x, got %#x", wantSubmit, got)
	}
	if got := binary.LittleEndian.Uint16(payload[4:6]); got != wantPort {
		t.Fatalf("port_service: want %#x, got %#x", wantPort, got)
	}
	if got := binary.LittleEndian.Uint16(payload[6:8]); got != wantTeam {
		t.Fatalf("team_id: want %#x, got %#x", wantTeam, got)
	}
}

func assertPayloadTail(t *testing.T, payload []byte, wantBody []byte) {
	t.Helper()

	offset := SizeTimestamp + SizePort + SizeTeamID
	want := append(append([]byte(nil), wantBody...), DelimiterBytes...)
	got := payload[offset:]
	if !bytes.Equal(got, want) {
		t.Fatalf("payload body+delimiter:\nwant %v\n got %v", want, got)
	}
}

func resetProtocolTestState(t *testing.T) {
	t.Helper()

	OnNewConfig = nil
	config.GetInstance().SetSharedConfig(sharedconfig.Shared{
		Services:   map[string]uint16{"original": 1},
		RegexFlag:  "ORIGINAL",
		Configured: false,
	})

	t.Cleanup(func() {
		OnNewConfig = nil
		config.GetInstance().SetSharedConfig(sharedconfig.Shared{
			Services: map[string]uint16{},
		})
	})
}
