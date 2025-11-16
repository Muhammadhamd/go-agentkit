# Agent Handoffs Guide

## Overview

Agent handoffs allow one agent to transfer control to another specialized agent. This enables building complex multi-agent workflows where different agents handle different aspects of a task.

## Table of Contents

1. [Basic Handoff Setup](#basic-handoff-setup)
2. [How Handoffs Work](#how-handoffs-work)
3. [System Instructions for Handoffs](#system-instructions-for-handoffs)
4. [Complete Examples](#complete-examples)
5. [Bidirectional Handoffs](#bidirectional-handoffs)
6. [Input Filtering](#input-filtering)
7. [Context Sharing](#context-sharing)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)

## Basic Handoff Setup

### Step 1: Create Multiple Agents

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	// Create specialized agents
	mathAgent := agent.NewAgent("Math Agent")
	mathAgent.SetModelProvider(provider)
	mathAgent.WithModel("gpt-4o-mini")
	mathAgent.SetSystemInstructions("You are a math specialist. Solve mathematical problems.")

	weatherAgent := agent.NewAgent("Weather Agent")
	weatherAgent.SetModelProvider(provider)
	weatherAgent.WithModel("gpt-4o-mini")
	weatherAgent.SetSystemInstructions("You provide weather information.")

	// Create a coordinator agent
	coordinatorAgent := agent.NewAgent("Coordinator")
	coordinatorAgent.SetModelProvider(provider)
	coordinatorAgent.WithModel("gpt-4o-mini")
	coordinatorAgent.SetSystemInstructions(`You coordinate requests by delegating to specialized agents.
For math problems, handoff to "Math Agent".
For weather questions, handoff to "Weather Agent".`)
```

### Step 2: Configure Handoffs

```go
	// Set up handoffs: coordinator can handoff to math and weather agents
	coordinatorAgent.WithHandoffs(mathAgent, weatherAgent)
```

### Step 3: Run the Coordinator

```go
	r := runner.NewRunner().WithDefaultProvider(provider)
	
	result, err := r.Run(context.Background(), coordinatorAgent, &runner.RunOptions{
		Input:    "What is 42 divided by 6?",
		MaxTurns: 20,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result.FinalOutput)
}
```

## How Handoffs Work

### The Handoff Process

1. **Agent Decides to Handoff**: The current agent decides it needs to delegate to another agent
2. **Handoff Tool Call**: The agent calls a special handoff tool (automatically available)
3. **Input Preparation**: The conversation history is prepared for the new agent
4. **Agent Switch**: Control transfers to the new agent
5. **Execution**: The new agent processes the task
6. **Return** (optional): The new agent can return results to the original agent

### Automatic Handoff Tools

When you call `WithHandoffs()`, the framework automatically creates handoff tools for each agent:

- Tool name format: `handoff_to_[AgentName]`
- Parameters: `{"input": "task description"}`
- Example: `handoff_to_Math_Agent` with `{"input": "Calculate 42 / 6"}`

## System Instructions for Handoffs

### For Coordinator Agents

The coordinator agent needs clear instructions on when and how to handoff:

```go
coordinatorAgent.SetSystemInstructions(`You are a coordinator that routes requests to specialized agents.

AVAILABLE AGENTS:
- Math Agent: For any mathematical calculations or problems
- Weather Agent: For weather information and forecasts

HANDOFF INSTRUCTIONS:
1. When you receive a request, determine which specialized agent should handle it
2. Use the handoff_to_[AgentName] tool to delegate
3. In the "input" parameter, provide clear instructions for the specialized agent
4. After the agent completes, provide a final response to the user

EXAMPLES:
- "What is 10 + 5?" → handoff to Math Agent with input "Calculate 10 plus 5"
- "What's the weather in Paris?" → handoff to Weather Agent with input "Get weather for Paris"`)
```

### For Specialized Agents

Specialized agents should focus on their domain:

```go
mathAgent.SetSystemInstructions(`You are a math specialist.

YOUR ROLE:
- Solve mathematical problems accurately
- Show your work when appropriate
- Provide clear, complete answers

IMPORTANT:
- Focus only on mathematical tasks
- If asked about non-math topics, explain that you only handle math`)
```

## Complete Examples

### Example 1: Simple Multi-Agent Workflow

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	// Create calculator tool
	calculatorTool := tool.NewFunctionTool(
		"calculate",
		"Perform a calculation",
		func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			operation := params["operation"].(string)
			a := params["a"].(float64)
			b := params["b"].(float64)

			switch operation {
			case "add":
				return a + b, nil
			case "subtract":
				return a - b, nil
			case "multiply":
				return a * b, nil
			case "divide":
				if b == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return a / b, nil
			default:
				return nil, fmt.Errorf("unknown operation")
			}
		},
	).WithSchema(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"add", "subtract", "multiply", "divide"},
				"description": "The operation to perform",
			},
			"a": map[string]interface{}{
				"type":        "number",
				"description": "First number",
			},
			"b": map[string]interface{}{
				"type":        "number",
				"description": "Second number",
			},
		},
		"required": []string{"operation", "a", "b"},
	})

	// Create Math Agent
	mathAgent := agent.NewAgent("Math Agent")
	mathAgent.SetModelProvider(provider)
	mathAgent.WithModel("gpt-4o-mini")
	mathAgent.SetSystemInstructions(`You are a math specialist.
Use the calculate tool to solve mathematical problems.
Always provide clear explanations of your calculations.`)
	mathAgent.WithTools(calculatorTool)

	// Create Coordinator Agent
	coordinatorAgent := agent.NewAgent("Coordinator")
	coordinatorAgent.SetModelProvider(provider)
	coordinatorAgent.WithModel("gpt-4o-mini")
	coordinatorAgent.SetSystemInstructions(`You coordinate requests.
For any math problems, handoff to "Math Agent" with clear instructions.
After receiving results, provide a friendly response to the user.`)
	coordinatorAgent.WithHandoffs(mathAgent)

	// Run
	r := runner.NewRunner().WithDefaultProvider(provider)
	result, err := r.Run(context.Background(), coordinatorAgent, &runner.RunOptions{
		Input:    "What is 42 divided by 6?",
		MaxTurns: 20,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result.FinalOutput)
}
```

### Example 2: Complex Multi-Agent System

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	// Support Agent
	supportAgent := agent.NewAgent("Support Agent")
	supportAgent.SetModelProvider(provider)
	supportAgent.WithModel("gpt-4o-mini")
	supportAgent.SetSystemInstructions(`You are a customer support agent.
You help customers with inquiries and can delegate to specialized agents.
For billing issues, handoff to "Billing Agent".
For technical issues, handoff to "Technical Agent".`)

	// Billing Agent
	billingAgent := agent.NewAgent("Billing Agent")
	billingAgent.SetModelProvider(provider)
	billingAgent.WithModel("gpt-4o-mini")
	billingAgent.SetSystemInstructions("You handle billing, refunds, and payment issues.")

	// Technical Agent
	technicalAgent := agent.NewAgent("Technical Agent")
	technicalAgent.SetModelProvider(provider)
	technicalAgent.WithModel("gpt-4o-mini")
	technicalAgent.SetSystemInstructions("You handle technical support and troubleshooting.")

	// Set up handoffs
	supportAgent.WithHandoffs(billingAgent, technicalAgent)

	// Run
	r := runner.NewRunner().WithDefaultProvider(provider)
	result, err := r.Run(context.Background(), supportAgent, &runner.RunOptions{
		Input:    "I need a refund for my order",
		MaxTurns: 20,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result.FinalOutput)
}
```

## Bidirectional Handoffs

Bidirectional handoffs allow agents to return to their delegator with results.

### Setup Bidirectional Handoffs

```go
// Method 1: Using WithBidirectionalHandoffs (recommended)
orchestratorAgent.WithBidirectionalHandoffs(workerAgent1, workerAgent2)

// Method 2: Manual setup
orchestratorAgent.WithHandoffs(workerAgent1, workerAgent2)
workerAgent1.WithHandoffs(orchestratorAgent)
workerAgent2.WithHandoffs(orchestratorAgent)
```

### Example: Task Delegation with Return

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	// Orchestrator Agent
	orchestratorAgent := agent.NewAgent("Orchestrator")
	orchestratorAgent.SetModelProvider(provider)
	orchestratorAgent.WithModel("gpt-4o-mini")
	orchestratorAgent.SetSystemInstructions(`You coordinate tasks by delegating to workers.
When a worker completes a task, they will return to you with results.
Analyze the results and provide a final response to the user.`)

	// Worker Agent
	workerAgent := agent.NewAgent("Worker")
	workerAgent.SetModelProvider(provider)
	workerAgent.WithModel("gpt-4o-mini")
	workerAgent.SetSystemInstructions(`You are a worker that processes tasks.
When you complete a task, return to "Orchestrator" with your results using handoff.
Include all relevant information in the handoff input.`)

	// Set up bidirectional handoffs
	orchestratorAgent.WithBidirectionalHandoffs(workerAgent)

	// Run
	r := runner.NewRunner().WithDefaultProvider(provider)
	result, err := r.Run(context.Background(), orchestratorAgent, &runner.RunOptions{
		Input:    "Process this data: [complex data]",
		MaxTurns: 20,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Result:", result.FinalOutput)
}
```

## Input Filtering

You can filter the conversation history passed to the new agent during handoffs.

### Custom Input Filter

```go
func customHandoffFilter(input interface{}) (interface{}, error) {
	// Convert to conversation array
	conversation, ok := input.([]interface{})
	if !ok {
		return input, nil
	}

	// Filter out tool calls and results, keep only messages
	filtered := make([]interface{}, 0)
	for _, item := range conversation {
		if msg, ok := item.(map[string]interface{}); ok {
			msgType := msg["type"]
			if msgType == "message" {
				filtered = append(filtered, item)
			}
		}
	}

	return filtered, nil
}

// Use the filter
result, err := r.Run(context.Background(), agent, &runner.RunOptions{
	Input: "Hello",
	RunConfig: &runner.RunConfig{
		HandoffInputFilter: customHandoffFilter,
	},
})
```

## Context Sharing

Context is automatically shared across handoffs through `RunContext`.

### Example: Shared Context Across Handoffs

```go
type SharedContext struct {
	UserID    string
	OrderID   string
	Metadata  map[string]interface{}
}

func main() {
	// Create shared context
	sharedCtx := &SharedContext{
		UserID:   "user_123",
		OrderID:  "order_456",
		Metadata: make(map[string]interface{}),
	}

	// Create agents
	supportAgent := agent.NewAgent("Support")
	supportAgent.SetModelProvider(provider)
	
	billingAgent := agent.NewAgent("Billing")
	billingAgent.SetModelProvider(provider)
	
	supportAgent.WithHandoffs(billingAgent)

	// Run with shared context
	result, err := r.Run(context.Background(), supportAgent, &runner.RunOptions{
		Input:   "I need help with my order",
		Context: sharedCtx, // Shared across all agents
		MaxTurns: 20,
	})

	// Access updated context
	if result.RunContext != nil {
		if runCtx, ok := result.RunContext.(*runner.RunContext); ok {
			if ctx, ok := runCtx.Context.(*SharedContext); ok {
				fmt.Printf("User ID: %s\n", ctx.UserID)
				fmt.Printf("Order ID: %s\n", ctx.OrderID)
			}
		}
	}
}
```

## Best Practices

### 1. Clear Agent Roles

Each agent should have a clear, focused role:

```go
// ✅ GOOD: Clear, focused role
mathAgent.SetSystemInstructions("You are a math specialist. Solve mathematical problems.")

// ❌ BAD: Vague, unfocused role
mathAgent.SetSystemInstructions("You help with various tasks.")
```

### 2. Explicit Handoff Instructions

Be explicit about when and how to handoff:

```go
coordinatorAgent.SetSystemInstructions(`For math problems, use handoff_to_Math_Agent.
For weather questions, use handoff_to_Weather_Agent.
Always provide clear instructions in the handoff input parameter.`)
```

### 3. Agent Names Must Match

The agent name in `WithHandoffs()` must match what you tell the LLM:

```go
// ✅ CORRECT
agent.WithHandoffs(mathAgent) // mathAgent.Name = "Math Agent"
// In instructions: "handoff to Math Agent"

// ❌ WRONG
agent.WithHandoffs(mathAgent) // mathAgent.Name = "Math Agent"
// In instructions: "handoff to MathBot" // Name doesn't match!
```

### 4. Limit Handoff Depth

Set `MaxTurns` to prevent infinite handoff loops:

```go
result, err := r.Run(context.Background(), agent, &runner.RunOptions{
	Input:    "Hello",
	MaxTurns: 20, // Prevents infinite loops
})
```

### 5. Use Tools Before Handoffs

If an agent can handle a task directly, use tools instead of handoffs:

```go
// ✅ GOOD: Use tool for simple tasks
agent.WithTools(calculatorTool)

// ✅ GOOD: Use handoff for complex, specialized tasks
agent.WithHandoffs(specializedAgent)
```

## Troubleshooting

### Problem: Handoff Not Working

**Symptoms**: Agent doesn't handoff even when instructed.

**Solutions**:
1. Check agent names match exactly:
   ```go
   // Agent name
   mathAgent := agent.NewAgent("Math Agent")
   
   // Must match in instructions
   "handoff to Math Agent" // ✅ Correct
   "handoff to math agent" // ❌ Wrong (case sensitive)
   ```

2. Verify handoffs are configured:
   ```go
   coordinatorAgent.WithHandoffs(mathAgent) // Must be called
   ```

3. Check system instructions mention handoffs:
   ```go
   coordinatorAgent.SetSystemInstructions("... handoff to Math Agent ...")
   ```

### Problem: Multiple Handoffs Detected

**Symptom**: Error message "Multiple handoffs detected, ignoring this one."

**Cause**: The LLM tried to handoff to multiple agents at once.

**Solution**: The framework automatically handles this by accepting only the first handoff. Make your instructions clearer:

```go
coordinatorAgent.SetSystemInstructions(`IMPORTANT: Only handoff to ONE agent at a time.
Choose the most appropriate agent for the task.`)
```

### Problem: Agent Doesn't Return

**Symptom**: Worker agent doesn't return to orchestrator.

**Solution**: 
1. Set up bidirectional handoffs:
   ```go
   orchestratorAgent.WithBidirectionalHandoffs(workerAgent)
   ```

2. Instruct worker to return:
   ```go
   workerAgent.SetSystemInstructions(`When you complete your task, 
   return to "Orchestrator" using handoff with your results.`)
   ```

### Problem: Context Not Shared

**Symptom**: Context data is lost during handoffs.

**Solution**: Use `RunContext`:

```go
result, err := r.Run(context.Background(), agent, &runner.RunOptions{
	Input:   "Hello",
	Context: sharedContext, // This is shared automatically
})
```

## Summary

- ✅ Use `WithHandoffs()` to configure which agents can be handed off to
- ✅ Agent names in instructions must match exactly
- ✅ Use `WithBidirectionalHandoffs()` for agents that need to return
- ✅ Set clear system instructions for when to handoff
- ✅ Context is automatically shared across handoffs
- ✅ Set `MaxTurns` to prevent infinite loops
- ✅ Use input filters to control what history is passed

For more examples, see:
- [examples/multi_agent_example](../examples/multi_agent_example)
- [examples/bidirectional_flow_example](../examples/bidirectional_flow_example)
- [examples/typescript_code_review_example](../examples/typescript_code_review_example)

