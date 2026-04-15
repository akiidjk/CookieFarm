package core

import (
	"context"
	"models"
	"testing"
	"time"

	"server/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Shared test helpers
// =============================================================================

// insertFlagWithAge inserts a flag whose response_time is set to
// (now - age), making it appear "age" old to DeleteFlagByTTL.
func insertFlagWithAge(t *testing.T, store *database.Store, code string, age time.Duration) {
	t.Helper()
	f := sampleFlag(code)
	f.ResponseTime = uint64(time.Now().Add(-age).Unix())
	f.Status = models.StatusAccepted
	insertFlag(t, store, f)
}

// countAllFlags returns the total number of flags in the store.
func countAllFlags(t *testing.T, store *database.Store) int {
	t.Helper()
	flags, err := store.Queries.GetAllFlags(context.Background())
	require.NoError(t, err)
	return len(flags)
}

// newRunnerWithStore creates a fresh store, config, and Runner, returning all
// three so callers that need direct store access can use it.
func newRunnerWithStore(t *testing.T) (*Runner, *database.Store) {
	t.Helper()
	store := newTestStore(t)
	cfg := newTestConfig(t)
	return NewRunner(store, cfg), store
}

// newCancelledCtx returns a context that has already been cancelled and a
// no-op cancel func (safe to call multiple times).
func newCancelledCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx, cancel
}

// newManagedCtx returns a background context whose cancel is registered with
// t.Cleanup so callers never need to remember to cancel.
func newManagedCtx(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx, cancel
}

// startTTLLoop launches ValidateFlagTTL in a goroutine and returns a channel
// that is closed when the function returns.
func startTTLLoop(r *Runner, ctx context.Context, flagTTL uint64, tickTime uint64) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		r.ValidateFlagTTL(ctx, flagTTL, tickTime)
	}()
	return done
}

// assertExitsWithin fails the test if the done channel is not closed within
// the given deadline.
func assertExitsWithin(t *testing.T, done <-chan struct{}, deadline time.Duration, msg string) {
	t.Helper()
	select {
	case <-done:
		// expected
	case <-time.After(deadline):
		t.Fatal(msg)
	}
}

// =============================================================================
// ValidateFlagTTL — context cancellation
// =============================================================================

func TestValidateFlagTTL_ContextCancelledBeforeFirstTick_ExitsCleanly(t *testing.T) {
	r := newTestRunner(t)
	// Cancel immediately so the ticker never fires.
	ctx, _ := newCancelledCtx()
	done := startTTLLoop(r, ctx, 1, 99999)
	assertExitsWithin(t, done, 2*time.Second, "ValidateFlagTTL did not exit after context cancellation")
}

func TestValidateFlagTTL_ContextCancelledMidRun_ExitsCleanly(t *testing.T) {
	r := newTestRunner(t)
	ctx, cancel := newManagedCtx(t)

	// Huge tick interval so the loop blocks on select; cancellation must unblock it.
	done := startTTLLoop(r, ctx, 1, 99999)

	time.Sleep(30 * time.Millisecond)
	cancel()

	assertExitsWithin(t, done, 2*time.Second, "ValidateFlagTTL did not exit after mid-run context cancellation")
}

// =============================================================================
// ValidateFlagTTL — deletion logic
// =============================================================================

func TestValidateFlagTTL_FlagsOlderThanTTL_AreDeleted(t *testing.T) {
	r, store := newRunnerWithStore(t)

	// Flag is 200 s old; TTL window = 1 × 5 = 5 s → must be deleted.
	insertFlagWithAge(t, store, "FLAG{ttl_old}", 200*time.Second)
	require.Equal(t, 1, countAllFlags(t, store))

	ctx, _ := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 5)

	waitFor(t, 6*time.Second, 20*time.Millisecond,
		"old flag was not deleted within TTL window",
		func() bool { return countAllFlags(t, store) == 0 },
	)
}

func TestValidateFlagTTL_FlagsNewerThanTTL_AreNotDeleted(t *testing.T) {
	r, store := newRunnerWithStore(t)

	// Fresh flag (0 s old); TTL window = 1 × 1 = 1 s → must be kept.
	insertFlagWithAge(t, store, "FLAG{ttl_fresh}", 0)
	require.Equal(t, 1, countAllFlags(t, store))

	ctx, cancel := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 1)

	time.Sleep(150 * time.Millisecond)
	cancel()

	assert.Equal(t, 1, countAllFlags(t, store), "fresh flag must not be deleted")
}

func TestValidateFlagTTL_MixedAges_OnlyOldFlagsDeleted(t *testing.T) {
	r, store := newRunnerWithStore(t)

	insertFlagWithAge(t, store, "FLAG{ttl_mix_old}", 300*time.Second)
	insertFlagWithAge(t, store, "FLAG{ttl_mix_fresh}", 0)
	require.Equal(t, 2, countAllFlags(t, store))

	// TTL window = 1 × 5 = 5 s. Old flag (300 s) → deleted; fresh (0 s) → kept.
	ctx, _ := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 5)

	waitFor(t, 6*time.Second, 20*time.Millisecond,
		"old flag was not deleted",
		func() bool { return countAllFlags(t, store) == 1 },
	)

	remaining, err := store.Queries.GetAllFlags(context.Background())
	require.NoError(t, err)
	require.Len(t, remaining, 1)
	assert.Equal(t, "FLAG{ttl_mix_fresh}", remaining[0].FlagCode)
}

func TestValidateFlagTTL_NoFlags_DoesNotPanic(t *testing.T) {
	r := newTestRunner(t)
	ctx, cancel := newManagedCtx(t)

	assert.NotPanics(t, func() {
		go r.ValidateFlagTTL(ctx, 1, 1)
		time.Sleep(150 * time.Millisecond)
		cancel()
	})
}

func TestValidateFlagTTL_ZeroAffectedRows_DoesNotReturnError(t *testing.T) {
	r, _ := newRunnerWithStore(t)
	ctx, cancel := newManagedCtx(t)

	done := startTTLLoop(r, ctx, 1, 1)

	// Let the loop tick a couple of times on an empty store, then stop.
	time.Sleep(150 * time.Millisecond)
	cancel()

	assertExitsWithin(t, done, 2*time.Second, "ValidateFlagTTL hung after zero-row deletion")
}

// =============================================================================
// ValidateFlagTTL — interval calculation
// =============================================================================

func TestValidateFlagTTL_IntervalIsProductOfTTLAndTickTime(t *testing.T) {
	r, store := newRunnerWithStore(t)

	// Interval = 1 × 1 = 1 s. Flag is 2 s old → deleted within the first tick.
	insertFlagWithAge(t, store, "FLAG{ttl_interval}", 2*time.Second)

	ctx, _ := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 1)

	waitFor(t, 2*time.Second, 20*time.Millisecond,
		"flag was not deleted within the expected interval",
		func() bool { return countAllFlags(t, store) == 0 },
	)
}

func TestValidateFlagTTL_MultipleTicksDeleteAccumulatedFlags(t *testing.T) {
	r, store := newRunnerWithStore(t)

	ctx, _ := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 1)

	// Insert the flag after the loop has started.
	time.Sleep(50 * time.Millisecond)
	insertFlagWithAge(t, store, "FLAG{ttl_tick2}", 10*time.Second)

	waitFor(t, 3*time.Second, 20*time.Millisecond,
		"flag inserted after loop start was not eventually deleted",
		func() bool { return countAllFlags(t, store) == 0 },
	)
}

// =============================================================================
// ValidateFlagTTL — TTL boundary conditions
// =============================================================================

func TestValidateFlagTTL_LargeTTLValue_DoesNotPanic(t *testing.T) {
	r := newTestRunner(t)
	ctx, cancel := newManagedCtx(t)

	assert.NotPanics(t, func() {
		go r.ValidateFlagTTL(ctx, ^uint64(0)>>1, 1)
		time.Sleep(30 * time.Millisecond)
		cancel()
	})
}

func TestValidateFlagTTL_TTLEqualToFlagAge_FlagIsDeleted(t *testing.T) {
	r, store := newRunnerWithStore(t)

	// TTL window = 1 × 2 = 2 s. Flag age = 3 s → must be deleted.
	insertFlagWithAge(t, store, "FLAG{ttl_boundary}", 3*time.Second)

	ctx, _ := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 2)

	waitFor(t, 5*time.Second, 20*time.Millisecond,
		"boundary-age flag was not deleted",
		func() bool { return countAllFlags(t, store) == 0 },
	)
}

// =============================================================================
// ValidateFlagTTL — concurrent runners
// =============================================================================

func TestValidateFlagTTL_TwoRunnersOnSameStore_NoPanic(t *testing.T) {
	store := newTestStore(t)
	cfg := newTestConfig(t)
	r1 := NewRunner(store, cfg)
	r2 := NewRunner(store, cfg)

	insertFlagWithAge(t, store, "FLAG{concurrent_ttl_1}", 200*time.Second)
	insertFlagWithAge(t, store, "FLAG{concurrent_ttl_2}", 200*time.Second)

	ctx1, cancel1 := newManagedCtx(t)
	ctx2, cancel2 := newManagedCtx(t)

	assert.NotPanics(t, func() {
		go r1.ValidateFlagTTL(ctx1, 1, 1)
		go r2.ValidateFlagTTL(ctx2, 1, 1)
		time.Sleep(200 * time.Millisecond)
		cancel1()
		cancel2()
	})
}

// =============================================================================
// ValidateFlagTTL — invalid configuration
// =============================================================================

func TestInvalidTickTime_DefaultsTo60Seconds(t *testing.T) {
	r, store := newRunnerWithStore(t)

	// Invalid tickTime (0) should default to 60 s; flag is 70 s old → deleted.
	insertFlagWithAge(t, store, "FLAG{invalid_tick}", 70*time.Second)
	require.Equal(t, 1, countAllFlags(t, store))

	ctx, _ := newManagedCtx(t)
	go r.ValidateFlagTTL(ctx, 1, 0)

	waitFor(t, 80*time.Second, 20*time.Millisecond,
		"flag was not deleted with invalid tickTime",
		func() bool { return countAllFlags(t, store) == 0 },
	)
}
