//go:build store
// +build store

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/muhammadhamd/go-agentkit/pkg/agent"
	"github.com/muhammadhamd/go-agentkit/pkg/model/providers/openai"
	"github.com/muhammadhamd/go-agentkit/pkg/runner"
	"github.com/muhammadhamd/go-agentkit/pkg/tool"
)

// TestContext holds shared context across agents
type TestContext struct {
	UserID    string
	OrderID   string
	SessionID string
	Metadata  map[string]interface{}
}

// Mock database/store for testing
var mockDB = map[string]interface{}{
	"user_123": map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"tier":  "premium",
	},
	"order_456": map[string]interface{}{
		"user_id":    "user_123",
		"amount":     99.99,
		"status":     "pending",
		"created_at": "2024-01-15",
	},
}

// Tool: Get user information
func getUserInfo(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	userID, ok := params["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id is required")
	}

	// Access shared context
	runCtxVal := ctx.Value("run_context")
	if runCtxVal != nil {
		if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil && runCtx.Context != nil {
			if testCtx, ok := runCtx.Context.(*TestContext); ok {
				testCtx.UserID = userID
			}
		}
	}

	userData, exists := mockDB[userID]
	if !exists {
		return map[string]interface{}{"error": "User not found"}, nil
	}

	return userData, nil
}

// Tool: Check order status
func checkOrderStatus(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	orderID, ok := params["order_id"].(string)
	if !ok {
		return nil, fmt.Errorf("order_id is required")
	}

	// Access shared context
	runCtxVal := ctx.Value("run_context")
	if runCtxVal != nil {
		if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil && runCtx.Context != nil {
			if testCtx, ok := runCtx.Context.(*TestContext); ok {
				testCtx.OrderID = orderID
			}
		}
	}

	orderData, exists := mockDB[orderID]
	if !exists {
		return map[string]interface{}{"error": "Order not found"}, nil
	}

	return orderData, nil
}

// Tool: Calculate refund amount
func calculateRefund(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	orderID, ok := params["order_id"].(string)
	if !ok {
		return nil, fmt.Errorf("order_id is required")
	}

	orderData, exists := mockDB[orderID]
	if !exists {
		return map[string]interface{}{"error": "Order not found"}, nil
	}

	order := orderData.(map[string]interface{})
	amount := order["amount"].(float64)
	refundAmount := amount * 0.9

	// Access shared context to store refund info
	runCtxVal := ctx.Value("run_context")
	if runCtxVal != nil {
		if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil && runCtx.Context != nil {
			if testCtx, ok := runCtx.Context.(*TestContext); ok {
				if testCtx.Metadata == nil {
					testCtx.Metadata = make(map[string]interface{})
				}
				testCtx.Metadata["refund_amount"] = refundAmount
				testCtx.Metadata["refund_calculated_at"] = time.Now().Format(time.RFC3339)
			}
		}
	}

	return map[string]interface{}{
		"order_id":          orderID,
		"original_amount":   amount,
		"refund_amount":     refundAmount,
		"refund_percentage": 90,
	}, nil
}

// Tool: Process refund
func processRefund(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	orderID, ok := params["order_id"].(string)
	if !ok {
		return nil, fmt.Errorf("order_id is required")
	}

	// Access shared context to get refund amount
	runCtxVal := ctx.Value("run_context")
	var refundAmount float64
	if runCtxVal != nil {
		if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil && runCtx.Context != nil {
			if testCtx, ok := runCtx.Context.(*TestContext); ok {
				if testCtx.Metadata != nil {
					if amount, ok := testCtx.Metadata["refund_amount"].(float64); ok {
						refundAmount = amount
					}
				}
			}
		}
	}

	if refundAmount == 0 {
		return map[string]interface{}{"error": "Refund amount not calculated. Please calculate refund first."}, nil
	}

	// Simulate refund processing
	orderData := mockDB[orderID].(map[string]interface{})
	orderData["status"] = "refunded"
	orderData["refund_amount"] = refundAmount
	orderData["refunded_at"] = time.Now().Format(time.RFC3339)

	return map[string]interface{}{
		"order_id":      orderID,
		"refund_amount": refundAmount,
		"status":        "refunded",
		"message":       "Refund processed successfully",
	}, nil
}

// Tool: Get conversation summary
func getConversationSummary(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	summary := map[string]interface{}{
		"tools_called": []string{},
		"context_data": map[string]interface{}{},
	}

	runCtxVal := ctx.Value("run_context")
	if runCtxVal != nil {
		if runCtx, ok := runCtxVal.(*runner.RunContext); ok && runCtx != nil {
			if runCtx.Context != nil {
				if testCtx, ok := runCtx.Context.(*TestContext); ok {
					summary["context_data"] = map[string]interface{}{
						"user_id":    testCtx.UserID,
						"order_id":   testCtx.OrderID,
						"session_id": testCtx.SessionID,
						"metadata":   testCtx.Metadata,
					}
				}
			}
			if runCtx.Usage != nil {
				summary["usage"] = map[string]interface{}{
					"requests":      runCtx.Usage.Requests,
					"input_tokens":  runCtx.Usage.InputTokens,
					"output_tokens": runCtx.Usage.OutputTokens,
					"total_tokens":  runCtx.Usage.TotalTokens,
				}
			}
		}
	}

	return summary, nil
}

func complexAgenticFlowTest() {
	ctx := context.Background()

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable is not set")
		fmt.Println("Please set it with: export OPENAI_API_KEY=your-api-key")
		os.Exit(1)
	}

	provider := openai.NewProvider(apiKey)
	provider.WithDefaultModel("gpt-4o-mini")

	testContext := &TestContext{
		SessionID: "session_" + fmt.Sprintf("%d", time.Now().Unix()),
		Metadata:  make(map[string]interface{}),
	}

	// Create tools
	getUserInfoTool := tool.NewFunctionTool("get_user_info", "Retrieves user information by user_id. Use this to get customer details.", getUserInfo)
	checkOrderStatusTool := tool.NewFunctionTool("check_order_status", "Checks the status of an order by order_id. Returns order details including amount and status.", checkOrderStatus)
	calculateRefundTool := tool.NewFunctionTool("calculate_refund", "Calculates the refund amount for an order (90% of original amount). Requires order_id.", calculateRefund)
	processRefundTool := tool.NewFunctionTool("process_refund", "Processes a refund for an order. Requires order_id. Make sure to calculate refund first.", processRefund)
	getConversationSummaryTool := tool.NewFunctionTool("get_conversation_summary", "Gets a summary of the conversation including all context data and usage statistics.", getConversationSummary)

	// Create Billing Agent
	billingAgent := agent.NewAgent("billing_agent")
	billingAgent.SetModelProvider(provider)
	billingAgent.WithModel("gpt-4o-mini")
	billingAgent.SetSystemInstructions(
		"You are a billing specialist. Your role is to handle refunds and billing inquiries. " +
			"WORKFLOW FOR PROCESSING REFUNDS (MUST FOLLOW IN ORDER): " +
			"1. FIRST: Check the order status using check_order_status tool with the order_id parameter (extract it from the conversation) " +
			"2. SECOND: Calculate the refund amount using calculate_refund tool with the order_id parameter " +
			"3. THIRD: Process the refund using process_refund tool with the order_id parameter " +
			"CRITICAL: You MUST call all three tools in this exact order. Do not skip any steps. " +
			"IMPORTANT: Always extract the order_id from the conversation text and pass it as a parameter to the tools. " +
			"For example, if the conversation mentions 'order_456', use {\"order_id\": \"order_456\"} as the parameter. " +
			"Always use the tools in this order. Be thorough and professional.",
	)
	billingAgent.WithTools(checkOrderStatusTool, calculateRefundTool, processRefundTool, getConversationSummaryTool)

	// Create Support Agent
	supportAgent := agent.NewAgent("support_agent")
	supportAgent.SetModelProvider(provider)
	supportAgent.WithModel("gpt-4o-mini")
	supportAgent.SetSystemInstructions(
		"You are a customer support agent. Your role is to help customers with their inquiries. " +
			"WORKFLOW FOR REFUND REQUESTS (MUST FOLLOW IN ORDER): " +
			"1. FIRST: Get user information using get_user_info tool with the user_id parameter (extract it from the conversation) " +
			"2. SECOND: Check order status using check_order_status tool with the order_id parameter (extract it from the conversation) " +
			"3. THIRD: After getting user info and order status, THEN handoff to billing_agent for refund processing " +
			"CRITICAL: You MUST call get_user_info and check_order_status BEFORE handing off to billing_agent. " +
			"Do NOT handoff to billing_agent until you have called both tools. " +
			"IMPORTANT: Always extract the order_id and user_id from the conversation text and pass them as parameters to the tools. " +
			"For example, if the conversation mentions 'user_123', use {\"user_id\": \"user_123\"} as the parameter. " +
			"If the conversation mentions 'order_456', use {\"order_id\": \"order_456\"} as the parameter. " +
			"Always be helpful and professional.",
	)
	supportAgent.WithTools(getUserInfoTool, checkOrderStatusTool, getConversationSummaryTool)
	supportAgent.WithHandoffs(billingAgent)

	// Create Receptionist Agent
	receptionistAgent := agent.NewAgent("receptionist")
	receptionistAgent.SetModelProvider(provider)
	receptionistAgent.WithModel("gpt-4o-mini")
	receptionistAgent.SetSystemInstructions(
		"You are a receptionist. Your role is to greet customers and route them to the appropriate agent. " +
			"You can handoff to support_agent for any customer inquiries. " +
			"Be friendly and welcoming. Always handoff to support_agent for detailed help.",
	)
	receptionistAgent.WithHandoffs(supportAgent)

	r := runner.NewRunner().WithDefaultProvider(provider)
	input := "Hi, I need a refund for my order. My order ID is order_456 and my user ID is user_123."

	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Println("COMPLEX AGENTIC FLOW TEST")
	fmt.Println("=" + strings.Repeat("=", 80))
	fmt.Printf("Input: %s\n\n", input)
	fmt.Printf("Expected Flow:\n")
	fmt.Printf("1. Receptionist greets customer and hands off to support_agent\n")
	fmt.Printf("2. Support_agent gets user info (user_123) and checks order status (order_456)\n")
	fmt.Printf("3. Support_agent hands off to billing_agent for refund\n")
	fmt.Printf("4. Billing_agent checks order, calculates refund, and processes it\n")
	fmt.Printf("5. Context (user_id, order_id, refund_amount) should be shared across all agents\n\n")
	fmt.Println("=" + strings.Repeat("=", 80))

	res, err := r.Run(ctx, receptionistAgent, &runner.RunOptions{
		Input:    input,
		Context:  testContext,
		MaxTurns: 20,
	})

	if err != nil {
		log.Fatalf("Run failed: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("RESULTS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\nFinal Output: %v\n\n", res.FinalOutput)
	fmt.Printf("Generated Items (%d):\n", len(res.NewItems))
	for i, item := range res.NewItems {
		itemJSON, _ := json.MarshalIndent(item, "", "  ")
		fmt.Printf("\n[%d] %s\n", i+1, string(itemJSON))
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("CONTEXT VERIFICATION")
	fmt.Println(strings.Repeat("=", 80))

	var runCtx *runner.RunContext
	if res.RunContext != nil {
		if rc, ok := res.RunContext.(*runner.RunContext); ok {
			runCtx = rc
		}
	}

	if runCtx != nil && runCtx.Context != nil {
		if testCtx, ok := runCtx.Context.(*TestContext); ok {
			if testCtx.UserID != "" {
				fmt.Printf("✓ User ID shared: %s\n", testCtx.UserID)
			} else {
				fmt.Printf("✗ User ID NOT shared\n")
			}

			if testCtx.OrderID != "" {
				fmt.Printf("✓ Order ID shared: %s\n", testCtx.OrderID)
			} else {
				fmt.Printf("✗ Order ID NOT shared\n")
			}

			if testCtx.Metadata["refund_amount"] != nil {
				fmt.Printf("✓ Refund amount calculated: $%.2f\n", testCtx.Metadata["refund_amount"].(float64))
			} else {
				fmt.Printf("✗ Refund amount NOT calculated\n")
			}
		}
	}

	if runCtx != nil && runCtx.Usage != nil {
		fmt.Printf("\nUsage Statistics:\n")
		fmt.Printf("  Requests: %d\n", runCtx.Usage.Requests)
		fmt.Printf("  Input Tokens: %d\n", runCtx.Usage.InputTokens)
		fmt.Printf("  Output Tokens: %d\n", runCtx.Usage.OutputTokens)
		fmt.Printf("  Total Tokens: %d\n", runCtx.Usage.TotalTokens)
	}

	fmt.Printf("\nLast Agent: %s\n", res.LastAgent.Name)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("EXPECTED RESULTS CHECKLIST")
	fmt.Println(strings.Repeat("=", 80))

	checks := []struct {
		name     string
		check    func() bool
		expected string
	}{
		{
			name: "User ID in context",
			check: func() bool {
				if runCtx != nil && runCtx.Context != nil {
					if testCtx, ok := runCtx.Context.(*TestContext); ok {
						return testCtx.UserID == "user_123"
					}
				}
				return false
			},
			expected: "user_123",
		},
		{
			name: "Order ID in context",
			check: func() bool {
				if runCtx != nil && runCtx.Context != nil {
					if testCtx, ok := runCtx.Context.(*TestContext); ok {
						return testCtx.OrderID == "order_456"
					}
				}
				return false
			},
			expected: "order_456",
		},
		{
			name: "Refund amount calculated",
			check: func() bool {
				if runCtx != nil && runCtx.Context != nil {
					if testCtx, ok := runCtx.Context.(*TestContext); ok {
						return testCtx.Metadata["refund_amount"] != nil
					}
				}
				return false
			},
			expected: "~$89.99 (90% of $99.99)",
		},
		{
			name: "Final output contains refund info",
			check: func() bool {
				if res.FinalOutput == nil {
					return false
				}
				outputStr := fmt.Sprintf("%v", res.FinalOutput)
				return strings.Contains(strings.ToLower(outputStr), "refund") ||
					strings.Contains(strings.ToLower(outputStr), "order")
			},
			expected: "Contains refund or order information",
		},
		{
			name: "Multiple agents involved",
			check: func() bool {
				for _, item := range res.NewItems {
					if item.GetType() == "handoff" {
						return true
					}
				}
				return false
			},
			expected: "At least one handoff occurred",
		},
		{
			name: "Tools were called",
			check: func() bool {
				for _, item := range res.NewItems {
					if item.GetType() == "tool_result" {
						return true
					}
				}
				return false
			},
			expected: "At least one tool was executed",
		},
	}

	allPassed := true
	for _, check := range checks {
		passed := check.check()
		status := "✓"
		if !passed {
			status = "✗"
			allPassed = false
		}
		fmt.Printf("%s %s (Expected: %s)\n", status, check.name, check.expected)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	if allPassed {
		fmt.Println("✓ ALL CHECKS PASSED - Agentic flow working correctly!")
	} else {
		fmt.Println("✗ SOME CHECKS FAILED - Review the results above")
	}
	fmt.Println(strings.Repeat("=", 80))
}

// main function for running the test directly with: go run -tags=store store.go
func main() {
	complexAgenticFlowTest()
}
