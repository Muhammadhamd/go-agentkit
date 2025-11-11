package tracing

import (
	"fmt"
	"time"
)

// Trace represents a trace in the OpenAI tracing API format
// Following the same structure as Python/TypeScript SDKs
type Trace struct {
	// Type is always "trace"
	Type string `json:"object"`

	// TraceID is the unique identifier for the trace
	TraceID string `json:"id"`

	// Name is the workflow name
	Name string `json:"workflow_name"`

	// GroupID is an optional grouping identifier
	GroupID string `json:"group_id,omitempty"`

	// Metadata is optional metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// StartedAt is when the trace started (ISO 8601)
	StartedAt string `json:"started_at,omitempty"`

	// EndedAt is when the trace ended (ISO 8601)
	EndedAt string `json:"ended_at,omitempty"`
}

// NewTrace creates a new trace
func NewTrace(name string, traceID string, groupID string, metadata map[string]interface{}) *Trace {
	if traceID == "" {
		traceID = generateTraceID()
	}

	return &Trace{
		Type:     "trace",
		TraceID:  traceID,
		Name:     name,
		GroupID:  groupID,
		Metadata: metadata,
	}
}

// Start marks the trace as started
func (t *Trace) Start() {
	t.StartedAt = time.Now().UTC().Format(time.RFC3339Nano)
}

// End marks the trace as ended
func (t *Trace) End() {
	t.EndedAt = time.Now().UTC().Format(time.RFC3339Nano)
}

// ToJSON converts the trace to JSON format for export
// Implements ExportableItem interface
func (t *Trace) ToJSON() (map[string]interface{}, error) {
	data := map[string]interface{}{
		"object":        t.Type,
		"id":            t.TraceID,
		"workflow_name": t.Name,
	}

	if t.GroupID != "" {
		data["group_id"] = t.GroupID
	}

	if len(t.Metadata) > 0 {
		data["metadata"] = t.Metadata
	}

	if t.StartedAt != "" {
		data["started_at"] = t.StartedAt
	}

	if t.EndedAt != "" {
		data["ended_at"] = t.EndedAt
	}

	return data, nil
}

// generateTraceID generates a new trace ID in the format trace_<32_hex_chars>
func generateTraceID() string {
	// Generate a proper hex ID (32 hex chars = 16 bytes)
	timestamp := time.Now().UnixNano()
	// Use timestamp to ensure uniqueness
	hexStr := fmt.Sprintf("%016x%016x", timestamp, timestamp>>32)
	// Ensure 32 hex characters
	if len(hexStr) > 32 {
		return fmt.Sprintf("trace_%s", hexStr[:32])
	}
	// Pad if needed
	padding := 32 - len(hexStr)
	return fmt.Sprintf("trace_%s%s", hexStr, fmt.Sprintf("%0*x", padding, 0))
}
