package runner

import (
	"sync"
)

// RunContext shares context across agent turns and handoffs
// Similar to OpenAI's RunContext, it tracks usage, approvals, and custom context
type RunContext struct {
	// Context is user-provided custom context
	Context interface{}

	// Usage tracks token usage across the run
	Usage *Usage

	// Approvals tracks tool approval states
	approvals map[string]*ApprovalRecord

	mu sync.RWMutex
}

// Usage tracks token usage
type Usage struct {
	Requests     int
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// ApprovalRecord tracks approval state for a tool call
type ApprovalRecord struct {
	Approved bool
	Rejected bool
	ToolName string
	CallID   string
}

// NewRunContext creates a new RunContext
func NewRunContext(context interface{}) *RunContext {
	return &RunContext{
		Context:   context,
		Usage:     &Usage{},
		approvals: make(map[string]*ApprovalRecord),
	}
}

// IsToolApproved checks if a tool call is approved
func (rc *RunContext) IsToolApproved(toolName, callID string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	key := rc.approvalKey(toolName, callID)
	record, exists := rc.approvals[key]
	return exists && record.Approved
}

// IsToolRejected checks if a tool call is rejected
func (rc *RunContext) IsToolRejected(toolName, callID string) bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	key := rc.approvalKey(toolName, callID)
	record, exists := rc.approvals[key]
	return exists && record.Rejected
}

// ApproveTool marks a tool call as approved
func (rc *RunContext) ApproveTool(toolName, callID string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	key := rc.approvalKey(toolName, callID)
	rc.approvals[key] = &ApprovalRecord{
		Approved: true,
		Rejected: false,
		ToolName: toolName,
		CallID:   callID,
	}
}

// RejectTool marks a tool call as rejected
func (rc *RunContext) RejectTool(toolName, callID string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	key := rc.approvalKey(toolName, callID)
	rc.approvals[key] = &ApprovalRecord{
		Approved: false,
		Rejected: true,
		ToolName: toolName,
		CallID:   callID,
	}
}

// approvalKey generates a unique key for a tool approval
func (rc *RunContext) approvalKey(toolName, callID string) string {
	return toolName + ":" + callID
}

// AddUsage adds token usage to the context
func (rc *RunContext) AddUsage(requests, inputTokens, outputTokens, totalTokens int) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.Usage.Requests += requests
	rc.Usage.InputTokens += inputTokens
	rc.Usage.OutputTokens += outputTokens
	rc.Usage.TotalTokens += totalTokens
}
