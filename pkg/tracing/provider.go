package tracing

import (
	"context"
	"os"
	"sync"
	"time"
)

// TraceProvider is responsible for creating traces and spans
// Following the same pattern as Python/TypeScript SDKs
type TraceProvider struct {
	processors []TracingProcessor // Support multiple processors (like Python/TypeScript)
	disabled   bool
	mu         sync.RWMutex
}

// NewTraceProvider creates a new trace provider
func NewTraceProvider(processors ...TracingProcessor) *TraceProvider {
	// Check environment variable (matching Python/TypeScript)
	disabled := false
	disableTracingEnv := os.Getenv("OPENAI_AGENTS_DISABLE_TRACING")
	if disableTracingEnv == "1" || disableTracingEnv == "true" {
		disabled = true
	}

	// If OPENAI_API_KEY is not set and we're using OpenAI exporter, disable tracing
	// This prevents creating traces/spans that can't be exported
	if !disabled && len(processors) > 0 {
		// Check if we have an OpenAI exporter
		hasOpenAIExporter := false
		for _, processor := range processors {
			if batchProcessor, ok := processor.(*BatchTraceProcessor); ok {
				if _, ok := batchProcessor.Exporter.(*OpenAITracingExporter); ok {
					hasOpenAIExporter = true
					break
				}
			}
		}

		// If we have OpenAI exporter but no API key, disable tracing
		if hasOpenAIExporter && os.Getenv("OPENAI_API_KEY") == "" {
			disabled = true
		}
	}

	return &TraceProvider{
		processors: processors,
		disabled:   disabled,
	}
}

// SetDisabled sets whether tracing is disabled
func (p *TraceProvider) SetDisabled(disabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disabled = disabled
}

// IsDisabled returns whether tracing is disabled
func (p *TraceProvider) IsDisabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.disabled
}

// RegisterProcessor adds a new trace processor (matching Python's add_trace_processor)
// This allows adding additional processors without replacing the default OpenAI exporter
func (p *TraceProvider) RegisterProcessor(processor TracingProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processors = append(p.processors, processor)
}

// SetProcessors sets the list of trace processors (matching Python's set_trace_processors)
// This replaces all existing processors, so you can send only to custom backends
func (p *TraceProvider) SetProcessors(processors []TracingProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processors = processors
}

// CreateTrace creates a new trace
func (p *TraceProvider) CreateTrace(name string, traceID string, groupID string, metadata map[string]interface{}) *Trace {
	if p.IsDisabled() {
		return nil
	}

	trace := NewTrace(name, traceID, groupID, metadata)
	trace.Start()

	// Notify all processors
	p.mu.RLock()
	processors := p.processors
	p.mu.RUnlock()

	for _, processor := range processors {
		_ = processor.OnTraceStart(context.Background(), trace)
	}

	return trace
}

// CreateSpan creates a new span
func (p *TraceProvider) CreateSpan(traceID string, parentID string, spanData SpanData) *Span {
	if p.IsDisabled() {
		return nil
	}

	span := NewSpan(traceID, parentID, spanData)
	span.Start()

	// Notify all processors
	p.mu.RLock()
	processors := p.processors
	p.mu.RUnlock()

	for _, processor := range processors {
		_ = processor.OnSpanStart(context.Background(), span)
	}

	return span
}

// FinishTrace finishes a trace
func (p *TraceProvider) FinishTrace(trace *Trace) {
	if trace == nil || p.IsDisabled() {
		return
	}

	trace.End()

	// Notify all processors
	p.mu.RLock()
	processors := p.processors
	p.mu.RUnlock()

	for _, processor := range processors {
		_ = processor.OnTraceEnd(context.Background(), trace)
	}
}

// FinishSpan finishes a span
func (p *TraceProvider) FinishSpan(span *Span) {
	if span == nil || p.IsDisabled() {
		return
	}

	span.End()

	// Notify all processors
	p.mu.RLock()
	processors := p.processors
	p.mu.RUnlock()

	for _, processor := range processors {
		_ = processor.OnSpanEnd(context.Background(), span)
	}
}

// ForceFlush forces a flush of all queued items
func (p *TraceProvider) ForceFlush() error {
	p.mu.RLock()
	processors := p.processors
	p.mu.RUnlock()

	var lastErr error
	for _, processor := range processors {
		if err := processor.ForceFlush(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Shutdown shuts down the provider
func (p *TraceProvider) Shutdown(timeout time.Duration) error {
	p.mu.RLock()
	processors := p.processors
	p.mu.RUnlock()

	var lastErr error
	for _, processor := range processors {
		if err := processor.Shutdown(timeout); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Global provider
var (
	globalProvider   *TraceProvider
	globalProviderMu sync.Mutex
)

// GetGlobalTraceProvider returns the global trace provider
func GetGlobalTraceProvider() *TraceProvider {
	globalProviderMu.Lock()
	defer globalProviderMu.Unlock()

	if globalProvider == nil {
		// Check if OPENAI_API_KEY is set - if not, skip tracing entirely (no local files)
		// This matches Python/TypeScript behavior - they don't create local files
		if os.Getenv("OPENAI_API_KEY") == "" {
			// No API key, create disabled provider (NO local files)
			globalProvider = NewTraceProvider()
			globalProvider.SetDisabled(true)
		} else {
			// Create default provider with OpenAI backend exporter (matching Python/TypeScript)
			// Python/TypeScript ONLY use backend exporter, no local files
			exporter := NewOpenAITracingExporter(nil)
			processor := NewBatchTraceProcessor(exporter, nil)
			processor.Start()
			globalProvider = NewTraceProvider(processor)
		}
	}

	return globalProvider
}

// SetGlobalTraceProvider sets the global trace provider
func SetGlobalTraceProvider(provider *TraceProvider) {
	globalProviderMu.Lock()
	defer globalProviderMu.Unlock()
	globalProvider = provider
}
