package memory

import (
	"context"
	"sync"

	"opita-sync-framework/internal/engine/events"
)

var _ events.EventLog = (*EventLog)(nil)

type EventLog struct {
	mu      sync.RWMutex
	records []events.Record
}

func NewEventLog() *EventLog {
	return &EventLog{records: []events.Record{}}
}

func (l *EventLog) Append(ctx context.Context, record events.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, record)
	return nil
}

func (l *EventLog) Records(ctx context.Context) []events.Record {
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]events.Record, len(l.records))
	copy(out, l.records)
	return out
}

func (l *EventLog) RecordsByExecution(ctx context.Context, executionID string) []events.Record {
	select {
	case <-ctx.Done():
		return nil
	default:
	}
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
