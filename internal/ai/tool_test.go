package ai

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockTool struct{}

func TestToolCallJSON(t *testing.T){
	data := []byte(`{"name":"calculator","id":"call_123","arguments":{"a":12,"b":8}}`)

	var example ToolCall

	err:=json.Unmarshal(data,&example)

	require.NoError(t,err)
	
	require.Equal(t, "call_123", example.ID)	
	require.Equal(t, "calculator", example.Name)
	require.JSONEq(t, `{"a":12,"b":8}`, string(example.Arguments))
}
func TestToolResultJSON(t *testing.T){
	data:=[]byte(`{"tool_call_id":"call_123","content":"HI I am an Agent","error":"Null"}`)

	var toJson ToolResult

	err:=json.Unmarshal(data,&toJson)

	require.NoError(t,err)


	st,er:=json.Marshal(toJson)

	require.NoError(t, er)
	require.JSONEq(t, `{"tool_call_id":"call_123","content":"HI I am an Agent","error":"Null"}`, string(st))
	
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
            wantErr: false,
        },
        {
            name: "empty id",
            call: ToolCall{
                ID:        "",
                Name:      "calculator",
                Arguments: json.RawMessage(`{"a":12,"b":8}`),
            },
            wantErr: true,
        },
        {
            name: "empty name",
            call: ToolCall{
                ID:        "call_123",
                Name:      "",
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
            Content:    "Hello from tool",
        }

        err := result.Validate()
        require.NoError(t, err)
    })

    t.Run("empty tool call id", func(t *testing.T) {
        result := ToolResult{
            ToolCallID: "   ",
            Content:    "Hello from tool",
        }

        err := result.Validate()
        require.Error(t, err)
    })
}

func (m MockTool)Name()string{

	return "mock"
}

func(m MockTool)Execute(t ToolCall)(ToolResult,error){
	var result ToolResult
	result.ToolCallID=t.ID
	result.Content="mock result"
	return result,nil
}


func TestMockTool(t *testing.T) {
    var tool Tool = MockTool{}

    require.Equal(t, "mock", tool.Name())

    call := ToolCall{
        ID:        "call_123",
        Name:      "mock",
        Arguments: json.RawMessage(`{}`),
    }

    result, err := tool.Execute(call)

    require.NoError(t, err)
    require.Equal(t, "call_123", result.ToolCallID)
    require.Equal(t, "mock result", result.Content)
}
