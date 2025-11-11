package tracing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// TracingExporter is the interface for exporting traces and spans
type TracingExporter interface {
	// Export exports the given traces and spans
	Export(items []ExportableItem) error
}

// ExportableItem represents an item that can be exported (Trace or Span)
type ExportableItem interface {
	ToJSON() (map[string]interface{}, error)
}

// OpenAITracingExporterOptions configures the OpenAI tracing exporter
type OpenAITracingExporterOptions struct {
	APIKey       string
	Organization string
	Project      string
	Endpoint     string
	MaxRetries   int
	BaseDelay    time.Duration
	MaxDelay     time.Duration
}

// DefaultOpenAITracingExporterOptions returns default options
func DefaultOpenAITracingExporterOptions() *OpenAITracingExporterOptions {
	return &OpenAITracingExporterOptions{
		Endpoint:   "https://api.openai.com/v1/traces/ingest",
		MaxRetries: 3,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
	}
}

// OpenAITracingExporter exports traces to OpenAI's tracing API
// Following the same pattern as Python's BackendSpanExporter and TypeScript's OpenAITracingExporter
type OpenAITracingExporter struct {
	options *OpenAITracingExporterOptions
	client  *http.Client
}

// NewOpenAITracingExporter creates a new OpenAI tracing exporter
func NewOpenAITracingExporter(options *OpenAITracingExporterOptions) *OpenAITracingExporter {
	if options == nil {
		options = DefaultOpenAITracingExporterOptions()
	}

	// Get API key from options or environment
	if options.APIKey == "" {
		options.APIKey = os.Getenv("OPENAI_API_KEY")
	}

	// Get organization from options or environment
	if options.Organization == "" {
		options.Organization = os.Getenv("OPENAI_ORG_ID")
	}

	// Get project from options or environment
	if options.Project == "" {
		options.Project = os.Getenv("OPENAI_PROJECT_ID")
	}

	return &OpenAITracingExporter{
		options: options,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAPIKey sets the OpenAI API key for the exporter
func (e *OpenAITracingExporter) SetAPIKey(apiKey string) {
	e.options.APIKey = apiKey
}

// Export exports traces and spans to OpenAI's backend
// Following the same pattern as Python/TypeScript SDKs
func (e *OpenAITracingExporter) Export(items []ExportableItem) error {
	if len(items) == 0 {
		return nil
	}

	if e.options.APIKey == "" {
		// Log warning but don't fail (non-fatal, like Python/TypeScript)
		fmt.Fprintf(os.Stderr, "[non-fatal] Tracing: OPENAI_API_KEY is not set, skipping trace export\n")
		return nil
	}

	// Convert items to JSON format
	data := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		jsonData, err := item.ToJSON()
		if err != nil {
			// Skip invalid items (non-fatal)
			continue
		}
		if jsonData != nil {
			data = append(data, jsonData)
		}
	}

	if len(data) == 0 {
		return nil
	}

	payload := map[string]interface{}{
		"data": data,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build headers (matching Python/TypeScript)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", e.options.APIKey),
		"Content-Type":  "application/json",
		"OpenAI-Beta":   "traces=v1",
		"User-Agent":    "go-agentkit",
	}

	if e.options.Organization != "" {
		headers["OpenAI-Organization"] = e.options.Organization
	}

	if e.options.Project != "" {
		headers["OpenAI-Project"] = e.options.Project
	}

	// Exponential backoff retry loop (matching Python/TypeScript)
	attempts := 0
	delay := e.options.BaseDelay

	for attempts < e.options.MaxRetries {
		req, err := http.NewRequest("POST", e.options.Endpoint, bytes.NewBuffer(payloadJSON))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := e.client.Do(req)
		if err != nil {
			// Network error - retry
			fmt.Fprintf(os.Stderr, "[non-fatal] Tracing: request failed: %v\n", err)
			attempts++
			if attempts >= e.options.MaxRetries {
				return fmt.Errorf("failed to export traces after %d attempts: %w", e.options.MaxRetries, err)
			}

			// Exponential backoff with jitter (10% jitter, matching Python/TypeScript)
			sleepTime := delay + time.Duration(float64(delay)*0.1*float64(attempts))
			time.Sleep(sleepTime)
			delay = time.Duration(float64(delay) * 2)
			if delay > e.options.MaxDelay {
				delay = e.options.MaxDelay
			}
			continue
		}

		// Read response body
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Success
		if resp.StatusCode < 300 {
			// Debug log (matching Python/TypeScript)
			if os.Getenv("DEBUG") == "1" {
				fmt.Fprintf(os.Stderr, "[Tracing] Exported %d items\n", len(data))
			}
			return nil
		}

		// Client error (4xx) - don't retry (matching Python/TypeScript)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			fmt.Fprintf(os.Stderr, "[non-fatal] Tracing client error %d: %s\n", resp.StatusCode, string(body))
			return nil // Don't retry on client errors
		}

		// Server error (5xx) - retry
		fmt.Fprintf(os.Stderr, "[non-fatal] Tracing: server error %d, retrying.\n", resp.StatusCode)
		attempts++
		if attempts >= e.options.MaxRetries {
			return fmt.Errorf("failed to export traces after %d attempts: server error %d", e.options.MaxRetries, resp.StatusCode)
		}

		// Exponential backoff with jitter
		sleepTime := delay + time.Duration(float64(delay)*0.1*float64(attempts))
		time.Sleep(sleepTime)
		delay = time.Duration(float64(delay) * 2)
		if delay > e.options.MaxDelay {
			delay = e.options.MaxDelay
		}
	}

	return fmt.Errorf("failed to export traces after %d attempts", e.options.MaxRetries)
}
