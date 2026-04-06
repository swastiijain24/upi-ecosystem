package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

type EventType string

const (
	EventAuthSuccess EventType = "auth_success"
	EventAuthFailure EventType = "auth_failure"
	EventKeyCreated  EventType = "key_created"
)

type AuditEvent struct {
	Timestamp time.Time      `json:"timestamp"`
	EventType EventType      `json:"event_type"`
	KeyID     string         `json:"key_id,omitempty"`
	IPAddress string         `json:"ip_address"`
	Resource  string         `json:"resource,omitempty"`
	Reason    string         `json:"reason,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type Logger struct {
	logger *slog.Logger
}

func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Log(ctx context.Context, event *AuditEvent) {
	eventJSON, _ := json.Marshal(event)
	l.logger.InfoContext(ctx, "audit_event",
		slog.String("event_type", string(event.EventType)),
		slog.String("key_id", event.KeyID),
		slog.String("ip_address", event.IPAddress),
		slog.String("event_data", string(eventJSON)),
	)
}

func (l *Logger) LogAuthSuccess(ctx context.Context, keyID, ipAddr, resource string) {
	l.Log(ctx, &AuditEvent{
		Timestamp: time.Now(),
		EventType: EventAuthSuccess,
		KeyID:     keyID,
		IPAddress: ipAddr,
		Resource:  resource,
	})
}

func (l *Logger) LogAuthFailure(ctx context.Context, keyID, reason, ipAddr string) {
	l.Log(ctx, &AuditEvent{
		Timestamp: time.Now(),
		EventType: EventAuthFailure,
		KeyID:     keyID,
		Reason:    reason,
		IPAddress: ipAddr,
	})
}

func (l *Logger) LogKeyCreated(ctx context.Context, keyID, ipAddr string) {
	l.Log(ctx, &AuditEvent{
		Timestamp: time.Now(),
		EventType: EventKeyCreated,
		KeyID:     keyID,
		IPAddress: ipAddr,
	})
}
