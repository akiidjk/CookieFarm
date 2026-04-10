/*
================================================================================
Category-Partition Methodology — CircuitBreaker (circuit_test.go)
Module:  client/websockets
Package: websockets  (white-box — same package as production code)
================================================================================

STEP 1 — Identify Independently Testable Units
────────────────────────────────────────────────────────────────────────────────
The following methods of *CircuitBreaker are the units under test.
Each method is small, mutex-protected, and has observable side-effects on
the struct's unexported fields (accessible here because we are in the same
package).

  U1. (*CircuitBreaker).RecordSuccess()
      Unconditionally resets failureCount to 0 and sets state to StateClosed.

  U2. (*CircuitBreaker).RecordFailure()
      Increments failureCount, records lastFailureTime, and conditionally
      opens the circuit:
        • If state == StateHalfOpen → state = StateOpen immediately.
        • If failureCount (after increment) >= failureThreshold (= 2)
          → state = StateOpen.
      Otherwise state is unchanged.

  U3. (*CircuitBreaker).IsAllowed() bool
      Returns whether a connection attempt should be permitted:
        • StateClosed   → always true.
        • StateHalfOpen → always true.
        • StateOpen     → false, UNLESS time.Since(lastFailureTime) > 30 s,
                          in which case state transitions to StateHalfOpen
                          and the method returns true.

  U4. Failure/Success interplay — end-to-end lifecycle:
      Two failures open the circuit; RecordSuccess resets it; IsAllowed
      must subsequently return true.

  U5. Concurrent access — all three methods invoked from multiple goroutines
      simultaneously must not produce a data race (validates mutex coverage).

────────────────────────────────────────────────────────────────────────────────

STEP 2 — Identify Categories for Each Unit
────────────────────────────────────────────────────────────────────────────────

  U1 — RecordSuccess
    C1: Circuit state immediately before the call
    C2: Value of failureCount immediately before the call

  U2 — RecordFailure
    C3: Circuit state immediately before the call
    C4: Relationship of failureCount (after the increment) to
        failureThreshold (= 2)

  U3 — IsAllowed
    C5: Circuit state at the moment of the call
    C6: Time elapsed since lastFailureTime (relevant only when StateOpen)

  U4 — Interplay
    C7: Specific ordered call sequence

  U5 — Concurrent access
    C8: Goroutine access pattern (variety of concurrent operations)

────────────────────────────────────────────────────────────────────────────────

STEP 3 — Partition Each Category into Choices
────────────────────────────────────────────────────────────────────────────────

  C1 (pre-call state for RecordSuccess):
    C1a — StateClosed   : normal operating state, no prior failures
    C1b — StateOpen     : circuit tripped after >= failureThreshold failures
    C1c — StateHalfOpen : timeout elapsed, probing whether server recovered

  C2 (failureCount before RecordSuccess):
    C2a — 0        : already at baseline (e.g., freshly created)
    C2b — positive : accumulated failures that must be cleared

  C3 (pre-call state for RecordFailure):
    C3a — StateClosed   : normal path; threshold check governs transition
    C3b — StateHalfOpen : any single failure must re-open immediately
    C3c — StateOpen     : already tripped; count continues to grow

  C4 (failureCount + 1 vs. failureThreshold after the increment):
    C4a — below threshold  : count → 1, still < 2  → state unchanged
    C4b — reaches threshold: count → 2, = 2         → state → Open
    C4c — above threshold  : count → 3+, already Open → state stays Open

  C5 (state for IsAllowed):
    C5a — StateClosed
    C5b — StateOpen with a recent failure (within 30 s)
    C5c — StateOpen with a stale  failure (more than 30 s ago)
    C5d — StateHalfOpen

  C6 (elapsed time since lastFailureTime — meaningful only for C5b / C5c):
    C6a — < resetTimeout (30 s) : paired with C5b → must return false
    C6b — > resetTimeout (30 s) : paired with C5c → must return true and
                                   transition state to HalfOpen

  C7 (interplay call sequence):
    C7a — RecordFailure × 2 → assert Open → RecordSuccess → assert Closed
          → assert IsAllowed = true

  C8 (concurrent access pattern):
    C8a — 100 goroutines all calling RecordFailure concurrently
    C8b — 50 goroutines calling RecordFailure interleaved with 50 goroutines
          calling RecordSuccess concurrently

────────────────────────────────────────────────────────────────────────────────

STEP 4 — Constraints, Eliminations, and Test Frames
────────────────────────────────────────────────────────────────────────────────

Constraints / eliminations:
  • C2 is dominated by C1 for RecordSuccess: the reset is unconditional, so
    crossing all C1 × C2 combinations (6 total) adds no new behaviour. One
    representative per C1 value is sufficient; C2 is sampled, not exhausted.
  • C3b (HalfOpen) overrides the threshold check entirely, so crossing C3b
    with each C4 value is redundant: one example with count deliberately
    below threshold (C4a) fully demonstrates the special-case rule.
  • C3c (Open): the threshold condition is already met; only C4c is reachable.
  • C5b ↔ C6a and C5c ↔ C6b are definitionally paired and cannot be split.
  • C8 tests must NOT assert on exact final counters: concurrent scheduling
    makes exact values non-deterministic. Their sole purpose is to expose
    data races under the race detector (go test -race).

Derived test frames:
  TF-01  RecordSuccess | Closed   (C1a, C2a)  → state=Closed, count=0
  TF-02  RecordSuccess | Open     (C1b, C2b)  → state=Closed, count=0
  TF-03  RecordSuccess | HalfOpen (C1c, C2b)  → state=Closed, count=0
  TF-04  RecordFailure ×1 from Closed  (C3a, C4a) → state=Closed, count=1
  TF-05  RecordFailure ×2 from Closed  (C3a, C4b) → state=Open,   count=2
  TF-06  RecordFailure from Open       (C3c, C4c) → state=Open,   count>threshold
  TF-07  RecordFailure from HalfOpen   (C3b)      → state=Open immediately
  TF-08  IsAllowed | StateClosed                  (C5a)      → true
  TF-09  IsAllowed | StateOpen + recent failure   (C5b, C6a) → false, stays Open
  TF-10  IsAllowed | StateOpen + stale  failure   (C5c, C6b) → true, → HalfOpen
  TF-11  IsAllowed | StateHalfOpen                (C5d)      → true
  TF-12  2×Failure → Open → RecordSuccess → Closed (C7a)     → IsAllowed=true
  TF-13  100 goroutines RecordFailure only  (C8a)            → no data race
  TF-14  50×Failure + 50×Success concurrent (C8b)            → no data race

================================================================================
*/

package websockets

import (
	"logger"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain configures the shared logger at "error" level before any test runs.
// This suppresses the debug / info / warn log lines that the production code
// emits during failure-path tests, keeping test output clean while still
// surfacing genuine error-level entries that may indicate real problems.
func TestMain(m *testing.M) {
	logger.Setup("error", false)
	os.Exit(m.Run())
}

// newCB returns a fresh, isolated *CircuitBreaker in StateClosed with zero
// failures. Every test that exercises CircuitBreaker MUST use this helper and
// MUST NOT reference the package-level circuitBreaker variable, which is
// shared mutable state across the entire test binary.
func newCB() *CircuitBreaker {
	return &CircuitBreaker{
		state:        StateClosed,
		failureCount: 0,
	}
}

// ---------------------------------------------------------------------------
// U1 — RecordSuccess                                         TF-01, 02, 03
// ---------------------------------------------------------------------------

// TestCircuitBreaker_RecordSuccess_TableDriven verifies TF-01, TF-02, and
// TF-03: RecordSuccess must unconditionally reset the circuit to StateClosed
// with failureCount = 0, regardless of the state it was in before the call.
func TestCircuitBreaker_RecordSuccess_TableDriven(t *testing.T) {
	cases := []struct {
		name         string // human-readable frame label
		initialState CircuitState
		initialCount int // C2: failureCount before the call
	}{
		// TF-01: C1a — success from the normal operating state (C2a: count=0)
		{name: "from_Closed_countZero", initialState: StateClosed, initialCount: 0},
		// TF-02: C1b — success resets an open circuit (C2b: count=5)
		{name: "from_Open_countPositive", initialState: StateOpen, initialCount: 5},
		// TF-03: C1c — success resets a half-open probe (C2b: count=3)
		{name: "from_HalfOpen_countPositive", initialState: StateHalfOpen, initialCount: 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cb := newCB()
			cb.state = tc.initialState
			cb.failureCount = tc.initialCount

			cb.RecordSuccess()

			assert.Equal(t, StateClosed, cb.state,
				"state must be StateClosed after RecordSuccess (was %v)", tc.initialState)
			assert.Equal(t, 0, cb.failureCount,
				"failureCount must be reset to 0 after RecordSuccess (was %d)", tc.initialCount)
		})
	}
}

// ---------------------------------------------------------------------------
// U2 — RecordFailure                                         TF-04..TF-07
// ---------------------------------------------------------------------------

// TestCircuitBreaker_RecordFailure_BelowThreshold_StaysClosed covers TF-04
// (C3a + C4a): a single failure from StateClosed increments failureCount to 1,
// which is below failureThreshold (= 2), so the circuit must remain Closed.
func TestCircuitBreaker_RecordFailure_BelowThreshold_StaysClosed(t *testing.T) {
	cb := newCB() // Closed, count = 0

	cb.RecordFailure()

	assert.Equal(t, StateClosed, cb.state,
		"circuit must remain Closed when failureCount has not yet reached failureThreshold")
	assert.Equal(t, 1, cb.failureCount,
		"failureCount must be 1 after the first failure")
	assert.False(t, cb.lastFailureTime.IsZero(),
		"lastFailureTime must be recorded on every failure")
}

// TestCircuitBreaker_RecordFailure_ReachesThreshold_OpensCircuit covers TF-05
// (C3a + C4b): a second consecutive failure from Closed brings failureCount to
// exactly failureThreshold (= 2) and must trip the circuit to StateOpen.
func TestCircuitBreaker_RecordFailure_ReachesThreshold_OpensCircuit(t *testing.T) {
	cb := newCB()

	cb.RecordFailure() // count → 1 (below threshold, Closed)
	cb.RecordFailure() // count → 2 = failureThreshold → Open

	assert.Equal(t, StateOpen, cb.state,
		"circuit must transition to StateOpen once failureCount reaches failureThreshold")
	assert.Equal(t, failureThreshold, cb.failureCount,
		"failureCount must equal failureThreshold after the threshold-reaching failure")
}

// TestCircuitBreaker_RecordFailure_AlreadyOpen_StaysOpen covers TF-06
// (C3c + C4c): a further failure when the circuit is already Open must keep
// it Open and continue incrementing failureCount beyond the threshold.
func TestCircuitBreaker_RecordFailure_AlreadyOpen_StaysOpen(t *testing.T) {
	cb := newCB()
	cb.state = StateOpen
	cb.failureCount = failureThreshold // already at the threshold

	cb.RecordFailure() // 3rd failure — count → failureThreshold + 1

	assert.Equal(t, StateOpen, cb.state,
		"additional failures while Open must not change state")
	assert.Greater(t, cb.failureCount, failureThreshold,
		"failureCount must continue to grow past the threshold")
}

// TestCircuitBreaker_RecordFailure_FromHalfOpen_OpensImmediately covers TF-07
// (C3b): a failure from StateHalfOpen must transition to StateOpen immediately,
// even when failureCount is below failureThreshold. This validates the special
// HalfOpen short-circuit in RecordFailure.
func TestCircuitBreaker_RecordFailure_FromHalfOpen_OpensImmediately(t *testing.T) {
	cb := newCB()
	cb.state = StateHalfOpen
	cb.failureCount = 0 // deliberately below failureThreshold

	cb.RecordFailure()

	assert.Equal(t, StateOpen, cb.state,
		"any failure from StateHalfOpen must immediately re-open the circuit, "+
			"regardless of failureCount (got count=%d, threshold=%d)",
		cb.failureCount, failureThreshold)
}

// ---------------------------------------------------------------------------
// U3 — IsAllowed                                            TF-08..TF-11
// ---------------------------------------------------------------------------

// TestCircuitBreaker_IsAllowed_TableDriven covers TF-08, TF-09, TF-10, and
// TF-11 in a single table. For each circuit state variant it checks both the
// boolean return value and any resulting state transition inside IsAllowed.
func TestCircuitBreaker_IsAllowed_TableDriven(t *testing.T) {
	cases := []struct {
		name        string
		setup       func() *CircuitBreaker
		wantAllowed bool
		wantState   CircuitState // state the CB must be in AFTER IsAllowed returns
		desc        string       // assertion message suffix
	}{
		{
			// TF-08: C5a — Closed always permits connections.
			name: "Closed_AlwaysPermits",
			setup: func() *CircuitBreaker {
				return newCB()
			},
			wantAllowed: true,
			wantState:   StateClosed,
			desc:        "StateClosed must always return true and remain Closed",
		},
		{
			// TF-09: C5b + C6a — Open with a recent failure blocks.
			name: "Open_RecentFailure_Blocked",
			setup: func() *CircuitBreaker {
				cb := newCB()
				cb.state = StateOpen
				cb.failureCount = failureThreshold
				cb.lastFailureTime = time.Now() // elapsed ≈ 0, well within resetTimeout
				return cb
			},
			wantAllowed: false,
			wantState:   StateOpen,
			desc:        "StateOpen with recent failure must return false and stay Open",
		},
		{
			// TF-10: C5c + C6b — Open with a stale failure transitions to HalfOpen.
			name: "Open_StaleFailure_TransitionsToHalfOpen",
			setup: func() *CircuitBreaker {
				cb := newCB()
				cb.state = StateOpen
				cb.failureCount = failureThreshold
				// 31 s exceeds the 30 s resetTimeout constant.
				cb.lastFailureTime = time.Now().Add(-31 * time.Second)
				return cb
			},
			wantAllowed: true,
			wantState:   StateHalfOpen,
			desc: "StateOpen with stale failure must return true and " +
				"transition state to StateHalfOpen",
		},
		{
			// TF-11: C5d — HalfOpen always permits (probe attempt).
			name: "HalfOpen_AlwaysPermits",
			setup: func() *CircuitBreaker {
				cb := newCB()
				cb.state = StateHalfOpen
				return cb
			},
			wantAllowed: true,
			wantState:   StateHalfOpen,
			desc:        "StateHalfOpen must always return true and remain HalfOpen",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cb := tc.setup()

			got := cb.IsAllowed()

			assert.Equal(t, tc.wantAllowed, got,
				"IsAllowed return value mismatch — %s", tc.desc)
			assert.Equal(t, tc.wantState, cb.state,
				"unexpected post-call state — %s", tc.desc)
		})
	}
}

// ---------------------------------------------------------------------------
// U4 — Failure/Success interplay                                      TF-12
// ---------------------------------------------------------------------------

// TestCircuitBreaker_Interplay_FailThenSuccess_ResetsCircuit covers TF-12
// (C7a): validates the complete trip-and-reset lifecycle end-to-end.
//
// Sequence:
//  1. Two consecutive failures open the circuit (asserted as precondition).
//  2. An immediate IsAllowed call returns false (open circuit blocks).
//  3. RecordSuccess resets state to Closed with zero failure count.
//  4. A subsequent IsAllowed call returns true.
func TestCircuitBreaker_Interplay_FailThenSuccess_ResetsCircuit(t *testing.T) {
	cb := newCB()

	// --- Phase 1: trip the circuit ---
	cb.RecordFailure()
	cb.RecordFailure()

	require.Equal(t, StateOpen, cb.state,
		"precondition violated: two failures must open the circuit")
	require.False(t, cb.IsAllowed(),
		"precondition violated: a freshly-opened circuit must block connections")

	// --- Phase 2: simulate a successful recovery ---
	cb.RecordSuccess()

	assert.Equal(t, StateClosed, cb.state,
		"RecordSuccess must return circuit to StateClosed")
	assert.Equal(t, 0, cb.failureCount,
		"RecordSuccess must reset failureCount to 0")
	assert.True(t, cb.IsAllowed(),
		"IsAllowed must return true after the circuit has been reset by RecordSuccess")
}

// ---------------------------------------------------------------------------
// U5 — Concurrent access                                    TF-13, TF-14
// ---------------------------------------------------------------------------
// These tests exercise the mutex protection under high goroutine concurrency.
// Run with:  go test -race ./...
//
// No assertions are made on exact counter values or state after concurrent
// execution, because goroutine scheduling makes those non-deterministic.
// The only meaningful guarantee is the absence of a data race.

// TestCircuitBreaker_ConcurrentRecordFailure_NoDataRace covers TF-13 (C8a):
// 100 goroutines all call RecordFailure on the same CircuitBreaker at the
// same time. The test passes iff the race detector reports no conflict.
func TestCircuitBreaker_ConcurrentRecordFailure_NoDataRace(t *testing.T) {
	cb := newCB()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			cb.RecordFailure()
		}()
	}
	wg.Wait()

	// Verify the final state is one of the three legal values. Under this
	// write-only workload the circuit will invariably be Open, but we use
	// a map-lookup guard rather than hard-coding StateOpen so that the
	// assertion survives any future threshold change.
	validStates := map[CircuitState]bool{
		StateClosed:   true,
		StateHalfOpen: true,
		StateOpen:     true,
	}
	assert.True(t, validStates[cb.state],
		"state %v after concurrent RecordFailure is not a valid CircuitState", cb.state)
}

// TestCircuitBreaker_ConcurrentMixedAccess_NoDataRace covers TF-14 (C8b):
// 50 goroutines call RecordFailure while 50 other goroutines simultaneously
// call RecordSuccess on the same CircuitBreaker. The test validates that the
// mutex correctly serialises all mixed reads and writes.
func TestCircuitBreaker_ConcurrentMixedAccess_NoDataRace(t *testing.T) {
	cb := newCB()
	const half = 50

	var wg sync.WaitGroup
	wg.Add(half * 2)

	for range half {
		go func() {
			defer wg.Done()
			cb.RecordFailure()
		}()
		go func() {
			defer wg.Done()
			cb.RecordSuccess()
		}()
	}
	wg.Wait()

	validStates := map[CircuitState]bool{
		StateClosed:   true,
		StateHalfOpen: true,
		StateOpen:     true,
	}
	assert.True(t, validStates[cb.state],
		"state %v after concurrent mixed access is not a valid CircuitState", cb.state)
}
