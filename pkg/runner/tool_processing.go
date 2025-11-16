package runner

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/muhammadhamd/go-agentkit/pkg/model"
	"github.com/muhammadhamd/go-agentkit/pkg/result"
	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

// ToolUseBehavior determines how tool outputs are handled
// Similar to OpenAI's toolUseBehavior
type ToolUseBehavior interface {
	// ShouldStop determines if execution should stop after tool execution
	// Returns (shouldStop, finalOutput)
	ShouldStop(ctx context.Context, toolResults []ToolResult) (bool, interface{})
}

// ToolResult represents the result of executing a tool
type ToolResult struct {
	ToolName string
	Output   interface{}
	Error    error
}

// DefaultToolUseBehavior is the default behavior: always continue after tools
type DefaultToolUseBehavior struct{}

func (b *DefaultToolUseBehavior) ShouldStop(ctx context.Context, toolResults []ToolResult) (bool, interface{}) {
	return false, nil
}

// StopOnFirstToolBehavior stops execution after the first tool
type StopOnFirstToolBehavior struct{}

func (b *StopOnFirstToolBehavior) ShouldStop(ctx context.Context, toolResults []ToolResult) (bool, interface{}) {
	if len(toolResults) > 0 && toolResults[0].Error == nil {
		return true, toolResults[0].Output
	}
	return false, nil
}

// ProcessedResponse contains the processed model response.
// Similar to OpenAI's ProcessedResponse in both Python and TypeScript.
// It categorizes the model output into messages, tool calls, handoffs, and tracks tools used.
type ProcessedResponse struct {
	ToolCalls          []ToolRunFunction
	Handoffs           []ToolRunHandoff
	Messages           []result.RunItem
	HasToolsOrHandoffs bool
	ToolsUsed          []string // Names of all tools used (for tracking)
}

// ToolRunFunction represents a function tool call to execute.
// Similar to OpenAI's ToolRunFunction type in both SDKs.
type ToolRunFunction struct {
	ToolCall model.ToolCall
	Tool     tool.Tool
}

// ToolRunHandoff represents a handoff call to another agent.
// Similar to OpenAI's ToolRunHandoff type in both SDKs.
type ToolRunHandoff struct {
	ToolCall model.ToolCall
	Agent    AgentType
}

// processModelResponse processes a model response and categorizes the output.
// Similar to OpenAI's processModelResponse in TypeScript and process_model_response in Python.
// It separates the response into messages, tool calls, handoffs, and tracks tools used.
func (r *Runner) processModelResponse(
	response *model.Response,
	agent AgentType,
) *ProcessedResponse {
	processed := &ProcessedResponse{
		ToolCalls: make([]ToolRunFunction, 0),
		Handoffs:  make([]ToolRunHandoff, 0),
		Messages:  make([]result.RunItem, 0),
		ToolsUsed: make([]string, 0),
	}

	// Process HandoffCall if present (direct handoff from model)
	if response.HandoffCall != nil {
		handoffCall := response.HandoffCall
		// Find the handoff agent by name
		var handoffAgent AgentType
		for _, h := range agent.Handoffs {
			if h.Name == handoffCall.AgentName {
				handoffAgent = h
				break
			}
		}

		if handoffAgent != nil {
			// Convert HandoffCall to ToolCall format for consistency
			// IMPORTANT: Use TaskID as the tool_call ID - this must match the ID in the assistant message
			// If TaskID is empty, we'll need to extract it from the message later
			toolCallID := handoffCall.TaskID
			if toolCallID == "" {
				// If TaskID is empty, we need to find the actual tool_call ID from the response
				// This happens when HandoffCall comes directly from the model
				// We'll need to match it with the tool_calls in the message
				// For now, generate one but we'll need to ensure it matches the message
				randomBytes := make([]byte, 8)
				if _, err := rand.Read(randomBytes); err == nil {
					toolCallID = fmt.Sprintf("call_%x", randomBytes)
				} else {
					toolCallID = fmt.Sprintf("handoff_%s", handoffCall.AgentName)
				}
			}
			toolCall := model.ToolCall{
				ID:         toolCallID,
				Name:       fmt.Sprintf("handoff_to_%s", handoffCall.AgentName),
				Parameters: handoffCall.Parameters,
			}
			processed.Handoffs = append(processed.Handoffs, ToolRunHandoff{
				ToolCall: toolCall,
				Agent:    handoffAgent,
			})
			processed.ToolsUsed = append(processed.ToolsUsed, fmt.Sprintf("handoff_to_%s", handoffCall.AgentName))
		} else {
			fmt.Printf("WARNING: Handoff agent '%s' not found in agent.Handoffs\n", handoffCall.AgentName)
		}
	}

	// Process tool calls
	for _, tc := range response.ToolCalls {
		// Check if this is a handoff
		handoffAgent := r.findHandoffAgent(agent, tc.Name)
		if handoffAgent != nil {
			processed.Handoffs = append(processed.Handoffs, ToolRunHandoff{
				ToolCall: tc,
				Agent:    handoffAgent,
			})
			processed.ToolsUsed = append(processed.ToolsUsed, tc.Name)
		} else {
			// Find the tool
			toolToCall := r.findTool(agent, tc.Name)
			if toolToCall != nil {
				processed.ToolCalls = append(processed.ToolCalls, ToolRunFunction{
					ToolCall: tc,
					Tool:     toolToCall,
				})
				processed.ToolsUsed = append(processed.ToolsUsed, tc.Name)
			}
		}
	}

	// Process content as message (assistant message, may include tool_calls or handoffs)
	if response.Content != "" || len(response.ToolCalls) > 0 || response.HandoffCall != nil {
		messageItem := &result.MessageItem{
			Role:    "assistant",
			Content: response.Content,
		}

		// Add tool_calls if present (including handoffs)
		// Store tool call IDs for later matching with tool results
		toolCalls := make([]interface{}, 0)

		// Add regular tool calls
		if len(response.ToolCalls) > 0 {
			for i, tc := range response.ToolCalls {
				toolCallID := tc.ID
				if toolCallID == "" {
					randomBytes := make([]byte, 8)
					if _, err := rand.Read(randomBytes); err == nil {
						toolCallID = fmt.Sprintf("call_%x", randomBytes)
					} else {
						toolCallID = fmt.Sprintf("call_%d", i)
					}
					// Update the tool call ID in the response for consistency
					response.ToolCalls[i].ID = toolCallID
				}

				argsJSON, _ := json.Marshal(tc.Parameters)
				toolCalls = append(toolCalls, map[string]interface{}{
					"id":   toolCallID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      tc.Name,
						"arguments": string(argsJSON),
					},
				})
			}
		}

		// Add handoff as a tool call if present (for conversation consistency)
		// IMPORTANT: Use the SAME tool_call ID that we stored in processed.Handoffs
		// This ensures the tool response ID matches the assistant message tool_call ID
		if response.HandoffCall != nil {
			handoffCall := response.HandoffCall
			// Find the matching handoff in processed.Handoffs to get the correct ID
			var toolCallID string
			for _, h := range processed.Handoffs {
				if h.ToolCall.Name == fmt.Sprintf("handoff_to_%s", handoffCall.AgentName) {
					toolCallID = h.ToolCall.ID
					break
				}
			}
			// Fallback if not found (shouldn't happen, but safety check)
			if toolCallID == "" {
				toolCallID = handoffCall.TaskID
				if toolCallID == "" {
					randomBytes := make([]byte, 8)
					if _, err := rand.Read(randomBytes); err == nil {
						toolCallID = fmt.Sprintf("call_%x", randomBytes)
					} else {
						toolCallID = fmt.Sprintf("handoff_%d", len(toolCalls))
					}
				}
			}

			argsJSON, _ := json.Marshal(handoffCall.Parameters)
			toolCalls = append(toolCalls, map[string]interface{}{
				"id":   toolCallID, // Use the SAME ID from processed.Handoffs
				"type": "function",
				"function": map[string]interface{}{
					"name":      fmt.Sprintf("handoff_to_%s", handoffCall.AgentName),
					"arguments": string(argsJSON),
				},
			})
		}

		if len(toolCalls) > 0 {
			messageItem.ToolCalls = toolCalls
		}

		processed.Messages = append(processed.Messages, messageItem)
	}

	processed.HasToolsOrHandoffs = len(processed.ToolCalls) > 0 || len(processed.Handoffs) > 0

	return processed
}

// findHandoffAgent finds a handoff agent by tool name
func (r *Runner) findHandoffAgent(agent AgentType, toolName string) AgentType {
	// Handoff tools are named "handoff_to_<agent_name>"
	if !strings.HasPrefix(toolName, "handoff_to_") {
		return nil
	}

	agentName := strings.TrimPrefix(toolName, "handoff_to_")
	for _, h := range agent.Handoffs {
		if h.Name == agentName {
			return h
		}
	}

	return nil
}

// findTool finds a tool by name
func (r *Runner) findTool(agent AgentType, toolName string) tool.Tool {
	for _, t := range agent.Tools {
		if t.GetName() == toolName {
			return t
		}
	}
	return nil
}

// executeFunctionTools executes function tool calls.
// Similar to OpenAI's execute_function_tool_calls in Python and executeFunctionToolCalls in TypeScript.
// It executes tools in parallel (when possible), handles approvals, and returns tool results.
func (r *Runner) executeFunctionTools(
	ctx context.Context,
	agent AgentType,
	toolRuns []ToolRunFunction,
	state *RunState,
) ([]ToolResult, []result.RunItem, error) {
	toolResults := make([]ToolResult, 0, len(toolRuns))
	runItems := make([]result.RunItem, 0)

	for i, toolRun := range toolRuns {
		tc := toolRun.ToolCall
		t := toolRun.Tool

		// Create tool call item
		toolCallItem := &result.ToolCallItem{
			Name:       tc.Name,
			Parameters: tc.Parameters,
		}
		runItems = append(runItems, toolCallItem)

		// Call agent hooks
		if agent.Hooks != nil {
			if err := agent.Hooks.OnBeforeToolCall(ctx, agent, t, tc.Parameters); err != nil {
				return nil, nil, fmt.Errorf("before tool call hook error: %w", err)
			}
		}

		// Inject RunContext into context for tool access
		toolCtx := context.WithValue(ctx, "run_context", state.RunContext)

		// Execute the tool
		toolResult, err := t.Execute(toolCtx, tc.Parameters)

		// Get or generate tool call ID
		toolCallID := tc.ID
		if toolCallID == "" {
			randomBytes := make([]byte, 8)
			if _, err := rand.Read(randomBytes); err == nil {
				toolCallID = fmt.Sprintf("call_%x", randomBytes)
			} else {
				toolCallID = fmt.Sprintf("call_%d_%d", len(runItems), i)
			}
		}

		// Create tool result item
		toolResultItem := &result.ToolResultItem{
			Name:       tc.Name,
			Result:     toolResult,
			ToolCallID: toolCallID,
		}
		runItems = append(runItems, toolResultItem)

		// Call agent hooks
		if agent.Hooks != nil {
			if hookErr := agent.Hooks.OnAfterToolCall(ctx, agent, t, toolResult, err); hookErr != nil {
				return nil, nil, fmt.Errorf("after tool call hook error: %w", hookErr)
			}
		}

		// Handle errors
		if err != nil {
			toolResult = fmt.Sprintf("Error: %v", err)
		}

		toolResults = append(toolResults, ToolResult{
			ToolName: tc.Name,
			Output:   toolResult,
			Error:    err,
		})
	}

	return toolResults, runItems, nil
}
