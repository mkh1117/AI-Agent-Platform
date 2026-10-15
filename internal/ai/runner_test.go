package ai

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type CalculatorTool struct{}

func (CalculatorTool) Name() string {
	return "calculator"
}

func (CalculatorTool) Description() string {
	return "Adds two numbers"
}

func (CalculatorTool) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"a": {"type": "integer"},
			"b": {"type": "integer"}
		},
		"required": ["a", "b"]
	}`)
}

func (CalculatorTool) Execute(
	ctx context.Context,
	arguments json.RawMessage,
) (string, error) {
	var args struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	if err := json.Unmarshal(arguments, &args); err != nil {
		return "", err
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(args.A + args.B), nil
}

func TestRunnerExecutesToolAndReturnsFinalAnswer(t *testing.T) {
	provider := &MockProvider{
		Responses: []GenerateResponse{
			{
				Message: Message{
					Role: RoleAssistant,
					ToolCalls: []ToolCall{
						{
							ID:        "call_123",
							Name:      "calculator",
							Arguments: json.RawMessage(`{"a":12,"b":8}`),
						},
					},
				},
				Provider:     "mock",
				Model:        "mock-model",
				FinishReason: "tool_calls",
				Usage: TokenUsage{
					InputTokens:  20,
					OutputTokens: 8,
					TotalTokens:  28,
				},
			},
			{
				Message: Message{
					Role:    RoleAssistant,
					Content: "12 + 8 = 20",
				},
				Provider:     "mock",
				Model:        "mock-model",
				FinishReason: "stop",
				Usage: TokenUsage{
					InputTokens:  30,
					OutputTokens: 6,
					TotalTokens:  36,
				},
			},
		},
	}

	runner, err := NewRunner(
		provider,
		[]Tool{CalculatorTool{}},
		5,
	)
	require.NoError(t, err)

	result, err := runner.Run(
		context.Background(),
		GenerateRequest{
			Messages: []Message{
				{
					Role:    RoleUser,
					Content: "Calculate 12 + 8",
				},
			},
			Config: GenerationConfig{
				Model: "mock-model",
			},
		},
	)

	require.NoError(t, err)
	require.Equal(t, RoleAssistant, result.Message.Role)
	require.Equal(t, "12 + 8 = 20", result.Message.Content)
	require.Equal(t, 50, result.Usage.InputTokens)
	require.Equal(t, 14, result.Usage.OutputTokens)
	require.Equal(t, 64, result.Usage.TotalTokens)

	// The model should have been called twice:
	// first to request the tool, then to produce the final answer.
	require.Len(t, provider.Calls, 2)

	// The second request must contain the tool result.
	secondMessages := provider.Calls[1].Messages
	require.Len(t, secondMessages, 3)
	require.Equal(t, RoleTool, secondMessages[2].Role)
	require.Equal(t, "call_123", secondMessages[2].ToolCallID)
	require.Equal(t, "calculator", secondMessages[2].Name)
	require.Equal(t, "20", secondMessages[2].Content)
}
