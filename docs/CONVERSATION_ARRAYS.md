# Conversation Arrays Guide

## Overview

Conversation arrays allow you to provide a complete conversation history to an agent, enabling it to continue from a previous context. This is essential for building chat applications, maintaining context across sessions, and implementing conversation persistence.

## Table of Contents

1. [Basic Usage](#basic-usage)
2. [Message Types](#message-types)
3. [Tool Calls and Results](#tool-calls-and-results)
4. [Complete Examples](#complete-examples)
5. [Best Practices](#best-practices)
6. [Common Pitfalls](#common-pitfalls)

## Basic Usage

### Simple Conversation

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

	agent := agent.NewAgent("Chat Assistant")
	agent.SetModelProvider(provider)
	agent.WithModel("gpt-4o-mini")
	agent.SetSystemInstructions("You are a helpful assistant.")

	r := runner.NewRunner().WithDefaultProvider(provider)

	// Simple conversation array
	conversation := []interface{}{
		map[string]interface{}{
			"type":    "message",
			"role":    "user",
			"content": "Hello!",
		},
		map[string]interface{}{
			"type":    "message",
			"role":    "assistant",
			"content": "Hi there! How can I help you?",
		},
		map[string]interface{}{
			"type":    "message",
			"role":    "user",
			"content": "What's the weather like?",
		},
	}

	result, err := r.Run(context.Background(), agent, &runner.RunOptions{
		Input: conversation,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.FinalOutput)
}
```

## Message Types

### System Messages

System messages set the behavior and context for the agent:

```go
conversation := []interface{}{
	map[string]interface{}{
		"type":    "message",
		"role":    "system",
		"content": "You are a helpful assistant that provides accurate information.",
	},
	// ... other messages
}
```

### User Messages

User messages represent input from the user:

```go
map[string]interface{}{
	"type":    "message",
	"role":    "user",
	"content": "Your question or input here",
}
```

### Assistant Messages

Assistant messages represent the agent's responses:

```go
map[string]interface{}{
	"type":    "message",
	"role":    "assistant",
	"content": "The assistant's response",
}
```

## Tool Calls and Results

### Critical: Tool Call Order

**IMPORTANT**: When including tool calls in conversation history, you MUST follow this order:

1. **Assistant message with `tool_calls`** (REQUIRED)
2. **Tool result** (response to the tool call)
3. **Assistant final response** (optional, after tool execution)

### Assistant Message with Tool Calls

```go
map[string]interface{}{
	"type":    "message",
	"role":    "assistant",
	"content": "", // Can be empty or contain text
	"tool_calls": []map[string]interface{}{
		{
			"id":   "call_unique_id_123",
			"type": "function",
			"function": map[string]interface{}{
				"name":      "tool_name",
				"arguments": `{"param1": "value1"}`, // JSON string
			},
		},
	},
}
```

### Tool Result

```go
map[string]interface{}{
	"type": "tool_result",
	"tool_call": map[string]interface{}{
		"id":   "call_unique_id_123", // MUST match the tool_call ID above
		"name": "tool_name",
	},
	"tool_result": map[string]interface{}{
		"content": "Tool execution result", // Can be string, number, object, etc.
	},
}
```

### Complete Tool Call Example

```go
conversation := []interface{}{
	// User asks a question
	map[string]interface{}{
		"type":    "message",
		"role":    "user",
		"content": "What time is it?",
	},
	
	// Assistant decides to call a tool
	map[string]interface{}{
		"type":    "message",
		"role":    "assistant",
		"content": "",
		"tool_calls": []map[string]interface{}{
			{
				"id":   "call_time_001",
				"type": "function",
				"function": map[string]interface{}{
					"name":      "get_current_time",
					"arguments": `{"format": "kitchen"}`,
				},
			},
		},
	},
	
	// Tool result (response from tool execution)
	map[string]interface{}{
		"type": "tool_result",
		"tool_call": map[string]interface{}{
			"id":   "call_time_001", // Must match ID above
			"name": "get_current_time",
		},
		"tool_result": map[string]interface{}{
			"content": "3:45PM",
		},
	},
	
	// Assistant's final response after using tool
	map[string]interface{}{
		"type":    "message",
		"role":    "assistant",
		"content": "The current time is 3:45PM.",
	},
	
	// User follow-up
	map[string]interface{}{
		"type":    "message",
		"role":    "user",
		"content": "What about in RFC3339 format?",
	},
}
```

## Complete Examples

### Example 1: Multi-Turn Conversation with Tools

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
	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

func main() {
	provider := openai.NewProvider("your-api-key")
	provider.SetDefaultModel("gpt-4o-mini")

	// Create a time tool
	timeTool := tool.NewFunctionTool(
		"get_current_time",
		"Get the current time in a specified format",
		func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			format := time.RFC3339
			if formatParam, ok := params["format"].(string); ok && formatParam != "" {
				switch formatParam {
				case "kitchen":
					format = time.Kitchen
				case "rfc3339":
					format = time.RFC3339
				}
			}
			return time.Now().Format(format), nil
		},
	).WithSchema(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"format": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"rfc3339", "kitchen"},
				"description": "Time format",
			},
		},
	})

	agent := agent.NewAgent("Time Assistant")
	agent.SetModelProvider(provider)
	agent.WithModel("gpt-4o-mini")
	agent.SetSystemInstructions("You help users with time-related questions.")
	agent.WithTools(timeTool)

	r := runner.NewRunner().WithDefaultProvider(provider)

	// Conversation with tool usage history
	conversation := []interface{}{
		map[string]interface{}{
			"type":    "message",
			"role":    "user",
			"content": "What time is it?",
		},
		map[string]interface{}{
			"type":    "message",
			"role":    "assistant",
			"content": "",
			"tool_calls": []map[string]interface{}{
				{
					"id":   "call_001",
					"type": "function",
					"function": map[string]interface{}{
						"name":      "get_current_time",
						"arguments": `{"format": "kitchen"}`,
					},
				},
			},
		},
		map[string]interface{}{
			"type": "tool_result",
			"tool_call": map[string]interface{}{
				"id":   "call_001",
				"name": "get_current_time",
			},
			"tool_result": map[string]interface{}{
				"content": "3:45PM",
			},
		},
		map[string]interface{}{
			"type":    "message",
			"role":    "assistant",
			"content": "The current time is 3:45PM.",
		},
		map[string]interface{}{
			"type":    "message",
			"role":    "user",
			"content": "Can you tell me the time in RFC3339 format?",
		},
	}

	result, err := r.Run(context.Background(), agent, &runner.RunOptions{
		Input: conversation,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", result.FinalOutput)
}
```

### Example 2: Conversation with Multiple Tool Calls

```go
conversation := []interface{}{
	map[string]interface{}{
		"type":    "message",
		"role":    "user",
		"content": "Get the weather and time",
	},
	map[string]interface{}{
		"type":    "message",
		"role":    "assistant",
		"content": "",
		"tool_calls": []map[string]interface{}{
			{
				"id":   "call_weather_001",
				"type": "function",
				"function": map[string]interface{}{
					"name":      "get_weather",
					"arguments": `{"city": "Paris"}`,
				},
			},
			{
				"id":   "call_time_001",
				"type": "function",
				"function": map[string]interface{}{
					"name":      "get_current_time",
					"arguments": `{"format": "rfc3339"}`,
				},
			},
		},
	},
	// Tool results must match the order of tool_calls
	map[string]interface{}{
		"type": "tool_result",
		"tool_call": map[string]interface{}{
			"id":   "call_weather_001",
			"name": "get_weather",
		},
		"tool_result": map[string]interface{}{
			"content": "Sunny, 72°F",
		},
	},
	map[string]interface{}{
		"type": "tool_result",
		"tool_call": map[string]interface{}{
			"id":   "call_time_001",
			"name": "get_current_time",
		},
		"tool_result": map[string]interface{}{
			"content": "2024-01-15T14:30:00Z",
		},
	},
	map[string]interface{}{
		"type":    "message",
		"role":    "assistant",
		"content": "The weather in Paris is Sunny, 72°F. The current time is 2024-01-15T14:30:00Z.",
	},
}
```

### Example 3: Loading Conversation from Database

```go
// Load conversation from database
func loadConversationFromDB(sessionID string) ([]interface{}, error) {
	// This is a mock - replace with your actual database query
	conversation := []interface{}{
		map[string]interface{}{
			"type":    "message",
			"role":    "system",
			"content": "You are a helpful assistant.",
		},
		// ... load messages from database
	}
	return conversation, nil
}

// Continue conversation
func continueConversation(sessionID string) {
	conversation, err := loadConversationFromDB(sessionID)
	if err != nil {
		log.Fatal(err)
	}

	// Add new user message
	conversation = append(conversation, map[string]interface{}{
		"type":    "message",
		"role":    "user",
		"content": "Continue our conversation",
	})

	result, err := r.Run(context.Background(), agent, &runner.RunOptions{
		Input: conversation,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Save new messages to database
	saveToDatabase(sessionID, result.NewItems)
}
```

## Best Practices

### 1. Always Include System Message

Include a system message at the beginning to set agent behavior:

```go
conversation := []interface{}{
	map[string]interface{}{
		"type":    "message",
		"role":    "system",
		"content": "You are a helpful assistant.",
	},
	// ... rest of conversation
}
```

### 2. Match Tool Call IDs

**CRITICAL**: Tool call IDs must match exactly between `tool_calls` and `tool_result`:

```go
// ✅ CORRECT
"tool_calls": []map[string]interface{}{
	{"id": "call_123", ...},
}
// Later...
"tool_call": map[string]interface{}{
	"id": "call_123", // ✅ Matches
}

// ❌ WRONG
"tool_calls": []map[string]interface{}{
	{"id": "call_123", ...},
}
// Later...
"tool_call": map[string]interface{}{
	"id": "call_456", // ❌ Doesn't match - will cause error
}
```

### 3. Maintain Correct Order

The order of messages is critical:

1. System message (optional, but recommended)
2. User message
3. Assistant message (with or without tool_calls)
4. Tool result (if assistant had tool_calls)
5. Assistant final response (after tool execution)
6. Repeat as needed

### 4. Tool Result Content Types

Tool result content can be various types:

```go
// String
"content": "Simple text result"

// Number
"content": 42

// Object (will be converted to JSON string)
"content": map[string]interface{}{
	"temperature": 72,
	"condition": "sunny",
}
```

### 5. Limit Conversation Length

For performance and cost reasons, consider limiting conversation history:

```go
func truncateConversation(conversation []interface{}, maxMessages int) []interface{} {
	if len(conversation) <= maxMessages {
		return conversation
	}
	// Keep system message and recent messages
	result := []interface{}{conversation[0]} // System message
	start := len(conversation) - maxMessages + 1
	return append(result, conversation[start:]...)
}
```

## Common Pitfalls

### Error: "messages with role 'tool' must be a response to a preceding message with 'tool_calls'"

**Problem**: You included a `tool_result` without a preceding assistant message with `tool_calls`.

**Solution**: Always include the assistant message with `tool_calls` before the `tool_result`:

```go
// ❌ WRONG - Missing assistant message with tool_calls
conversation := []interface{}{
	map[string]interface{}{
		"type":    "message",
		"role":    "user",
		"content": "What time is it?",
	},
	map[string]interface{}{
		"type": "tool_result", // ❌ No preceding tool_calls
		// ...
	},
}

// ✅ CORRECT
conversation := []interface{}{
	map[string]interface{}{
		"type":    "message",
		"role":    "user",
		"content": "What time is it?",
	},
	map[string]interface{}{
		"type":    "message",
		"role":    "assistant",
		"content": "",
		"tool_calls": []map[string]interface{}{...}, // ✅ Required
	},
	map[string]interface{}{
		"type": "tool_result", // ✅ Now it's valid
		// ...
	},
}
```

### Error: Tool Call ID Mismatch

**Problem**: The `tool_call.id` in the tool result doesn't match the `id` in the assistant's `tool_calls`.

**Solution**: Ensure IDs match exactly:

```go
// ✅ CORRECT
"tool_calls": []map[string]interface{}{
	{"id": "call_abc123", ...},
}
// Later...
"tool_call": map[string]interface{}{
	"id": "call_abc123", // ✅ Exact match
}
```

### Error: Invalid Message Format

**Problem**: Missing required fields like `type` or `role`.

**Solution**: Always include all required fields:

```go
// ❌ WRONG
map[string]interface{}{
	"role":    "user",
	"content": "Hello",
	// Missing "type"
}

// ✅ CORRECT
map[string]interface{}{
	"type":    "message", // Required
	"role":    "user",    // Required
	"content": "Hello",   // Required
}
```

## Advanced: Converting from Other Formats

### From OpenAI Format

If you have messages in OpenAI's format:

```go
// OpenAI format
openaiMessages := []map[string]interface{}{
	{"role": "user", "content": "Hello"},
	{"role": "assistant", "content": "Hi!"},
}

// Convert to go-agentkit format
conversation := make([]interface{}, len(openaiMessages))
for i, msg := range openaiMessages {
	conversation[i] = map[string]interface{}{
		"type":    "message",
		"role":    msg["role"],
		"content": msg["content"],
	}
}
```

### From Database Records

```go
type MessageRecord struct {
	Role    string
	Content string
	ToolCalls string // JSON string
}

func convertDBMessages(records []MessageRecord) []interface{} {
	conversation := make([]interface{}, 0, len(records))
	for _, record := range records {
		msg := map[string]interface{}{
			"type":    "message",
			"role":    record.Role,
			"content": record.Content,
		}
		
		// Add tool_calls if present
		if record.ToolCalls != "" {
			var toolCalls []map[string]interface{}
			json.Unmarshal([]byte(record.ToolCalls), &toolCalls)
			msg["tool_calls"] = toolCalls
		}
		
		conversation = append(conversation, msg)
	}
	return conversation
}
```

## Summary

- ✅ Always use `"type": "message"` for regular messages
- ✅ Include assistant message with `tool_calls` before `tool_result`
- ✅ Match tool call IDs exactly between `tool_calls` and `tool_result`
- ✅ Maintain correct message order
- ✅ Include system message for consistent behavior
- ✅ Limit conversation length for performance
- ✅ Validate message format before sending

For more examples, see the [examples directory](../examples).

