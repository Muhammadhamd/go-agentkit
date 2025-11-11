package runner

import (
	"context"
	"fmt"

	"github.com/muhammadhamd/go-agentkit/pkg/model"
	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// executeHandoffs processes handoff calls.
// Similar to OpenAI's execute_handoffs in Python and executeHandoffs in TypeScript.
// It handles multiple handoffs (rejecting all but the first), applies input filters,
// and prepares the handoff to the new agent.
func (r *Runner) executeHandoffs(
	ctx context.Context,
	agent AgentType,
	originalInput interface{},
	preStepItems []result.RunItem,
	newStepItems []result.RunItem,
	response *model.Response,
	handoffs []ToolRunHandoff,
	state *RunState,
	opts *RunOptions,
) (*TurnResult, error) {
	if len(handoffs) == 0 {
		return nil, nil
	}

	// Handle multiple handoffs (like Python: reject all but first)
	// If there is more than one handoff, add tool responses that reject those handoffs
	multipleHandoffs := len(handoffs) > 1
	if multipleHandoffs {
		outputMessage := "Multiple handoffs detected, ignoring this one."
		for _, rejectedHandoff := range handoffs[1:] {
			// Add tool result for rejected handoffs (required by OpenAI API)
			rejectionToolResult := &result.ToolResultItem{
				Name:       rejectedHandoff.ToolCall.Name,
				Result:     outputMessage,
				ToolCallID: rejectedHandoff.ToolCall.ID,
			}
			newStepItems = append(newStepItems, rejectionToolResult)
		}
	}

	// Process the first handoff
	handoff := handoffs[0]
	handoffCall := handoff.ToolCall
	handoffAgent := handoff.Agent

	// Extract handoff input from tool call parameters
	var handoffInputStr string
	if inputVal, ok := handoffCall.Parameters["input"]; ok {
		if str, ok := inputVal.(string); ok {
			handoffInputStr = str
		} else {
			handoffInputStr = fmt.Sprintf("%v", inputVal)
		}
	}

	// Create handoff tool result item FIRST (required by OpenAI API)
	// This must be added to newStepItems BEFORE creating HandoffInputData
	// so it's included in the conversation history for the new agent
	transferMessage := fmt.Sprintf(`{"assistant":"%s"}`, handoffAgent.Name)
	handoffToolResult := &result.ToolResultItem{
		Name:       handoffCall.Name,
		Result:     transferMessage,
		ToolCallID: handoffCall.ID,
	}

	// Add tool result to newStepItems BEFORE filtering
	// This ensures the tool response is in the conversation history
	newStepItemsWithToolResult := make([]result.RunItem, len(newStepItems))
	copy(newStepItemsWithToolResult, newStepItems)
	newStepItemsWithToolResult = append(newStepItemsWithToolResult, handoffToolResult)

	// Create HandoffInputData for filtering (now includes tool result)
	handoffInputData := &HandoffInputData{
		InputHistory:    originalInput,
		PreHandoffItems: preStepItems,
		NewItems:        newStepItemsWithToolResult, // Include tool result in filtering
		RunContext:      state.RunContext,
	}

	// Apply handoff input filter if available
	var filteredInput interface{}
	if opts.RunConfig != nil && opts.RunConfig.HandoffInputFilter != nil {
		filtered, err := opts.RunConfig.HandoffInputFilter(handoffInputData.AllItems())
		if err != nil {
			return nil, fmt.Errorf("handoff input filter error: %w", err)
		}
		filteredInput = filtered
	} else {
		// Default: use all items + handoff input string
		allItems := handoffInputData.AllItems()
		allItems = append(allItems, map[string]interface{}{
			"type":    "message",
			"role":    "user",
			"content": handoffInputStr,
		})
		filteredInput = allItems
	}

	// Create handoff call item (for internal tracking)
	handoffCallItem := &result.HandoffItem{
		AgentName: handoffAgent.Name,
		Input:     filteredInput,
	}

	// Build final items list (tool result already added to newStepItemsWithToolResult above)
	newItems := make([]result.RunItem, len(newStepItemsWithToolResult))
	copy(newItems, newStepItemsWithToolResult)
	newItems = append(newItems, handoffCallItem) // Internal handoff tracking item

	// Create next step
	nextStep := &NextStepHandoff{
		NewAgent: handoffAgent,
		Input:    filteredInput,
	}

	// Call both run hooks and agent hooks on handoff (like Python)
	// Python calls: hooks.on_handoff and agent.hooks.on_handoff
	if opts.Hooks != nil {
		if err := opts.Hooks.OnHandoff(ctx, agent, handoffAgent); err != nil {
			return nil, fmt.Errorf("run handoff hook error: %w", err)
		}
	}
	if agent.Hooks != nil {
		if err := agent.Hooks.OnBeforeHandoff(ctx, agent, handoffAgent); err != nil {
			return nil, fmt.Errorf("agent handoff hook error: %w", err)
		}
	}

	return NewTurnResult(originalInput, newItems, nextStep, response), nil
}
