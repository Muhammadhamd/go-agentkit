package mcp

import (
	"context"
)

// Transport defines the interface for MCP server communication
type Transport interface {
	// Connect establishes connection to the MCP server
	Connect(ctx context.Context) error

	// SendRequest sends a JSON-RPC request and returns the response
	SendRequest(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error)

	// SendNotification sends a JSON-RPC notification (no response)
	SendNotification(ctx context.Context, notif *JSONRPCNotification) error

	// Close closes the transport connection
	Close() error

	// IsConnected returns whether the transport is currently connected
	IsConnected() bool
}

// BaseTransport provides common functionality for transports
type BaseTransport struct {
	connected bool
}

// IsConnected returns the connection status
func (t *BaseTransport) IsConnected() bool {
	return t.connected
}

// SetConnected sets the connection status
func (t *BaseTransport) SetConnected(connected bool) {
	t.connected = connected
}
