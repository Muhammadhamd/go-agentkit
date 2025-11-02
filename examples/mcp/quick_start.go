package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Muhammadhamd/go-agentkit/pkg/agent"
	"github.com/Muhammadhamd/go-agentkit/pkg/mcp"
	"github.com/Muhammadhamd/go-agentkit/pkg/mcp/local"
	"github.com/Muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/Muhammadhamd/go-agentkit/pkg/runner"
)

// Quick Start Example: Using Local MCP Server
func main() {
	ctx := context.Background()

	// Step 1: Setup OpenAI provider
	provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
	provider.WithDefaultModel("gpt-4o-mini")

	// Step 2: Create stdio transport for local MCP server
	// Example: Assuming you have an MCP server script at ./mcp-server.js
	stdioTransport := local.NewStdioServer(local.StdioServerConfig{
		Command: "node",
		Args:    []string{"./mcp-server.js"},
	})

	// Step 3: Get MCP tools and convert them to SDK tools
	mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
		Transports:             []mcp.Transport{stdioTransport},
		ConvertSchemasToStrict: true, // Strongly recommended
	})
	if err != nil {
		log.Fatalf("Failed to get MCP tools: %v", err)
	}

	fmt.Printf("✓ Connected to MCP server\n")
	fmt.Printf("✓ Found %d MCP tools\n", len(mcpTools))

	// List available tools
	for _, tool := range mcpTools {
		fmt.Printf("  - %s: %s\n", tool.GetName(), tool.GetDescription())
	}

	// Step 4: Create agent with MCP tools
	agent := agent.NewAgent("MCP Assistant")
	agent.Description = "An assistant with access to MCP tools"
	agent.SetSystemInstructions(`You are a helpful assistant with access to MCP tools.
Use the available tools when appropriate to help the user.`)
	agent.WithModel("gpt-4o-mini")
	agent.SetModelProvider(provider)
	agent.WithTools(mcpTools...)

	// Step 5: Create runner and execute
	r := runner.NewRunner()
	r.WithDefaultProvider(provider)

	fmt.Println("\n🤖 Running agent...")
	result, err := r.Run(ctx, agent, &runner.RunOptions{
		Input: "Hello! What tools do you have available?",
	})
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	fmt.Println("\n📝 Agent Response:")
	fmt.Println(result.FinalOutput)
}
