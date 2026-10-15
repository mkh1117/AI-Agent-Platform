package ai

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockProviderReturnsResponse(t *testing.T) {
	expected := GenerateResponse{
		Message: Message{
			Role:    RoleAssistant,
			Content: "Hello!",
		},
		Provider:     "mock",
		Model:        "mock-model",
		FinishReason: "stop",
		Usage: TokenUsage{
			InputTokens:  10,
			OutputTokens: 5,
			TotalTokens:  15,
		},
	}

	provider := &MockProvider{
		Responses: []GenerateResponse{expected},
	}

	request := GenerateRequest{
		Messages: []Message{
			{Role: RoleUser, Content: "Hi!"},
		},
	}

	got, err := provider.Generate(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Len(t, provider.Calls, 1)
	require.Equal(t, request, provider.Calls[0])
}

func TestMockProviderReturnsResponsesInOrder(t *testing.T) {
	first := GenerateResponse{
		Message: Message{
			Role:    RoleAssistant,
			Content: "First response",
		},
	}

	second := GenerateResponse{
		Message: Message{
			Role:    RoleAssistant,
			Content: "Second response",
		},
	}

	provider := &MockProvider{
		Responses: []GenerateResponse{first, second},
	}

	gotFirst, err := provider.Generate(
		context.Background(),
		GenerateRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, first, gotFirst)

	gotSecond, err := provider.Generate(
		context.Background(),
		GenerateRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, second, gotSecond)
}

func TestMockProviderReturnsConfiguredError(t *testing.T) {
	expectedErr := errors.New("simulated provider failure")

	provider := &MockProvider{
		Err: expectedErr,
	}

	got, err := provider.Generate(
		context.Background(),
		GenerateRequest{},
	)

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, GenerateResponse{}, got)
}

func TestMockProviderWithoutResponses(t *testing.T) {
	provider := &MockProvider{}

	_, err := provider.Generate(
		context.Background(),
		GenerateRequest{},
	)

	require.ErrorIs(t, err, ErrNoMockResponse)
}

func TestMockProviderRespectsCanceledContext(t *testing.T) {
	provider := &MockProvider{
		Responses: []GenerateResponse{
			{
				Message: Message{
					Role:    RoleAssistant,
					Content: "Hello!",
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := provider.Generate(ctx, GenerateRequest{})

	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, provider.Calls)
}
