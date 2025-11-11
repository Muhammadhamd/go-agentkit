package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/mcp"
	"github.com/muhammadhamd/go-agentkit/pkg/mcp/local"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
)

func main() {
	ctx := context.Background()

	// Create OpenAI provider
	provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
	provider.WithDefaultModel("gpt-3.5-turbo")

	// Create stdio MCP server transport
	stdioTransport := local.NewStdioServer(local.StdioServerConfig{
		Command: "node",
		Args:    []string{"path/to/mcp-server.js"},
	})

	// Fetch and convert MCP tools
	mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
		Transports:             []mcp.Transport{stdioTransport},
		ConvertSchemasToStrict: true,
	})
	if err != nil {
		log.Fatalf("Failed to get MCP tools: %v", err)
	}

	// Create specialized agent with MCP tools
	knowledgeAgent := agent.NewAgent("Knowledge Agent")
	knowledgeAgent.Description = "Specialized agent that uses MCP tools for knowledge retrieval"
	knowledgeAgent.SetSystemInstructions("You are a specialized knowledge agent. Use the available MCP tools to retrieve information.")
	knowledgeAgent.WithModel("gpt-3.5-turbo")
	knowledgeAgent.SetModelProvider(provider)
	knowledgeAgent.WithTools(mcpTools...)

	// Create triage agent that can handoff to knowledge agent
	triageAgent := agent.NewAgent("Triage Agent")
	triageAgent.Description = "Main agent that coordinates requests and delegates to specialized agents"
	triageAgent.SetSystemInstructions(`You are a triage agent. When a user needs information retrieval, 
hand off to the Knowledge Agent. Otherwise, handle the request yourself.`)
	triageAgent.WithModel("gpt-3.5-turbo")
	triageAgent.SetModelProvider(provider)
	triageAgent.WithHandoffs(knowledgeAgent)

	// Create runner
	r := runner.NewRunner()
	r.WithDefaultProvider(provider)

	// Run the triage agent
	result, err := r.Run(ctx, triageAgent, &runner.RunOptions{
		Input: "I need information about Model Context Protocol. Please use the Knowledge Agent.",
	})
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	fmt.Println("\nFinal response:")
	fmt.Println(result.FinalOutput)

	// Print handoff history
	fmt.Println("\nHandoff items:")
	for _, item := range result.NewItems {
		fmt.Printf("- %+v\n", item)
	}
}
