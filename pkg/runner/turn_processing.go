package runner

import (
	"context"
	"fmt"

	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// processSingleTurn processes a single turn of the agent loop.
// Similar to OpenAI's _run_single_turn in Python and processSingleTurn in TypeScript.
// It handles agent start hooks, model invocation, response processing, tool execution, and handoffs.
func (r *Runner) processSingleTurn(
	ctx context.Context,
	state *RunState,
	opts *RunOptions,
) (*TurnResult, error) {
	// Increment turn
	state.IncrementTurn()

	// Check max turns
	if state.CurrentTurn > state.MaxTurns {
		return nil, fmt.Errorf("max turns (%d) exceeded", state.MaxTurns)
	}

	// Run agent start hooks if needed (first turn or after handoff)
	// Similar to Python's should_run_agent_start_hooks logic
	if state.ShouldRunAgentStartHooks {
		// Call both run hooks and agent hooks (like Python does)
		if opts.Hooks != nil {
			if err := opts.Hooks.OnAgentStart(ctx, state.CurrentAgent, state.OriginalInput); err != nil {
				return nil, fmt.Errorf("run agent start hook error: %w", err)
			}
		}
		if state.CurrentAgent.Hooks != nil {
			if err := state.CurrentAgent.Hooks.OnAgentStart(ctx, state.CurrentAgent, state.OriginalInput); err != nil {
				return nil, fmt.Errorf("agent start hook error: %w", err)
			}
		}
		state.ShouldRunAgentStartHooks = false
	}

	// Get turn input (combines originalInput + generatedItems)
	turnInput := state.GetTurnInput()

	// Execute model request
	response, err := r.executeModelRequest(
		ctx,
		state.CurrentAgent,
		turnInput,
		state.ConsecutiveToolCalls,
		opts,
		state.CurrentTurn,
		state.ToolUseTracker,
	)

	if err != nil {
		return nil, fmt.Errorf("model request error: %w", err)
	}

	// Store raw response
	state.AddRawResponse(*response)

	// Update usage tracking
	if response.Usage != nil && state.RunContext != nil {
		state.RunContext.AddUsage(
			1, // requests
			response.Usage.PromptTokens,
			response.Usage.CompletionTokens,
			response.Usage.TotalTokens,
		)
	}

	// Process the response
	processedResponse := r.processModelResponse(response, state.CurrentAgent)

	// Track items generated in this step
	newStepItems := make([]result.RunItem, 0)

	// Add message items (always add, even if empty content - for tracking)
	if len(processedResponse.Messages) > 0 {
		newStepItems = append(newStepItems, processedResponse.Messages...)
	} else if response.Content != "" || response.HandoffCall != nil {
		// If no messages were created but we have content or a handoff, create one
		// This ensures handoffs are tracked even if there's no content
		newStepItems = append(newStepItems, &result.MessageItem{
			Role:    "assistant",
			Content: response.Content,
		})
	}

	// Handle structured output
	if state.CurrentAgent.OutputType != nil {
		// Parse and validate structured output
		parsedOutput, err := r.parseStructuredOutput(response.Content, state.CurrentAgent.OutputType)
		if err != nil {
			return nil, fmt.Errorf("structured output parsing failed: %w", err)
		}

		return NewTurnResult(
			state.OriginalInput,
			newStepItems,
			&NextStepFinalOutput{Output: parsedOutput},
			response,
		), nil
	}

	// Handle handoffs
	if len(processedResponse.Handoffs) > 0 {
		// Get pre-step items (everything before this turn)
		preStepItems := make([]result.RunItem, len(state.GeneratedItems))
		copy(preStepItems, state.GeneratedItems)

		turnResult, err := r.executeHandoffs(
			ctx,
			state.CurrentAgent,
			state.OriginalInput,
			preStepItems,
			newStepItems,
			response,
			processedResponse.Handoffs,
			state,
			opts,
		)
		if err != nil {
			return nil, err
		}
		if turnResult != nil {
			return turnResult, nil
		}
	}

	// Handle tool calls
	if len(processedResponse.ToolCalls) > 0 {
		// Execute tools
		toolResults, toolItems, err := r.executeFunctionTools(
			ctx,
			state.CurrentAgent,
			processedResponse.ToolCalls,
			state,
		)
		if err != nil {
			return nil, err
		}

		// Add tool items to new step items (tool calls + tool results)
		newStepItems = append(newStepItems, toolItems...)

		// Track tool usage for reset_tool_choice logic (similar to Python)
		if state.ToolUseTracker != nil {
			toolNames := make([]string, len(processedResponse.ToolCalls))
			for i, tc := range processedResponse.ToolCalls {
				toolNames[i] = tc.ToolCall.Name
			}
			state.ToolUseTracker.AddToolUse(state.CurrentAgent.Name, toolNames)
		}

		// Update consecutive tool calls counter
		if len(processedResponse.ToolCalls) == 1 {
			state.ConsecutiveToolCalls++
		} else {
			state.ConsecutiveToolCalls = 0
		}

		// Check tool use behavior
		toolUseBehavior := r.getToolUseBehavior(state.CurrentAgent)
		shouldStop, finalOutput := toolUseBehavior.ShouldStop(ctx, toolResults)
		if shouldStop {
			return NewTurnResult(
				state.OriginalInput,
				newStepItems,
				&NextStepFinalOutput{Output: finalOutput},
				response,
			), nil
		}

		// Add prompt if too many consecutive tool calls
		// This will be included in the next turn's input via GetTurnInput
		if state.ConsecutiveToolCalls >= 3 {
			newStepItems = append(newStepItems, &result.MessageItem{
				Role:    "user",
				Content: "Now that you have the information from the tool(s), please provide a complete response to my original question.",
			})
		}

		// Continue loop - the assistant message with tool_calls and tool results
		// are already in newStepItems, and GetTurnInput will convert them to input format
		return NewTurnResult(
			state.OriginalInput,
			newStepItems,
			&NextStepRunAgain{},
			response,
		), nil
	}

	// If we have content but no tools/handoffs, it's a final output
	if response.Content != "" && !processedResponse.HasToolsOrHandoffs {
		return NewTurnResult(
			state.OriginalInput,
			newStepItems,
			&NextStepFinalOutput{Output: response.Content},
			response,
		), nil
	}

	// If no content and no tools/handoffs, something went wrong - treat as final output
	if response.Content == "" && !processedResponse.HasToolsOrHandoffs {
		// This shouldn't happen, but if it does, stop the loop
		return NewTurnResult(
			state.OriginalInput,
			newStepItems,
			&NextStepFinalOutput{Output: ""},
			response,
		), nil
	}

	// Default: continue (we have tools/handoffs to process)
	return NewTurnResult(
		state.OriginalInput,
		newStepItems,
		&NextStepRunAgain{},
		response,
	), nil
}

// getToolUseBehavior gets the tool use behavior for an agent
func (r *Runner) getToolUseBehavior(agent AgentType) ToolUseBehavior {
	if agent.ToolUseBehavior == nil {
		return &DefaultToolUseBehavior{}
	}

	// If it's already a ToolUseBehavior instance, return it
	if behavior, ok := agent.ToolUseBehavior.(ToolUseBehavior); ok {
		return behavior
	}

	// If it's a string, convert to behavior
	if behaviorStr, ok := agent.ToolUseBehavior.(string); ok {
		switch behaviorStr {
		case "stop_on_first_tool":
			return &StopOnFirstToolBehavior{}
		case "run_llm_again":
			fallthrough
		default:
			return &DefaultToolUseBehavior{}
		}
	}

	// Default fallback
	return &DefaultToolUseBehavior{}
}
