package hosted

import (
	"context"
	"strings"
)

// ApprovalRequirement defines when approval is required for tool execution
type ApprovalRequirement string

const (
	// ApprovalNever means approval is never required
	ApprovalNever ApprovalRequirement = "never"

	// ApprovalAlways means approval is always required
	ApprovalAlways ApprovalRequirement = "always"

	// ApprovalOnTool means approval is required based on tool-specific logic
	ApprovalOnTool ApprovalRequirement = "on_tool"
)

// OnApprovalCallback is called when approval is requested
// Returns true if approved, false if denied, and an error if something went wrong
type OnApprovalCallback func(ctx context.Context, toolName string, params map[string]interface{}) (bool, error)

// ApprovalRequest represents an approval request
type ApprovalRequest struct {
	ToolName string
	Params   map[string]interface{}
}

// ApprovalResponse represents an approval response
type ApprovalResponse struct {
	Approved bool
	Reason   string
}

// DefaultApprovalChecker provides default approval logic
type DefaultApprovalChecker struct{}

// ShouldRequireApproval checks if approval should be required for a tool call
func (c *DefaultApprovalChecker) ShouldRequireApproval(toolName string, params map[string]interface{}) bool {
	// Default: require approval for tools that might modify state
	// This is a simple heuristic - can be customized

	// Check for common write operations in tool name
	name := strings.ToLower(toolName)
	if strings.Contains(name, "write") ||
		strings.Contains(name, "delete") ||
		strings.Contains(name, "update") ||
		strings.Contains(name, "create") ||
		strings.Contains(name, "modify") {
		return true
	}

	// Check for write operations in params
	if operation, ok := params["operation"].(string); ok {
		op := strings.ToLower(operation)
		if strings.Contains(op, "write") ||
			strings.Contains(op, "delete") ||
			strings.Contains(op, "update") ||
			strings.Contains(op, "create") {
			return true
		}
	}

	return false
}
