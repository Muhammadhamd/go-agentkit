package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

// GetAllMCPToolsConfig configures tool fetching from MCP servers
type GetAllMCPToolsConfig struct {
	Transports             []Transport
	ConvertSchemasToStrict bool
	ToolFilter             func(string) bool // Optional tool name filter
}

// GetAllMCPTools fetches and converts all tools from MCP servers
func GetAllMCPTools(ctx context.Context, config GetAllMCPToolsConfig) ([]tool.Tool, error) {
	var allTools []tool.Tool
	toolNames := make(map[string]bool)

	// Create clients for each transport
	clients := make([]*Client, len(config.Transports))
	for idx, transport := range config.Transports {
		client := NewClient(ClientConfig{
			Transport:       transport,
			ProtocolVersion: DefaultProtocolVersion,
		})

		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to MCP server %d: %w", idx, err)
		}

		clients[idx] = client
	}

	// Fetch tools from each server
	for _, client := range clients {
		mcpTools, err := client.ListTools(ctx)
		if err != nil {
			// Continue with other servers even if one fails
			continue
		}

		for _, mcpTool := range mcpTools {
			// Check for duplicate tool names
			if toolNames[mcpTool.Name] {
				return nil, fmt.Errorf("duplicate tool name: %s", mcpTool.Name)
			}

			// Apply tool filter if provided
			if config.ToolFilter != nil && !config.ToolFilter(mcpTool.Name) {
				continue
			}

			// Convert MCP tool to SDK tool with client reference
			sdkTool, err := ConvertMCPToolToSDKTool(mcpTool, client, config.ConvertSchemasToStrict)
			if err != nil {
				return nil, fmt.Errorf("failed to convert tool %s: %w", mcpTool.Name, err)
			}

			allTools = append(allTools, sdkTool)
			toolNames[mcpTool.Name] = true
		}
	}

	return allTools, nil
}

// ConvertMCPToolToSDKTool converts an MCP tool to an SDK tool
func ConvertMCPToolToSDKTool(mcpTool MCPTool, client *Client, convertSchemasToStrict bool) (tool.Tool, error) {
	// Convert schema format
	schema := convertMCPSchemaToOpenAIFormat(mcpTool.InputSchema, convertSchemasToStrict)

	// Create adapter that executes via MCP client
	adapter := &MCPToolAdapter{
		client:   client,
		toolName: mcpTool.Name,
	}

	// Create function tool that delegates to adapter
	ft := tool.NewFunctionTool(
		mcpTool.Name,
		mcpTool.Description,
		func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return adapter.Execute(ctx, params)
		},
	)

	ft.WithSchema(schema)
	return ft, nil
}

// convertMCPSchemaToOpenAIFormat converts MCP JSON schema to OpenAI-compatible format
func convertMCPSchemaToOpenAIFormat(mcpSchema map[string]interface{}, strict bool) map[string]interface{} {
	if mcpSchema == nil {
		return map[string]interface{}{
			"type":       "object",
			"properties": make(map[string]interface{}),
		}
	}

	result := make(map[string]interface{})
	result["type"] = "object"

	// Copy properties
	if properties, ok := mcpSchema["properties"].(map[string]interface{}); ok {
		convertedProps := make(map[string]interface{})
		for key, value := range properties {
			convertedProps[key] = convertSchemaProperty(value, strict)
		}
		result["properties"] = convertedProps
	} else {
		result["properties"] = make(map[string]interface{})
	}

	// Handle required fields
	if required, ok := mcpSchema["required"].([]interface{}); ok {
		requiredStrs := make([]string, 0, len(required))
		for _, req := range required {
			if str, ok := req.(string); ok {
				requiredStrs = append(requiredStrs, str)
			}
		}
		if len(requiredStrs) > 0 {
			result["required"] = requiredStrs
		}
	} else if requiredStrs, ok := mcpSchema["required"].([]string); ok {
		result["required"] = requiredStrs
	}

	// If strict mode, ensure additionalProperties is false
	if strict {
		result["additionalProperties"] = false
	}

	return result
}

// convertSchemaProperty converts a single schema property
func convertSchemaProperty(prop interface{}, strict bool) map[string]interface{} {
	propMap, ok := prop.(map[string]interface{})
	if !ok {
		return map[string]interface{}{"type": "string"}
	}

	result := make(map[string]interface{})

	// Convert type
	if typeVal, ok := propMap["type"].(string); ok {
		result["type"] = typeVal
	} else {
		result["type"] = "string"
	}

	// Copy description if present
	if desc, ok := propMap["description"].(string); ok {
		result["description"] = desc
	}

	// Handle enum
	if enum, ok := propMap["enum"]; ok {
		result["enum"] = enum
	}

	// Handle items for arrays
	if items, ok := propMap["items"]; ok {
		result["items"] = convertSchemaProperty(items, strict)
	}

	// Handle properties for objects
	if properties, ok := propMap["properties"].(map[string]interface{}); ok {
		convertedProps := make(map[string]interface{})
		for key, value := range properties {
			convertedProps[key] = convertSchemaProperty(value, strict)
		}
		result["properties"] = convertedProps
	}

	// Handle required for objects
	if required, ok := propMap["required"].([]interface{}); ok {
		requiredStrs := make([]string, 0, len(required))
		for _, req := range required {
			if str, ok := req.(string); ok {
				requiredStrs = append(requiredStrs, str)
			}
		}
		if len(requiredStrs) > 0 {
			result["required"] = requiredStrs
		}
	}

	// If strict mode and object type, set additionalProperties
	if strict && result["type"] == "object" {
		result["additionalProperties"] = false
	}

	return result
}

// MCPToolAdapter wraps an MCP client and tool name for execution
type MCPToolAdapter struct {
	client   *Client
	toolName string
}

// NewMCPToolAdapter creates a new MCP tool adapter
func NewMCPToolAdapter(client *Client, toolName string) *MCPToolAdapter {
	return &MCPToolAdapter{
		client:   client,
		toolName: toolName,
	}
}

// Execute executes the tool via MCP client
func (a *MCPToolAdapter) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	toolCall := &MCPToolCall{
		Name:      a.toolName,
		Arguments: params,
	}

	result, err := a.client.CallTool(ctx, toolCall)
	if err != nil {
		return nil, err
	}

	// Convert MCP result to string representation
	if result.IsError {
		return nil, fmt.Errorf("tool execution error: %v", result.Content)
	}

	// Extract text content from result
	var content strings.Builder
	for _, c := range result.Content {
		if c.Type == "text" {
			content.WriteString(c.Text)
		}
	}

	return content.String(), nil
}
