# Lifecycle Hooks Guide

## Overview

Hooks allow you to intercept and customize agent behavior at various points in the execution lifecycle. There are two types of hooks: **Agent Hooks** (attached to specific agents) and **Run Hooks** (attached to runs).

## Table of Contents

1. [Agent Hooks](#agent-hooks)
2. [Run Hooks](#run-hooks)
3. [Complete Examples](#complete-examples)
4. [Use Cases](#use-cases)
5. [Best Practices](#best-practices)

## Agent Hooks

Agent hooks are attached to specific agents and are called during that agent's execution.

### Hook Methods

```go
type Hooks interface {
	// Called when the agent starts processing
	OnAgentStart(ctx context.Context, agent *Agent, input interface{}) error

	// Called before the model is called
	OnBeforeModelCall(ctx context.Context, agent *Agent, request *model.Request) error

	// Called after the model is called
	OnAfterModelCall(ctx context.Context, agent *Agent, response *model.Response) error

	// Called before a tool is called
	OnBeforeToolCall(ctx context.Context, agent *Agent, tool tool.Tool, params map[string]interface{}) error

	// Called after a tool is called
	OnAfterToolCall(ctx context.Context, agent *Agent, tool tool.Tool, result interface{}, err error) error

	// Called before a handoff to another agent
	OnBeforeHandoff(ctx context.Context, agent *Agent, handoffAgent *Agent) error

	// Called after a handoff to another agent
	OnAfterHandoff(ctx context.Context, agent *Agent, handoffAgent *Agent, result interface{}) error

	// Called when the agent finishes processing
	OnAgentEnd(ctx context.Context, agent *Agent, result interface{}) error
}
```

### Example: Logging Agent Hooks

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

// LoggingHooks implements agent.Hooks for logging
type LoggingHooks struct{}

func (h *LoggingHooks) OnAgentStart(ctx context.Context, agent *agent.Agent, input interface{}) error {
	fmt.Printf("[HOOK] Agent %s started with input: %v\n", agent.Name, input)
	return nil
}

func (h *LoggingHooks) OnBeforeModelCall(ctx context.Context, agent *agent.Agent, request *model.Request) error {
	fmt.Printf("[HOOK] Before model call for agent %s\n", agent.Name)
	return nil
}

func (h *LoggingHooks) OnAfterModelCall(ctx context.Context, agent *agent.Agent, response *model.Response) error {
	fmt.Printf("[HOOK] After model call for agent %s, content: %s\n", agent.Name, response.Content)
	return nil
}

func (h *LoggingHooks) OnBeforeToolCall(ctx context.Context, agent *agent.Agent, tool tool.Tool, params map[string]interface{}) error {
	fmt.Printf("[HOOK] Before tool call: %s with params: %v\n", tool.GetName(), params)
	return nil
}

func (h *LoggingHooks) OnAfterToolCall(ctx context.Context, agent *agent.Agent, tool tool.Tool, result interface{}, err error) error {
	if err != nil {
		fmt.Printf("[HOOK] Tool call %s failed: %v\n", tool.GetName(), err)
	} else {
		fmt.Printf("[HOOK] Tool call %s succeeded: %v\n", tool.GetName(), result)
	}
	return nil
}

func (h *LoggingHooks) OnBeforeHandoff(ctx context.Context, agent *agent.Agent, handoffAgent *agent.Agent) error {
	fmt.Printf("[HOOK] Agent %s handing off to %s\n", agent.Name, handoffAgent.Name)
	return nil
}

func (h *LoggingHooks) OnAfterHandoff(ctx context.Context, agent *agent.Agent, handoffAgent *agent.Agent, result interface{}) error {
	fmt.Printf("[HOOK] Handoff from %s to %s completed\n", agent.Name, handoffAgent.Name)
	return nil
}

func (h *LoggingHooks) OnAgentEnd(ctx context.Context, agent *agent.Agent, result interface{}) error {
	fmt.Printf("[HOOK] Agent %s finished\n", agent.Name)
	return nil
}

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	// Create agent with hooks
	agent := agent.NewAgent("Test Agent")
	agent.SetModelProvider(provider)
	agent.WithModel("gpt-4o-mini")
	agent.SetSystemInstructions("You are a helpful assistant.")
	agent.WithHooks(&LoggingHooks{})

	// Add a simple tool
	agent.WithTools(tool.NewFunctionTool(
		"get_time",
		"Get current time",
		func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return "3:45PM", nil
		},
	))

	r := runner.NewRunner().WithDefaultProvider(provider)
	result, err := r.Run(context.Background(), agent, &runner.RunOptions{
		Input: "What time is it?",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result.FinalOutput)
}
```

### Example: Rate Limiting with Hooks

```go
type RateLimitingHooks struct {
	lastCall    time.Time
	minInterval time.Duration
}

func NewRateLimitingHooks(minInterval time.Duration) *RateLimitingHooks {
	return &RateLimitingHooks{
		minInterval: minInterval,
	}
}

func (h *RateLimitingHooks) OnBeforeModelCall(ctx context.Context, agent *agent.Agent, request *model.Request) error {
	now := time.Now()
	if !h.lastCall.IsZero() && now.Sub(h.lastCall) < h.minInterval {
		return fmt.Errorf("rate limit: must wait %v between calls", h.minInterval)
	}
	h.lastCall = now
	return nil
}

// Implement other hook methods...
```

## Run Hooks

Run hooks are attached to runs and are called during the entire run lifecycle, including handoffs between agents.

### Hook Methods

```go
type RunHooks interface {
	// Called when the run starts
	OnRunStart(ctx context.Context, agent *agent.Agent, input interface{}) error

	// Called when an agent starts (first turn or after handoff)
	OnAgentStart(ctx context.Context, agent AgentType, input interface{}) error

	// Called when a turn starts
	OnTurnStart(ctx context.Context, agent *agent.Agent, turn int) error

	// Called when a turn ends
	OnTurnEnd(ctx context.Context, agent *agent.Agent, turn int, result *SingleTurnResult) error

	// Called when the run ends
	OnRunEnd(ctx context.Context, result *result.RunResult) error

	// Called when a handoff occurs
	OnHandoff(ctx context.Context, fromAgent AgentType, toAgent AgentType) error

	// Called before a handoff (deprecated, use OnHandoff)
	OnBeforeHandoff(ctx context.Context, agent AgentType, handoffAgent AgentType) error

	// Called after a handoff completes
	OnAfterHandoff(ctx context.Context, agent AgentType, handoffAgent AgentType, result interface{}) error
}
```

### Example: Performance Monitoring Hooks

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
)

type PerformanceHooks struct {
	startTime    time.Time
	turnTimes    []time.Duration
	agentStarts  map[string]time.Time
}

func NewPerformanceHooks() *PerformanceHooks {
	return &PerformanceHooks{
		agentStarts: make(map[string]time.Time),
	}
}

func (h *PerformanceHooks) OnRunStart(ctx context.Context, agent *agent.Agent, input interface{}) error {
	h.startTime = time.Now()
	fmt.Printf("[PERF] Run started at %s\n", h.startTime.Format(time.RFC3339))
	return nil
}

func (h *PerformanceHooks) OnAgentStart(ctx context.Context, agent runner.AgentType, input interface{}) error {
	agentName := agent.GetName()
	h.agentStarts[agentName] = time.Now()
	fmt.Printf("[PERF] Agent %s started\n", agentName)
	return nil
}

func (h *PerformanceHooks) OnTurnStart(ctx context.Context, agent *agent.Agent, turn int) error {
	fmt.Printf("[PERF] Turn %d started for agent %s\n", turn, agent.Name)
	return nil
}

func (h *PerformanceHooks) OnTurnEnd(ctx context.Context, agent *agent.Agent, turn int, result *runner.SingleTurnResult) error {
	turnTime := time.Since(h.startTime)
	h.turnTimes = append(h.turnTimes, turnTime)
	fmt.Printf("[PERF] Turn %d ended in %v\n", turn, turnTime)
	return nil
}

func (h *PerformanceHooks) OnHandoff(ctx context.Context, fromAgent runner.AgentType, toAgent runner.AgentType) error {
	fromName := fromAgent.GetName()
	toName := toAgent.GetName()
	if startTime, ok := h.agentStarts[fromName]; ok {
		duration := time.Since(startTime)
		fmt.Printf("[PERF] Handoff: %s -> %s (after %v)\n", fromName, toName, duration)
	}
	return nil
}

func (h *PerformanceHooks) OnRunEnd(ctx context.Context, result *runner.RunResult) error {
	totalTime := time.Since(h.startTime)
	avgTurnTime := time.Duration(0)
	if len(h.turnTimes) > 0 {
		sum := time.Duration(0)
		for _, t := range h.turnTimes {
			sum += t
		}
		avgTurnTime = sum / time.Duration(len(h.turnTimes))
	}

	fmt.Printf("[PERF] Run completed in %v\n", totalTime)
	fmt.Printf("[PERF] Average turn time: %v\n", avgTurnTime)
	fmt.Printf("[PERF] Total turns: %d\n", len(h.turnTimes))
	return nil
}

func (h *PerformanceHooks) OnBeforeHandoff(ctx context.Context, agent runner.AgentType, handoffAgent runner.AgentType) error {
	return nil // Use OnHandoff instead
}

func (h *PerformanceHooks) OnAfterHandoff(ctx context.Context, agent runner.AgentType, handoffAgent runner.AgentType, result interface{}) error {
	return nil
}

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	agent1 := agent.NewAgent("Agent 1")
	agent1.SetModelProvider(provider)
	agent1.WithModel("gpt-4o-mini")

	agent2 := agent.NewAgent("Agent 2")
	agent2.SetModelProvider(provider)
	agent2.WithModel("gpt-4o-mini")

	agent1.WithHandoffs(agent2)

	r := runner.NewRunner().WithDefaultProvider(provider)
	result, err := r.Run(context.Background(), agent1, &runner.RunOptions{
		Input: "Test",
		Hooks: NewPerformanceHooks(),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result.FinalOutput)
}
```

## Complete Examples

### Example: Validation Hooks

```go
type ValidationHooks struct{}

func (h *ValidationHooks) OnBeforeModelCall(ctx context.Context, agent *agent.Agent, request *model.Request) error {
	// Validate input before sending to model
	if inputStr, ok := request.Input.(string); ok {
		if len(inputStr) > 10000 {
			return fmt.Errorf("input too long: %d characters", len(inputStr))
		}
	}
	return nil
}

func (h *ValidationHooks) OnAfterModelCall(ctx context.Context, agent *agent.Agent, response *model.Response) error {
	// Validate response
	if len(response.Content) > 5000 {
		return fmt.Errorf("response too long: %d characters", len(response.Content))
	}
	return nil
}
```

### Example: Cost Tracking Hooks

```go
type CostTrackingHooks struct {
	totalTokens int
	modelCosts  map[string]float64
}

func NewCostTrackingHooks() *CostTrackingHooks {
	return &CostTrackingHooks{
		modelCosts: map[string]float64{
			"gpt-4o-mini": 0.00015, // per 1K tokens
			"gpt-4":       0.03,    // per 1K tokens
		},
	}
}

func (h *CostTrackingHooks) OnAfterModelCall(ctx context.Context, agent *agent.Agent, response *model.Response) error {
	// Access usage from context
	runCtxVal := ctx.Value("run_context")
	if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx.Usage != nil {
		h.totalTokens += runCtx.Usage.TotalTokens
		
		modelName := agent.Model.(string)
		costPer1K := h.modelCosts[modelName]
		if costPer1K > 0 {
			cost := float64(runCtx.Usage.TotalTokens) / 1000.0 * costPer1K
			fmt.Printf("[COST] Model: %s, Tokens: %d, Cost: $%.4f\n", 
				modelName, runCtx.Usage.TotalTokens, cost)
		}
	}
	return nil
}

func (h *CostTrackingHooks) OnRunEnd(ctx context.Context, result *runner.RunResult) error {
	fmt.Printf("[COST] Total tokens used: %d\n", h.totalTokens)
	return nil
}
```

## Use Cases

### 1. Logging and Debugging

Use hooks to log all agent activities for debugging:

```go
agent.WithHooks(&LoggingHooks{})
```

### 2. Rate Limiting

Prevent too many API calls:

```go
agent.WithHooks(NewRateLimitingHooks(1 * time.Second))
```

### 3. Cost Tracking

Monitor token usage and costs:

```go
runOptions.Hooks = NewCostTrackingHooks()
```

### 4. Input/Output Validation

Validate data before and after processing:

```go
agent.WithHooks(&ValidationHooks{})
```

### 5. Performance Monitoring

Track execution times:

```go
runOptions.Hooks = NewPerformanceHooks()
```

## Best Practices

1. **Don't Block**: Hooks should return quickly. Avoid long-running operations.

2. **Error Handling**: Return errors from hooks to stop execution if needed.

3. **Context Usage**: Use context for cancellation and timeouts.

4. **State Management**: Use struct fields to maintain state across hook calls.

5. **Selective Implementation**: You don't need to implement all hook methods. Use `DefaultAgentHooks` or `DefaultRunHooks` as base.

## Summary

- ✅ Agent hooks are attached to specific agents
- ✅ Run hooks are attached to runs and track the entire execution
- ✅ Hooks allow you to intercept and customize behavior
- ✅ Use hooks for logging, validation, rate limiting, and monitoring
- ✅ Return errors from hooks to stop execution if needed

