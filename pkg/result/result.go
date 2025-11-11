package result

import (
	"fmt"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model"
)

// RunItem represents an item generated during a run
type RunItem interface {
	// GetType returns the type of the item
	GetType() string

	// ToInputItem converts the item to an input item
	ToInputItem() interface{}
}

// MessageItem represents a message item
type MessageItem struct {
	Role      string
	Content   string
	ToolCalls []interface{} // Optional: for assistant messages with tool calls
}

// GetType returns the type of the item
func (i *MessageItem) GetType() string {
	return "message"
}

// ToInputItem converts the item to an input item
func (i *MessageItem) ToInputItem() interface{} {
	item := map[string]interface{}{
		"type":    "message",
		"role":    i.Role,
		"content": i.Content,
	}
	// Add tool_calls if present (for assistant messages)
	if len(i.ToolCalls) > 0 {
		item["tool_calls"] = i.ToolCalls
	}
	return item
}

// ToolCallItem represents a tool call item
type ToolCallItem struct {
	Name       string
	Parameters map[string]interface{}
}

// GetType returns the type of the item
func (i *ToolCallItem) GetType() string {
	return "tool_call"
}

// ToInputItem converts the item to an input item
func (i *ToolCallItem) ToInputItem() interface{} {
	return map[string]interface{}{
		"type":       "tool_call",
		"name":       i.Name,
		"parameters": i.Parameters,
	}
}

// ToolResultItem represents a tool result item
type ToolResultItem struct {
	Name       string
	Result     interface{}
	ToolCallID string // ID of the tool call this result corresponds to
}

// GetType returns the type of the item
func (i *ToolResultItem) GetType() string {
	return "tool_result"
}

// ToInputItem converts the item to an input item
// Returns format: { type: "tool_result", tool_call: {...}, tool_result: {...} }
func (i *ToolResultItem) ToInputItem() interface{} {
	return map[string]interface{}{
		"type": "tool_result",
		"tool_call": map[string]interface{}{
			"name": i.Name,
			"id":   i.ToolCallID,
		},
		"tool_result": map[string]interface{}{
			"content": fmt.Sprintf("%v", i.Result),
		},
	}
}

// HandoffItem represents a handoff item
type HandoffItem struct {
	AgentName string
	Input     interface{}
}

// GetType returns the type of the item
func (i *HandoffItem) GetType() string {
	return "handoff"
}

// ToInputItem converts the item to an input item
func (i *HandoffItem) ToInputItem() interface{} {
	return map[string]interface{}{
		"type":       "handoff",
		"agent_name": i.AgentName,
		"input":      i.Input,
	}
}

// RunResult contains the result of an agent run
type RunResult struct {
	// Input is the original input to the run
	Input interface{}

	// NewItems are the new items generated during the run
	NewItems []RunItem

	// RawResponses are the raw LLM responses generated during the run
	RawResponses []model.Response

	// FinalOutput is the output of the last agent
	FinalOutput interface{}

	// InputGuardrailResults are guardrail results for the input
	InputGuardrailResults []GuardrailResult

	// OutputGuardrailResults are guardrail results for the final output
	OutputGuardrailResults []GuardrailResult

	// LastAgent is the last agent that was run
	LastAgent *agent.Agent

	// RunContext contains the shared context and usage statistics
	RunContext interface{}
}

// GuardrailResult represents the result of a guardrail check
type GuardrailResult struct {
	// Name is the name of the guardrail
	Name string

	// Passed indicates whether the guardrail check passed
	Passed bool

	// Message is a message from the guardrail
	Message string

	// Error is an error from the guardrail
	Error error
}

// ToInputList converts the result to an input list
func (r *RunResult) ToInputList() []interface{} {
	var result []interface{}

	// Add the original input
	if input, ok := r.Input.(string); ok {
		messageItem := &MessageItem{
			Role:    "user",
			Content: input,
		}
		result = append(result, messageItem.ToInputItem())
	} else if inputList, ok := r.Input.([]interface{}); ok {
		result = append(result, inputList...)
	}

	// Add the new items
	for _, item := range r.NewItems {
		result = append(result, item.ToInputItem())
	}

	return result
}
