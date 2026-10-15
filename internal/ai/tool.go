package ai

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

type Tool interface{
	Name()string
	Description() string
	Parameters() json.RawMessage
	Execute(ctx context.Context,arguments json.RawMessage,) (string, error)
}


type ToolCall struct{
	Name string `json:"name"`
	ID string `json:"id"`
	Arguments json.RawMessage `json:"arguments"`
}

type ToolResult struct{
	ToolCallID string `json:"tool_call_id"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Error      string `json:"error,omitempty"`
}

func (call ToolCall)Validate()error{
	idTrim:=strings.TrimSpace(call.ID)
	nameTrim:=strings.TrimSpace(call.Name)
	argumentCheck:=json.Valid(call.Arguments)

	if idTrim==""{
		return errors.New("tool call id is empty")
	}
	if nameTrim==""{
		return  errors.New("tool name is empty")
	}
	if !argumentCheck{
		return errors.New("tool arguments is invalid or empty")	
	}
	return nil
}

func (call ToolResult)Validate()error{
	callidTrim:=strings.TrimSpace(call.ToolCallID)

	if callidTrim==""{
		return errors.New("tool call id is empty")
	}
	return nil
}