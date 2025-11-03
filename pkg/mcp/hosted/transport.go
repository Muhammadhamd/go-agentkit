package hosted

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Muhammadhamd/go-agentkit/pkg/mcp"
)

// HTTPSSETransport implements MCP transport over HTTP/SSE
type HTTPSSETransport struct {
	base      mcp.BaseTransport
	url       string
	headers   map[string]string
	client    *http.Client
	sseReader *SSEReader
	mu        sync.Mutex
}

// HTTPSSETransportConfig configures HTTP/SSE transport
type HTTPSSETransportConfig struct {
	URL     string
	Headers map[string]string
	Timeout time.Duration
}

var _ mcp.Transport = (*HTTPSSETransport)(nil)

// NewHTTPSSETransport creates a new HTTP/SSE transport
func NewHTTPSSETransport(config HTTPSSETransportConfig) mcp.Transport {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &HTTPSSETransport{
		url:     config.URL,
		headers: config.Headers,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Connect establishes HTTP/SSE connection
func (t *HTTPSSETransport) Connect(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.base.IsConnected() {
		return fmt.Errorf("already connected")
	}

	// For SSE, we'll use HTTP POST for requests and SSE for responses
	// This is a simplified implementation - actual MCP over HTTP may vary

	t.base.SetConnected(true)
	return nil
}

// SendRequest sends a JSON-RPC request via HTTP POST
func (t *HTTPSSETransport) SendRequest(ctx context.Context, req *mcp.JSONRPCRequest) (*mcp.JSONRPCResponse, error) {
	if !t.base.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, mcp.NewParseError(err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", t.url, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, mcp.NewTransportError(err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range t.headers {
		httpReq.Header.Set(key, value)
	}

	if os.Getenv("DEBUG") == "1" {
		fmt.Println("MCP POST", t.url, "method:", req.Method)
		for key, value := range t.headers {
			fmt.Printf("Header %s: %s\n", key, value)
		}
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, mcp.NewTransportError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, mcp.NewTransportError(fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	var mcpResp mcp.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return nil, mcp.NewParseError(err)
	}

	return &mcpResp, nil
}

// SendNotification sends a JSON-RPC notification
func (t *HTTPSSETransport) SendNotification(ctx context.Context, notif *mcp.JSONRPCNotification) error {
	if !t.base.IsConnected() {
		return fmt.Errorf("not connected")
	}

	notifJSON, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", t.url, bytes.NewReader(notifJSON))
	if err != nil {
		return mcp.NewTransportError(err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range t.headers {
		httpReq.Header.Set(key, value)
	}

	// Fire and forget for notifications
	_, err = t.client.Do(httpReq)
	return err
}

// Close closes the transport connection
func (t *HTTPSSETransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.base.SetConnected(false)
	if t.sseReader != nil {
		t.sseReader.Close()
		t.sseReader = nil
	}
	return nil
}

// IsConnected implements Transport interface
func (t *HTTPSSETransport) IsConnected() bool {
	return t.base.IsConnected()
}

// SSEReader reads Server-Sent Events
type SSEReader struct {
	reader io.ReadCloser
	closed bool
}

// NewSSEReader creates a new SSE reader
func NewSSEReader(reader io.ReadCloser) *SSEReader {
	return &SSEReader{
		reader: reader,
		closed: false,
	}
}

// ReadEvent reads the next SSE event
func (r *SSEReader) ReadEvent() ([]byte, error) {
	if r.closed {
		return nil, io.EOF
	}

	var builder strings.Builder
	buf := make([]byte, 1)

	for {
		n, err := r.reader.Read(buf)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			continue
		}

		if buf[0] == '\n' {
			line := builder.String()
			builder.Reset()

			if strings.HasPrefix(line, "data: ") {
				return []byte(line[6:]), nil
			}
			if line == "" {
				// Empty line indicates end of event
				continue
			}
		} else {
			builder.WriteByte(buf[0])
		}
	}
}

// Close closes the SSE reader
func (r *SSEReader) Close() error {
	r.closed = true
	if r.reader != nil {
		return r.reader.Close()
	}
	return nil
}
