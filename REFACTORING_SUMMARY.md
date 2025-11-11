# Agentic Loop Refactoring Summary

## Overview
This document summarizes the refactoring of the Go agentic loop to match OpenAI's implementation pattern.

## Key Changes

### 1. RunState Structure (`pkg/runner/run_state.go`)
- **Created**: `RunState` struct that maintains state similar to OpenAI's `RunState`
- **Key Features**:
  - `OriginalInput`: Never mutated, preserved separately
  - `GeneratedItems`: Accumulated over time (messages, tool calls, results)
  - `CurrentAgent`: Tracks the active agent
  - `CurrentStep`: Uses NextStep pattern for state machine
  - `GetTurnInput()`: Efficiently combines originalInput + generatedItems

### 2. NextStep Pattern (`pkg/runner/next_step.go`)
- **Created**: Discriminated union pattern for next steps
- **Types**:
  - `NextStepRunAgain`: Continue with another turn
  - `NextStepFinalOutput`: Agent produced final output
  - `NextStepHandoff`: Handoff to another agent
  - `NextStepInterruption`: Interruption (e.g., tool approvals)

### 3. RunContext (`pkg/runner/run_context.go`)
- **Created**: Context sharing across handoffs
- **Features**:
  - Custom context storage
  - Token usage tracking
  - Tool approval state management

### 4. Main Loop Refactoring (`pkg/runner/runner.go`)
- **Before**: Used local variables, mutated `currentInput` each turn
- **After**: Uses `RunState` and `NextStep` pattern
- **Pattern**: 
  ```go
  for {
    switch step := state.CurrentStep.(type) {
    case *NextStepRunAgain:
      // Process turn
    case *NextStepFinalOutput:
      // Return result
    case *NextStepHandoff:
      // Switch agent
    }
  }
  ```

### 5. Tool Processing (`pkg/runner/tool_processing.go`)
- **Created**: `ProcessedResponse` to categorize model outputs
- **Created**: `executeFunctionTools` to execute tools and create RunItems
- **Added**: `ToolUseBehavior` interface for configurable tool behavior
- **Fixed**: Tool call IDs properly tracked and matched with results

### 6. Handoff Processing (`pkg/runner/handoff_processing.go`)
- **Created**: `HandoffInputData` structure
- **Added**: Support for `HandoffInputFilter` (from RunConfig)
- **Fixed**: Handoffs now pass full conversation history + handoff input

### 7. Turn Processing (`pkg/runner/turn_processing.go`)
- **Created**: `processSingleTurn` function
- **Pattern**: Matches OpenAI's `_run_single_turn` / `processSingleTurn`
- **Flow**:
  1. Increment turn, check max turns
  2. Get turn input (originalInput + generatedItems)
  3. Execute model request
  4. Process response (categorize into tools/handoffs/messages)
  5. Execute tools/handoffs
  6. Determine next step

### 8. Result Items (`pkg/result/result.go`)
- **Enhanced**: `MessageItem` now supports `ToolCalls` field
- **Enhanced**: `ToolResultItem` now stores `ToolCallID` for proper matching
- **Fixed**: `ToInputItem()` methods return proper format for model providers

## Architecture Comparison

### OpenAI Pattern
```
RunState {
  _originalInput (never mutated)
  _generatedItems (accumulated)
  _currentStep (NextStep)
  _currentAgent
}

Loop:
  while (true) {
    switch (state._currentStep.type) {
      case 'next_step_run_again':
        turnResult = processSingleTurn()
        state._currentStep = turnResult.nextStep
      case 'next_step_final_output':
        return
      case 'next_step_handoff':
        switch agent
    }
  }
```

### Go Pattern (After Refactoring)
```
RunState {
  OriginalInput (never mutated)
  GeneratedItems (accumulated)
  CurrentStep (NextStep interface)
  CurrentAgent
}

Loop:
  for {
    switch step := state.CurrentStep.(type) {
    case *NextStepRunAgain:
      turnResult = processSingleTurn()
      state.CurrentStep = turnResult.NextStep
    case *NextStepFinalOutput:
      return
    case *NextStepHandoff:
      switch agent
    }
  }
```

## Key Improvements

1. **State Management**: Original input preserved, generated items accumulated
2. **Input History**: Efficient `GetTurnInput()` instead of reconstructing each time
3. **State Machine**: Explicit NextStep pattern instead of boolean flags
4. **Tool Use Behavior**: Extensible interface for custom tool behavior
5. **Handoff Filtering**: Support for input filters during handoffs
6. **Context Sharing**: RunContext maintains state across handoffs

## Remaining Work

1. **Tool Approvals**: Implement `NextStepInterruption` handling
2. **Structured Output**: Complete structured output parsing
3. **Guardrails**: Integrate input/output guardrails
4. **Testing**: Comprehensive testing to verify OpenAI compatibility

## Files Created

- `pkg/runner/run_state.go`
- `pkg/runner/next_step.go`
- `pkg/runner/run_context.go`
- `pkg/runner/handoff_input_data.go`
- `pkg/runner/turn_result.go`
- `pkg/runner/tool_processing.go`
- `pkg/runner/handoff_processing.go`
- `pkg/runner/turn_processing.go`

## Files Modified

- `pkg/runner/runner.go` (main loop refactored)
- `pkg/result/result.go` (enhanced MessageItem and ToolResultItem)

