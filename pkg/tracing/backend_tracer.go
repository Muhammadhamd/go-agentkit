package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// BackendTracer is a tracer that uses the OpenAI backend exporter
// It bridges the old Event-based API to the new Trace/Span system
type BackendTracer struct {
	provider     *TraceProvider
	currentTrace *Trace
	currentSpan  *Span
	traceID      string
	agentName    string
	mu           sync.Mutex
}

// NewBackendTracer creates a new backend tracer for an agent
func NewBackendTracer(agentName string) (*BackendTracer, error) {
	provider := GetGlobalTraceProvider()
	if provider.IsDisabled() {
		return nil, fmt.Errorf("tracing is disabled")
	}

	// Create a trace for this agent run
	traceID := generateTraceID()
	trace := provider.CreateTrace(
		fmt.Sprintf("Agent: %s", agentName),
		traceID,
		"",
		map[string]interface{}{
			"agent_name": agentName,
		},
	)

	return &BackendTracer{
		provider:     provider,
		currentTrace: trace,
		traceID:      traceID,
		agentName:    agentName,
	}, nil
}

// RecordEvent records an event (bridges old Event API to new Trace/Span system)
func (t *BackendTracer) RecordEvent(ctx context.Context, event Event) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.provider.IsDisabled() || t.currentTrace == nil {
		return
	}

	switch event.Type {
	case EventTypeAgentStart:
		// Already handled in NewBackendTracer
		return

	case EventTypeAgentEnd:
		// Finish the trace
		if t.currentTrace != nil {
			t.provider.FinishTrace(t.currentTrace)
		}
		return

	case EventTypeToolCall:
		// Create a function span for tool call
		inputJSON, _ := json.Marshal(event.Details["parameters"])
		spanData := &FunctionSpanData{
			Name:  fmt.Sprintf("%v", event.Details["tool_name"]),
			Input: string(inputJSON),
		}
		span := t.provider.CreateSpan(t.traceID, "", spanData)
		t.currentSpan = span
		return

	case EventTypeToolResult:
		// Finish the tool span
		if t.currentSpan != nil {
			outputJSON, _ := json.Marshal(event.Details["result"])
			if funcData, ok := t.currentSpan.SpanData.(*FunctionSpanData); ok {
				funcData.Output = string(outputJSON)
			}
			if event.Error != nil {
				t.currentSpan.SetError(event.Error.Error(), nil)
			}
			t.provider.FinishSpan(t.currentSpan)
			t.currentSpan = nil
		}
		return

	case EventTypeModelRequest:
		// Create a generation span for model request
		spanData := &GenerationSpanData{
			Model: fmt.Sprintf("%v", event.Details["model"]),
		}
		if prompt, ok := event.Details["prompt"]; ok {
			// Convert prompt to input format
			if promptList, ok := prompt.([]interface{}); ok {
				input := make([]map[string]interface{}, 0, len(promptList))
				for _, p := range promptList {
					if pm, ok := p.(map[string]interface{}); ok {
						input = append(input, pm)
					}
				}
				spanData.Input = input
			}
		}
		span := t.provider.CreateSpan(t.traceID, "", spanData)
		t.currentSpan = span
		return

	case EventTypeModelResponse:
		// Finish the generation span
		if t.currentSpan != nil {
			if genData, ok := t.currentSpan.SpanData.(*GenerationSpanData); ok {
				if response, ok := event.Details["response"]; ok {
					// Convert response to output format
					if responseList, ok := response.([]interface{}); ok {
						output := make([]map[string]interface{}, 0, len(responseList))
						for _, r := range responseList {
							if rm, ok := r.(map[string]interface{}); ok {
								output = append(output, rm)
							}
						}
						genData.Output = output
					}
				}
			}
			if event.Error != nil {
				t.currentSpan.SetError(event.Error.Error(), nil)
			}
			t.provider.FinishSpan(t.currentSpan)
			t.currentSpan = nil
		}
		return

	case EventTypeHandoff:
		// Create a handoff span
		spanData := &HandoffSpanData{
			FromAgent: t.agentName,
			ToAgent:   fmt.Sprintf("%v", event.Details["to_agent"]),
		}
		span := t.provider.CreateSpan(t.traceID, "", spanData)
		span.End()
		t.provider.FinishSpan(span)
		return

	default:
		// For other event types, create custom spans
		spanData := &CustomSpanData{
			Name: event.Type,
			Data: event.Details,
		}
		span := t.provider.CreateSpan(t.traceID, "", spanData)
		span.End()
		t.provider.FinishSpan(span)
		return
	}
}

// Flush flushes any buffered events
func (t *BackendTracer) Flush() error {
	return t.provider.ForceFlush()
}

// Close closes the tracer
func (t *BackendTracer) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Finish trace if still active
	if t.currentTrace != nil {
		t.provider.FinishTrace(t.currentTrace)
		t.currentTrace = nil
	}

	// Flush before closing
	return t.provider.ForceFlush()
}

// CustomSpanData represents custom span data
type CustomSpanData struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

func (d *CustomSpanData) Type() string {
	return "custom"
}

func (d *CustomSpanData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type": "custom",
		"name": d.Name,
		"data": d.Data,
	}
}
