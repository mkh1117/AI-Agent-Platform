
package ai

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockTool struct{}

func (m MockTool) Name() string {
	return "mock"
}

func (m MockTool) Description() string {
	return "A mock tool for testing"
}

func (m MockTool) Parameters() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {}
	}`)
}

func (m MockTool) Execute(
	ctx context.Context,
	arguments json.RawMessage,
) (string, error) {
	return "mock result", nil
}

func TestToolCallJSON(t *testing.T) {
	data := []byte(`{
		"name": "calculator",
		"id": "call_123",
		"arguments": {"a": 12, "b": 8}
	}`)

	var call ToolCall

	err := json.Unmarshal(data, &call)

	require.NoError(t, err)
	require.Equal(t, "call_123", call.ID)
	require.Equal(t, "calculator", call.Name)
	require.JSONEq(t, `{"a":12,"b":8}`, string(call.Arguments))
}

func TestToolResultJSON(t *testing.T) {
	result := ToolResult{
		ToolCallID: "call_123",
		Name:       "calculator",
		Content:    "20",
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	var decoded ToolResult
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	require.Equal(t, result.ToolCallID, decoded.ToolCallID)
	require.Equal(t, result.Name, decoded.Name)
	require.Equal(t, result.Content, decoded.Content)
	require.Empty(t, decoded.Error)
}

func TestToolCallValidate(t *testing.T) {
	tests := []struct {
		name    string
		call    ToolCall
		wantErr bool
	}{
		{
			name: "valid tool call",
			call: ToolCall{
				ID:        "call_123",
				Name:      "calculator",
				Arguments: json.RawMessage(`{"a":12,"b":8}`),
			},
		},
		{
			name: "empty id",
			call: ToolCall{
				Name:      "calculator",
				Arguments: json.RawMessage(`{"a":12,"b":8}`),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			call: ToolCall{
				ID:        "call_123",
				Arguments: json.RawMessage(`{"a":12,"b":8}`),
			},
			wantErr: true,
		},
		{
			name: "invalid arguments",
			call: ToolCall{
				ID:        "call_123",
				Name:      "calculator",
				Arguments: json.RawMessage(`{"a":}`),
			},
			wantErr: true,
		},
		{
			name: "empty arguments",
			call: ToolCall{
				ID:        "call_123",
				Name:      "calculator",
				Arguments: nil,
			},
			wantErr: true,
		},
		{
			name: "whitespace id",
			call: ToolCall{
				ID:        "   ",
				Name:      "calculator",
				Arguments: json.RawMessage(`{}`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call.Validate()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToolResultValidate(t *testing.T) {
	t.Run("valid tool result", func(t *testing.T) {
		result := ToolResult{
			ToolCallID: "call_123",
			Name:       "mock",
			Content:    "Hello from tool",
		}

		require.NoError(t, result.Validate())
	})

	t.Run("empty tool call id", func(t *testing.T) {
		result := ToolResult{
			ToolCallID: "   ",
			Name:       "mock",
			Content:    "Hello from tool",
		}

		require.Error(t, result.Validate())
	})
}

func TestMockTool(t *testing.T) {
	var tool Tool = MockTool{}

	require.Equal(t, "mock", tool.Name())
	require.NotEmpty(t, tool.Description())
	require.True(t, json.Valid(tool.Parameters()))

	call := ToolCall{
		ID:        "call_123",
		Name:      "mock",
		Arguments: json.RawMessage(`{}`),
	}

	require.NoError(t, call.Validate())

	content, err := tool.Execute(
		context.Background(),
		call.Arguments,
	)

	require.NoError(t, err)
	require.Equal(t, "mock result", content)
}
