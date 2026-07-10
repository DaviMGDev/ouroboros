package main

import (
	"context"
)

// LLM abstracts a language model provider with chat completion and
// single-turn text completion methods.
type LLM interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	Complete(ctx context.Context, prompt string) (string, error)
}

// MessageRole identifies the sender of a chat message.
type MessageRole string

const (
	// RoleSystem indicates a system-level instruction message.
	RoleSystem    MessageRole = "system"
	// RoleUser indicates a message from the end user.
	RoleUser      MessageRole = "user"
	// RoleAssistant indicates a message from the AI assistant.
	RoleAssistant MessageRole = "assistant"
)

// Message represents a single message in a chat conversation.
type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

// ChatRequest contains the parameters for a chat completion request.
type ChatRequest struct {
	Messages      []Message `json:"messages"`
	Model         string    `json:"model"`
	Temperature   float64   `json:"temperature"`
	MaxTokens     int       `json:"max_tokens"`
	StopSequences []string  `json:"stop_sequences"`
}

// ChatResponse contains the result of a chat completion call.
type ChatResponse struct {
	Message      Message      `json:"message"`
	Model        string       `json:"model"`
	Usage        UsageStats   `json:"usage"`
	FinishReason FinishReason `json:"finish_reason"`
}

// UsageStats contains token counts for an LLM API call.
type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FinishReason explains why the model stopped generating tokens.
type FinishReason string

const (
	// FinishReasonStop indicates the model finished naturally.
	FinishReasonStop          FinishReason = "stop"
	// FinishReasonLength indicates the response was cut off by max_tokens.
	FinishReasonLength        FinishReason = "length"
	// FinishReasonError indicates an error occurred during generation.
	FinishReasonError         FinishReason = "error"
	// FinishReasonContentFilter indicates the response was flagged by a content filter.
	FinishReasonContentFilter FinishReason = "content_filter"
)

// MockLLM is an echo implementation of LLM that returns the user's input back.
// It is useful for unit testing code that depends on the LLM interface.
type MockLLM struct{}

var _ LLM = (*MockLLM)(nil)

func (m *MockLLM) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	content := ""
	if len(req.Messages) > 0 {
		content = req.Messages[len(req.Messages)-1].Content
	}
	return &ChatResponse{
		Message: Message{
			Role:    RoleAssistant,
			Content: content,
		},
		Model: req.Model,
		Usage: UsageStats{
			PromptTokens:     len(content),
			CompletionTokens: len(content),
			TotalTokens:      len(content) * 2,
		},
		FinishReason: FinishReasonStop,
	}, nil
}

func (m *MockLLM) Complete(ctx context.Context, prompt string) (string, error) {
	return prompt, nil
}
