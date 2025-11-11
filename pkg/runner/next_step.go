package runner

import (
	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// NextStep represents the next step in the agent loop.
// Similar to OpenAI's NextStep discriminated union in both Python and TypeScript.
// This follows OpenAI's discriminated union pattern for controlling loop flow.
type NextStep interface {
	StepType() string
}

// NextStepRunAgain indicates the loop should continue with another turn
type NextStepRunAgain struct{}

func (n *NextStepRunAgain) StepType() string {
	return "next_step_run_again"
}

// NextStepFinalOutput indicates the agent has produced a final output
type NextStepFinalOutput struct {
	Output interface{}
}

func (n *NextStepFinalOutput) StepType() string {
	return "next_step_final_output"
}

// NextStepHandoff indicates a handoff to another agent
type NextStepHandoff struct {
	NewAgent AgentType
	Input    interface{}
}

func (n *NextStepHandoff) StepType() string {
	return "next_step_handoff"
}

// NextStepInterruption indicates the run was interrupted (e.g., for tool approval)
type NextStepInterruption struct {
	Interruptions []result.RunItem // Items that require approval/intervention
}

func (n *NextStepInterruption) StepType() string {
	return "next_step_interruption"
}
