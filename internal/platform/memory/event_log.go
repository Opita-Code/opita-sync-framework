package memory

import (
	"sync"

	"opita-sync-framework/internal/engine/events"
)

type EventLog struct {
	mu      sync.RWMutex
	records []events.Record
}

func NewEventLog() *EventLog {
	return &EventLog{records: []events.Record{}}
}

func (l *EventLog) Append(record events.Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, record)
	return nil
}

func (l *EventLog) Records() []events.Record {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]events.Record, len(l.records))
	copy(out, l.records)
	return out
}

func (l *EventLog) RecordsByExecution(executionID string) []events.Record {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]events.Record, 0)
	for _, record := range l.records {
		if record.ExecutionID == executionID {
			out = append(out, record)
		}
	}
	return out
}
