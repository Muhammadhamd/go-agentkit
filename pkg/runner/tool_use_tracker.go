package runner

// AgentToolUseTracker tracks which tools each agent has used
// Similar to OpenAI's AgentToolUseTracker
type AgentToolUseTracker struct {
	agentToTools map[string][]string // Maps agent name to list of tool names used
}

// NewAgentToolUseTracker creates a new tool use tracker
func NewAgentToolUseTracker() *AgentToolUseTracker {
	return &AgentToolUseTracker{
		agentToTools: make(map[string][]string),
	}
}

// AddToolUse adds tool names to the tracker for an agent
func (t *AgentToolUseTracker) AddToolUse(agentName string, toolNames []string) {
	if existing, ok := t.agentToTools[agentName]; ok {
		t.agentToTools[agentName] = append(existing, toolNames...)
	} else {
		t.agentToTools[agentName] = append([]string{}, toolNames...)
	}
}

// HasUsedTools checks if an agent has used any tools
func (t *AgentToolUseTracker) HasUsedTools(agentName string) bool {
	tools, ok := t.agentToTools[agentName]
	return ok && len(tools) > 0
}

