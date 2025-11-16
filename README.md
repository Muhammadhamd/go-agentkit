<p align="center">
  <img src="./go-agentkit-header.gif" alt="Agent SDK Go">
</p>

<div align="center">
  <p><strong>Build, deploy, and scale AI agents with ease</strong></p>
  
  <a href="https://go-agent.org"><img src="https://img.shields.io/badge/website-go--agent.org-blue?style=for-the-badge" alt="Website" /></a>
  <a href="https://go-agent.org/#waitlist"><img src="https://img.shields.io/badge/Cloud_Waitlist-Sign_Up-4285F4?style=for-the-badge" alt="Cloud Waitlist" /></a>
  
</div>

<p align="center">
  Agent SDK Go is an open-source framework for building powerful AI agents with Go that supports multiple LLM providers, function calling, agent handoffs, and more.
</p>

<p align="center">
    <a href="https://github.com/muhammadhamd/go-agentkit/actions/workflows/code-quality.yml"><img src="https://github.com/muhammadhamd/go-agentkit/actions/workflows/code-quality.yml/badge.svg" alt="Code Quality"></a>
    <a href="https://goreportcard.com/report/github.com/muhammadhamd/go-agentkit/"><img src="https://goreportcard.com/badge/github.com/muhammadhamd/go-agentkit/" alt="Go Report Card"></a>
    <a href="https://github.com/muhammadhamd/go-agentkit/blob/main/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/muhammadhamd/go-agentkit" alt="Go Version"></a>
    <a href="https://pkg.go.dev/github.com/muhammadhamd/go-agentkit/"><img src="https://pkg.go.dev/badge/github.com/muhammadhamd/go-agentkit/.svg" alt="PkgGoDev"></a><br>
    <a href="https://github.com/muhammadhamd/go-agentkit/actions/workflows/codeql-analysis.yml"><img src="https://github.com/muhammadhamd/go-agentkit/actions/workflows/codeql-analysis.yml/badge.svg" alt="CodeQL"></a>
    <a href="https://github.com/muhammadhamd/go-agentkit/blob/main/LICENSE"><img src="https://img.shields.io/github/license/muhammadhamd/go-agentkit" alt="License"></a>
    <a href="https://github.com/muhammadhamd/go-agentkit/stargazers"><img src="https://img.shields.io/github/stars/muhammadhamd/go-agentkit" alt="Stars"></a>
    <a href="https://github.com/muhammadhamd/go-agentkit/graphs/contributors"><img src="https://img.shields.io/github/contributors/muhammadhamd/go-agentkit" alt="Contributors"></a>
    <a href="https://github.com/muhammadhamd/go-agentkit/commits/main"><img src="https://img.shields.io/github/last-commit/muhammadhamd/go-agentkit" alt="Last Commit"></a>
</p>

<p align="center">
  <a href="https://go-agent.org/#waitlist">☁️ Cloud Waitlist</a> •
  <a href="https://github.com/muhammadhamd/go-agentkit/blob/main/LICENSE">📜 License</a>
</p>

<p align="center">
  Created by <a href="https://linkedin.com/in/MuhammadHamd">Muhammad Hamd</a> • 
  <a href="https://github.com/muhammadhamd">GitHub</a> • 
  <a href="https://linkedin.com/in/MuhammadHamd">LinkedIn</a>
</p>

<p align="center">
  Inspired by <a href="https://platform.openai.com/docs/assistants/overview">OpenAI's Assistants API</a> and <a href="https://github.com/openai/openai-agents-python">OpenAI's Python Agent SDK</a>.
</p>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Provider Setup](#-provider-setup)
- [Key Components](#-key-components)
  - [Agent](#agent)
  - [Runner](#runner)
  - [Tools](#tools)
  - [Model Providers](#model-providers)
- [Advanced Features](#-advanced-features)
  - [Agentic Loop & Context Sharing](#agentic-loop--context-sharing)
  - [Multi-Agent Workflows](#multi-agent-workflows)
  - [Bidirectional Agent Flow](#bidirectional-agent-flow)
  - [Tracing](#tracing)
  - [Structured Output](#structured-output)
  - [Streaming](#streaming)
  - [OpenAI Tool Definitions](#openai-tool-definitions)
  - [Workflow State Management](#workflow-state-management)
- [Examples](#-examples)
- [Cloud Support](#-cloud-support)
- [Development](#-development)
- [Contributing](#-contributing)
- [License](#-license)
- [Acknowledgements](#-acknowledgements)

---

## 🔍 What is This?

**Go Agent SDK** is a powerful, easy-to-use framework for building AI agents in Go. Think of it as your toolkit for creating intelligent applications that can:

- **Talk to AI models** - Use OpenAI, Anthropic Claude, or local models
- **Call functions** - Let AI agents use your Go functions as tools
- **Work together** - Create teams of specialized agents that hand off tasks to each other
- **Share context** - Agents can share data and remember information across conversations
- **Stream responses** - Get real-time updates as agents work
- **Debug easily** - Built-in tracing helps you see exactly what your agents are doing

### Why Use Go Agent SDK?

- ✅ **Simple** - Easy to learn, powerful to use
- ✅ **Production-Ready** - Built with best practices and error handling
- ✅ **Flexible** - Works with multiple AI providers and models
- ✅ **Compatible** - Matches OpenAI's Python and TypeScript SDKs behavior
- ✅ **Well-Documented** - Comprehensive examples and guides

**Visit [go-agent.org](https://go-agent.org) for more documentation, examples, and cloud service waitlist.**

## 🌟 Complete Feature List

The Go Agent SDK is a powerful, production-ready framework for building AI agents. Here are all the features available:

### 🤖 Core Agent Features
- ✅ **Multiple LLM Providers** - OpenAI (GPT-3.5, GPT-4, GPT-4o), Anthropic Claude (Haiku, Sonnet, Opus), LM Studio (local models), and custom base URLs (DeepSeek, etc.)
- ✅ **Agent Configuration** - System instructions, model settings (temperature, max tokens), output types, and tool use behavior
- ✅ **Agentic Loop** - Turn-based execution with automatic state management, tool calling, and handoff handling
- ✅ **Max Turns Control** - Prevent infinite loops by limiting the number of agent turns
- ✅ **Backward Compatibility** - All old agent creation methods still work (method chaining, direct field access)

### 🛠️ Tool Features
- ✅ **Function Tools** - Convert any Go function into a tool that agents can call
- ✅ **OpenAI-Compatible Tools** - Use OpenAI tool definitions directly
- ✅ **Tool Schema Generation** - Automatic JSON schema generation for tool parameters
- ✅ **Tool Parameter Validation** - Automatic validation of tool parameters
- ✅ **Tool Error Handling** - Custom error handling for tool failures
- ✅ **Tool Approval** - Human-in-the-loop approval for sensitive tool calls
- ✅ **Tool Use Behavior** - Control how agents handle tool outputs (run_llm_again, stop_on_first_tool, custom)

### 👥 Multi-Agent Features
- ✅ **Agent Handoffs** - Transfer control between specialized agents
- ✅ **Bidirectional Flow** - Agents can delegate tasks and receive results back
- ✅ **Task Delegation** - Track and manage delegated tasks with unique IDs
- ✅ **Input Filtering** - Filter conversation history during handoffs
- ✅ **Context Sharing** - Share custom data, usage stats, and tool approvals across agents

### 📊 Data & Output Features
- ✅ **Structured Output** - Parse LLM responses into Go structs with automatic validation
- ✅ **JSON Schema Validation** - Validate structured outputs against schemas
- ✅ **Context Sharing** - RunContext for sharing data across agents and turns
- ✅ **Usage Tracking** - Automatic token usage tracking (input, output, total)

### 🔄 Streaming & Real-time
- ✅ **Streaming Responses** - Get real-time streaming responses from agents
- ✅ **Stream Events** - Content, tool calls, handoffs, and completion events
- ✅ **AsyncIterable Pattern** - Easy-to-use streaming interface

### 🔍 Observability & Debugging
- ✅ **OpenAI Backend Tracing** - Automatic tracing sent to OpenAI dashboard (matches Python/TypeScript SDKs)
- ✅ **Environment Variable Control** - Disable tracing via `OPENAI_AGENTS_DISABLE_TRACING`
- ✅ **Per-Run Tracing Config** - Enable/disable tracing per run
- ✅ **No Local Files** - No trace files created locally (matches Python/TypeScript behavior)

### 🔌 Integration Features
- ✅ **MCP Support** - Model Context Protocol for local and hosted MCP servers
- ✅ **Custom Base URLs** - Connect to any OpenAI-compatible API (DeepSeek, etc.)
- ✅ **Lifecycle Hooks** - Agent hooks and run hooks for custom logic
- ✅ **Workflow State Management** - State persistence, retry, recovery, and checkpointing

### ⚙️ Advanced Features
- ✅ **Error Handling** - Comprehensive error types (tool errors, model errors, max turns, guardrails)
- ✅ **Guardrails** - Input and output validation (framework ready)
- ✅ **Session Management** - Framework ready for conversation history persistence
- ✅ **Consecutive Tool Call Tracking** - Prevent infinite tool call loops

## 📦 Installation

There are several ways to add this module to your project:

### Option 1: Using `go get` (Recommended)

```bash
go get github.com/muhammadhamd/go-agentkit/
```

### Option 2: Add to your imports and use `go mod tidy`

1. Add imports to your Go files:
   ```go
   import (
       "github.com/muhammadhamd/go-agentkit/pkg/agent"
       "github.com/muhammadhamd/go-agentkit/pkg/model/providers/lmstudio"
       "github.com/muhammadhamd/go-agentkit/pkg/runner"
       "github.com/muhammadhamd/go-agentkit/pkg/tool"
       // Import other packages as needed
   )
   ```

2. Run `go mod tidy` to automatically fetch dependencies:
   ```bash
   go mod tidy
   ```

### Option 3: Manually edit your `go.mod` file

Add the following line to your `go.mod` file:
```
require github.com/muhammadhamd/go-agentkit/ latest
```

Then run:
```bash
go mod tidy
```

### New Project Setup

If you're starting a new project:

1. Create and navigate to your project directory:
   ```bash
   mkdir my-agent-project
   cd my-agent-project
   ```

2. Initialize a new Go module:
   ```bash
   go mod init github.com/yourusername/my-agent-project
   ```

3. Install the Agent SDK:
   ```bash
   go get github.com/muhammadhamd/go-agentkit/
   ```

### Troubleshooting

- If you encounter version conflicts, you can specify a version:
  ```bash
  go get github.com/muhammadhamd/go-agentkit/@v0.1.0  # Replace with desired version
  ```

- For private repositories or local development, consider using Go workspaces or replace directives in your go.mod file.

> **Note:** Requires Go 1.23 or later.

## 🚀 Quick Start

Let's build your first AI agent in 5 minutes! This example shows how to create an agent that can call a function.

### Step 1: Install the SDK

```bash
go get github.com/muhammadhamd/go-agentkit/
```

### Step 2: Set Your API Key

```bash
# For OpenAI
export OPENAI_API_KEY=your-api-key-here

# Or for Anthropic
export ANTHROPIC_API_KEY=your-api-key-here
```

### Step 3: Create Your First Agent

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/muhammadhamd/go-agentkit/pkg/agent"
    "github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
    "github.com/muhammadhamd/go-agentkit/pkg/runner"
    "github.com/muhammadhamd/go-agentkit/pkg/tool"
)

func main() {
    // Step 1: Get your API key from environment
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("Please set OPENAI_API_KEY environment variable")
    }

    // Step 2: Create a provider (this connects to OpenAI)
    provider := openai.NewProvider(apiKey)
    provider.SetDefaultModel("gpt-4o-mini") // Use a cheaper model for testing

    // Step 3: Create a tool (a function the agent can call)
    getWeather := tool.NewFunctionTool(
        "get_weather",
        "Get the current weather for a city",
        func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
            city := params["city"].(string)
            // In a real app, you'd call a weather API here
            return fmt.Sprintf("The weather in %s is 72°F and sunny.", city), nil
        },
    ).WithSchema(map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "city": map[string]interface{}{
                "type":        "string",
                "description": "The city name to get weather for",
            },
        },
        "required": []string{"city"},
    })

    // Step 4: Create an agent
    assistant := agent.NewAgent("Weather Assistant", "You are a helpful weather assistant.")
    assistant.SetModelProvider(provider)
    assistant.WithModel("gpt-4o-mini")
    assistant.WithTools(getWeather)

    // Step 5: Create a runner and run the agent
    r := runner.NewRunner()
    r.WithDefaultProvider(provider)

    result, err := r.RunSync(assistant, &runner.RunOptions{
        Input: "What's the weather in Tokyo?",
    })
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    // Step 6: Print the result
    fmt.Println("Agent Response:", result.FinalOutput)
}
```

### What Just Happened?

1. **Provider** - Connects to OpenAI's API
2. **Tool** - A function the agent can call (get_weather)
3. **Agent** - The AI assistant with instructions and tools
4. **Runner** - Executes the agent and handles the conversation
5. **Result** - Contains the agent's final response

The agent automatically:
- ✅ Understood the user's question
- ✅ Decided to call the `get_weather` tool
- ✅ Called the tool with `city: "Tokyo"`
- ✅ Used the tool's result to answer the question

### Try Different Providers

**Anthropic Claude:**
```go
import "github.com/muhammadhamd/go-agentkit/pkg/model/providers/anthropic"

provider := anthropic.NewProvider(os.Getenv("ANTHROPIC_API_KEY"))
provider.SetDefaultModel("claude-3-haiku-20240307")
```

**LM Studio (Local Model):**
```go
import "github.com/muhammadhamd/go-agentkit/pkg/model/providers/lmstudio"

provider := lmstudio.NewProvider()
provider.SetBaseURL("http://127.0.0.1:1234/v1")
provider.SetDefaultModel("gemma-3-4b-it")
```

**DeepSeek (Custom Base URL):**
```go
deepseekProvider := openai.NewProvider(os.Getenv("DEEPSEEK_API_KEY"))
deepseekProvider.SetBaseURL("https://api.deepseek.com/v1")
deepseekProvider.SetDefaultModel("deepseek-chat")
```

See the [DeepSeek example](./examples/deepseek_example) for a complete example.

## 💡 Common Examples

### Example 1: Multi-Agent Workflow

Create specialized agents that work together:

```go
// Create a support agent
supportAgent := agent.NewAgent("Support Agent", "You handle customer support.")
supportAgent.SetModelProvider(provider)
supportAgent.WithTools(getUserInfoTool, checkOrderTool)

// Create a billing agent
billingAgent := agent.NewAgent("Billing Agent", "You handle billing and refunds.")
billingAgent.SetModelProvider(provider)
billingAgent.WithTools(calculateRefundTool, processRefundTool)

// Connect them: support agent can hand off to billing agent
supportAgent.WithHandoffs(billingAgent)

// Run with shared context
type MyContext struct {
    UserID  string
    OrderID string
}

ctx := &MyContext{}
result, err := r.RunSync(supportAgent, &runner.RunOptions{
    Input:   "I need a refund for order_123",
    Context: ctx, // Shared across all agents
})
```

### Example 2: Context Sharing

Share data between agents and tools:

```go
// Define your context
type AppContext struct {
    UserID    string
    SessionID string
    Metadata  map[string]interface{}
}

// Create context
appCtx := &AppContext{
    UserID:    "user_123",
    SessionID: "session_456",
    Metadata:  make(map[string]interface{}),
}

// Access context in tools
myTool := tool.NewFunctionTool(
    "update_user",
    "Update user information",
    func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
        // Get RunContext from context
        runCtxVal := ctx.Value("run_context")
        if runCtx, ok := runCtxVal.(*runner.RunContext); ok {
            if appCtx, ok := runCtx.Context.(*AppContext); ok {
                // Modify shared context
                appCtx.Metadata["last_updated"] = time.Now()
            }
        }
        return "Updated", nil
    },
)

// Run with context
result, err := r.RunSync(agent, &runner.RunOptions{
    Input:   "Update my profile",
    Context: appCtx, // Shared across agents and tools
})

// Access updated context from result
if result.RunContext != nil {
    // Context is automatically updated
}
```

### Example 3: Streaming Responses

Get real-time updates:

```go
streamResult, err := r.RunStreaming(ctx, agent, &runner.RunOptions{
    Input: "Tell me a story",
})

for event := range streamResult.Stream {
    switch event.Type {
    case model.StreamEventTypeContent:
        fmt.Print(event.Content) // Print as it streams
    case model.StreamEventTypeToolCall:
        fmt.Printf("\n🔧 Calling tool: %s\n", event.ToolCall.Name)
    case model.StreamEventTypeDone:
        fmt.Println("\n✅ Done!")
    }
}
```

### Example 4: Structured Output

Get responses as Go structs:

```go
type WeatherReport struct {
    City        string  `json:"city"`
    Temperature float64 `json:"temperature"`
    Condition   string  `json:"condition"`
}

agent := agent.NewAgent("Weather Agent", "You provide weather reports.")
agent.SetOutputType(WeatherReport{}) // Specify output type

result, err := r.RunSync(agent, &runner.RunOptions{
    Input: "What's the weather in Paris?",
})

// Parse structured output
var report WeatherReport
if err := json.Unmarshal([]byte(result.FinalOutput.(string)), &report); err == nil {
    fmt.Printf("City: %s, Temp: %.1f°F, Condition: %s\n", 
        report.City, report.Temperature, report.Condition)
}
```

See more examples in the [examples directory](./examples).

## 🖥️ Provider Setup

### OpenAI Setup

To use the OpenAI provider:

1. **Get an API Key**
   - Sign up at [OpenAI](https://platform.openai.com/)
   - Create an API key in your account settings

2. **Configure the Provider**
   ```go
   provider := openai.NewProvider()
   provider.SetAPIKey("your-openai-api-key")
   provider.SetDefaultModel("gpt-3.5-turbo")  // or any other OpenAI model
   ```

### Anthropic Setup

<details>
<summary>Click to expand setup instructions</summary>

To use the Anthropic provider:

1. **Get an API Key**
   - Sign up at [Anthropic Console](https://console.anthropic.com/)
   - Create an API key in your account settings

2. **Configure the Provider**
   ```go
   provider := anthropic.NewProvider("your-anthropic-api-key")
   provider.SetDefaultModel("claude-3-haiku-20240307")  // or claude-3-sonnet/opus
   
   // Optional rate limiting configuration
   provider.WithRateLimit(40, 80000) // 40 requests/min, 80,000 tokens/min
   
   // Optional retry configuration
   provider.WithRetryConfig(3, 2*time.Second) // 3 retries with exponential backoff
   ```

</details>

### LM Studio Setup

<details>
<summary>Click to expand setup instructions</summary>

To use the LM Studio provider:

1. **Install LM Studio**
   - Download from [lmstudio.ai](https://lmstudio.ai/)
   - Install and run the application

2. **Load a Model**
   - Download a model in LM Studio (Like Gemma-3-4B-It, Llama3, or other compatible models)
   - Load the model

3. **Start the Server**
   - Go to the "Local Server" tab
   - Click "Start Server"
   - Note the server URL (default: http://127.0.0.1:1234)

4. **Configure the Provider**
   ```go
   provider := lmstudio.NewProvider()
   provider.SetBaseURL("http://127.0.0.1:1234/v1")
   provider.SetDefaultModel("gemma-3-4b-it") // Replace with your model
   ```

</details>

## 🧩 Key Components

### Agent

The Agent is the core component that encapsulates the LLM with instructions, tools, and other configuration.

```go
// Create a new agent
agent := agent.NewAgent("Assistant")
agent.SetSystemInstructions("You are a helpful assistant.")
agent.WithModel("gemma-3-4b-it")
agent.WithTools(tool1, tool2) // Add multiple tools at once
```

### Runner

The Runner executes agents, handling the agent loop, tool calls, and handoffs.

```go
// Create a runner
runner := runner.NewRunner()
runner.WithDefaultProvider(provider)

// Run the agent
result, err := runner.RunSync(agent, &runner.RunOptions{
    Input: "Hello, world!",
    MaxTurns: 10, // Optional: limit the number of turns
})
```

### Tools

Tools allow agents to perform actions using your Go functions.

```go
// Create a function tool
tool := tool.NewFunctionTool(
    "get_weather",
    "Get the weather for a city",
    func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
        city := params["city"].(string)
        return fmt.Sprintf("The weather in %s is sunny.", city), nil
    },
).WithSchema(map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "city": map[string]interface{}{
            "type": "string",
            "description": "The city to get weather for",
        },
    },
    "required": []string{"city"},
})
```

### Model Providers

Model providers allow you to use different LLM providers.

```go
// Create a provider for OpenAI
openaiProvider := openai.NewProvider("your-openai-api-key")
openaiProvider.SetDefaultModel("gpt-4")

// Create a provider for Anthropic Claude
anthropicProvider := anthropic.NewProvider("your-anthropic-api-key")
anthropicProvider.SetDefaultModel("claude-3-haiku-20240307")

// Create a provider for LM Studio
lmStudioProvider := lmstudio.NewProvider()
lmStudioProvider.SetBaseURL("http://127.0.0.1:1234/v1")
lmStudioProvider.SetDefaultModel("gemma-3-4b-it")

// Set a provider as the default provider
runner := runner.NewRunner()
runner.WithDefaultProvider(openaiProvider) // or anthropicProvider or lmStudioProvider
```

## 🔧 Advanced Features

### Agentic Loop & Context Sharing

<details>
<summary>Understand how the agentic loop works and share context between agents</summary>

The agentic loop is the core execution engine that powers agent interactions. It manages turn-based execution, tool calls, handoffs, and context sharing across agents.

#### How the Agentic Loop Works

The agentic loop follows a turn-based execution model:

1. **Turn Execution**: Each turn involves:
   - Getting input (combines original input + generated items from previous turns)
   - Calling the model with the current conversation history
   - Processing the model response (messages, tool calls, handoffs)
   - Executing tools if needed
   - Determining the next step (continue, handoff, or final output)

2. **State Management**: The `RunState` tracks:
   - Current agent and turn number
   - Generated items (messages, tool calls, tool results)
   - Conversation history
   - Usage statistics

3. **Next Step Types**:
   - `NextStepRunAgain`: Continue the loop (after tool execution)
   - `NextStepHandoff`: Transfer control to another agent
   - `NextStepFinalOutput`: Agent has produced final output
   - `NextStepInterruption`: Pause for approvals or user input

#### Context Sharing with RunContext

`RunContext` allows you to share custom data, usage statistics, and tool approval states across agent turns and handoffs:

```go
// Define your custom context type
type MyContext struct {
    UserID    string
    OrderID   string
    Metadata  map[string]interface{}
}

// Create context instance
myContext := &MyContext{
    UserID:   "user_123",
    OrderID:  "order_456",
    Metadata: make(map[string]interface{}),
}

// Pass context to the runner
result, err := runner.Run(ctx, agent, &runner.RunOptions{
    Input:   "Process my order",
    Context: myContext, // Shared across all agents
    MaxTurns: 20,
})

// Access context from result
if result.RunContext != nil {
    if rc, ok := result.RunContext.(*runner.RunContext); ok {
        if myCtx, ok := rc.Context.(*MyContext); ok {
            fmt.Printf("User ID: %s\n", myCtx.UserID)
            fmt.Printf("Usage: %d tokens\n", rc.Usage.TotalTokens)
        }
    }
}
```

#### Accessing Context in Tools

Tools can access and modify the shared context:

```go
func myTool(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // Access RunContext from the context
    runCtxVal := ctx.Value("run_context")
    if runCtxVal != nil {
        if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil {
            if myCtx, ok := runCtx.Context.(*MyContext); ok {
                // Read from context
                userID := myCtx.UserID
                
                // Modify context (shared across all agents)
                myCtx.Metadata["last_tool_called"] = "myTool"
                myCtx.Metadata["called_at"] = time.Now()
            }
        }
    }
    
    // Tool implementation
    return result, nil
}
```

#### Complex Multi-Agent Flow Example

Here's a complete example demonstrating context sharing across multiple agents:

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

// Shared context across agents
type TestContext struct {
    UserID    string
    OrderID   string
    SessionID string
    Metadata  map[string]interface{}
}

// Tool that updates shared context
func getUserInfo(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    userID := params["user_id"].(string)
    
    // Access and update shared context
    runCtxVal := ctx.Value("run_context")
    if runCtxVal != nil {
        if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil {
            if testCtx, ok := runCtx.Context.(*TestContext); ok {
                testCtx.UserID = userID // Shared across all agents
            }
        }
    }
    
    return map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    }, nil
}

func main() {
    ctx := context.Background()
    provider := openai.NewProvider("your-api-key")
    
    // Create shared context
    testContext := &TestContext{
        SessionID: fmt.Sprintf("session_%d", time.Now().Unix()),
        Metadata:  make(map[string]interface{}),
    }
    
    // Create tools
    getUserInfoTool := tool.NewFunctionTool(
        "get_user_info",
        "Get user information by user_id",
        getUserInfo,
    )
    
    // Create specialized agents
    supportAgent := agent.NewAgent("support_agent")
    supportAgent.SetModelProvider(provider)
    supportAgent.SetSystemInstructions(
        "You are a support agent. Get user info first, then handoff to billing.",
    )
    supportAgent.WithTools(getUserInfoTool)
    
    billingAgent := agent.NewAgent("billing_agent")
    billingAgent.SetModelProvider(provider)
    billingAgent.SetSystemInstructions("You handle billing and refunds.")
    
    // Set up handoffs
    supportAgent.WithHandoffs(billingAgent)
    
    // Create runner
    r := runner.NewRunner().WithDefaultProvider(provider)
    
    // Run with shared context
    result, err := r.Run(ctx, supportAgent, &runner.RunOptions{
        Input:    "I need a refund for order_456, user_123",
        Context:  testContext, // Shared across all agents
        MaxTurns: 20,
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Access updated context
    if result.RunContext != nil {
        if rc, ok := result.RunContext.(*runner.RunContext); ok {
            if testCtx, ok := rc.Context.(*TestContext); ok {
                fmt.Printf("User ID in context: %s\n", testCtx.UserID)
                fmt.Printf("Usage: %d tokens\n", rc.Usage.TotalTokens)
            }
        }
    }
}
```

#### Key Features

- **Turn-based Execution**: Each turn processes one model interaction
- **State Persistence**: Conversation history and generated items are maintained
- **Context Sharing**: Custom context is shared across all agents via `RunContext`
- **Usage Tracking**: Token usage is automatically tracked and available in `RunContext`
- **Tool Context Access**: Tools can read and modify shared context
- **Handoff Context**: Context is preserved when agents hand off to each other

#### Best Practices

1. **Context Design**: Keep your context type simple and focused on shared data
2. **Thread Safety**: The `RunContext` is thread-safe, but your custom context should handle concurrent access if needed
3. **Context Size**: Keep context data small to avoid token overhead
4. **Tool Modifications**: Tools should modify context carefully to avoid race conditions
5. **Error Handling**: Always check for nil when accessing context values

See the complete example in [store.go](./store.go) (run with `go run -tags=store store.go`).

</details>

### Multi-Agent Workflows

<details>
<summary>Create specialized agents that collaborate on complex tasks</summary>

```go
// Create specialized agents
mathAgent := agent.NewAgent("Math Agent")
mathAgent.SetModelProvider(provider)
mathAgent.WithModel("gemma-3-4b-it")
mathAgent.SetSystemInstructions("You are a specialized math agent.")
mathAgent.WithTools(calculatorTool)

weatherAgent := agent.NewAgent("Weather Agent")
weatherAgent.SetModelProvider(provider)
weatherAgent.WithModel("gemma-3-4b-it")
weatherAgent.SetSystemInstructions("You provide weather information.")
weatherAgent.WithTools(weatherTool)

// Create a frontend agent that coordinates tasks
frontendAgent := agent.NewAgent("Frontend Agent")
frontendAgent.SetModelProvider(provider)
frontendAgent.WithModel("gemma-3-4b-it")
frontendAgent.SetSystemInstructions(`You coordinate requests by delegating to specialized agents.
For math calculations, delegate to the Math Agent.
For weather information, delegate to the Weather Agent.`)
frontendAgent.WithHandoffs(mathAgent, weatherAgent)

// Run the frontend agent
result, err := runner.RunSync(frontendAgent, &runner.RunOptions{
    Input: "What is 42 divided by 6 and what's the weather in Paris?",
    MaxTurns: 20,
})
```

See the complete example in [examples/multi_agent_example](./examples/multi_agent_example).

</details>

### Bidirectional Agent Flow

<details>
<summary>Create agents that can hand off tasks and receive results back</summary>

Bidirectional agent flow allows agents to delegate tasks to other agents and receive results back once the tasks are complete. This enables more complex workflows with proper task context management.

```go
// Create specialized agents
orchestratorAgent := agent.NewAgent("Orchestrator")
orchestratorAgent.SetModelProvider(provider)
orchestratorAgent.WithModel("gpt-4")
orchestratorAgent.SetSystemInstructions("You coordinate tasks and analyze results.")

workerAgent := agent.NewAgent("Worker")
workerAgent.SetModelProvider(provider)
workerAgent.WithModel("gpt-3.5-turbo")
workerAgent.SetSystemInstructions("You process data and return results.")
workerAgent.WithTools(processingTool)

// Set up bidirectional handoffs
orchestratorAgent.WithHandoffs(workerAgent)
workerAgent.WithHandoffs(orchestratorAgent)  // Allow worker to return to orchestrator

// Run the orchestrator agent
result, err := runner.RunSync(orchestratorAgent, &runner.RunOptions{
    Input: "Analyze this data: [complex data]",
    MaxTurns: 10,
})
```

Key components of bidirectional flow:
- `TaskID`: Unique identifier for tracking tasks across agents
- `ReturnToAgent`: Specifies which agent to return to after task completion
- `IsTaskComplete`: Flag indicating whether the task is complete

See the complete example in [examples/bidirectional_flow_example](./examples/bidirectional_flow_example).

</details>

### Tracing

<details>
<summary>Debug your agent workflows with tracing</summary>

Tracing is enabled by default. You can disable tracing in two ways:

1. **Environment variable** (matches Python/TypeScript SDKs):
   ```bash
   export OPENAI_AGENTS_DISABLE_TRACING=1
   # or
   export OPENAI_AGENTS_DISABLE_TRACING=true
   ```

2. **Per-run configuration**:
   ```go
   // Run with tracing disabled
   result, err := runner.RunSync(agent, &runner.RunOptions{
       Input: "Hello, world!",
       RunConfig: &runner.RunConfig{
           TracingDisabled: true,
       },
   })
   ```

```go
// Run with tracing enabled and custom configuration
result, err := runner.RunSync(agent, &runner.RunOptions{
    Input: "Hello, world!",
    RunConfig: &runner.RunConfig{
        TracingDisabled: false,
        TracingConfig: &runner.TracingConfig{
            WorkflowName: "my_workflow",
        },
    },
})
```

By default, traces are written to local log files (`trace_<agent_name>.log`). The environment variable `OPENAI_AGENTS_DISABLE_TRACING` is checked first, matching the behavior of the Python and TypeScript SDKs.

</details>

### Structured Output

<details>
<summary>Parse responses into Go structs</summary>

```go
// Define an output type
type WeatherReport struct {
    City        string  `json:"city"`
    Temperature float64 `json:"temperature"`
    Condition   string  `json:"condition"`
}

// Create an agent with structured output
agent := agent.NewAgent("Weather Agent")
agent.SetSystemInstructions("You provide weather reports")
agent.SetOutputType(reflect.TypeOf(WeatherReport{}))
```

</details>

### Streaming

<details>
<summary>Get real-time streaming responses</summary>

```go
// Run the agent with streaming
streamedResult, err := runner.RunStreaming(context.Background(), agent, &runner.RunOptions{
    Input: "Hello, world!",
})
if err != nil {
    log.Fatalf("Error running agent: %v", err)
}

// Process streaming events
for event := range streamedResult.Stream {
    switch event.Type {
    case model.StreamEventTypeContent:
        fmt.Print(event.Content)
    case model.StreamEventTypeToolCall:
        fmt.Printf("\nCalling tool: %s\n", event.ToolCall.Name)
    case model.StreamEventTypeDone:
        fmt.Println("\nDone!")
    }
}
```

</details>

### OpenAI Tool Definitions

<details>
<summary>Work with OpenAI-compatible tool definitions</summary>

```go
// Auto-generate OpenAI-compatible tool definitions from Go functions
getCurrentTimeTool := tool.NewFunctionTool(
    "get_current_time",
    "Get the current time in a specified format",
    func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
        return time.Now().Format(time.RFC3339), nil
    },
)

// Convert it to OpenAI format (handled automatically when added to an agent)
openAITool := tool.ToOpenAITool(getCurrentTimeTool)

// Add an OpenAI-compatible tool definition directly to an agent
agent := agent.NewAgent("My Agent")
agent.AddToolFromDefinition(openAITool)

// Add multiple tool definitions at once
toolDefinitions := []map[string]interface{}{
    tool.ToOpenAITool(tool1),
    tool.ToOpenAITool(tool2),
}

agent.AddToolsFromDefinitions(toolDefinitions)
```

</details>

### Workflow State Management

<details>
<summary>Manage state between agent executions</summary>

```go
// Create a state store
stateStore := mocks.NewInMemoryStateStore()

// Create workflow configuration
workflowConfig := &runner.WorkflowConfig{
    RetryConfig: &runner.RetryConfig{
        MaxRetries:         2,
        RetryDelay:        time.Second,
        RetryBackoffFactor: 2.0,
    },
    StateManagement: &runner.StateManagementConfig{
        PersistState:        true,
        StateStore:          stateStore,
        CheckpointFrequency: time.Second * 5,
    },
    ValidationConfig: &runner.ValidationConfig{
        PreHandoffValidation: []runner.ValidationRule{
            {
                Name:         "StateValidation",
                Validate:     func(data interface{}) (bool, error) {
                    state, ok := data.(*runner.WorkflowState)
                    return ok && state != nil, nil
                },
                ErrorMessage: "Invalid workflow state",
                Severity:     runner.ValidationWarning,
            },
        },
    },
}

// Create workflow runner
workflowRunner := runner.NewWorkflowRunner(baseRunner, workflowConfig)

// Initialize workflow state
state := &runner.WorkflowState{
    CurrentPhase:    "",
    CompletedPhases: make([]string, 0),
    Artifacts:       make(map[string]interface{}),
    LastCheckpoint:  time.Now(),
    Metadata:        make(map[string]interface{}),
}

// Run workflow with state management
result, err := workflowRunner.RunWorkflow(context.Background(), agent, &runner.RunOptions{
    MaxTurns:       10,
    RunConfig:      runConfig,
    WorkflowConfig: workflowConfig,
    Input:         state,
})
```

See the complete example in [examples/workflow_example](./examples/workflow_example).
</details>

## 📚 Examples

The repository includes several examples to help you get started:

| Example | Description |
|---------|-------------|
| [Multi-Agent Example](./examples/multi_agent_example) | Demonstrates how to create a system of specialized agents that can collaborate on complex tasks using a local LLM via LM Studio |
| [OpenAI Example](./examples/openai_example) | Shows how to use the OpenAI provider with function calling capabilities |
| [OpenAI Multi-Agent Example](./examples/openai_multi_agent_example) | Illustrates multi-agent functionality using OpenAI models, with proper tool calling and streaming support |
| [Anthropic Example](./examples/anthropic_example) | Demonstrates how to use the Anthropic Claude API with tool calling capabilities |
| [Anthropic Handoff Example](./examples/anthropic_handoff_example) | Shows how to implement agent handoffs with Anthropic Claude models |
| [Bidirectional Flow Example](./examples/bidirectional_flow_example) | Demonstrates bidirectional agent communication with task delegation and return handoffs |
| [TypeScript Code Review Example](./examples/typescript_code_review_example) | Shows a practical application with specialized code review agents that collaborate using bidirectional handoffs |
| [Workflow Example](./examples/workflow_example) | Demonstrates advanced workflow management with state persistence between agent executions |
| [Complex Agentic Flow Test](./store.go) | Comprehensive test demonstrating multi-agent workflows with context sharing, tool execution, and handoffs (run with `go run -tags=store store.go`) |

### Running Examples with a Local LLM

1. Make sure LM Studio is running with a server at `http://127.0.0.1:1234/v1`
2. Navigate to the example directory
   ```bash
   cd examples/multi_agent_example # or any other example using LM Studio
   ```
3. Run the example
   ```bash
   go run .
   ```

### Running Examples with OpenAI

1. Set your OpenAI API key as an environment variable
   ```bash
   export OPENAI_API_KEY=your-api-key
   ```
2. Navigate to the example directory
   ```bash
   cd examples/openai_example # or openai_multi_agent_example
   ```
3. Run the example
   ```bash
   go run .
   ```

### Running Examples with Anthropic

1. Set your Anthropic API key as an environment variable
   ```bash
   export ANTHROPIC_API_KEY=your-anthropic-api-key
   ```
2. Navigate to the example directory
   ```bash
   cd examples/anthropic_example # or anthropic_handoff_example
   ```
3. Run the example
   ```bash
   go run .
   ```

### Debugging

You can enable debug output for various components by setting the appropriate environment variable:

For general debugging (runner and core components):
```bash
DEBUG=1 go run examples/bidirectional_flow_example/main.go
```

For provider-specific debugging:
```bash
# OpenAI provider debugging
OPENAI_DEBUG=1 go run examples/openai_multi_agent_example/main.go

# Anthropic provider debugging
ANTHROPIC_DEBUG=1 go run examples/anthropic_example/main.go

# LM Studio provider debugging
LMSTUDIO_DEBUG=1 go run examples/multi_agent_example/main.go
```

You can also combine multiple debug flags:
```bash
DEBUG=1 OPENAI_DEBUG=1 go run examples/typescript_code_review_example/main.go
```

## 🛠️ Development

<details>
<summary>Development setup and workflows</summary>

### Requirements

- Go 1.23 or later

### Setup

1. Clone the repository
2. Run the setup script to install required tools:

```bash
./scripts/ci_setup.sh
```

### Development Workflow

The project includes several scripts to help with development:

- `./scripts/lint.sh`: Runs formatting and linting checks
- `./scripts/security_check.sh`: Runs security checks with gosec
- `./scripts/check_all.sh`: Runs all checks including tests
- `./scripts/version.sh`: Helps with versioning (run with `bump` argument to bump version)

### Running Tests

Tests are located in the `test` directory and can be run with:

```bash
cd test && make test
```

Or use the check_all script to run all checks including tests:

```bash
./scripts/check_all.sh
```

### CI/CD

The project uses GitHub Actions for CI/CD. The workflow is defined in `.github/workflows/ci.yml`.

</details>

## 👥 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/muhammadhamd/go-agentkit/blob/main/LICENSE) file for details.

## 👨‍💻 Creator

**Go Agent SDK** is created and maintained by **[Muhammad Hamd](https://linkedin.com/in/MuhammadHamd)**.

- 🔗 **LinkedIn**: [linkedin.com/in/MuhammadHamd](https://linkedin.com/in/MuhammadHamd)
- 💻 **GitHub**: [github.com/muhammadhamd](https://github.com/muhammadhamd)

## 🙏 Acknowledgements

This project is inspired by [OpenAI's Assistants API](https://platform.openai.com/docs/assistants/overview) and [OpenAI's Python Agent SDK](https://github.com/openai/openai-agents-py), with the goal of providing similar capabilities in Go while being compatible with local LLMs and multiple providers.

Special thanks to the OpenAI team for their excellent Python and TypeScript SDKs, which served as the reference implementation for this Go version.

## ☁️ Cloud Support

For production deployments, we're developing a fully managed cloud service. Join our waitlist to be among the first to access:

- **Managed Agent Deployment** - Deploy agents without infrastructure hassle
- **Horizontal Scaling** - Handle any traffic volume
- **Observability & Monitoring** - Track performance and usage
- **Cost Optimization** - Pay only for what you use
- **Enterprise Security** - SOC2 compliance and data protection

**[Sign up for the Cloud Waitlist →](https://go-agent.org/#waitlist)**

## 👥 Community & Support

- **Website**: [go-agent.org](https://go-agent.org)
- **GitHub Issues**: [Report bugs or request features](https://github.com/muhammadhamd/go-agentkit/issues)
- **Discussions**: [Join the conversation](https://github.com/muhammadhamd/go-agentkit/discussions)
- **Waitlist**: [Join the cloud service waitlist](https://go-agent.org/#waitlist) 