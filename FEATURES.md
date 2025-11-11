# Complete Feature List - Go Agent SDK

## Core Features

1. **Multiple LLM Provider Support**
   - OpenAI (GPT-3.5, GPT-4, GPT-4o, etc.)
   - Anthropic Claude (Haiku, Sonnet, Opus)
   - LM Studio (Local models)
   - Custom Base URLs (DeepSeek, other OpenAI-compatible APIs)

2. **Tool Integration**
   - Function Tools (Convert Go functions to tools)
   - OpenAI-Compatible Tool Definitions
   - Tool Schema Generation
   - Tool Parameter Validation
   - Tool Error Handling

3. **Agent Handoffs**
   - Multi-Agent Workflows
   - Bidirectional Agent Flow
   - Task Delegation
   - Return to Delegator
   - Input Filtering for Handoffs

4. **Structured Output**
   - Parse LLM responses into Go structs
   - JSON Schema Validation
   - Type-safe output handling

5. **Streaming**
   - Real-time streaming responses
   - Stream events (content, tool calls, handoffs, done)
   - AsyncIterable pattern

6. **Tracing & Monitoring**
   - OpenAI Backend Tracing (sends to OpenAI dashboard)
   - Environment variable control (OPENAI_AGENTS_DISABLE_TRACING)
   - Per-run tracing configuration
   - No local file creation (matches Python/TypeScript)

7. **Context Sharing**
   - RunContext for shared data across agents
   - Usage statistics tracking
   - Tool approval states
   - Custom context types

8. **Agentic Loop**
   - Turn-based execution
   - State management (RunState)
   - Next step types (RunAgain, Handoff, FinalOutput, Interruption)
   - Max turns control
   - Consecutive tool call tracking

9. **Workflow State Management**
   - State persistence
   - Retry configuration
   - Recovery functions
   - Checkpointing
   - Validation rules

10. **Lifecycle Hooks**
    - Agent hooks (before/after model call, before/after tool call)
    - Run hooks (run start/end, agent start, turn start/end, handoff)
    - Custom hook implementations

11. **Tool Use Behavior**
    - run_llm_again (default - continue after tools)
    - stop_on_first_tool (stop after first tool call)
    - Custom tool use behavior
    - Reset tool choice (prevents infinite loops)

12. **MCP Support (Model Context Protocol)**
    - Local MCP servers (stdio transport)
    - Hosted MCP servers (HTTP/SSE transport)
    - MCP tool conversion to SDK tools
    - MCP tool filtering

13. **Custom Base URLs**
    - Support for OpenAI-compatible APIs
    - DeepSeek integration example
    - Custom API endpoints

14. **Error Handling**
    - Tool call errors
    - Model behavior errors
    - Max turns exceeded
    - Guardrail errors (input/output)
    - User errors

15. **Usage Tracking**
    - Token usage (input/output/total)
    - Per-agent usage
    - Available in RunContext

16. **Tool Approval (Human-in-the-Loop)**
    - Tool approval states
    - Approve/reject tools
    - Interruption support

17. **Input Filtering**
    - Filter conversation history during handoffs
    - Remove tools from history
    - Custom filter functions

18. **Agent Configuration**
    - System instructions
    - Model settings (temperature, max tokens, etc.)
    - Output type specification
    - Tool use behavior configuration
    - Reset tool choice option

19. **Backward Compatibility**
    - Old agent creation methods still work
    - Method chaining support
    - Direct field access support

20. **Environment Variables**
    - OPENAI_API_KEY (for OpenAI provider)
    - OPENAI_AGENTS_DISABLE_TRACING (disable tracing)
    - DEEPSEEK_API_KEY (for DeepSeek example)
    - ANTHROPIC_API_KEY (for Anthropic provider)

