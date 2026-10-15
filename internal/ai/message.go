package ai

import (
	"errors"
	"strings"
)


type Role string

const (
	RoleSystem Role="system"
	RoleUser Role="user"
	RoleAssistant Role="assistant"
	RoleTool      Role = "tool"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

func (m Message) Validate() error {
	switch m.Role {
	case RoleSystem, RoleUser:
		if strings.TrimSpace(m.Content) == "" {
			return errors.New("message content is empty")
		}

	case RoleAssistant:
		// پاسخ مدل می‌تواند متن یا درخواست اجرای ابزار باشد.
		if strings.TrimSpace(m.Content) == "" &&
			len(m.ToolCalls) == 0 {
			return errors.New(
				"assistant message must contain content or tool calls",
			)
		}

		for _, call := range m.ToolCalls {
			if err := call.Validate(); err != nil {
				return err
			}
		}

	case RoleTool:
		if strings.TrimSpace(m.ToolCallID) == "" {
			return errors.New("tool message call ID is empty")
		}

	default:
		return errors.New("invalid message role")
	}

	return nil
}