package ai

import (
	"context"
)


type GenerateRequest struct {
	Messages []Message       `json:"messages"`
	Config   GenerationConfig `json:"config"`
	Tools    []ToolDefinition `json:"tools,omitempty"`
}

type GenerationConfig struct {
	Model           string   `json:"model"`
	Temperature     *float64 `json:"temperature,omitempty"`
	MaxOutputTokens int      `json:"max_output_tokens,omitempty"`
}

type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  []byte `json:"parameters"`
}

type GenerateResponse struct {
	Message      Message    `json:"message"`
	Usage        TokenUsage `json:"usage"`
	Provider     string     `json:"provider"`
	Model        string     `json:"model"`
	FinishReason string     `json:"finish_reason"`
}

type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type ModelProvider interface{
	Generate(ctx context.Context,request GenerateRequest)(GenerateResponse,error)
}


