package tracing

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// TracingProcessor processes traces and spans
type TracingProcessor interface {
	// OnTraceStart is called when a trace starts
	OnTraceStart(ctx context.Context, trace *Trace) error

	// OnTraceEnd is called when a trace ends
	OnTraceEnd(ctx context.Context, trace *Trace) error

	// OnSpanStart is called when a span starts
	OnSpanStart(ctx context.Context, span *Span) error

	// OnSpanEnd is called when a span ends
	OnSpanEnd(ctx context.Context, span *Span) error

	// ForceFlush forces an immediate flush of all queued items
	ForceFlush() error

	// Shutdown shuts down the processor
	Shutdown(timeout time.Duration) error
}

// BatchTraceProcessorOptions configures the batch trace processor
type BatchTraceProcessorOptions struct {
	MaxQueueSize       int
	MaxBatchSize       int
	ScheduleDelay      time.Duration
	ExportTriggerRatio float64
}

// DefaultBatchTraceProcessorOptions returns default options
func DefaultBatchTraceProcessorOptions() *BatchTraceProcessorOptions {
	return &BatchTraceProcessorOptions{
		MaxQueueSize:       8192,
		MaxBatchSize:       128,
		ScheduleDelay:      5 * time.Second,
		ExportTriggerRatio: 0.7,
	}
}

// BatchTraceProcessor batches traces and spans before exporting
// Following the same pattern as Python's BatchTraceProcessor and TypeScript's BatchTraceProcessor
type BatchTraceProcessor struct {
	Exporter          TracingExporter // Exported for custom.go access
	options           *BatchTraceProcessorOptions
	queue             chan ExportableItem
	exportTriggerSize int
	shutdown          chan struct{}
	wg                sync.WaitGroup
	mu                sync.Mutex
	started           bool
}

// NewBatchTraceProcessor creates a new batch trace processor
func NewBatchTraceProcessor(exporter TracingExporter, options *BatchTraceProcessorOptions) *BatchTraceProcessor {
	if options == nil {
		options = DefaultBatchTraceProcessorOptions()
	}

	exportTriggerSize := int(float64(options.MaxQueueSize) * options.ExportTriggerRatio)
	if exportTriggerSize < 1 {
		exportTriggerSize = 1
	}

	return &BatchTraceProcessor{
		Exporter:          exporter,
		options:           options,
		queue:             make(chan ExportableItem, options.MaxQueueSize),
		exportTriggerSize: exportTriggerSize,
		shutdown:          make(chan struct{}),
	}
}

// Start starts the background worker thread
func (p *BatchTraceProcessor) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.started {
		return
	}

	p.started = true
	p.wg.Add(1)
	go p.run()
}

// OnTraceStart queues a trace when it starts
func (p *BatchTraceProcessor) OnTraceStart(ctx context.Context, trace *Trace) error {
	p.Start() // Ensure worker is running

	select {
	case p.queue <- trace:
		return nil
	default:
		// Queue is full - drop trace (non-fatal, matching Python/TypeScript)
		return nil
	}
}

// OnTraceEnd is called when a trace ends (no-op, we send on start)
func (p *BatchTraceProcessor) OnTraceEnd(ctx context.Context, trace *Trace) error {
	return nil
}

// OnSpanStart is called when a span starts (no-op, we send on end)
func (p *BatchTraceProcessor) OnSpanStart(ctx context.Context, span *Span) error {
	return nil
}

// OnSpanEnd queues a span when it ends
func (p *BatchTraceProcessor) OnSpanEnd(ctx context.Context, span *Span) error {
	p.Start() // Ensure worker is running

	select {
	case p.queue <- span:
		return nil
	default:
		// Queue is full - drop span (non-fatal, matching Python/TypeScript)
		return nil
	}
}

// ForceFlush forces an immediate flush of all queued items
func (p *BatchTraceProcessor) ForceFlush() error {
	return p.exportBatches(true)
}

// Shutdown shuts down the processor
func (p *BatchTraceProcessor) Shutdown(timeout time.Duration) error {
	close(p.shutdown)
	p.wg.Wait()

	// Final flush
	return p.exportBatches(true)
}

// run is the background worker that processes the queue
func (p *BatchTraceProcessor) run() {
	defer p.wg.Done()

	ticker := time.NewTicker(200 * time.Millisecond) // Check every 200ms
	defer ticker.Stop()

	nextExportTime := time.Now().Add(p.options.ScheduleDelay)

	for {
		select {
		case <-p.shutdown:
			// Final drain
			p.exportBatches(true)
			return

		case <-ticker.C:
			currentTime := time.Now()
			queueSize := len(p.queue)

			// Export if it's time for scheduled flush or queue is above trigger threshold
			if !currentTime.Before(nextExportTime) || queueSize >= p.exportTriggerSize {
				p.exportBatches(false)
				nextExportTime = time.Now().Add(p.options.ScheduleDelay)
			}
		}
	}
}

// exportBatches exports items in batches
func (p *BatchTraceProcessor) exportBatches(force bool) error {
	for {
		items := make([]ExportableItem, 0, p.options.MaxBatchSize)

		// Gather a batch
		for len(items) < p.options.MaxBatchSize {
			select {
			case item := <-p.queue:
				items = append(items, item)
			default:
				// Queue is empty or we've gathered enough
				goto export
			}
		}

	export:
		if len(items) == 0 {
			break
		}

		// Export the batch
		if err := p.Exporter.Export(items); err != nil {
			// Log error but continue (non-fatal)
			if os.Getenv("DEBUG") == "1" {
				fmt.Fprintf(os.Stderr, "[Tracing] Export error: %v\n", err)
			}
		}

		// If not forcing, only export one batch per call
		if !force {
			break
		}
	}

	return nil
}
