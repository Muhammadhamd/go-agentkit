package local

// ToolFilter defines an interface for filtering MCP tools
type ToolFilter interface {
	// ShouldInclude returns true if the tool should be included
	ShouldInclude(toolName string) bool
}

// StaticToolFilter implements a static allow/deny list filter
type StaticToolFilter struct {
	allowedTools map[string]bool
	deniedTools  map[string]bool
	allowAll     bool
}

// NewStaticToolFilter creates a new static tool filter
func NewStaticToolFilter() *StaticToolFilter {
	return &StaticToolFilter{
		allowedTools: make(map[string]bool),
		deniedTools:  make(map[string]bool),
		allowAll:     true, // Default: allow all
	}
}

// AllowAll sets the filter to allow all tools
func (f *StaticToolFilter) AllowAll() *StaticToolFilter {
	f.allowAll = true
	f.allowedTools = make(map[string]bool)
	return f
}

// AllowOnly sets the filter to allow only the specified tools
func (f *StaticToolFilter) AllowOnly(toolNames ...string) *StaticToolFilter {
	f.allowAll = false
	f.allowedTools = make(map[string]bool)
	for _, name := range toolNames {
		f.allowedTools[name] = true
	}
	return f
}

// Deny sets the filter to deny specific tools
func (f *StaticToolFilter) Deny(toolNames ...string) *StaticToolFilter {
	for _, name := range toolNames {
		f.deniedTools[name] = true
	}
	return f
}

// ShouldInclude implements ToolFilter interface
func (f *StaticToolFilter) ShouldInclude(toolName string) bool {
	// Check deny list first
	if f.deniedTools[toolName] {
		return false
	}

	// If allow all, return true (unless denied)
	if f.allowAll {
		return true
	}

	// Check allow list
	return f.allowedTools[toolName]
}
