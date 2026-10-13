package ai

import (
	"encoding/json"
	"errors"
	"strings"
)

type Tool interface{
	Name()string
	Execute(t ToolCall)(r ToolResult,e error)
}


type ToolCall struct{
	Name string `json:"name"`
	ID string `json:"id"`
	Arguments json.RawMessage `json:"arguments"`
}

type ToolResult struct{
	ToolCallID string `json:"tool_call_id"`
	Content string `json:"content"`
	Error string `json:"error"`
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