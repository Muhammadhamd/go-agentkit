package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/mcp/hosted"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
	ctx := context.Background()

	// Create OpenAI provider
	provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
	provider.WithDefaultModel("gpt-3.5-turbo")

	// Create hosted MCP tool
	hostedTool, err := hosted.NewHostedMCPTool(hosted.HostedMCPToolConfig{
		ServerLabel: "local-mcp",
		ServerURL:   "http://127.0.0.1:8000/mcp/", // Adjust to your MCP JSON-RPC endpoint
		// Leave AllowedTools empty to avoid filtering until tool names are known
		Headers: map[string]string{
			"Authorization": "Bearer " + os.Getenv("MCP_API_KEY"),
		},
		RequireApproval: hosted.ApprovalNever,
		Description:     "Access to local MCP server",
	})
	if err != nil {
		log.Fatalf("Failed to create hosted MCP tool: %v", err)
	}

	// Eagerly connect to MCP to verify connectivity and trigger POST initialize
	if ht, ok := hostedTool.(*hosted.HostedMCPTool); ok {
		if err := ht.Connect(ctx); err != nil {
			log.Fatalf("Failed to connect to MCP: %v", err)
		}
	}

	// Create an agent with hosted MCP tool
	agent := agent.NewAgent("Knowledge Agent")
	agent.Description = "An agent that uses hosted MCP tools for knowledge retrieval"
	agent.SetSystemInstructions("You are a helpful agent with access to hosted MCP tools. Use them to retrieve information when needed.")
	agent.WithModel("gpt-3.5-turbo")
	agent.SetModelProvider(provider)
	agent.WithTools(hostedTool)

	// Create runner
	r := runner.NewRunner()
	r.WithDefaultProvider(provider)

	// Run the agent
	result, err := r.Run(ctx, agent, &runner.RunOptions{
		Input: "Use the MCP tools to search for information about AI agents",
	})
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	fmt.Println("\nAgent response:")
	fmt.Println(result.FinalOutput)
}
