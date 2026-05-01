package database

import (
	"context"
	"errors"
	"logger"
	"sync"
	"time"
)

const (
	defaultMaxBufferSize = 500
	minBufferSize        = 100
	maxBufferSizeLimit   = 5000
	flushInterval        = 10 * time.Second
	flushTimeout         = 10 * time.Second
)

// FlagCollector manages the collection and flushing of flags to the database
type FlagCollector struct {
	buffer        []Flag         // Buffer for storing flags
	mutex         sync.Mutex     // Mutex for thread-safe access
	flushTimer    *time.Timer    // Timer for periodic flushing
	stopChan      chan struct{}  // Channel to signal stop
	running       bool           // Indicates if the collector is running
	flushCond     *sync.Cond     // Condition variable for flushing
	stats         CollectorStats // Statistics about the collector
	store         *Store         // Reference to the database store
	maxBufferSize int
}

// CollectorStats holds statistics about the flag collector
type CollectorStats struct {
	TotalFlagsReceived  int
	TotalFlushes        int
	SuccessfulFlushes   int
	FailedFlushes       int
	LastFlushTime       time.Time
	LastSuccessfulFlush time.Time
	TotalFlagsFlushed   int
	LastError           error
}

// Singleton instance of the FlagCollector
var (
	collector *FlagCollector
	once      sync.Once
)

// GetCollector return the instance of FlagCollector singleton
func GetCollector() *FlagCollector {
	once.Do(func() {
		c := &FlagCollector{
			buffer:        make([]Flag, 0, defaultMaxBufferSize),
			stopChan:      make(chan struct{}),
			maxBufferSize: defaultMaxBufferSize,
		}
		c.flushCond = sync.NewCond(&c.mutex)
		collector = c
	})
	return collector
}

// SetStore sets the database store on the collector.
// This must be called before the collector flushes any flags.
func (fc *FlagCollector) SetStore(s *Store) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	fc.store = s
}

func (fc *FlagCollector) GetMaxBufferSize() int {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()
	return fc.maxBufferSize
}

func (fc *FlagCollector) adjustBufferSize() {
	stats := fc.stats

	totalFlushes := stats.SuccessfulFlushes + stats.FailedFlushes
	if totalFlushes < 5 {
		return
	}
	failRate := float64(stats.FailedFlushes) / float64(totalFlushes)

	switch {
	case failRate > 0.3:
		fc.maxBufferSize = max(minBufferSize, fc.maxBufferSize/2)
		logger.Log.Warn().
			Float64("fail_rate", failRate).
			Int("new_max", fc.maxBufferSize).
			Msg("High flush failure rate, reducing buffer size")

	case failRate < 0.05 && stats.TotalFlagsFlushed > 1000:
		fc.maxBufferSize = min(maxBufferSizeLimit, fc.maxBufferSize*2)
		logger.Log.Debug().
			Int("new_max", fc.maxBufferSize).
			Msg("Healthy flush rate, increasing buffer size")
	}
}

// Start the collector to begin collecting flags
func (fc *FlagCollector) Start() {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	if fc.running {
		return
	}

	fc.running = true
	fc.flushTimer = time.NewTimer(flushInterval)

	go func() {
		for {
			select {
			case <-fc.flushTimer.C:
				logger.Log.Debug().Msg("Flushing flags on timer")
				ctx, cancel := context.WithTimeout(context.Background(), flushTimeout)
				err := fc.FlushWithContext(ctx)
				cancel()

				if err != nil {
					logger.Log.Error().Err(err).Msg("Error flushing flag buffer on timer")
				}

				fc.mutex.Lock()
				if fc.running {
					fc.flushTimer.Reset(flushInterval)
				}
				fc.mutex.Unlock()
			case <-fc.stopChan:
				fc.mutex.Lock()
				if fc.flushTimer != nil {
					fc.flushTimer.Stop()
				}
				fc.running = false
				fc.flushCond.Broadcast()
				fc.mutex.Unlock()
				return
			}
		}
	}()

	logger.Log.Info().
		Int("max_buffer", defaultMaxBufferSize).
		Dur("flush_interval", flushInterval).
		Msg("Flag collector started")
}

// Stop the collector and flush any remaining flags
func (fc *FlagCollector) Stop() error {
	fc.mutex.Lock()

	if !fc.running {
		fc.mutex.Unlock()
		return nil
	}

	close(fc.stopChan)

	ctx, cancel := context.WithTimeout(context.Background(), flushTimeout)
	defer cancel()

	fc.mutex.Unlock()

	err := fc.FlushWithContext(ctx)
	fc.running = false

	// Take a locked snapshot of the stats so we don't race with the
	// background goroutine that may still be writing them (Issue 4.8).
	stats := fc.GetStats()
	bufferSize := fc.GetBufferSize()
	logger.Log.Info().
		Int("flags_processed", stats.TotalFlagsFlushed).
		Int("buffer_remaining", bufferSize).
		Int("successful_flushes", stats.SuccessfulFlushes).
		Int("failed_flushes", stats.FailedFlushes).
		Msg("Flag collector stopped")

	return err
}

// AddFlag adds a flag to the collector's buffer
func (fc *FlagCollector) AddFlag(flag Flag) error {
	fc.mutex.Lock()

	if !fc.running && fc.stopChan == nil {
		fc.mutex.Unlock()
		return nil
	}

	if !fc.running {
		fc.mutex.Unlock()
		fc.Start()
		fc.mutex.Lock()
	}

	fc.stats.TotalFlagsReceived++
	fc.buffer = append(fc.buffer, flag)

	if len(fc.buffer) >= fc.maxBufferSize {
		logger.Log.Debug().
			Int("buffer_size", len(fc.buffer)).
			Int("max_buffer_size", fc.maxBufferSize).
			Msg("Flushing flag buffer due to size limit")

		flagsToInsert := make([]Flag, len(fc.buffer))
		copy(flagsToInsert, fc.buffer)
		fc.buffer = fc.buffer[:0]
		fc.adjustBufferSize()

		fc.mutex.Unlock()
		ctx := context.Background()
		var err error

		err = fc.store.BulkInsertThings(ctx, flagsToInsert)

		fc.mutex.Lock()

		fc.updateFlushStats(err, len(flagsToInsert))

		if err != nil {
			if len(fc.buffer)+len(flagsToInsert) <= defaultMaxBufferSize {
				fc.buffer = append(fc.buffer, flagsToInsert...)
			} else {
				logger.Log.Error().
					Int("dropped_flags", len(flagsToInsert)).
					Msg("Buffer overflow, dropped flags due to database error")
			}
			fc.mutex.Unlock()
			return err
		}
	}

	fc.mutex.Unlock()
	return nil
}

// Flush sends all accumulated flags to the database
func (fc *FlagCollector) Flush() error {
	return fc.FlushWithContext(context.Background())
}

// FlushWithContext sends all accumulated flags to the database using the provided context
func (fc *FlagCollector) FlushWithContext(ctx context.Context) error {
	fc.mutex.Lock()

	if len(fc.buffer) == 0 {
		fc.mutex.Unlock()
		return nil
	}

	flagsToInsert := make([]Flag, len(fc.buffer))
	copy(flagsToInsert, fc.buffer)
	fc.buffer = fc.buffer[:0]
	fc.stats.TotalFlushes++

	if fc.store == nil {
		fc.buffer = append(flagsToInsert, fc.buffer...)
		nilErr := errors.New("flag collector has no store set: call SetStore before flushing")
		fc.updateFlushStats(nilErr, 0)
		fc.mutex.Unlock()
		return nilErr
	}

	fc.adjustBufferSize()
	fc.mutex.Unlock()

	err := fc.store.BulkInsertThings(ctx, flagsToInsert)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error flushing flag buffer to database")
	}
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.updateFlushStats(err, len(flagsToInsert))

	if err != nil {
		if len(fc.buffer)+len(flagsToInsert) <= defaultMaxBufferSize {
			fc.buffer = append(fc.buffer, flagsToInsert...)
		} else {
			logger.Log.Error().
				Int("dropped_flags", len(flagsToInsert)).
				Msg("Buffer overflow, dropped flags due to database error")
		}
		return err
	}

	fc.flushCond.Broadcast()

	return nil
}

// updateFlushStats updates the statistics after a flush operation
func (fc *FlagCollector) updateFlushStats(err error, flagCount int) {
	fc.stats.LastFlushTime = time.Now()

	if err != nil {
		fc.stats.FailedFlushes++
		fc.stats.LastError = err
	} else {
		fc.stats.SuccessfulFlushes++
		fc.stats.LastSuccessfulFlush = fc.stats.LastFlushTime
		fc.stats.TotalFlagsFlushed += flagCount

		logger.Log.Debug().
			Int("flag_count", flagCount).
			Time("flush_time", fc.stats.LastFlushTime).
			Int("total_flushed", fc.stats.TotalFlagsFlushed).
			Msg("Successfully flushed flags to database")
	}
}

// GetBufferSize returns the current size of the buffer
func (fc *FlagCollector) GetBufferSize() int {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	return len(fc.buffer)
}

// GetStats returns the current statistics of the collector
func (fc *FlagCollector) GetStats() CollectorStats {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	return fc.stats
}

// IsRunning checks if the collector is currently running
func (fc *FlagCollector) IsRunning() bool {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	return fc.running
}
