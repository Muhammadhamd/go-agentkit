# MCP Usage Examples

This guide shows how to use Model Context Protocol (MCP) with the Agent SDK Go.

## Table of Contents

1. [Local MCP (stdio)](#local-mcp-stdio)
2. [Hosted MCP (HTTP/SSE)](#hosted-mcp-httpsse)
3. [MCP with Agent Handoffs](#mcp-with-agent-handoffs)
4. [Tool Filtering](#tool-filtering)
5. [Approval System](#approval-system)

## Local MCP (stdio)

Connect to a local MCP server running as a subprocess.

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Muhammadhamd/agent-sdk-go/pkg/agent"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/mcp"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/mcp/local"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/model/providers/openai"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/runner"
)

func main() {
    ctx := context.Background()

    // 1. Create OpenAI provider
    provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
    provider.WithDefaultModel("gpt-3.5-turbo")

    // 2. Create stdio MCP server transport
    stdioTransport := local.NewStdioServer(local.StdioServerConfig{
        Command: "node",
        Args:    []string{"path/to/your-mcp-server.js"},
        // Optional: Env: []string{"NODE_ENV=production"},
        // Optional: Dir: "/path/to/working/directory",
    })

    // 3. Fetch and convert MCP tools
    mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
        Transports:             []mcp.Transport{stdioTransport},
        ConvertSchemasToStrict: true, // Strongly recommended
    })
    if err != nil {
        log.Fatalf("Failed to get MCP tools: %v", err)
    }

    fmt.Printf("Found %d MCP tools\n", len(mcpTools))

    // 4. Create agent with MCP tools
    agent := agent.NewAgent("MCP Agent")
    agent.SetSystemInstructions("You are a helpful agent with access to MCP tools. Use them when appropriate.")
    agent.WithModel("gpt-3.5-turbo")
    agent.SetModelProvider(provider)
    agent.WithTools(mcpTools...)

    // 5. Create runner and run
    r := runner.NewRunner()
    r.WithDefaultProvider(provider)

    result, err := r.Run(ctx, agent, &runner.RunOptions{
        Input: "Use the MCP tools to help me with my task",
    })
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(result.FinalOutput)
}
```

### Multiple MCP Servers

```go
// Connect to multiple MCP servers
transports := []mcp.Transport{
    local.NewStdioServer(local.StdioServerConfig{
        Command: "node",
        Args:    []string{"mcp-server-1.js"},
    }),
    local.NewStdioServer(local.StdioServerConfig{
        Command: "python",
        Args:    []string{"mcp-server-2.py"},
    }),
}

mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
    Transports:             transports,
    ConvertSchemasToStrict: true,
})
```

## Hosted MCP (HTTP/SSE)

Connect to a remote MCP server via HTTP/SSE.

### Basic Example

```go
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

    // 1. Create OpenAI provider
    provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
    provider.WithDefaultModel("gpt-3.5-turbo")

    // 2. Create hosted MCP tool
    hostedTool, err := hosted.NewHostedMCPTool(hosted.HostedMCPToolConfig{
        ServerLabel:  "deepwiki",
        ServerURL:    "https://mcp.example.com/mcp",
        AllowedTools: []string{"read", "search"}, // Optional: restrict which tools
        Headers: map[string]string{
            "Authorization": "Bearer " + os.Getenv("MCP_API_KEY"),
        },
        RequireApproval: hosted.ApprovalNever,
        Description:     "Access to DeepWiki knowledge base via MCP",
    })
    if err != nil {
        log.Fatalf("Failed to create hosted MCP tool: %v", err)
    }

    // 3. Create agent with hosted tool
    agent := agent.NewAgent("Knowledge Agent")
    agent.Description = "An agent that uses hosted MCP tools"
    agent.SetSystemInstructions("Use the MCP tools to retrieve information.")
    agent.WithModel("gpt-3.5-turbo")
    agent.SetModelProvider(provider)
    agent.WithTools(hostedTool)

    // 4. Run agent
    r := runner.NewRunner()
    r.WithDefaultProvider(provider)

    result, err := r.Run(ctx, agent, &runner.RunOptions{
        Input: "Search for information about AI agents",
    })
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(result.FinalOutput)
}
```

## MCP with Agent Handoffs

Combine MCP tools with agent handoffs for complex workflows.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Muhammadhamd/agent-sdk-go/pkg/agent"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/mcp"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/mcp/local"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/model/providers/openai"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/runner"
)

func main() {
    ctx := context.Background()
    provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
    provider.WithDefaultModel("gpt-3.5-turbo")

    // 1. Get MCP tools
    stdioTransport := local.NewStdioServer(local.StdioServerConfig{
        Command: "node",
        Args:    []string{"knowledge-mcp-server.js"},
    })

    mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
        Transports:             []mcp.Transport{stdioTransport},
        ConvertSchemasToStrict: true,
    })
    if err != nil {
        log.Fatalf("Failed to get MCP tools: %v", err)
    }

    // 2. Create specialized agent with MCP tools
    knowledgeAgent := agent.NewAgent("Knowledge Agent")
    knowledgeAgent.Description = "Specialized agent with MCP tools for knowledge retrieval"
    knowledgeAgent.SetSystemInstructions("Use MCP tools to retrieve information.")
    knowledgeAgent.WithModel("gpt-3.5-turbo")
    knowledgeAgent.SetModelProvider(provider)
    knowledgeAgent.WithTools(mcpTools...)

    // 3. Create triage agent that delegates to knowledge agent
    triageAgent := agent.NewAgent("Triage Agent")
    triageAgent.Description = "Main agent that coordinates and delegates"
    triageAgent.SetSystemInstructions(`When the user needs information retrieval, 
hand off to the Knowledge Agent. Otherwise, handle requests yourself.`)
    triageAgent.WithModel("gpt-3.5-turbo")
    triageAgent.SetModelProvider(provider)
    triageAgent.WithHandoffs(knowledgeAgent)

    // 4. Run
    r := runner.NewRunner()
    r.WithDefaultProvider(provider)

    result, err := r.Run(ctx, triageAgent, &runner.RunOptions{
        Input: "I need information about Model Context Protocol",
    })
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(result.FinalOutput)
}
```

## Tool Filtering

Filter which MCP tools to include.

```go
// Using tool filter function
mcpTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
    Transports:             []mcp.Transport{stdioTransport},
    ConvertSchemasToStrict: true,
    ToolFilter: func(toolName string) bool {
        // Only include specific tools
        allowed := []string{"read", "search", "write"}
        for _, name := range allowed {
            if toolName == name {
                return true
            }
        }
        return false
    },
})

// Or use static filter (for local MCP)
filter := local.NewStaticToolFilter()
filter.AllowOnly("read", "search") // Only these tools
// filter.Deny("write", "delete")  // Or deny specific tools
```

## Approval System

Require approval for hosted MCP tools.

### With Approval Callback

```go
hostedTool, err := hosted.NewHostedMCPTool(hosted.HostedMCPToolConfig{
    ServerLabel:  "production-api",
    ServerURL:    "https://api.example.com/mcp",
    AllowedTools: []string{"write", "delete"},
    Headers: map[string]string{
        "Authorization": "Bearer " + os.Getenv("MCP_API_KEY"),
    },
    RequireApproval: hosted.ApprovalAlways, // Always require approval
    OnApproval: func(ctx context.Context, toolName string, params map[string]interface{}) (bool, error) {
        // Your approval logic
        fmt.Printf("Approval requested for tool: %s\n", toolName)
        fmt.Printf("Parameters: %+v\n", params)
        
        // Example: auto-approve read operations, require manual approval for writes
        if toolName == "read" {
            return true, nil
        }
        
        // For writes, you could prompt user, check permissions, etc.
        // For now, auto-approve (in production, implement proper logic)
        return true, nil
    },
    Description: "Production API access with approval",
})
```

### Approval Requirements

```go
// Never require approval (default)
RequireApproval: hosted.ApprovalNever

// Always require approval
RequireApproval: hosted.ApprovalAlways

// Require approval based on tool/operation
RequireApproval: hosted.ApprovalOnTool
// ApprovalOnTool automatically requires approval for operations
// containing "write", "delete", "update", "create" in name or params
```

## Complete Example: Local MCP with File Operations

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Muhammadhamd/agent-sdk-go/pkg/agent"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/mcp"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/mcp/local"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/model/providers/openai"
    "github.com/Muhammadhamd/agent-sdk-go/pkg/runner"
)

func main() {
    ctx := context.Background()

    // Setup
    provider := openai.NewProvider(os.Getenv("OPENAI_API_KEY"))
    provider.WithDefaultModel("gpt-3.5-turbo")

    // Connect to file system MCP server
    fsTransport := local.NewStdioServer(local.StdioServerConfig{
        Command: "npx",
        Args:    []string{"-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/directory"},
    })

    // Get file system tools
    fsTools, err := mcp.GetAllMCPTools(ctx, mcp.GetAllMCPToolsConfig{
        Transports:             []mcp.Transport{fsTransport},
        ConvertSchemasToStrict: true,
        ToolFilter: func(name string) bool {
            // Only allow read operations for safety
            return name == "read_file" || name == "list_directory"
        },
    })
    if err != nil {
        log.Fatalf("Failed to get FS tools: %v", err)
    }

    // Create agent
    agent := agent.NewAgent("File Assistant")
    agent.SetSystemInstructions("You can read files and list directories. Be helpful and safe.")
    agent.WithModel("gpt-3.5-turbo")
    agent.SetModelProvider(provider)
    agent.WithTools(fsTools...)

    // Run
    r := runner.NewRunner()
    r.WithDefaultProvider(provider)

    result, err := r.Run(ctx, agent, &runner.RunOptions{
        Input: "Read the README.md file and summarize it",
    })
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(result.FinalOutput)
}
```

## Environment Variables

```bash
# For OpenAI
export OPENAI_API_KEY="sk-your-key-here"

# For MCP servers
export MCP_API_KEY="your-mcp-api-key"
export NODE_ENV="production"
```

## Best Practices

1. **Always use `ConvertSchemasToStrict: true`** - Ensures better schema validation
2. **Filter tools** - Only include tools your agent needs
3. **Use approval for hosted MCP** - Protect against unauthorized operations
4. **Handle errors gracefully** - MCP servers may be unavailable
5. **Close connections** - Properly clean up MCP client connections
6. **Test locally first** - Validate with local MCP servers before production

## Troubleshooting

### "Connection failed"
- Check MCP server is running
- Verify command/args are correct
- Check environment variables

### "Duplicate tool name"
- Tool names must be unique across all MCP servers
- Use tool filtering to exclude duplicates

### "Tool execution failed"
- Check MCP server logs
- Verify tool parameters match schema
- Ensure server has required permissions

### "Transport error"
- For stdio: verify subprocess starts correctly
- For HTTP/SSE: check network connectivity and authentication

