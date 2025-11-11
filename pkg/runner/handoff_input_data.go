package runner

import (
	"github.com/muhammadhamd/go-agentkit/pkg/result"
)

// HandoffInputData contains the data passed to a handoff input filter
// This follows OpenAI's HandoffInputData pattern
type HandoffInputData struct {
	// InputHistory is the original input to the run
	InputHistory interface{}

	// PreHandoffItems are items generated before the handoff
	PreHandoffItems []result.RunItem

	// NewItems are items generated in the current step (including the handoff call)
	NewItems []result.RunItem

	// RunContext is the shared run context
	RunContext *RunContext
}

// AllItems returns all items that should be passed to the next agent
// This combines inputHistory, preHandoffItems, and newItems
func (d *HandoffInputData) AllItems() []interface{} {
	var items []interface{}

	// Add input history
	if inputStr, ok := d.InputHistory.(string); ok {
		items = append(items, map[string]interface{}{
			"type":    "message",
			"role":    "user",
			"content": inputStr,
		})
	} else if inputList, ok := d.InputHistory.([]interface{}); ok {
		items = append(items, inputList...)
	}

	// Add pre-handoff items
	for _, item := range d.PreHandoffItems {
		items = append(items, item.ToInputItem())
	}

	// Add new items
	for _, item := range d.NewItems {
		items = append(items, item.ToInputItem())
	}

	return items
}
