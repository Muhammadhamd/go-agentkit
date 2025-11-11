package runner

import (
	"context"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// SingleTurnResult contains the result of a single turn
type SingleTurnResult struct {
	Agent    *agent.Agent
	Response interface{}
	Output   interface{}
}

// RunHooks defines lifecycle hooks for a run
type RunHooks interface {
	// OnRunStart is called when the run starts
	OnRunStart(ctx context.Context, agent *agent.Agent, input interface{}) error

	// OnAgentStart is called when an agent starts (first turn or after handoff)
	// Similar to Python's hooks.on_agent_start
	OnAgentStart(ctx context.Context, agent AgentType, input interface{}) error

	// OnTurnStart is called when a turn starts
	OnTurnStart(ctx context.Context, agent *agent.Agent, turn int) error

	// OnTurnEnd is called when a turn ends
	OnTurnEnd(ctx context.Context, agent *agent.Agent, turn int, result *SingleTurnResult) error

	// OnRunEnd is called when the run ends
	OnRunEnd(ctx context.Context, result *result.RunResult) error

	// OnHandoff is called when a handoff occurs (from_agent -> to_agent)
	// Similar to Python's hooks.on_handoff
	OnHandoff(ctx context.Context, fromAgent AgentType, toAgent AgentType) error

	// OnBeforeHandoff is called before a handoff occurs (deprecated, use OnHandoff)
	OnBeforeHandoff(ctx context.Context, agent AgentType, handoffAgent AgentType) error

	// OnAfterHandoff is called after a handoff completes
	OnAfterHandoff(ctx context.Context, agent AgentType, handoffAgent AgentType, result interface{}) error
}

// DefaultRunHooks provides a default implementation of RunHooks
type DefaultRunHooks struct{}

// OnRunStart is called when the run starts
func (h *DefaultRunHooks) OnRunStart(ctx context.Context, agent *agent.Agent, input interface{}) error {
	return nil
}

// OnAgentStart is called when an agent starts (first turn or after handoff)
func (h *DefaultRunHooks) OnAgentStart(ctx context.Context, agent AgentType, input interface{}) error {
	return nil
}

// OnTurnStart is called when a turn starts
func (h *DefaultRunHooks) OnTurnStart(ctx context.Context, agent *agent.Agent, turn int) error {
	return nil
}

// OnTurnEnd is called when a turn ends
func (h *DefaultRunHooks) OnTurnEnd(ctx context.Context, agent *agent.Agent, turn int, result *SingleTurnResult) error {
	return nil
}

// OnRunEnd is called when the run ends
func (h *DefaultRunHooks) OnRunEnd(ctx context.Context, result *result.RunResult) error {
	return nil
}

// OnHandoff is called when a handoff occurs
func (h *DefaultRunHooks) OnHandoff(ctx context.Context, fromAgent AgentType, toAgent AgentType) error {
	return nil
}

// OnBeforeHandoff is called before a handoff occurs (deprecated, use OnHandoff)
func (h *DefaultRunHooks) OnBeforeHandoff(ctx context.Context, agent AgentType, handoffAgent AgentType) error {
	return nil
}

// OnAfterHandoff is called after a handoff completes
func (h *DefaultRunHooks) OnAfterHandoff(ctx context.Context, agent AgentType, handoffAgent AgentType, result interface{}) error {
	return nil
}
