package hosted

import (
	"context"
	"fmt"
	"strings"

	"github.com/Muhammadhamd/agent-sdk-go/pkg/mcp"
	"github.com/Muhammadhamd/agent-sdk-go/pkg/tool"
)

// HostedMCPTool represents a tool that executes on a remote MCP server
type HostedMCPTool struct {
	name            string
	description     string
	schema          map[string]interface{}
	serverLabel     string
	serverURL       string
	allowedTools    []string
	headers         map[string]string
	requireApproval ApprovalRequirement
	onApproval      OnApprovalCallback
	transport       mcp.Transport
	client          *mcp.Client
}

// HostedMCPToolConfig configures a hosted MCP tool
type HostedMCPToolConfig struct {
	ServerLabel     string
	ServerURL       string
	AllowedTools    []string
	Headers         map[string]string
	RequireApproval ApprovalRequirement
	OnApproval      OnApprovalCallback
	Description     string
	Schema          map[string]interface{}
}

// NewHostedMCPTool creates a new hosted MCP tool
func NewHostedMCPTool(config HostedMCPToolConfig) (tool.Tool, error) {
	if config.ServerURL == "" {
		return nil, fmt.Errorf("server URL is required")
	}

	transport := NewHTTPSSETransport(HTTPSSETransportConfig{
		URL:     config.ServerURL,
		Headers: config.Headers,
	})

	client := mcp.NewClient(mcp.ClientConfig{
		Transport:       transport,
		ProtocolVersion: mcp.DefaultProtocolVersion,
	})

	// Determine tool name
	toolName := config.ServerLabel
	if toolName == "" {
		toolName = "hosted_mcp_tool"
	}

	// Determine description
	description := config.Description
	if description == "" {
		description = fmt.Sprintf("Access to MCP server at %s", config.ServerURL)
	}

	return &HostedMCPTool{
		name:            toolName,
		description:     description,
		schema:          config.Schema,
		serverLabel:     config.ServerLabel,
		serverURL:       config.ServerURL,
		allowedTools:    config.AllowedTools,
		headers:         config.Headers,
		requireApproval: config.RequireApproval,
		onApproval:      config.OnApproval,
		transport:       transport,
		client:          client,
	}, nil
}

// GetName implements tool.Tool interface
func (t *HostedMCPTool) GetName() string {
	return t.name
}

// GetDescription implements tool.Tool interface
func (t *HostedMCPTool) GetDescription() string {
	return t.description
}

// GetParametersSchema implements tool.Tool interface
func (t *HostedMCPTool) GetParametersSchema() map[string]interface{} {
	if t.schema != nil {
		return t.schema
	}
	return map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}
}

// Execute implements tool.Tool interface
func (t *HostedMCPTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Check if approval is required
	if t.requireApproval == ApprovalAlways ||
		(t.requireApproval == ApprovalOnTool && t.needsApproval(params)) {
		if t.onApproval != nil {
			approved, err := t.onApproval(ctx, t.name, params)
			if err != nil {
				return nil, err
			}
			if !approved {
				return nil, fmt.Errorf("tool execution not approved")
			}
		} else {
			return nil, fmt.Errorf("approval required but no callback provided")
		}
	}

	// Connect client if not connected
	if !t.client.IsInitialized() {
		if err := t.client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to MCP server: %w", err)
		}
	}

	// Extract tool name from params or use server label
	toolName := t.serverLabel
	if name, ok := params["tool_name"].(string); ok && name != "" {
		toolName = name
	}

	// Check if tool is allowed
	if len(t.allowedTools) > 0 {
		allowed := false
		for _, allowedTool := range t.allowedTools {
			if allowedTool == toolName {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("tool %s is not in allowed tools list", toolName)
		}
	}

	// Execute tool call
	toolCall := &mcp.MCPToolCall{
		Name:      toolName,
		Arguments: params,
	}

	result, err := t.client.CallTool(ctx, toolCall)
	if err != nil {
		return nil, err
	}

	// Convert result to string
	var content strings.Builder
	for _, c := range result.Content {
		if c.Type == "text" {
			content.WriteString(c.Text)
		}
	}

	return content.String(), nil
}

// needsApproval determines if approval is needed based on params
func (t *HostedMCPTool) needsApproval(params map[string]interface{}) bool {
	// Simple heuristic: check for write/destructive operations
	// This can be enhanced based on specific requirements
	if operation, ok := params["operation"].(string); ok {
		operation = strings.ToLower(operation)
		return strings.Contains(operation, "write") ||
			strings.Contains(operation, "delete") ||
			strings.Contains(operation, "update") ||
			strings.Contains(operation, "create")
	}
	return false
}

// Close closes the hosted tool connection
func (t *HostedMCPTool) Close() error {
	if t.client != nil {
		return t.client.Close()
	}
	return nil
}

// CreateHostedMCPTool is a helper function to create a hosted MCP tool (similar to TypeScript version)
func CreateHostedMCPTool(config HostedMCPToolConfig) (tool.Tool, error) {
	return NewHostedMCPTool(config)
}
