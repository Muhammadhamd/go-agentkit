# DeepSeek Integration Example

This example shows how to connect to DeepSeek API using a custom base URL with the Go Agent SDK.

## DeepSeek API

DeepSeek provides an OpenAI-compatible API, so we can use the OpenAI provider with DeepSeek's base URL.

## Setup

1. Get your DeepSeek API key from [DeepSeek Platform](https://platform.deepseek.com/)

2. Set the API key as an environment variable:
   ```bash
   export DEEPSEEK_API_KEY=your-deepseek-api-key
   ```

3. Run the example:
   ```bash
   go run examples/deepseek_example/main.go
   ```

## Code Example

```go
package main

import (
    "context"
    "os"
    
    "github.com/muhammadhamd/go-agentkit/pkg/agent"
    "github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
    "github.com/muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
    // Get API key
    apiKey := os.Getenv("DEEPSEEK_API_KEY")
    
    // Create provider with DeepSeek base URL
    deepseekProvider := openai.NewOpenAIProvider(apiKey).
        SetBaseURL("https://api.deepseek.com/v1").
        WithDefaultModel("deepseek-chat")
    
    // Get the model
    model, _ := deepseekProvider.GetModel("deepseek-chat")
    
    // Create agent
    agent := agent.NewAgent("DeepSeek Assistant", "You are a helpful assistant.")
    agent.SetModelProvider(deepseekProvider)
    agent.WithModel(model)
    
    // Create runner
    r := runner.NewRunner()
    r.WithDefaultProvider(deepseekProvider)
    
    // Run
    result, _ := r.Run(context.Background(), agent, &runner.RunOptions{
        Input: "Hello!",
    })
    
    fmt.Println(result.FinalOutput)
}
```

## DeepSeek Configuration

- **Base URL**: `https://api.deepseek.com/v1`
- **Available Models**:
  - `deepseek-chat` - Main chat model
  - `deepseek-coder` - Code-focused model

## Alternative: Using Method Chaining

You can also configure it using method chaining:

```go
deepseekProvider := openai.NewOpenAIProvider(apiKey).
    SetBaseURL("https://api.deepseek.com/v1").
    WithDefaultModel("deepseek-chat").
    WithRateLimit(100, 50000) // Optional: Set rate limits
```

## Notes

- DeepSeek uses OpenAI-compatible API format
- All OpenAI provider features work with DeepSeek (tools, streaming, etc.)
- Make sure to use the correct model name for DeepSeek

