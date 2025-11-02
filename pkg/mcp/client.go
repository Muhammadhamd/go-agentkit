package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Client manages connection to an MCP server
type Client struct {
	transport        Transport
	serverInfo       *MCPServerInfo
	capabilities     *MCPServerCapabilities
	protocolVersion  string
	initialized      bool
	mu               sync.RWMutex
	requestIDCounter int64
	pendingRequests  map[interface{}]chan *JSONRPCResponse
	pendingMu        sync.Mutex
}

// ClientConfig configures an MCP client
type ClientConfig struct {
	Transport       Transport
	ProtocolVersion string
	ClientInfo      MCPClientInfo
	Capabilities    MCPCapabilities
}

// NewClient creates a new MCP client
func NewClient(config ClientConfig) *Client {
	return &Client{
		transport:       config.Transport,
		protocolVersion: config.ProtocolVersion,
		initialized:     false,
		pendingRequests: make(map[interface{}]chan *JSONRPCResponse),
	}
}

// Connect establishes connection and initializes the MCP server
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Connect transport
	if err := c.transport.Connect(ctx); err != nil {
		return NewConnectionError(err)
	}

	// Initialize the server
	if err := c.initialize(ctx); err != nil {
		c.transport.Close()
		return fmt.Errorf("failed to initialize MCP server: %w", err)
	}

	return nil
}

// initialize performs the MCP initialization handshake
func (c *Client) initialize(ctx context.Context) error {
	params := MCPInitializeParams{
		ProtocolVersion: c.protocolVersion,
		Capabilities: MCPCapabilities{
			Experimental: make(map[string]interface{}),
		},
		ClientInfo: MCPClientInfo{
			Name:    "agent-sdk-go",
			Version: "1.0.0",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      c.getNextRequestID(),
		Method:  "initialize",
		Params:  paramsJSON,
	}

	resp, err := c.transport.SendRequest(ctx, req)
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf("MCP initialize error: code=%d, message=%s", resp.Error.Code, resp.Error.Message)
	}

	var result MCPInitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("failed to unmarshal initialize result: %w", err)
	}

	c.serverInfo = &result.ServerInfo
	c.capabilities = &result.Capabilities
	c.protocolVersion = result.ProtocolVersion
	c.initialized = true

	// Send initialized notification
	notif := &JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}

	return c.transport.SendNotification(ctx, notif)
}

// ListTools fetches all available tools from the MCP server
func (c *Client) ListTools(ctx context.Context) ([]MCPTool, error) {
	if !c.isInitialized() {
		return nil, fmt.Errorf("client not initialized")
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      c.getNextRequestID(),
		Method:  "tools/list",
	}

	resp, err := c.transport.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("MCP tools/list error: code=%d, message=%s", resp.Error.Code, resp.Error.Message)
	}

	var result MCPToolsListResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tools/list result: %w", err)
	}

	return result.Tools, nil
}

// CallTool executes a tool call on the MCP server
func (c *Client) CallTool(ctx context.Context, toolCall *MCPToolCall) (*MCPToolResult, error) {
	if !c.isInitialized() {
		return nil, fmt.Errorf("client not initialized")
	}

	params := map[string]interface{}{
		"name":      toolCall.Name,
		"arguments": toolCall.Arguments,
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      c.getNextRequestID(),
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp, err := c.transport.SendRequest(ctx, req)
	if err != nil {
		return nil, NewToolExecutionError(toolCall.Name, err)
	}

	if resp.Error != nil {
		return nil, NewToolExecutionError(toolCall.Name, fmt.Errorf("MCP error: code=%d, message=%s", resp.Error.Code, resp.Error.Message))
	}

	var result MCPToolResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tools/call result: %w", err)
	}

	return &result, nil
}

// ListResources fetches all available resources from the MCP server
func (c *Client) ListResources(ctx context.Context) ([]MCPResource, error) {
	if !c.isInitialized() {
		return nil, fmt.Errorf("client not initialized")
	}

	req := &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      c.getNextRequestID(),
		Method:  "resources/list",
	}

	resp, err := c.transport.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("MCP resources/list error: code=%d, message=%s", resp.Error.Code, resp.Error.Message)
	}

	var result MCPResourcesListResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resources/list result: %w", err)
	}

	return result.Resources, nil
}

// Close closes the MCP client connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.initialized = false
	return c.transport.Close()
}

// GetServerInfo returns the server information
func (c *Client) GetServerInfo() *MCPServerInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.serverInfo
}

// IsInitialized returns whether the client is initialized
func (c *Client) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

// Helper methods
func (c *Client) isInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

func (c *Client) getNextRequestID() interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requestIDCounter++
	return c.requestIDCounter
}

// DefaultProtocolVersion is the default MCP protocol version
const DefaultProtocolVersion = "2024-11-05"
