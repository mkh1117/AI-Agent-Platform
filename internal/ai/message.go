package ai

import (
	"errors"
	"strings"
)


type Role string

const (
	RoleSystem Role="system"
	RoleUser Role="user"
	RoleAssistant="assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func (message Message) Validate() error {

	if message.Role != RoleSystem && message.Role != RoleUser && message.Role != RoleAssistant {
		return errors.New("Role (who messages) Not Validate")
	}


	messageTrim:=strings.TrimSpace(message.Content)
	if messageTrim ==""{
		return errors.New("Message is empty")
	}


	return nil
}