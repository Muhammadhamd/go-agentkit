package runner

import (
	"github.com/muhammadhamd/go-agentkit/pkg/model"
	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// TurnResult represents the result of processing a single turn.
// Similar to OpenAI's SingleStepResult in both Python and TypeScript.
// It contains the generated items, next step, and raw model response for a turn.
type TurnResult struct {
	// OriginalInput is the original input (may be mutated by handoff filters)
	OriginalInput interface{}

	// GeneratedItems are all items generated during this turn
	GeneratedItems []result.RunItem

	// NextStep determines what should happen next
	NextStep NextStep

	// ModelResponse is the raw model response for this turn
	ModelResponse *model.Response
}

// NewTurnResult creates a new TurnResult
func NewTurnResult(originalInput interface{}, generatedItems []result.RunItem, nextStep NextStep, modelResponse *model.Response) *TurnResult {
	return &TurnResult{
		OriginalInput:  originalInput,
		GeneratedItems: generatedItems,
		NextStep:       nextStep,
		ModelResponse:  modelResponse,
	}
}
