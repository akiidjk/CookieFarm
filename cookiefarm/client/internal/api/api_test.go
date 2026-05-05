package api

/*
Step 1 – Identify Parameters and Environment Conditions
-------------------------------------------------------
Unit under test:
- Helpers: doJSON, checkStatus, getCookie
- HTTP transport: (*Client).doRequest + wrappers get/postJSON/postForm
- High-level API calls: Login, GetConfig, SubmitBatchDirect, SubmitFlag

Inputs and environment:
- HTTP status code + response body
- Response cookies (token present/missing)
- JSON payload validity (valid/malformed)
- Auth mode (AUTHED / NOTAUTHED)
- Config state from client/config singleton: host, port, token
- Server reachability (up/down)
- Request body categories (empty/non-empty)
- Endpoints and content-types

Returns:
- error (nil / non-nil)
- parsed Shared config
- side-effect: token saved into config singleton after Login

Step 2 – Define Categories
--------------------------
A) Status categories:
   - A1: 200 OK
   - A2: non-200 (400/401/403/500)

B) JSON categories:
   - B1: valid Shared JSON
   - B2: malformed JSON
   - B3: valid JSON but for different shape (decode issue for target)

C) Cookie categories:
   - C1: token cookie present
   - C2: token cookie missing

D) Auth/header categories:
   - D1: AUTHED with host configured
   - D2: AUTHED with host empty (function currently checks host and returns "missing auth token")
   - D3: NOTAUTHED

E) Network categories:
   - E1: server reachable
   - E2: server unreachable

F) Payload categories:
   - F1: empty flags slice
   - F2: non-empty flags slice
   - F3: single flag payload

Step 3 – Define Constraints
---------------------------
- Login requires endpoint /api/v1/auth/login and a reachable server (E1) to test status/cookie branches.
- GetConfig/Submit* are AUTHED by implementation path; for deterministic success tests use D1.
- D2 triggers early error in doRequest and bypasses network.
- C2 branch is meaningful only when status is A1 and endpoint returns 200.
- For malformed JSON branch (B2), status must be A1 so parsing is reached.
- "Unreachable server" (E2) short-circuits before status/body checks.

Step 4 – Generate Test Frames
-----------------------------
1) checkStatus
   - A1 => nil
   - A2 => error contains status code

2) doJSON
   - B1 => decoded struct correct
   - B2 => json decode error

3) getCookie
   - C1 => returns token
   - C2 => returns not found error

4) doRequest
   - D2 => returns missing auth token without HTTP call
   - D1 + E1 + body/content-type => returns response body and headers observed by server
   - D3 + E1 => request sent without auth cookie
   - E2 => network error

5) Login
   - A1 + C1 => nil and token stored
   - A2 => error
   - A1 + C2 => error

6) GetConfig
   - A1 + B1 => parsed Shared
   - A1 + B2 => decode error
   - A2 => status error

7) SubmitBatchDirect
   - F1 + A1 => nil
   - F2 + A2 => error

8) SubmitFlag
   - F3 + A1 => nil
   - F3 + A2 => error

Step 5 – Implement Tests
------------------------
Convention: should_<expected>_when_<condition>
All tests are independent; singleton client/config state reset per test.
*/

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/database"
	"strconv"
	"strings"
	"sync"
	"testing"

	clientconfig "client/config"

	sharedconfig "sharedconfig"
)

func resetClientSingletonForTest(host string, port uint16, token string) {
	once = sync.Once{}
	instance = nil

	cm := clientconfig.GetInstance()
	cm.SetHost(host)
	cm.SetPort(port)
	cm.SetToken(token)
}

func parseHostPort(t *testing.T, rawURL string) (string, uint16) {
	t.Helper()

	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}

	host, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatalf("split host/port: %v", err)
	}

	p, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("atoi port: %v", err)
	}

	return host, uint16(p)
}

func makeFlag(i int) database.Flag {
	return database.Flag{
		FlagCode:    fmt.Sprintf("FLAG_%d", i),
		ServiceName: "svc",
		PortService: 1337,
		SubmitTime:  1,
		Status:      0,
		TeamID:      1,
		Username:    "tester",
		ExploitName: "exp.py",
		Msg:         "ok",
	}
}

func TestCheckStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    int
		body    []byte
		wantErr bool
	}{
		// Category A1
		{name: "should_return_nil_when_status_is_200", code: 200, body: []byte("ok"), wantErr: false},
		// Category A2
		{name: "should_return_error_when_status_is_non_200", code: 403, body: []byte("forbidden"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := checkStatus(tt.code, tt.body)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if tt.wantErr && !strings.Contains(err.Error(), strconv.Itoa(tt.code)) {
				t.Fatalf("expected error to contain status code %d, got %q", tt.code, err.Error())
			}
		})
	}
}

func TestDoJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    []byte
		wantErr bool
	}{
		// Category B1
		{
			name:    "should_decode_shared_when_json_is_valid",
			body:    []byte(`{"services":{"http":80},"regex_flag":"FLAG{.*}","configured":true}`),
			wantErr: false,
		},
		// Category B2
		{
			name:    "should_return_error_when_json_is_malformed",
			body:    []byte(`{"services":`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var out sharedconfig.Shared
			err := doJSON(tt.body, &out)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if !tt.wantErr {
				if out.Services["http"] != 80 {
					t.Fatalf("expected services[http]=80, got %d", out.Services["http"])
				}
				if !out.Configured {
					t.Fatalf("expected configured=true")
				}
			}
		})
	}
}

func TestGetCookie(t *testing.T) {
	t.Parallel()

	t.Run("should_return_cookie_value_when_cookie_exists", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			Header: http.Header{
				"Set-Cookie": []string{"token=abc123; Path=/; HttpOnly"},
			},
		}

		got, err := getCookie(resp, "token")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "abc123" {
			t.Fatalf("expected token abc123, got %q", got)
		}
	})

	t.Run("should_return_error_when_cookie_does_not_exist", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			Header: http.Header{
				"Set-Cookie": []string{"session=session-id; Path=/"},
			},
		}

		_, err := getCookie(resp, "token")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "token not found") {
			t.Fatalf("expected token not found error, got %v", err)
		}
	})
}

func TestDoRequest_AuthedEmptyHost(t *testing.T) {
	t.Run("should_return_error_when_authed_and_host_is_empty", func(t *testing.T) {
		resetClientSingletonForTest("", 0, "tok")
		c := &Client{
			baseURL: "http://127.0.0.1:1",
			http:    &http.Client{},
		}

		resp, _, err := c.doRequest(http.MethodGet, "/x", nil, AUTHED, "")
		if resp != nil {
			_ = resp.Body.Close()
		}
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "missing auth token") {
			t.Fatalf("expected missing auth token, got %v", err)
		}
	})
}

func TestDoRequest_AuthedSendCookie(t *testing.T) {
	t.Run("should_send_auth_cookie_and_content_type_when_authed", func(t *testing.T) {
		var gotCookie string
		var gotCT string
		var gotMethod string
		var gotPath string
		var gotBody string

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotCookie = r.Header.Get("Cookie")
			gotCT = r.Header.Get("Content-Type")
			gotMethod = r.Method
			gotPath = r.URL.Path
			b, _ := io.ReadAll(r.Body)
			gotBody = string(b)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`ok`))
		}))
		defer srv.Close()

		host, port := parseHostPort(t, srv.URL)
		resetClientSingletonForTest(host, port, "mytoken")

		c := getClient()
		resp, body, err := c.doRequest(http.MethodPost, "/abc", []byte(`{"a":1}`), AUTHED, "application/json")
		if resp != nil {
			_ = resp.Body.Close()
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(body) != "ok" {
			t.Fatalf("expected body ok, got %q", string(body))
		}
		if gotCookie != "token=mytoken" {
			t.Fatalf("expected auth cookie token=mytoken, got %q", gotCookie)
		}
		if gotCT != "application/json" {
			t.Fatalf("expected content-type application/json, got %q", gotCT)
		}
		if gotMethod != http.MethodPost || gotPath != "/abc" {
			t.Fatalf("unexpected request method/path: %s %s", gotMethod, gotPath)
		}
		if gotBody != `{"a":1}` {
			t.Fatalf("unexpected request body: %q", gotBody)
		}
	})
}

func TestDoRequest_NotAuthedNoCookie(t *testing.T) {
	t.Run("should_not_send_auth_cookie_when_not_authed", func(t *testing.T) {
		var gotCookie string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotCookie = r.Header.Get("Cookie")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`ok`))
		}))
		defer srv.Close()

		host, port := parseHostPort(t, srv.URL)
		resetClientSingletonForTest(host, port, "present-but-not-used")

		c := getClient()
		resp, _, err := c.doRequest(http.MethodGet, "/noauth", nil, NOTAUTHED, "")
		if resp != nil {
			_ = resp.Body.Close()
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotCookie != "" {
			t.Fatalf("expected no cookie in NOTAUTHED request, got %q", gotCookie)
		}
	})
}

func TestDoRequest_NetworkError(t *testing.T) {
	t.Run("should_return_network_error_when_server_unreachable", func(t *testing.T) {
		// Reserve and close a port to maximize unreachable determinism.
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen failed: %v", err)
		}
		addr := ln.Addr().String()
		_ = ln.Close()

		host, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			t.Fatalf("split host/port: %v", err)
		}
		p, _ := strconv.Atoi(portStr)
		resetClientSingletonForTest(host, uint16(p), "tok")

		c := getClient()
		resp, _, err := c.doRequest(http.MethodGet, "/x", nil, NOTAUTHED, "")
		if resp != nil {
			_ = resp.Body.Close()
		}
		if err == nil {
			t.Fatalf("expected network error, got nil")
		}
		if !strings.Contains(err.Error(), "do request") {
			t.Fatalf("expected wrapped do request error, got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("should_store_token_when_login_successful_and_cookie_present", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/auth/login" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			_ = r.ParseForm()
			if r.Form.Get("username") != "user1" || r.Form.Get("password") != "pass1" {
				t.Fatalf("unexpected form values: username=%q password=%q", r.Form.Get("username"), r.Form.Get("password"))
			}

			http.SetCookie(w, &http.Cookie{Name: "token", Value: "jwt-token", Path: "/"})
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`ok`))
		}))
		defer srv.Close()

		host, port := parseHostPort(t, srv.URL)
		resetClientSingletonForTest(host, port, "")

		if err := Login("user1", "pass1"); err != nil {
			t.Fatalf("unexpected login error: %v", err)
		}

		if got := clientconfig.GetInstance().GetToken(); got != "jwt-token" {
			t.Fatalf("expected token saved in config, got %q", got)
		}
	})

	loginCases := []struct {
		name       string
		status     int
		body       string
		errMessage string
	}{
		{"should_return_error_when_login_status_is_non_200", http.StatusUnauthorized, `unauthorized`, "status 401"},
		{"should_return_error_when_cookie_missing_even_if_status_200", http.StatusOK, `ok`, "token not found"},
	}

	for _, tc := range loginCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			host, port := parseHostPort(t, srv.URL)
			resetClientSingletonForTest(host, port, "")

			err := Login("u", "p")
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.errMessage) {
				t.Fatalf("expected %q in error, got %v", tc.errMessage, err)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	t.Run("should_return_shared_config_when_status_200_and_json_valid", func(t *testing.T) {
		const payload = `{
			"services":{"http":80,"https":443},
			"regex_flag":"FLAG\\{.*\\}",
			"format_ip_teams":"10.10.{}.1",
			"my_team_id":1,
			"url_flag_ids":"http://flags.local",
			"nop_team":99,
			"range_ip_teams":30,
			"configured":true
		}`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/v1/config" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.Header.Get("Cookie") != "token=my-token" {
				t.Fatalf("expected auth cookie token=my-token, got %q", r.Header.Get("Cookie"))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
		}))
		defer srv.Close()

		host, port := parseHostPort(t, srv.URL)
		resetClientSingletonForTest(host, port, "my-token")

		cfg, err := GetConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Services["http"] != 80 || cfg.Services["https"] != 443 {
			t.Fatalf("unexpected services map: %#v", cfg.Services)
		}
		if cfg.RangeIPTeams != 30 || !cfg.Configured {
			t.Fatalf("unexpected parsed config: %+v", cfg)
		}
	})

	configCases := []struct {
		name       string
		status     int
		body       string
		errMessage string
	}{
		{"should_return_error_when_status_non_200", http.StatusForbidden, `forbidden`, "status 403"},
		{"should_return_error_when_json_is_invalid", http.StatusOK, `{"services":`, "json decode error"},
	}

	for _, tc := range configCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			host, port := parseHostPort(t, srv.URL)
			resetClientSingletonForTest(host, port, "tok")

			_, err := GetConfig()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.errMessage) {
				t.Fatalf("expected %q in error, got %v", tc.errMessage, err)
			}
		})
	}
}

func TestSubmitBatchDirect(t *testing.T) {
	t.Run("should_return_nil_when_empty_batch_and_status_200", func(t *testing.T) {
		var gotContentType string
		var gotCookie string
		var gotBody string

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotContentType = r.Header.Get("Content-Type")
			gotCookie = r.Header.Get("Cookie")
			bodyBytes, _ := io.ReadAll(r.Body)
			gotBody = string(bodyBytes)

			if r.URL.Path != "/api/v1/submit-flags-standalone" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`ok`))
		}))
		defer srv.Close()

		host, port := parseHostPort(t, srv.URL)
		resetClientSingletonForTest(host, port, "tok-empty")

		err := SubmitBatchDirect([]database.Flag{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if gotContentType != "application/json" {
			t.Fatalf("expected JSON content type, got %q", gotContentType)
		}
		if gotCookie != "token=tok-empty" {
			t.Fatalf("expected auth cookie, got %q", gotCookie)
		}
		if !strings.Contains(gotBody, `"flags":[]`) {
			t.Fatalf("expected empty flags payload, got %s", gotBody)
		}
	})

	t.Run("should_return_error_when_server_returns_500_for_non_empty_batch", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`boom`))
		}))
		defer srv.Close()

		host, port := parseHostPort(t, srv.URL)
		resetClientSingletonForTest(host, port, "tok")

		err := SubmitBatchDirect([]database.Flag{makeFlag(1), makeFlag(2)})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "status 500") {
			t.Fatalf("expected status 500 error, got %v", err)
		}
	})
}

func TestSubmitFlag(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		// Category F3 + A1
		{name: "should_return_nil_when_single_flag_submit_status_200", statusCode: 200, wantErr: false},
		// Category F3 + A2
		{name: "should_return_error_when_single_flag_submit_status_400", statusCode: 400, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var seenPath string
			var seenBody string

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seenPath = r.URL.Path
				bodyBytes, _ := io.ReadAll(r.Body)
				seenBody = string(bodyBytes)

				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(`resp`))
			}))
			defer srv.Close()

			host, port := parseHostPort(t, srv.URL)
			resetClientSingletonForTest(host, port, "tok")

			err := SubmitFlag(makeFlag(10))
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if seenPath != "/api/v1/submit-flag" {
				t.Fatalf("unexpected path: %s", seenPath)
			}
			if !strings.Contains(seenBody, `"flag"`) || !strings.Contains(seenBody, `"flag_code":"FLAG_10"`) {
				t.Fatalf("unexpected payload: %s", seenBody)
			}
		})
	}
}

func TestGetConfig_ShouldReturnErrorWhenServerUnreachable(t *testing.T) {
	t.Parallel()

	// Extra E2 coverage at high-level API.
	resetClientSingletonForTest("127.0.0.1", 1, "tok")
	_, err := GetConfig()
	if err == nil {
		t.Fatalf("expected error for unreachable server, got nil")
	}
	if !strings.Contains(err.Error(), "do request") && !errors.Is(err, err) {
		t.Fatalf("expected wrapped request error, got %v", err)
	}
}
