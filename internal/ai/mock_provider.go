package ai

import (
	"context"
	"errors"
	"sync"
)

var ErrNoMockResponse = errors.New("mock provider: no response configured")

type MockProvider struct {
	mu        sync.Mutex
	Responses []GenerateResponse
	Err       error
	Calls     []GenerateRequest
}

func (p *MockProvider) Generate(ctx context.Context,request GenerateRequest,) (GenerateResponse, error) {
	if err := ctx.Err(); err != nil {
		return GenerateResponse{}, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.Calls = append(p.Calls, request)

	if p.Err != nil {
		return GenerateResponse{}, p.Err
	}

	if len(p.Responses) == 0 {
		return GenerateResponse{}, ErrNoMockResponse
	}

	response := p.Responses[0]
	p.Responses = p.Responses[1:]

	return response, nil
}
