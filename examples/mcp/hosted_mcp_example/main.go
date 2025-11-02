package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Muhammadhamd/agent-sdk-go/pkg/agent"
	"github.com/Muhammadhamd/agent-sdk-go/pkg/mcp/hosted"
	"github.com/Muhammadhamd/agent-sdk-go/pkg/model/providers/openai"
	"github.com/Muhammadhamd/agent-sdk-go/pkg/runner"
)

func main() {
	ctx := context.Background()

	// Create OpenAI provider
	provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
	provider.WithDefaultModel("gpt-3.5-turbo")

	// Create hosted MCP tool
	hostedTool, err := hosted.NewHostedMCPTool(hosted.HostedMCPToolConfig{
		ServerLabel:  "deepwiki",
		ServerURL:    "https://mcp.example.com/mcp",
		AllowedTools: []string{"read", "search"},
		Headers: map[string]string{
			"Authorization": "Bearer " + os.Getenv("MCP_API_KEY"),
		},
		RequireApproval: hosted.ApprovalNever,
		Description:     "Access to DeepWiki knowledge base via MCP",
	})
	if err != nil {
		log.Fatalf("Failed to create hosted MCP tool: %v", err)
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
