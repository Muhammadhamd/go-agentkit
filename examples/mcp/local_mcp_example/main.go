package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Muhammadhamd/go-agentkit/pkg/agent"
	"github.com/Muhammadhamd/go-agentkit/pkg/mcp"
	"github.com/Muhammadhamd/go-agentkit/pkg/mcp/local"
	"github.com/Muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/Muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
	ctx := context.Background()

	// Create OpenAI provider
	provider := openai.NewProvider("your-api-key")
	provider.WithDefaultModel("gpt-3.5-turbo")

	// Create stdio MCP server transport
	// Example: connecting to a local MCP server
	stdioTransport := local.NewStdioServer(local.StdioServerConfig{
		Command: "node",
		Args:    []string{"path/to/mcp-server.js"},
		// Optional: Env: []string{"NODE_ENV=production"},
		// Optional: Dir: "/path/to/working/directory",
	})

	// Fetch and convert MCP tools
	mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
		Transports:             []mcp.Transport{stdioTransport},
		ConvertSchemasToStrict: true, // Strongly recommended
	})
	if err != nil {
		log.Fatalf("Failed to get MCP tools: %v", err)
	}

	fmt.Printf("Found %d MCP tools\n", len(mcpTools))

	// Create an agent with MCP tools
	agent := agent.NewAgent("MCP Agent")
	agent.Description = "An agent that uses MCP tools"
	agent.SetSystemInstructions("You are a helpful agent with access to MCP tools. Use them when appropriate.")
	agent.WithModel("gpt-3.5-turbo")
	agent.SetModelProvider(provider)
	agent.WithTools(mcpTools...)

	// Create runner
	r := runner.NewRunner()
	r.WithDefaultProvider(provider)

	// Run the agent
	result, err := r.Run(ctx, agent, &runner.RunOptions{
		Input: "Use the MCP tools to help me with my task",
	})
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	fmt.Println("\nAgent response:")
	fmt.Println(result.FinalOutput)
}
