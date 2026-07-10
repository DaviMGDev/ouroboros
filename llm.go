package main

import (
	"context"
)

type LLM interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	Complete(ctx context.Context, prompt string) (string, error)
}

type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

type ChatRequest struct {
	Messages      []Message `json:"messages"`
	Model         string    `json:"model"`
	Temperature   float64   `json:"temperature"`
	MaxTokens     int       `json:"max_tokens"`
	StopSequences []string  `json:"stop_sequences"`
}

type ChatResponse struct {
	Message      Message      `json:"messages"`
	Model        string       `json:"model"`
	Usage        UsageStats   `json:"usage"`
	FinishReason FinishReason `json:"finish_reason"`
}

type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonError         FinishReason = "error"
	FinishReasonContentFilter FinishReason = "content_filter"
)

type MockLLM struct{}

func (m *MockLLM) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	return &ChatResponse{
		Message: Message{
			Role:    RoleAssistant,
			Content: "This is a mock response.",
		},
		Model: "mock-model",
		Usage: UsageStats{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
		FinishReason: FinishReasonStop,
	}, nil
}
