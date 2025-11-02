package mcp

import "fmt"

// MCPError represents an MCP-specific error
type MCPError struct {
	Code    int
	Message string
	Details interface{}
}

func (e *MCPError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("MCP error %d: %s (%v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("MCP error %d: %s", e.Code, e.Message)
}

// Standard MCP error codes
const (
	// JSON-RPC 2.0 error codes
	ErrCodeParseError     = -32700
	ErrCodeInvalidRequest = -32600
	ErrCodeMethodNotFound = -32601
	ErrCodeInvalidParams  = -32602
	ErrCodeInternalError  = -32603

	// MCP-specific error codes (4000-4999)
	ErrCodeServerInitialization = -32001
	ErrCodeConnectionFailed     = -32002
	ErrCodeToolNotFound         = -32003
	ErrCodeToolExecutionFailed  = -32004
	ErrCodeTransportError       = -32005
)

// Error constructors
func NewParseError(details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeParseError,
		Message: "Parse error",
		Details: details,
	}
}

func NewInvalidRequestError(details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeInvalidRequest,
		Message: "Invalid Request",
		Details: details,
	}
}

func NewMethodNotFoundError(method string) *MCPError {
	return &MCPError{
		Code:    ErrCodeMethodNotFound,
		Message: "Method not found",
		Details: method,
	}
}

func NewInvalidParamsError(details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeInvalidParams,
		Message: "Invalid params",
		Details: details,
	}
}

func NewInternalError(details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeInternalError,
		Message: "Internal error",
		Details: details,
	}
}

func NewConnectionError(details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeConnectionFailed,
		Message: "Connection failed",
		Details: details,
	}
}

func NewToolNotFoundError(toolName string) *MCPError {
	return &MCPError{
		Code:    ErrCodeToolNotFound,
		Message: "Tool not found",
		Details: toolName,
	}
}

func NewToolExecutionError(toolName string, details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeToolExecutionFailed,
		Message: "Tool execution failed",
		Details: fmt.Sprintf("tool: %s, error: %v", toolName, details),
	}
}

func NewTransportError(details interface{}) *MCPError {
	return &MCPError{
		Code:    ErrCodeTransportError,
		Message: "Transport error",
		Details: details,
	}
}
