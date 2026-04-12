package submitter

/*
Step 1 – Identify Parameters and Environment Conditions

Unit under test:
- SubmitFlags(flagsChan <-chan database.Flag)
- SubmitFlag(flag database.Flag) (string, error)

Inputs:
1) flagsChan:
   - channel lifecycle (open/closed)
   - number of flags sent
   - payload values inside each flag
2) flag (for SubmitFlag):
   - any database.Flag value

Outputs/observable behavior:
- SubmitFlags: no direct return value; expected to drain channel and terminate after close.
- SubmitFlag: returns success message + nil error, or empty string + error.

Relevant environment:
- submitter delegates network operations to api.SubmitBatchDirect / api.SubmitFlag.
- In tests, we avoid network assumptions; we validate control-flow behavior:
  channel draining, termination, and return contract.

Step 2 – Define Categories

A) Channel cardinality for SubmitFlags
- A1: empty channel (closed immediately)
- A2: below batch threshold (1, 49)
- A3: exactly threshold (50)
- A4: above threshold by 1 (51)
- A5: multiple full batches (100)

B) Channel closure
- B1: closed immediately
- B2: closed after writes

C) SubmitFlag result categories
- C1: success path (api accepts flag) -> success message
- C2: error path (api rejects/unreachable) -> empty message + error

Step 3 – Define Constraints

- SubmitFlags only terminates when channel is closed.
- We can assert deterministic behavior without introspecting private batch state by checking:
  function returns after close for each cardinality category.
- For SubmitFlag:
  - C1/C2 depend on external API behavior.
  - To keep tests deterministic without patching dependencies, use an invariant test:
    whenever error is non-nil, message must be empty; whenever error is nil, message must equal
    "Flag submitted successfully".

Excluded/redundant combinations:
- Payload variations inside database.Flag do not affect batching logic here.
- Re-testing same cardinality class with equivalent values is redundant.

Step 4 – Generate Test Frames

TF1: SubmitFlags with empty channel (A1+B1)
- Input: closed channel
- Expected: returns quickly
- Path: boundary/normal

TF2: SubmitFlags with 1 flag (A2+B2)
- Input: 1 flag then close
- Expected: returns after draining
- Path: normal (below threshold)

TF3: SubmitFlags with 49 flags (A2+B2)
- Input: 49 flags then close
- Expected: returns
- Path: boundary (just below threshold)

TF4: SubmitFlags with 50 flags (A3+B2)
- Input: 50 flags then close
- Expected: returns
- Path: boundary (exact threshold)

TF5: SubmitFlags with 51 flags (A4+B2)
- Input: 51 flags then close
- Expected: returns
- Path: boundary (just above threshold)

TF6: SubmitFlags with 100 flags (A5+B2)
- Input: 100 flags then close
- Expected: returns
- Path: normal (multiple batches)

TF7: SubmitFlag return contract invariant (C1/C2)
- Input: representative flag
- Expected:
  - if err == nil => message == success literal
  - if err != nil => message == ""
- Path: success or error depending on runtime environment
*/

import (
	"fmt"
	"server/database"
	"sync"
	"testing"
	"time"
)

func makeFlag(i int) database.Flag {
	return database.Flag{
		FlagCode:    fmt.Sprintf("FLAG{%d}", i),
		ServiceName: "svc",
		PortService: 1337,
		SubmitTime:  uint64(time.Now().Unix()),
		Status:      0,
		TeamID:      int64(i % 10),
		Username:    "tester",
		ExploitName: "exp.py",
		Msg:         "ok",
	}
}

func runSubmitFlagsAndWait(t *testing.T, n int) {
	t.Helper()

	ch := make(chan database.Flag, n)
	done := make(chan struct{})

	go func() {
		SubmitFlags(ch)
		close(done)
	}()

	for i := range n {
		ch <- makeFlag(i)
	}
	close(ch)

	select {
	case <-done:
		// expected
	case <-time.After(2 * time.Second):
		t.Fatalf("SubmitFlags did not return for n=%d", n)
	}
}

func TestSubmitFlagsCategoryPartitions(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		n    int
	}{
		{
			name: "should_return_when_channel_closed_immediately",
			n:    0, // A1 + B1
		},
		{
			name: "should_return_when_one_flag_sent_and_channel_closed",
			n:    1, // A2 + B2
		},
		{
			name: "should_return_when_49_flags_sent_below_batch_threshold",
			n:    49, // A2 + B2 boundary
		},
		{
			name: "should_return_when_50_flags_sent_exact_batch_threshold",
			n:    50, // A3 + B2 boundary
		},
		{
			name: "should_return_when_51_flags_sent_above_batch_threshold",
			n:    51, // A4 + B2 boundary
		},
		{
			name: "should_return_when_100_flags_sent_multiple_full_batches",
			n:    100, // A5 + B2
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Category-partition: batching behavior by channel cardinality.
			runSubmitFlagsAndWait(t, tc.n)
		})
	}
}

func TestSubmitFlagsShouldHandleConcurrentProducersWhenChannelEventuallyClosed(t *testing.T) {
	t.Parallel()

	const producers = 4
	const perProducer = 25 // total 100 (multiple batches)

	ch := make(chan database.Flag, producers*perProducer)
	done := make(chan struct{})

	go func() {
		SubmitFlags(ch)
		close(done)
	}()

	var wg sync.WaitGroup
	wg.Add(producers)
	for p := range producers {
		go func() {
			defer wg.Done()
			start := p * perProducer
			end := start + perProducer
			for i := start; i < end; i++ {
				ch <- makeFlag(i)
			}
		}()
	}

	wg.Wait()
	close(ch)

	select {
	case <-done:
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("SubmitFlags did not return with concurrent producers")
	}
}

func TestSubmitFlagShouldRespectReturnContractWhenApiOutcomeVaries(t *testing.T) {
	t.Parallel()

	// Category-partition (C1/C2):
	// This test is environment-agnostic and validates the contract regardless of API availability.
	msg, err := SubmitFlag(makeFlag(1))

	if err == nil {
		// Success partition
		if msg != "Flag submitted successfully" {
			t.Fatalf("expected success message, got %q", msg)
		}
		return
	}

	// Error partition
	if msg != "" {
		t.Fatalf("expected empty message on error, got %q (err=%v)", msg, err)
	}
}
