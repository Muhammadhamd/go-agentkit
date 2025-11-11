package tracing

import (
	"time"
)

// AddTraceProcessor adds a new trace processor (matching Python's add_trace_processor)
// This allows adding additional processors without replacing the default OpenAI exporter
// Traces will be sent to both OpenAI backend AND your custom processor
func AddTraceProcessor(processor TracingProcessor) {
	GetGlobalTraceProvider().RegisterProcessor(processor)
}

// SetTraceProcessors sets the list of trace processors (matching Python's set_trace_processors)
// This replaces all existing processors, so you can send only to custom backends
// If you want to keep OpenAI backend, include it in the processors list
func SetTraceProcessors(processors []TracingProcessor) {
	GetGlobalTraceProvider().SetProcessors(processors)
}

// SetTracingDisabled sets whether tracing is globally disabled (matching Python/TypeScript)
func SetTracingDisabled(disabled bool) {
	GetGlobalTraceProvider().SetDisabled(disabled)
}

// SetTracingExportAPIKey sets the OpenAI API key for the backend exporter
// (matching Python's set_tracing_export_api_key)
func SetTracingExportAPIKey(apiKey string) {
	// Get the default OpenAI exporter and set its API key
	provider := GetGlobalTraceProvider()
	provider.mu.Lock()
	defer provider.mu.Unlock()

	// Find the OpenAI exporter in the processors and update its API key
	for _, processor := range provider.processors {
		if batchProcessor, ok := processor.(*BatchTraceProcessor); ok {
			if exporter, ok := batchProcessor.Exporter.(*OpenAITracingExporter); ok {
				exporter.SetAPIKey(apiKey)
				return
			}
		}
	}
}

// ForceFlush forces an immediate flush of all queued traces/spans
func ForceFlush() error {
	return GetGlobalTraceProvider().ForceFlush()
}

// Shutdown shuts down the tracing system
func Shutdown(timeout time.Duration) error {
	return GetGlobalTraceProvider().Shutdown(timeout)
}
