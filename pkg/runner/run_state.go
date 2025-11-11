package runner

import (
	"github.com/muhammadhamd/go-agentkit/pkg/model"
	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// RunState tracks the state of an agent run.
// Similar to OpenAI's RunState in both Python and TypeScript.
// It maintains originalInput separately from generatedItems, following OpenAI's pattern.
type RunState struct {
	// OriginalInput is the initial input to the run, never mutated
	OriginalInput interface{}

	// GeneratedItems are items generated during the run (messages, tool calls, etc.)
	// These accumulate over time
	GeneratedItems []result.RunItem

	// CurrentAgent is the agent currently handling the conversation
	CurrentAgent AgentType

	// CurrentStep determines what should happen next in the loop
	CurrentStep NextStep

	// CurrentTurn is the current turn number
	CurrentTurn int

	// MaxTurns is the maximum number of turns allowed
	MaxTurns int

	// ConsecutiveToolCalls tracks consecutive calls to the same tool
	ConsecutiveToolCalls int

	// RawResponses stores all raw model responses
	RawResponses []model.Response

	// LastTurnResponse stores the last model response for interruption handling
	LastTurnResponse *model.Response

	// RunContext shares context across handoffs (usage, approvals, custom context)
	RunContext *RunContext

	// ToolUseTracker tracks which tools each agent has used
	ToolUseTracker *AgentToolUseTracker

	// ShouldRunAgentStartHooks controls when agent start hooks should run
	// Set to true on first turn or after handoff
	ShouldRunAgentStartHooks bool
}

// NewRunState creates a new RunState
func NewRunState(agent AgentType, input interface{}, maxTurns int, runContext *RunContext) *RunState {
	if runContext == nil {
		runContext = NewRunContext(nil)
	}

	return &RunState{
		OriginalInput:            input,
		GeneratedItems:           make([]result.RunItem, 0),
		CurrentAgent:             agent,
		CurrentStep:              &NextStepRunAgain{},
		CurrentTurn:              0,
		MaxTurns:                 maxTurns,
		ConsecutiveToolCalls:     0,
		RawResponses:             make([]model.Response, 0),
		LastTurnResponse:         nil,
		RunContext:               runContext,
		ToolUseTracker:           NewAgentToolUseTracker(),
		ShouldRunAgentStartHooks: true, // Run on first turn
	}
}

// GetTurnInput combines originalInput and generatedItems to create the input for the next turn
// This follows OpenAI's getTurnInput pattern, filtering out internal items like approvals
func (s *RunState) GetTurnInput() []interface{} {
	// Convert original input to list format
	var originalItems []interface{}
	if inputStr, ok := s.OriginalInput.(string); ok {
		originalItems = []interface{}{
			map[string]interface{}{
				"type":    "message",
				"role":    "user",
				"content": inputStr,
			},
		}
	} else if inputList, ok := s.OriginalInput.([]interface{}); ok {
		originalItems = make([]interface{}, len(inputList))
		copy(originalItems, inputList)
	} else {
		originalItems = []interface{}{}
	}

	// Convert generated items to input format, filtering out internal items
	generatedInputItems := make([]interface{}, 0, len(s.GeneratedItems))
	for _, item := range s.GeneratedItems {
		// Filter out approval items (similar to OpenAI's pattern)
		if item.GetType() == "tool_approval" {
			continue
		}
		// Filter out tool_call items - tool calls are already in the assistant message
		// Only include tool results and messages
		if item.GetType() == "tool_call" {
			continue
		}
		// Filter out handoff items - they're internal tracking items
		// The tool result for handoff is what matters for conversation history
		if item.GetType() == "handoff" {
			continue
		}
		generatedInputItems = append(generatedInputItems, item.ToInputItem())
	}

	// Combine original and generated items
	return append(originalItems, generatedInputItems...)
}

// AddGeneratedItem adds a new item to the generated items list
func (s *RunState) AddGeneratedItem(item result.RunItem) {
	s.GeneratedItems = append(s.GeneratedItems, item)
}

// AddGeneratedItems adds multiple items to the generated items list
func (s *RunState) AddGeneratedItems(items []result.RunItem) {
	s.GeneratedItems = append(s.GeneratedItems, items...)
}

// IncrementTurn increments the turn counter
func (s *RunState) IncrementTurn() {
	s.CurrentTurn++
}

// AddRawResponse adds a raw model response
func (s *RunState) AddRawResponse(response model.Response) {
	s.RawResponses = append(s.RawResponses, response)
	s.LastTurnResponse = &response
}
