package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNilProvider     = errors.New("runner: model provider is nil")
	ErrMaxIterations   = errors.New("runner: maximum iterations reached")
	ErrInvalidToolCall = errors.New("runner: invalid tool call")
	ErrInvalidResponse = errors.New("runner: invalid model response")
)

type RunResult struct {
	Message  Message
	Usage    TokenUsage
	Provider string
	Model    string
}

type Runner struct {
	provider      ModelProvider
	tools         map[string]Tool
	maxIterations int
}

func NewRunner(provider ModelProvider,tools []Tool,maxIterations int,) (*Runner, error) {
	if provider == nil {
		return nil, ErrNilProvider
	}

	if maxIterations <= 0 {
		maxIterations = 10
	}

	toolRegistry := make(map[string]Tool, len(tools))

	for _, tool := range tools {
		if tool == nil {
			continue
		}

		name := strings.TrimSpace(tool.Name())
		if name == "" {
			return nil, errors.New("runner: tool name cannot be empty")
		}

		if _, exists := toolRegistry[name]; exists {
			return nil, fmt.Errorf("runner: duplicate tool %q", name)
		}

		toolRegistry[name] = tool
	}

	return &Runner{
		provider:      provider,
		tools:         toolRegistry,
		maxIterations: maxIterations,
	}, nil
}

func (r *Runner) Run(ctx context.Context,request GenerateRequest,) (RunResult, error) {

	if len(request.Messages) == 0 {
		return RunResult{}, errors.New("runner: messages cannot be empty")
	}

	// Do not modify the caller's underlying message slice.
	messages := append([]Message(nil), request.Messages...)

	toolDefinitions := make([]ToolDefinition, 0, len(r.tools))
	for _, tool := range r.tools {
		toolDefinitions = append(toolDefinitions, ToolDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  tool.Parameters(),
		})
	}

	request.Tools = toolDefinitions

	var totalUsage TokenUsage

	for iteration := 0; iteration < r.maxIterations; iteration++ {
		if err := ctx.Err(); err != nil {
			return RunResult{}, err
		}

		request.Messages = messages

		response, err := r.provider.Generate(ctx, request)
		if err != nil {
			return RunResult{}, fmt.Errorf(
				"runner: model generation failed: %w",
				err,
			)
		}

		totalUsage.InputTokens += response.Usage.InputTokens
		totalUsage.OutputTokens += response.Usage.OutputTokens
		totalUsage.TotalTokens += response.Usage.TotalTokens

		if response.Message.Role != RoleAssistant {
			return RunResult{}, fmt.Errorf(
				"%w: expected assistant message, got %q",
				ErrInvalidResponse,
				response.Message.Role,
			)
		}

		messages = append(messages, response.Message)

		if len(response.Message.ToolCalls) == 0 {
			return RunResult{
				Message:  response.Message,
				Usage:    totalUsage,
				Provider: response.Provider,
				Model:    response.Model,
			}, nil
		}

		for _, call := range response.Message.ToolCalls {
			if err := ctx.Err(); err != nil {
				return RunResult{}, err
			}

			if strings.TrimSpace(call.ID) == "" ||
				strings.TrimSpace(call.Name) == "" ||
				!json.Valid(call.Arguments) {
				return RunResult{}, fmt.Errorf(
					"%w: missing ID/name or invalid arguments",
					ErrInvalidToolCall,
				)
			}

			tool, exists := r.tools[call.Name]

			var content string

			if !exists {
				content = fmt.Sprintf(
					"Tool %q is not available.",
					call.Name,
				)
			} else {
				output, executeErr := tool.Execute(ctx, call.Arguments)
				if executeErr != nil {
					content = fmt.Sprintf(
						"Tool execution failed: %s",
						executeErr.Error(),
					)
				} else {
					content = output
				}
			}

			messages = append(messages, Message{
				Role:       RoleTool,
				Name:       call.Name,
				ToolCallID: call.ID,
				Content:    content,
			})
		}
	}

	return RunResult{}, ErrMaxIterations
}
