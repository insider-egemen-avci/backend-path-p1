package processing

import "sync/atomic"

type StatsTracker struct {
	CompletedTransactions uint64
	FailedTransactions    uint64
}

func NewStatsTracker() *StatsTracker {
	return &StatsTracker{}
}

func (tracker *StatsTracker) IncrementCompletedTransactions() {
	atomic.AddUint64(&tracker.CompletedTransactions, 1)
}

func (tracker *StatsTracker) IncrementFailedTransactions() {
	atomic.AddUint64(&tracker.FailedTransactions, 1)
}

type Stats struct {
	Completed uint64 `json:"completed"`
	Failed    uint64 `json:"failed"`
}

func (tracker *StatsTracker) GetStats() Stats {
	return Stats{
		Completed: atomic.LoadUint64(&tracker.CompletedTransactions),
		Failed:    atomic.LoadUint64(&tracker.FailedTransactions),
	}
}
