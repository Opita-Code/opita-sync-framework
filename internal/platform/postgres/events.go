package postgres

import (
	"encoding/json"
	"fmt"

	"opita-sync-framework/internal/engine/events"
)

type EventLog struct {
	store *Store
}

func NewEventLog(store *Store) *EventLog {
	return &EventLog{store: store}
}

func (l *EventLog) Append(record events.Record) error {
	raw, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal event record: %w", err)
	}
	_, err = l.store.DB.ExecContext(contextBackground(), `
		insert into event_records (event_id, execution_id, event_type, trace_id, tenant_id, occurred_at, payload)
		values ($1, $2, $3, $4, $5, $6, $7)
	`, record.EventID, record.ExecutionID, record.EventType, record.TraceID, record.TenantID, record.OccurredAt, raw)
	if err != nil {
		return fmt.Errorf("insert event record: %w", err)
	}
	return nil
}

func (l *EventLog) RecordsByExecution(executionID string) []events.Record {
	rows, err := l.store.DB.QueryContext(contextBackground(), `select payload from event_records where execution_id = $1 order by occurred_at asc`, executionID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := make([]events.Record, 0)
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var record events.Record
		if err := json.Unmarshal(raw, &record); err != nil {
			continue
		}
		out = append(out, record)
	}
	return out
}
