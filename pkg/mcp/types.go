package mcp

import (
	"encoding/json"
)

// MCP protocol types based on JSON-RPC 2.0 specification

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"` // Can be string or number
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCNotification represents a JSON-RPC 2.0 notification (no ID)
type JSONRPCNotification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// MCPInitializeParams represents the initialize request parameters
type MCPInitializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    MCPCapabilities `json:"capabilities"`
	ClientInfo      MCPClientInfo   `json:"clientInfo"`
}

// MCPCapabilities represents client capabilities
type MCPCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
}

// MCPClientInfo represents client information
type MCPClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPInitializeResult represents the initialize response
type MCPInitializeResult struct {
	ProtocolVersion string                `json:"protocolVersion"`
	Capabilities    MCPServerCapabilities `json:"capabilities"`
	ServerInfo      MCPServerInfo         `json:"serverInfo"`
}

// MCPServerCapabilities represents server capabilities
type MCPServerCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
}

// MCPServerInfo represents server information
type MCPServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPTool represents an MCP tool definition
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPResource represents an MCP resource
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// MCPToolCall represents a tool call request
type MCPToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// MCPToolResult represents a tool call result
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent represents content in a tool result
type MCPContent struct {
	Type string `json:"type"` // "text" or "resource"
	Text string `json:"text,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// MCPToolsListResult represents the result of listing tools
type MCPToolsListResult struct {
	Tools []MCPTool `json:"tools"`
}

// MCPResourcesListResult represents the result of listing resources
type MCPResourcesListResult struct {
	Resources []MCPResource `json:"resources"`
}

// MCPServer represents an MCP server connection (deprecated: use Transport instead)
type MCPServer interface {
	// Connect establishes connection and initializes the server
	Connect(ctx interface{}) error

	// SendRequest sends a JSON-RPC request and waits for response
	SendRequest(ctx interface{}, req *JSONRPCRequest) (*JSONRPCResponse, error)

	// SendNotification sends a JSON-RPC notification (no response expected)
	SendNotification(ctx interface{}, notif *JSONRPCNotification) error

	// Close closes the connection
	Close() error
}
