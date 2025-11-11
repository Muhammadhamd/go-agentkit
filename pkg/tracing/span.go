package tracing

import (
	"fmt"
	"time"
)

// SpanData represents the data payload of a span
// Following the same structure as Python/TypeScript SDKs
type SpanData interface {
	// Type returns the type of the span data
	Type() string

	// ToMap converts the span data to a map for JSON export
	ToMap() map[string]interface{}
}

// AgentSpanData represents agent span data
type AgentSpanData struct {
	Name       string   `json:"name"`
	Handoffs   []string `json:"handoffs,omitempty"`
	Tools      []string `json:"tools,omitempty"`
	OutputType string   `json:"output_type,omitempty"`
}

func (d *AgentSpanData) Type() string {
	return "agent"
}

func (d *AgentSpanData) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"type": "agent",
		"name": d.Name,
	}
	if len(d.Handoffs) > 0 {
		m["handoffs"] = d.Handoffs
	}
	if len(d.Tools) > 0 {
		m["tools"] = d.Tools
	}
	if d.OutputType != "" {
		m["output_type"] = d.OutputType
	}
	return m
}

// GenerationSpanData represents generation span data
type GenerationSpanData struct {
	Input       []map[string]interface{} `json:"input,omitempty"`
	Output      []map[string]interface{} `json:"output,omitempty"`
	Model       string                   `json:"model,omitempty"`
	ModelConfig map[string]interface{}   `json:"model_config,omitempty"`
	Usage       map[string]interface{}   `json:"usage,omitempty"`
}

func (d *GenerationSpanData) Type() string {
	return "generation"
}

func (d *GenerationSpanData) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"type": "generation",
	}
	if len(d.Input) > 0 {
		m["input"] = d.Input
	}
	if len(d.Output) > 0 {
		m["output"] = d.Output
	}
	if d.Model != "" {
		m["model"] = d.Model
	}
	if len(d.ModelConfig) > 0 {
		m["model_config"] = d.ModelConfig
	}
	if len(d.Usage) > 0 {
		m["usage"] = d.Usage
	}
	return m
}

// FunctionSpanData represents function/tool call span data
type FunctionSpanData struct {
	Name    string `json:"name"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	MCPData string `json:"mcp_data,omitempty"`
}

func (d *FunctionSpanData) Type() string {
	return "function"
}

func (d *FunctionSpanData) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"type":   "function",
		"name":   d.Name,
		"input":  d.Input,
		"output": d.Output,
	}
	if d.MCPData != "" {
		m["mcp_data"] = d.MCPData
	}
	return m
}

// HandoffSpanData represents handoff span data
type HandoffSpanData struct {
	FromAgent string `json:"from_agent,omitempty"`
	ToAgent   string `json:"to_agent,omitempty"`
}

func (d *HandoffSpanData) Type() string {
	return "handoff"
}

func (d *HandoffSpanData) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"type": "handoff",
	}
	if d.FromAgent != "" {
		m["from_agent"] = d.FromAgent
	}
	if d.ToAgent != "" {
		m["to_agent"] = d.ToAgent
	}
	return m
}

// Span represents a span in the OpenAI tracing API format
// Following the same structure as Python/TypeScript SDKs
type Span struct {
	// Type is always "span"
	Type string `json:"object"`

	// SpanID is the unique identifier for the span
	SpanID string `json:"id"`

	// TraceID is the trace this span belongs to
	TraceID string `json:"trace_id"`

	// ParentID is the parent span ID (if any)
	ParentID string `json:"parent_id,omitempty"`

	// SpanData contains the span data
	SpanData SpanData `json:"-"`

	// StartedAt is when the span started (ISO 8601)
	StartedAt string `json:"started_at,omitempty"`

	// EndedAt is when the span ended (ISO 8601)
	EndedAt string `json:"ended_at,omitempty"`

	// Error contains error information if the span failed
	Error *SpanError `json:"error,omitempty"`
}

// SpanError represents an error in a span
type SpanError struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// NewSpan creates a new span
func NewSpan(traceID string, parentID string, spanData SpanData) *Span {
	return &Span{
		Type:     "span",
		SpanID:   generateSpanID(),
		TraceID:  traceID,
		ParentID: parentID,
		SpanData: spanData,
	}
}

// Start marks the span as started
func (s *Span) Start() {
	s.StartedAt = time.Now().UTC().Format(time.RFC3339Nano)
}

// End marks the span as ended
func (s *Span) End() {
	s.EndedAt = time.Now().UTC().Format(time.RFC3339Nano)
}

// SetError sets an error on the span
func (s *Span) SetError(message string, data map[string]interface{}) {
	s.Error = &SpanError{
		Message: message,
		Data:    data,
	}
}

// ToJSON converts the span to JSON format for export
// Implements ExportableItem interface
func (s *Span) ToJSON() (map[string]interface{}, error) {
	if s.SpanData == nil {
		return nil, fmt.Errorf("span data is nil")
	}

	data := map[string]interface{}{
		"object":   s.Type,
		"id":       s.SpanID,
		"trace_id": s.TraceID,
	}

	if s.ParentID != "" {
		data["parent_id"] = s.ParentID
	}

	// Add span data
	spanDataMap := s.SpanData.ToMap()
	for k, v := range spanDataMap {
		data[k] = v
	}

	if s.StartedAt != "" {
		data["started_at"] = s.StartedAt
	}

	if s.EndedAt != "" {
		data["ended_at"] = s.EndedAt
	}

	if s.Error != nil {
		errorData := map[string]interface{}{
			"message": s.Error.Message,
		}
		if len(s.Error.Data) > 0 {
			errorData["data"] = s.Error.Data
		}
		data["error"] = errorData
	}

	return data, nil
}

// generateSpanID generates a new span ID in the format span_<24_hex_chars>
func generateSpanID() string {
	// Generate a proper hex ID (24 hex chars = 12 bytes)
	// Format: span_ + 24 hex characters
	timestamp := time.Now().UnixNano()
	// Use timestamp to ensure uniqueness, pad to ensure 24 chars after "span_"
	hexStr := fmt.Sprintf("%016x%08x", timestamp, timestamp>>32)
	// Take first 24 characters (after "span_")
	if len(hexStr) > 24 {
		return fmt.Sprintf("span_%s", hexStr[:24])
	}
	// Pad if needed
	padding := 24 - len(hexStr)
	return fmt.Sprintf("span_%s%s", hexStr, fmt.Sprintf("%0*x", padding, 0))
}
