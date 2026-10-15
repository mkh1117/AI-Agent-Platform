
package ai

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockProvider struct {
	Response GenerateResponse
	Err      error
}

func (p MockProvider) Generate(
	ctx context.Context,
	request GenerateRequest,
) (GenerateResponse, error) {
	if p.Err != nil {
		return GenerateResponse{}, p.Err
	}

	return p.Response, nil
}

func TestMockProvider(t *testing.T) {
	provider := MockProvider{
		Response: GenerateResponse{
			Message: Message{
				Role:    RoleAssistant,
				Content: "Hello!",
			},
			Usage: TokenUsage{
				InputTokens:  10,
				OutputTokens: 5,
				TotalTokens:  15,
			},
			Provider:     "mock",
			Model:        "mock-model",
			FinishReason: "stop",
		},
	}

	var model ModelProvider = provider

	request := GenerateRequest{
		Messages: []Message{
			{
				Role:    RoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    RoleUser,
				Content: "Say hello.",
			},
		},
		Config: GenerationConfig{
			Model: "mock-model",
		},
	}

	response, err := model.Generate(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, RoleAssistant, response.Message.Role)
	require.Equal(t, "Hello!", response.Message.Content)
	require.Equal(t, "mock", response.Provider)
	require.Equal(t, "mock-model", response.Model)
	require.Equal(t, "stop", response.FinishReason)
	require.Equal(t, 10, response.Usage.InputTokens)
	require.Equal(t, 5, response.Usage.OutputTokens)
	require.Equal(t, 15, response.Usage.TotalTokens)
	require.Len(t, request.Messages, 2)
}

func TestMockProviderToolCall(t *testing.T) {
	provider := MockProvider{
		Response: GenerateResponse{
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
			Usage: TokenUsage{
				InputTokens:  20,
				OutputTokens: 8,
				TotalTokens:  28,
			},
			Provider:     "mock",
			Model:        "mock-model",
			FinishReason: "tool_calls",
		},
	}

	response, err := provider.Generate(
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
	require.Equal(t, "tool_calls", response.FinishReason)
	require.Len(t, response.Message.ToolCalls, 1)

	call := response.Message.ToolCalls[0]
	require.Equal(t, "call_123", call.ID)
	require.Equal(t, "calculator", call.Name)
	require.JSONEq(t, `{"a":12,"b":8}`, string(call.Arguments))
	require.Equal(t, 28, response.Usage.TotalTokens)
}

func TestMockProviderError(t *testing.T) {
	provider := MockProvider{
		Err: context.DeadlineExceeded,
	}

	_, err := provider.Generate(
		context.Background(),
		GenerateRequest{},
	)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}
