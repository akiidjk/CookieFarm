// =============================================================================
// Category Partition Methodology — config_test.go
// =============================================================================
//
// STEP 1 — IDENTIFY UNITS UNDER TEST
// ------------------------------------
// The following behaviours of ConfigManager are exercised:
//
//   1.  SetHost / GetHost
//   2.  SetPort / GetPort
//   3.  SetHTTPS / GetHTTPS
//   4.  SetUsername / GetUsername
//   5.  SetToken / GetToken
//   6.  MapServiceToPort
//   7.  MapPortToService
//   8.  SetSharedConfig   (data preservation + map-copy isolation)
//   9.  Reset             (file creation with default template)
//  10.  GetSession        (file present vs. absent)
//  11.  Logout            (file present vs. absent)
//  12.  WriteLocal + read (YAML write → disk → unmarshal round-trip)
//  13.  Concurrent SetHost (data-race safety on atomic.Value)
//  14.  Get snapshot immutability (previously obtained pointer not mutated)
//
// STEP 2 — PARAMETERS AND ENVIRONMENTAL CONDITIONS
// -------------------------------------------------
//  SetHost      : host     string
//  SetPort      : port     uint16
//  SetHTTPS     : https    bool
//  SetUsername  : username string
//  SetToken     : token    string
//
//  MapServiceToPort : service string    [env: Services map contents]
//  MapPortToService : port    uint16    [env: Services map contents]
//
//  SetSharedConfig  : sc sharedconfig.Shared
//                     [env: pointer aliasing risk — caller mutates map after call]
//
//  Reset      : [env: DefaultPath writable / pre-existing / non-existent]
//  GetSession : [env: session file present or absent]
//  Logout     : [env: session file present or absent]
//  WriteLocal : LocalConfig struct  [env: DefaultPath writable]
//
//  Concurrency : N goroutines × SetHost  [env: atomic.Value under concurrent load]
//
// STEP 3 — CATEGORIES FOR EACH PARAMETER
// ----------------------------------------
//
//  host
//    A1 — typical hostname      ("localhost", "example.com")
//    A2 — IP address literal    ("10.0.0.1", "192.168.1.100")
//    A3 — empty string          (boundary: zero value, valid Go string)
//
//  port (uint16)
//    B1 — typical well-known value  (8080, 443)
//    B2 — zero                      (lower boundary)
//    B3 — 65535                     (upper boundary, type maximum)
//
//  https (bool)
//    C1 — true
//    C2 — false  (test both toggle directions)
//
//  username
//    D1 — non-empty string  ("admin", "cookieguest")
//    D2 — empty string      (boundary)
//
//  token
//    E1 — non-empty string  (round-trip preserved)
//    E2 — empty string      (overwrite / reset boundary)
//
//  MapServiceToPort — service key
//    F1 — key present  → expected non-zero port
//    F2 — key absent   → 0  (Go map zero-value semantics)
//
//  MapPortToService — port value
//    G1 — port found in map   → service name
//    G2 — port not in map     → ""
//    G3 — multiple entries    → only the unique matching name returned
//
//  SetSharedConfig — Shared struct / Services map
//    H1 — non-nil map with entries  → all fields and entries preserved after set
//    H2 — nil Services field        → stored as non-nil empty map, no panic
//    H3 — map isolation             → mutating original map after Set does not
//                                     affect the stored copy (copyMap semantics)
//
//  Reset — DefaultPath filesystem state
//    I1 — path already exists  → writes client.yml containing default keys
//    I2 — path does not exist  → MkdirAll creates dir tree, then writes file
//
//  GetSession — session file
//    J1 — file present  → returns exact content, nil error
//    J2 — file absent   → returns "", non-nil error
//
//  Logout — session file
//    K1 — file present  → file removed from disk, success message returned
//    K2 — file absent   → returns non-nil error
//
//  WriteLocal + read round-trip — LocalConfig values
//    L1 — fully populated LocalConfig  → YAML round-trip identical
//    L2 — zero-value LocalConfig       → round-trips to zero struct
//
//  Concurrency
//    M1 — 50 goroutines calling SetHost simultaneously → no data race
//         (detected by `go test -race ./...`)
//
//  Snapshot immutability
//    N1 — pointer from Get() before mutation is not affected by later SetHost
//
// STEP 4 — CONSTRAINTS AND FRAME SELECTION
// -----------------------------------------
//  * uint16 bounds (0–65535) are enforced by the type system; no overflow test needed.
//  * Go map zero-value for uint16 is 0, making F2 an inherent language behaviour.
//  * DefaultPath MUST be replaced with t.TempDir() for every file-I/O test so
//    tests are hermetic and never touch the real ~/.config/cookiefarm.
//  * newTestManager() bypasses the package-level singleton for test isolation.
//  * GetInstance() singleton is intentionally absent from these tests; the
//    singleton pattern itself does not require functional re-testing here.
//  * copyMap(nil) ⇒ make(map[string]uint16, 0): non-nil empty map (supports H2).
//  * ShowLocalConfigContent reads "config.yml"; Reset writes "client.yml" — they
//    are distinct files. ShowLocalConfigContent is not separately tested here
//    because the underlying read() generic path is fully covered by L1/L2.
//  * WriteShared uses the same write[T] generic function as WriteLocal; it is
//    not duplicated since the tested invariant is identical.
//  * Concurrency test uses sync.WaitGroup; no time.Sleep assertions.
//  * Eliminated redundant cases: SetHTTPS(true)→SetHTTPS(true) (duplicate C1);
//    non-empty token followed by non-empty token (covered by E1 independently).
//
// STEP 5 — Test implementation follows.
// =============================================================================

package config

import (
	"os"
	"path/filepath"
	"sharedconfig"
	"strings"
	"sync"
	"testing"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// newTestManager constructs a fresh ConfigManager that bypasses the
// package-level singleton, giving each test an independent starting state.
func newTestManager() *ConfigManager {
	cm := &ConfigManager{}
	cm.state.Store(&RuntimeConfig{
		Local:  LocalConfig{},
		Shared: sharedconfig.Shared{Services: make(map[string]uint16)},
		Token:  "",
	})
	return cm
}

// writeSessionFile writes content to DefaultPath/session.
// It is only valid to call after DefaultPath has been redirected to a temp dir.
func writeSessionFile(t *testing.T, content string) {
	t.Helper()
	p := filepath.Join(DefaultPath, "session")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeSessionFile: %v", err)
	}
}

// overrideDefaultPath redirects DefaultPath to a fresh temp directory for the
// duration of one test and restores the original value via t.Cleanup.
func overrideDefaultPath(t *testing.T) {
	t.Helper()
	orig := DefaultPath
	DefaultPath = t.TempDir()
	t.Cleanup(func() { DefaultPath = orig })
}

// ---------------------------------------------------------------------------
// A — Host  (categories A1, A2, A3)
// ---------------------------------------------------------------------------

func TestSetGetHost(t *testing.T) {
	cases := []struct {
		name string
		host string
	}{
		// A1: typical hostname
		{name: "A1_localhost", host: "localhost"},
		{name: "A1_fqdn", host: "example.com"},
		// A2: IP address literal
		{name: "A2_ipv4", host: "10.0.0.1"},
		{name: "A2_private_ip", host: "192.168.1.100"},
		// A3: boundary — empty string
		{name: "A3_empty", host: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := newTestManager()
			cm.SetHost(tc.host)
			if got := cm.GetHost(); got != tc.host {
				t.Errorf("GetHost() = %q; want %q", got, tc.host)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// B — Port  (categories B1, B2, B3)
// ---------------------------------------------------------------------------

func TestSetGetPort(t *testing.T) {
	cases := []struct {
		name string
		port uint16
	}{
		// B1: typical values
		{name: "B1_http_alt", port: 8080},
		{name: "B1_https", port: 443},
		// B2: lower boundary
		{name: "B2_zero", port: 0},
		// B3: upper boundary (uint16 max)
		{name: "B3_max", port: 65535},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := newTestManager()
			cm.SetPort(tc.port)
			if got := cm.GetPort(); got != tc.port {
				t.Errorf("GetPort() = %d; want %d", got, tc.port)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// C — HTTPS  (categories C1, C2)
// ---------------------------------------------------------------------------

func TestSetGetHTTPS(t *testing.T) {
	cases := []struct {
		name  string
		value bool
	}{
		{name: "C1_true", value: true},
		{name: "C2_false", value: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := newTestManager()
			cm.SetHTTPS(tc.value)
			if got := cm.GetHTTPS(); got != tc.value {
				t.Errorf("GetHTTPS() = %v; want %v", got, tc.value)
			}
		})
	}
}

// TestSetHTTPS_Toggle exercises the C1↔C2 transition in both directions.
func TestSetHTTPS_Toggle(t *testing.T) {
	cm := newTestManager()

	cm.SetHTTPS(true)
	if !cm.GetHTTPS() {
		t.Fatal("expected HTTPS=true after SetHTTPS(true)")
	}

	cm.SetHTTPS(false)
	if cm.GetHTTPS() {
		t.Fatal("expected HTTPS=false after SetHTTPS(false)")
	}
}

// ---------------------------------------------------------------------------
// D — Username  (categories D1, D2)
// ---------------------------------------------------------------------------

func TestSetGetUsername(t *testing.T) {
	cases := []struct {
		name     string
		username string
	}{
		// D1: non-empty
		{name: "D1_admin", username: "admin"},
		{name: "D1_cookieguest", username: "cookieguest"},
		// D2: boundary
		{name: "D2_empty", username: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := newTestManager()
			cm.SetUsername(tc.username)
			if got := cm.GetUsername(); got != tc.username {
				t.Errorf("GetUsername() = %q; want %q", got, tc.username)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// E — Token  (categories E1, E2)
// ---------------------------------------------------------------------------

func TestSetGetToken(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		// E1: non-empty round-trip
		{name: "E1_jwt_like", token: "eyJhbGciOiJIUzI1NiJ9.payload.sig"},
		{name: "E1_simple", token: "secret-token-42"},
		// E2: boundary
		{name: "E2_empty", token: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := newTestManager()
			cm.SetToken(tc.token)
			if got := cm.GetToken(); got != tc.token {
				t.Errorf("GetToken() = %q; want %q", got, tc.token)
			}
		})
	}
}

// TestSetToken_Overwrite covers the E1 → E2 transition: replacing a non-empty
// token with an empty one must zero out the stored value.
func TestSetToken_Overwrite(t *testing.T) {
	cm := newTestManager()
	cm.SetToken("initial-token")
	cm.SetToken("")
	if got := cm.GetToken(); got != "" {
		t.Errorf("expected empty token after overwrite; got %q", got)
	}
}

// ---------------------------------------------------------------------------
// F — MapServiceToPort  (categories F1, F2)
// ---------------------------------------------------------------------------

func TestMapServiceToPort(t *testing.T) {
	cm := newTestManager()
	cm.SetSharedConfig(sharedconfig.Shared{
		Services: map[string]uint16{
			"http":  80,
			"https": 443,
			"ssh":   22,
		},
	})

	cases := []struct {
		name     string
		service  string
		wantPort uint16
	}{
		// F1: key present
		{name: "F1_http", service: "http", wantPort: 80},
		{name: "F1_https", service: "https", wantPort: 443},
		{name: "F1_ssh", service: "ssh", wantPort: 22},
		// F2: key absent → Go zero-value for uint16
		{name: "F2_unknown_service", service: "telnet", wantPort: 0},
		{name: "F2_empty_key", service: "", wantPort: 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cm.MapServiceToPort(tc.service)
			if got != tc.wantPort {
				t.Errorf("MapServiceToPort(%q) = %d; want %d", tc.service, got, tc.wantPort)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// G — MapPortToService  (categories G1, G2, G3)
// ---------------------------------------------------------------------------

func TestMapPortToService(t *testing.T) {
	cm := newTestManager()
	cm.SetSharedConfig(sharedconfig.Shared{
		Services: map[string]uint16{
			"http":  80,
			"https": 443,
			"ssh":   22,
		},
	})

	cases := []struct {
		name        string
		port        uint16
		wantService string
	}{
		// G1: port found
		{name: "G1_port_80", port: 80, wantService: "http"},
		{name: "G1_port_443", port: 443, wantService: "https"},
		{name: "G1_port_22", port: 22, wantService: "ssh"},
		// G2: port absent
		{name: "G2_unknown_port", port: 9999, wantService: ""},
		{name: "G2_zero_port", port: 0, wantService: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cm.MapPortToService(tc.port)
			if got != tc.wantService {
				t.Errorf("MapPortToService(%d) = %q; want %q", tc.port, got, tc.wantService)
			}
		})
	}
}

// TestMapPortToService_MultipleEntries covers G3: with many services all having
// distinct ports, each lookup must return the uniquely matching name.
func TestMapPortToService_MultipleEntries(t *testing.T) {
	entries := map[string]uint16{
		"alpha":   1001,
		"beta":    1002,
		"gamma":   1003,
		"delta":   1004,
		"epsilon": 1005,
	}

	cm := newTestManager()
	cm.SetSharedConfig(sharedconfig.Shared{Services: entries})

	for wantName, port := range entries {
		t.Run("G3_"+wantName, func(t *testing.T) {
			got := cm.MapPortToService(port)
			if got != wantName {
				t.Errorf("MapPortToService(%d) = %q; want %q", port, got, wantName)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// H — SetSharedConfig  (categories H1, H2, H3)
// ---------------------------------------------------------------------------

// TestSetSharedConfig_DataPreserved covers H1: all fields of a fully populated
// Shared struct survive the copy performed inside SetSharedConfig.
func TestSetSharedConfig_DataPreserved(t *testing.T) {
	cm := newTestManager()
	input := sharedconfig.Shared{
		Services:      map[string]uint16{"web": 8080, "db": 5432},
		RegexFlag:     `FLAG\{[^}]+\}`,
		FormatIPTeams: "10.60.%d.1",
		MyTeamID:      7,
		NOPTeam:       1,
		RangeIPTeams:  255,
		URLFlagIds:    "http://flags.example.com",
		Configured:    true,
	}

	cm.SetSharedConfig(input)
	got := cm.Get().Shared

	if got.Services["web"] != 8080 {
		t.Errorf("Services[web] = %d; want 8080", got.Services["web"])
	}
	if got.Services["db"] != 5432 {
		t.Errorf("Services[db] = %d; want 5432", got.Services["db"])
	}
	if got.RegexFlag != input.RegexFlag {
		t.Errorf("RegexFlag = %q; want %q", got.RegexFlag, input.RegexFlag)
	}
	if got.FormatIPTeams != input.FormatIPTeams {
		t.Errorf("FormatIPTeams = %q; want %q", got.FormatIPTeams, input.FormatIPTeams)
	}
	if got.MyTeamID != input.MyTeamID {
		t.Errorf("MyTeamID = %d; want %d", got.MyTeamID, input.MyTeamID)
	}
	if got.NOPTeam != input.NOPTeam {
		t.Errorf("NOPTeam = %d; want %d", got.NOPTeam, input.NOPTeam)
	}
	if got.RangeIPTeams != input.RangeIPTeams {
		t.Errorf("RangeIPTeams = %d; want %d", got.RangeIPTeams, input.RangeIPTeams)
	}
	if got.URLFlagIds != input.URLFlagIds {
		t.Errorf("URLFlagIds = %q; want %q", got.URLFlagIds, input.URLFlagIds)
	}
	if !got.Configured {
		t.Error("Configured = false; want true")
	}
}

// TestSetSharedConfig_NilServicesMap covers H2: passing a nil Services map
// must not panic and must result in a non-nil empty map being stored
// (copyMap(nil) ⇒ make(map[string]uint16, 0)).
func TestSetSharedConfig_NilServicesMap(t *testing.T) {
	cm := newTestManager()

	// Must not panic.
	cm.SetSharedConfig(sharedconfig.Shared{Services: nil})

	stored := cm.Get().Shared.Services
	if stored == nil {
		t.Error("Services map is nil after SetSharedConfig(nil map); want non-nil empty map")
	}
	if l := len(stored); l != 0 {
		t.Errorf("Services map len = %d; want 0", l)
	}
}

// TestSetSharedConfig_EmptyServicesMap covers H2 with an explicitly empty
// (non-nil but zero-length) map — same expectations as nil.
func TestSetSharedConfig_EmptyServicesMap(t *testing.T) {
	cm := newTestManager()
	cm.SetSharedConfig(sharedconfig.Shared{Services: map[string]uint16{}})

	stored := cm.Get().Shared.Services
	if stored == nil {
		t.Error("Services map is nil; want non-nil empty map")
	}
	if l := len(stored); l != 0 {
		t.Errorf("Services map len = %d; want 0", l)
	}
}

// TestSetSharedConfig_MapIsolation covers H3: mutations to the caller's map
// after SetSharedConfig must not affect the stored copy.
func TestSetSharedConfig_MapIsolation(t *testing.T) {
	cm := newTestManager()

	originalMap := map[string]uint16{"svc": 9000}
	cm.SetSharedConfig(sharedconfig.Shared{Services: originalMap})

	// Mutate the caller's reference — must not bleed into the stored copy.
	originalMap["svc"] = 1111
	originalMap["new_entry"] = 2222

	stored := cm.Get().Shared.Services
	if stored["svc"] != 9000 {
		t.Errorf("stored Services[svc] = %d after external mutation; want 9000 (isolation broken)", stored["svc"])
	}
	if _, exists := stored["new_entry"]; exists {
		t.Error("stored Services contains 'new_entry' after external mutation; map isolation broken")
	}
}

// ---------------------------------------------------------------------------
// I — Reset  (categories I1, I2)
// ---------------------------------------------------------------------------

// TestReset_CreatesFileWithDefaults covers I1: Reset on an already-existing
// DefaultPath must write client.yml containing all expected default keys.
func TestReset_CreatesFileWithDefaults(t *testing.T) {
	overrideDefaultPath(t) // I1: path already exists via TempDir

	cm := newTestManager()
	if err := cm.Reset(); err != nil {
		t.Fatalf("Reset() error: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(DefaultPath, "client.yml"))
	if err != nil {
		t.Fatalf("reading client.yml after Reset: %v", err)
	}

	content := string(raw)
	for _, want := range []string{"localhost", "8080", "cookieguest"} {
		if !strings.Contains(content, want) {
			t.Errorf("client.yml missing expected default value %q\nfull content:\n%s", want, content)
		}
	}
}

// TestReset_CreatesDirectoryIfAbsent covers I2: Reset when DefaultPath does not
// yet exist must call os.MkdirAll and then write the file successfully.
func TestReset_CreatesDirectoryIfAbsent(t *testing.T) {
	orig := DefaultPath
	// Point to a nested path that does not exist yet.
	DefaultPath = filepath.Join(t.TempDir(), "nested", "config", "cookiefarm")
	t.Cleanup(func() { DefaultPath = orig })

	cm := newTestManager()
	if err := cm.Reset(); err != nil {
		t.Fatalf("Reset() on non-existent path: %v", err)
	}

	if _, err := os.Stat(DefaultPath); os.IsNotExist(err) {
		t.Error("Reset() did not create DefaultPath directory")
	}

	if _, err := os.Stat(filepath.Join(DefaultPath, "client.yml")); os.IsNotExist(err) {
		t.Error("Reset() did not create client.yml")
	}
}

// ---------------------------------------------------------------------------
// J — GetSession  (categories J1, J2)
// ---------------------------------------------------------------------------

// TestGetSession_FilePresent covers J1: existing session file → content returned.
func TestGetSession_FilePresent(t *testing.T) {
	overrideDefaultPath(t)
	const wantToken = "my-session-token-abc123"
	writeSessionFile(t, wantToken)

	cm := newTestManager()
	got, err := cm.GetSession()
	if err != nil {
		t.Fatalf("GetSession() unexpected error: %v", err)
	}
	if got != wantToken {
		t.Errorf("GetSession() = %q; want %q", got, wantToken)
	}
}

// TestGetSession_FileAbsent covers J2: no session file → non-nil error, empty string.
func TestGetSession_FileAbsent(t *testing.T) {
	overrideDefaultPath(t)

	cm := newTestManager()
	got, err := cm.GetSession()
	if err == nil {
		t.Errorf("GetSession() expected error when session absent; got %q", got)
	}
	if got != "" {
		t.Errorf("GetSession() = %q on error path; want empty string", got)
	}
}

// ---------------------------------------------------------------------------
// K — Logout  (categories K1, K2)
// ---------------------------------------------------------------------------

// TestLogout_FilePresent covers K1: existing session file is removed and a
// non-empty success message is returned.
func TestLogout_FilePresent(t *testing.T) {
	overrideDefaultPath(t)
	writeSessionFile(t, "some-valid-token")

	cm := newTestManager()
	msg, err := cm.Logout()
	if err != nil {
		t.Fatalf("Logout() unexpected error: %v", err)
	}
	if msg == "" {
		t.Error("Logout() returned empty success message; want descriptive string")
	}

	sessionPath := filepath.Join(DefaultPath, "session")
	if _, statErr := os.Stat(sessionPath); !os.IsNotExist(statErr) {
		t.Error("session file still present on disk after Logout()")
	}
}

// TestLogout_FileAbsent covers K2: absent session file → non-nil error.
func TestLogout_FileAbsent(t *testing.T) {
	overrideDefaultPath(t)

	cm := newTestManager()
	_, err := cm.Logout()
	if err == nil {
		t.Error("Logout() expected error when session file absent; got nil")
	}
}

// ---------------------------------------------------------------------------
// L — WriteLocal + read round-trip  (categories L1, L2)
// ---------------------------------------------------------------------------

// TestWriteLocal_RoundTrip_Populated covers L1: a fully populated LocalConfig
// survives a YAML marshal–write–read–unmarshal cycle identically.
func TestWriteLocal_RoundTrip_Populated(t *testing.T) {
	overrideDefaultPath(t)

	cm := newTestManager()
	want := LocalConfig{
		Host:     "192.168.0.10",
		Username: "alice",
		Port:     9443,
		HTTPS:    true,
	}
	cm.SetHost(want.Host)
	cm.SetPort(want.Port)
	cm.SetHTTPS(want.HTTPS)
	cm.SetUsername(want.Username)

	if err := cm.WriteLocal(); err != nil {
		t.Fatalf("WriteLocal() error: %v", err)
	}

	var got LocalConfig
	if err := read(&got, "client.yml"); err != nil {
		t.Fatalf("read() after WriteLocal error: %v", err)
	}

	if got.Host != want.Host {
		t.Errorf("Host: got %q; want %q", got.Host, want.Host)
	}
	if got.Port != want.Port {
		t.Errorf("Port: got %d; want %d", got.Port, want.Port)
	}
	if got.HTTPS != want.HTTPS {
		t.Errorf("HTTPS: got %v; want %v", got.HTTPS, want.HTTPS)
	}
	if got.Username != want.Username {
		t.Errorf("Username: got %q; want %q", got.Username, want.Username)
	}
}

// TestWriteLocal_RoundTrip_ZeroValue covers L2: a zero-value LocalConfig
// round-trips without corruption.
func TestWriteLocal_RoundTrip_ZeroValue(t *testing.T) {
	overrideDefaultPath(t)

	cm := newTestManager() // all LocalConfig fields are zero by construction

	if err := cm.WriteLocal(); err != nil {
		t.Fatalf("WriteLocal() error: %v", err)
	}

	var got LocalConfig
	if err := read(&got, "client.yml"); err != nil {
		t.Fatalf("read() after WriteLocal error: %v", err)
	}

	zero := LocalConfig{}
	if got != zero {
		t.Errorf("zero LocalConfig round-trip = %+v; want %+v", got, zero)
	}
}

// ---------------------------------------------------------------------------
// M — Concurrency  (category M1)
// ---------------------------------------------------------------------------

// TestSetHost_Concurrency verifies that 50 goroutines calling SetHost
// simultaneously do not produce a data race on the underlying atomic.Value.
// Always run with: go test -race ./...
func TestSetHost_Concurrency(t *testing.T) {
	cm := newTestManager()

	candidates := []string{
		"host-alpha", "host-beta", "host-gamma", "host-delta",
		"10.0.0.1", "192.168.1.1", "localhost", "example.com",
	}
	validSet := make(map[string]bool, len(candidates))
	for _, h := range candidates {
		validSet[h] = true
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		host := candidates[i%len(candidates)]
		go func(h string) {
			defer wg.Done()
			cm.SetHost(h)
		}(host)
	}

	wg.Wait()

	// After all writes, GetHost must return one of the written values.
	if final := cm.GetHost(); !validSet[final] {
		t.Errorf("GetHost() = %q after concurrent writes; not in expected value set", final)
	}
}

// ---------------------------------------------------------------------------
// N — Snapshot immutability  (category N1)
// ---------------------------------------------------------------------------

// TestGet_SnapshotImmutability verifies that a *RuntimeConfig obtained from
// Get() before a mutation still reflects the pre-mutation state, confirming
// that update() stores a new pointer rather than modifying in place.
func TestGet_SnapshotImmutability(t *testing.T) {
	cm := newTestManager()
	cm.SetHost("before")

	snapshot := cm.Get() // obtain pointer to current state

	cm.SetHost("after") // triggers atomic swap; snapshot must remain untouched

	if snapshot.Local.Host != "before" {
		t.Errorf("snapshot.Local.Host = %q after subsequent SetHost; want %q (snapshot mutated)", snapshot.Local.Host, "before")
	}
	if cm.GetHost() != "after" {
		t.Errorf("GetHost() = %q; want %q (live state not updated)", cm.GetHost(), "after")
	}
}
