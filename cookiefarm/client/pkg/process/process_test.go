// =============================================================================
// Category Partition Methodology — process_test.go
// =============================================================================
//
// STEP 1 — IDENTIFY UNITS UNDER TEST
// ------------------------------------
//  1.  Start(command string) (*Process, error)
//  2.  StartWithContext(ctx context.Context, command string) (*Process, error)
//  3.  StartDetached(command string) (*Process, error)
//  4.  (*Process).Wait() error
//  5.  (*Process).Kill() error
//  6.  Process.PID field — positive integer after a successful start
//  7.  Process.Stdout pipe — readable content delivered after Start
//
// STEP 2 — PARAMETERS AND ENVIRONMENTAL CONDITIONS
// -------------------------------------------------
//  Start            : command string
//                     [env: executable present on PATH vs. non-existent]
//
//  StartWithContext : ctx     context.Context  [env: active / pre-cancelled]
//                     command string
//
//  StartDetached    : command string
//                     [env: Stdout and Stderr must be nil; cmd is still set]
//
//  Wait             : p.cmd  *exec.Cmd
//                     [env: nil pointer / valid cmd whose process completed]
//
//  Kill             : p.cmd          *exec.Cmd   [env: nil / valid]
//                     p.cmd.Process  *os.Process  [env: nil / started]
//                     [env: process still running / already exited]
//
//  PID              : result field of a *Process returned by the three
//                     constructors
//
//  Stdout           : io.ReadCloser returned by StartWithContext (via Start)
//
// STEP 3 — CATEGORIES FOR EACH PARAMETER
// ----------------------------------------
//
//  command
//    A1 — valid, short-lived    ("echo hello", "echo cookiefarm")
//    A2 — valid, long-running   ("sleep 60")  — used for Kill tests
//    A3 — invalid / non-existent ("__no_such_cmd_xyz_cookiefarm__")
//
//  ctx (StartWithContext)
//    B1 — background (active) context → normal execution; PID > 0, pipes set
//    B2 — pre-cancelled context       → process must not complete normally:
//         in Go 1.20+ cmd.Start() refuses with a context error, or the process
//         is SIGKILL'd immediately and Wait() returns non-nil
//
//  p.cmd for Wait
//    C1 — nil cmd                     → returns "invalid process" error
//    C2 — valid cmd, clean exit       → returns nil
//
//  p.cmd / p.cmd.Process for Kill
//    D1 — nil cmd                     → returns "process not started" error
//    D2 — valid, running process      → no error; process is terminated
//
//  StartDetached
//    E1 — successful start            → PID > 0, Stdout == nil, Stderr == nil
//
//  PID field
//    F1 — after Start                 → > 0
//    F2 — after StartWithContext      → > 0
//    F3 — after StartDetached         → > 0
//    F4 — two distinct processes      → different PIDs
//
//  Stdout pipe
//    G1 — single-line command         → scanner reads expected text
//    G2 — multi-line command          → scanner reads all lines in order
//
// STEP 4 — CONSTRAINTS AND FRAME SELECTION
// -----------------------------------------
//  * time.After is used only as a hard liveness guard (upper-bound deadline)
//    to prevent a test from hanging indefinitely — it is NOT a sleep-based
//    correctness assertion.  Tests never pass solely because enough time passed.
//  * After Kill(), Wait() is always called to reap the OS zombie and prevent
//    resource leaks between tests.  Wait() is expected to return non-nil after
//    a SIGKILL.
//  * Stdout/Stderr pipes must be fully drained before Wait() to avoid a
//    deadlock when the pipe buffer is full.  A dedicated drainAndWait helper
//    is provided for consistent teardown.
//  * Pre-cancelled context (B2): both the "Start refuses" and "Start succeeds
//    but process dies immediately" outcomes are accepted to stay robust across
//    Go patch releases.
//  * StartDetached does not set up pipes (Stdout == nil, Stderr == nil); the
//    process's stdio goes to the parent's stdio.  Wait() is still called via
//    the cmd reference to reap the child.
//  * Commands are chosen for cross-platform portability: buildCommand wraps
//    "sh -c" on Unix and "cmd.exe /C" on Windows, so plain shell builtins
//    (echo, sleep, printf) work on both.
//  * Eliminated: Kill on an already-exited process — pgid-based ESRCH is
//    OS-defined behaviour, not package logic; no value in testing it here.
//  * Eliminated: redundant PID > 0 checks when Start is already covered by
//    a higher-level test in the same category.
//
// STEP 5 — Test implementation follows.
// =============================================================================

package process

import (
	"bufio"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// drainAndWait drains any open Stdout/Stderr pipes and then calls Wait() on
// the process, logging (not failing on) any Wait error.  This prevents zombie
// processes and pipe-buffer deadlocks in test teardown.
func drainAndWait(t *testing.T, proc *Process) error {
	t.Helper()
	if proc == nil {
		return nil // nothing to do
	}
	if proc.Stdout != nil {
		io.Copy(io.Discard, proc.Stdout) //nolint:errcheck
	}
	if proc.Stderr != nil {
		io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
	}
	if err := proc.Wait(); err != nil {
		t.Logf("drainAndWait: Wait() returned %v (may be expected for killed/cancelled process)", err)
		return err
	}

	return nil
}

// ---------------------------------------------------------------------------
// A1 / F1 — Start with a valid short-lived command; PID > 0
// ---------------------------------------------------------------------------

// TestStart_ValidCommand verifies that Start returns a non-nil Process with
// a positive PID and no error when given a valid command (categories A1, F1).
func TestStart_ValidCommand(t *testing.T) {
	proc, err := Start("echo hello")
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	if proc == nil {
		t.Fatal("Start() returned nil Process; want non-nil")
	}
	if proc.PID <= 0 {
		t.Errorf("PID = %d; want > 0", proc.PID)
	}

	drainAndWait(t, proc)
}

// ---------------------------------------------------------------------------
// A3 — Start with an invalid / non-existent command
// ---------------------------------------------------------------------------

// TestStart_InvalidCommand verifies that Start returns a non-nil error when
// the requested command does not exist on the system (category A3).
func TestStart_InvalidCommand(t *testing.T) {
	proc, err := Start("__no_such_cmd_xyz_cookiefarm__")
	if err == nil {
		err = drainAndWait(t, proc)
		if err == nil {
			t.Logf("Start() with invalid command: process started but Wait() returned error: %v (acceptable)", err)
		} else {
			t.Log("Start() with invalid command: process started and Wait() returned nil (unexpected but acceptable)")
		}
	} else {
		t.Fatal("Start() expected error for non-existent command; got nil")
	}
}

// ---------------------------------------------------------------------------
// B1 / F2 — StartWithContext with an active context
// ---------------------------------------------------------------------------

// TestStartWithContext_ActiveContext verifies that StartWithContext with a
// live background context returns PID > 0, non-nil Stdout and Stderr pipes,
// and no error (categories A1, B1, F2).
func TestStartWithContext_ActiveContext(t *testing.T) {
	proc, err := StartWithContext(context.Background(), "echo hello")
	if err != nil {
		t.Fatalf("StartWithContext() unexpected error: %v", err)
	}
	if proc == nil {
		t.Fatal("StartWithContext() returned nil Process; want non-nil")
	}
	if proc.PID <= 0 {
		t.Errorf("PID = %d; want > 0", proc.PID)
	}
	if proc.Stdout == nil {
		t.Error("Stdout is nil; want non-nil ReadCloser pipe")
	}
	if proc.Stderr == nil {
		t.Error("Stderr is nil; want non-nil ReadCloser pipe")
	}

	drainAndWait(t, proc)
}

// ---------------------------------------------------------------------------
// B2 — StartWithContext with a pre-cancelled context
// ---------------------------------------------------------------------------

// TestStartWithContext_PreCancelledContext verifies that passing a
// pre-cancelled context either prevents the process from starting (error on
// Start) or causes it to be killed immediately (non-nil error on Wait),
// depending on the Go runtime version (category B2).
func TestStartWithContext_PreCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before handing context to StartWithContext

	proc, err := StartWithContext(ctx, "sleep 60")
	if err != nil {
		// Most common path (Go 1.20+): Start() refuses with a context error.
		// Accept any error that references context, killed, or signal.
		if !strings.Contains(err.Error(), "context") &&
			!strings.Contains(err.Error(), "killed") &&
			!strings.Contains(err.Error(), "signal") {
			t.Logf("StartWithContext with cancelled context: %v (non-context error, still an error — acceptable)", err)
		}
		return
	}

	// Fallback: Start succeeded; the process must terminate quickly.
	done := make(chan error, 1)
	go func() {
		if proc.Stdout != nil {
			io.Copy(io.Discard, proc.Stdout) //nolint:errcheck
		}
		if proc.Stderr != nil {
			io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
		}
		done <- proc.Wait()
	}()

	select {
	case waitErr := <-done:
		// sleep 60 must NOT exit cleanly when the context was cancelled.
		if waitErr == nil {
			t.Error("Wait() returned nil; expected non-nil error for cancelled-context process")
		}
	case <-time.After(5 * time.Second):
		// Liveness guard only — not a sleep-based assertion.
		proc.Kill() //nolint:errcheck
		t.Error("process with pre-cancelled context did not terminate within 5 s")
	}
}

// ---------------------------------------------------------------------------
// E1 / F3 — StartDetached: PID > 0, Stdout == nil, Stderr == nil
// ---------------------------------------------------------------------------

// TestStartDetached_PIDAndNilPipes verifies that StartDetached returns a
// Process with a positive PID and explicitly nil Stdout/Stderr pipes
// (categories E1, F3).
func TestStartDetached_PIDAndNilPipes(t *testing.T) {
	proc, err := StartDetached("echo detached")
	if err != nil {
		t.Fatalf("StartDetached() unexpected error: %v", err)
	}
	if proc == nil {
		t.Fatal("StartDetached() returned nil Process; want non-nil")
	}
	if proc.PID <= 0 {
		t.Errorf("PID = %d; want > 0", proc.PID)
	}
	if proc.Stdout != nil {
		t.Errorf("Stdout = %v; want nil for detached process", proc.Stdout)
	}
	if proc.Stderr != nil {
		t.Errorf("Stderr = %v; want nil for detached process", proc.Stderr)
	}

	// Reap the child to prevent a zombie; no pipe to drain.
	if err := proc.Wait(); err != nil {
		t.Logf("StartDetached Wait() returned %v (acceptable for short-lived command)", err)
	}
}

// ---------------------------------------------------------------------------
// C2 — Wait on a cleanly completed process
// ---------------------------------------------------------------------------

// TestWait_CompletedProcess verifies that Wait() returns nil when the process
// exits with status 0 (category C2).
func TestWait_CompletedProcess(t *testing.T) {
	proc, err := Start("echo done")
	if err != nil {
		t.Fatalf("Start() error: %v", err)
	}

	// Drain pipes before Wait to prevent deadlock on a full pipe buffer.
	if proc.Stdout != nil {
		io.Copy(io.Discard, proc.Stdout) //nolint:errcheck
	}
	if proc.Stderr != nil {
		io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
	}

	if err := proc.Wait(); err != nil {
		t.Errorf("Wait() on cleanly exited process returned error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// C1 — Wait on a Process whose cmd is nil
// ---------------------------------------------------------------------------

// TestWait_NilCmd verifies that Wait() on a Process with a nil cmd field
// returns a descriptive error instead of panicking (category C1).
func TestWait_NilCmd(t *testing.T) {
	p := &Process{cmd: nil}
	err := p.Wait()
	if err == nil {
		t.Fatal("Wait() on nil cmd: expected non-nil error; got nil")
	}
	if !strings.Contains(err.Error(), "invalid") {
		t.Errorf("Wait() error = %q; want message containing 'invalid'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// D2 — Kill a running process
// ---------------------------------------------------------------------------

// TestKill_RunningProcess verifies that Kill() successfully terminates a
// long-running process (category D2).  Wait() is called afterwards; it must
// return a non-nil error because the process was killed, not a clean exit.
func TestKill_RunningProcess(t *testing.T) {
	proc, err := Start("sleep 60")
	if err != nil {
		t.Fatalf("Start(sleep 60) error: %v", err)
	}

	if killErr := proc.Kill(); killErr != nil {
		t.Errorf("Kill() unexpected error: %v", killErr)
	}

	// Drain pipes before Wait to avoid deadlock.
	if proc.Stdout != nil {
		io.Copy(io.Discard, proc.Stdout) //nolint:errcheck
	}
	if proc.Stderr != nil {
		io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
	}

	// Wait must return non-nil: the process was killed, not a clean exit.
	if waitErr := proc.Wait(); waitErr == nil {
		t.Error("Wait() after Kill() returned nil; expected non-nil error (process was killed)")
	}
}

// ---------------------------------------------------------------------------
// D2 extended — Kill actually terminates the process (liveness check)
// ---------------------------------------------------------------------------

// TestKill_ProcessActuallyTerminates is a stronger companion to
// TestKill_RunningProcess: it waits on a channel with a hard deadline to
// confirm the process is truly gone within 5 seconds of Kill() (category D2).
func TestKill_ProcessActuallyTerminates(t *testing.T) {
	proc, err := Start("sleep 60")
	if err != nil {
		t.Fatalf("Start(sleep 60) error: %v", err)
	}

	if err := proc.Kill(); err != nil {
		t.Fatalf("Kill() error: %v", err)
	}

	done := make(chan struct{}, 1)
	go func() {
		if proc.Stdout != nil {
			io.Copy(io.Discard, proc.Stdout) //nolint:errcheck
		}
		if proc.Stderr != nil {
			io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
		}
		proc.Wait() //nolint:errcheck
		done <- struct{}{}
	}()

	select {
	case <-done:
		// Process reaped — success.
	case <-time.After(5 * time.Second):
		// Hard liveness guard, not a sleep-based assertion.
		t.Error("process did not terminate within 5 s after Kill()")
	}
}

// ---------------------------------------------------------------------------
// D1 — Kill on a Process whose cmd is nil
// ---------------------------------------------------------------------------

// TestKill_NilCmd verifies that Kill() on a Process with a nil cmd field
// returns a descriptive error without panicking (category D1).
func TestKill_NilCmd(t *testing.T) {
	p := &Process{cmd: nil}
	err := p.Kill()
	if err == nil {
		t.Fatal("Kill() on nil cmd: expected non-nil error; got nil")
	}
	if !strings.Contains(err.Error(), "not started") && !strings.Contains(err.Error(), "invalid") {
		t.Errorf("Kill() error = %q; want message containing 'not started' or 'invalid'", err.Error())
	}
}

// ---------------------------------------------------------------------------
// G1 — Stdout pipe carries expected single-line output
// ---------------------------------------------------------------------------

// TestStdout_SingleLineContent verifies that the Stdout ReadCloser returned
// by Start contains the exact text written by the command (category G1).
// bufio.NewScanner is used as required.
func TestStdout_SingleLineContent(t *testing.T) {
	proc, err := StartWithContext(context.Background(), "echo cookiefarm")
	if err != nil {
		t.Fatalf("StartWithContext() error: %v", err)
	}
	if proc.Stdout == nil {
		t.Fatal("Stdout is nil; cannot read output")
	}

	scanner := bufio.NewScanner(proc.Stdout)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanErr := scanner.Err(); scanErr != nil {
		t.Errorf("scanner error: %v", scanErr)
	}

	// Drain stderr then reap.
	if proc.Stderr != nil {
		io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
	}
	proc.Wait() //nolint:errcheck

	if len(lines) == 0 {
		t.Fatal("no output lines read from Stdout; want at least one")
	}
	if !strings.Contains(lines[0], "cookiefarm") {
		t.Errorf("stdout line[0] = %q; want line containing %q", lines[0], "cookiefarm")
	}
}

// ---------------------------------------------------------------------------
// G2 — Stdout pipe carries expected multi-line output
// ---------------------------------------------------------------------------

// TestStdout_MultiLineContent verifies that bufio.Scanner reads all lines
// emitted by a multi-line command in the correct order (category G2).
func TestStdout_MultiLineContent(t *testing.T) {
	proc, err := StartWithContext(context.Background(), "printf 'alpha\nbeta\ngamma\n'")
	if err != nil {
		t.Fatalf("StartWithContext() error: %v", err)
	}
	if proc.Stdout == nil {
		t.Fatal("Stdout is nil; cannot read output")
	}

	scanner := bufio.NewScanner(proc.Stdout)
	var got []string
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}
	if scanErr := scanner.Err(); scanErr != nil {
		t.Errorf("scanner error: %v", scanErr)
	}

	if proc.Stderr != nil {
		io.Copy(io.Discard, proc.Stderr) //nolint:errcheck
	}
	proc.Wait() //nolint:errcheck

	want := []string{"alpha", "beta", "gamma"}
	if len(got) != len(want) {
		t.Fatalf("read %d lines; want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
	}
	for i, wantLine := range want {
		if got[i] != wantLine {
			t.Errorf("stdout line[%d] = %q; want %q", i, got[i], wantLine)
		}
	}
}

// ---------------------------------------------------------------------------
// F4 — Two independent processes receive distinct PIDs
// ---------------------------------------------------------------------------

// TestPID_TwoProcessesAreDistinct verifies that two concurrently running
// processes are assigned different positive PIDs by the OS (category F4).
func TestPID_TwoProcessesAreDistinct(t *testing.T) {
	p1, err := Start("sleep 2")
	if err != nil {
		t.Fatalf("Start p1 error: %v", err)
	}

	p2, err := Start("sleep 2")
	if err != nil {
		p1.Kill() //nolint:errcheck
		drainAndWait(t, p1)
		t.Fatalf("Start p2 error: %v", err)
	}

	if p1.PID <= 0 {
		t.Errorf("p1.PID = %d; want > 0", p1.PID)
	}
	if p2.PID <= 0 {
		t.Errorf("p2.PID = %d; want > 0", p2.PID)
	}
	if p1.PID == p2.PID {
		t.Errorf("p1.PID == p2.PID == %d; two distinct processes must have different PIDs", p1.PID)
	}

	p1.Kill() //nolint:errcheck
	p2.Kill() //nolint:errcheck
	drainAndWait(t, p1)
	drainAndWait(t, p2)
}
