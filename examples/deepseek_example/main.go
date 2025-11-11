package main

import (
	"context"
	"fmt"
	"os"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
	// Get DeepSeek API key from environment
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: DEEPSEEK_API_KEY environment variable is not set")
		fmt.Println("Please set it with: export DEEPSEEK_API_KEY=your-api-key")
		os.Exit(1)
	}

	// Create OpenAI provider with DeepSeek base URL
	// DeepSeek uses OpenAI-compatible API, so we can use the OpenAI provider
	deepseekProvider := openai.NewOpenAIProvider(apiKey).
		SetBaseURL("https://api.deepseek.com/v1"). // DeepSeek's base URL
		WithDefaultModel("deepseek-chat")          // DeepSeek model name

	// Get the DeepSeek model
	deepseekModel, err := deepseekProvider.GetModel("deepseek-chat")
	if err != nil {
		fmt.Printf("Error getting model: %v\n", err)
		os.Exit(1)
	}

	// Create an agent with DeepSeek model
	deepseekAgent := agent.NewAgent("DeepSeek Assistant", "You are a helpful assistant powered by DeepSeek.")
	deepseekAgent.SetModelProvider(deepseekProvider)
	deepseekAgent.WithModel(deepseekModel)

	// Create runner with the provider
	r := runner.NewRunner()
	r.WithDefaultProvider(deepseekProvider)

	// Run the agent
	ctx := context.Background()
	result, err := r.Run(ctx, deepseekAgent, &runner.RunOptions{
		Input: "Hello! Can you tell me a fun fact about AI?",
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", result.FinalOutput)
}
