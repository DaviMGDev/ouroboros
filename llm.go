package main

import (
	"context"
	"fmt"
	"time"
)

// LLM abstracts a language model provider with chat completion,
// single-turn text completion, and streaming chat completion methods.
type LLM interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	Complete(ctx context.Context, prompt string) (string, error)
	// StreamChat returns a ChatStream that yields chunks incrementally.
	// The caller MUST call Close() when finished, regardless of whether
	// iteration completes naturally.
	StreamChat(ctx context.Context, req *ChatRequest) (ChatStream, error)
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

// ChatStream is a streaming iterator over chat completion chunks.
// The caller MUST call Close() when finished, regardless of whether
// iteration completes naturally.
type ChatStream interface {
	// Next advances to the next chunk. Returns false when the stream
	// is exhausted (call Err() to check for errors).
	Next() bool
	// Current returns the most recently yielded chunk. Only valid
	// after a true return from Next().
	Current() ChatChunk
	// Err returns the first error encountered during streaming, if any.
	Err() error
	// Close releases any resources held by the stream. Safe to call
	// more than once.
	Close() error
}

// ChatChunk is one incremental delta from a streaming chat response.
type ChatChunk struct {
	Content     string          `json:"content"`
	Role        MessageRole     `json:"role"`
	ToolCalls   []ToolCallDelta `json:"tool_calls,omitempty"`
	FinishReason FinishReason   `json:"finish_reason,omitempty"`
	Usage       *UsageStats     `json:"usage,omitempty"`
}

// ToolCallDelta carries incremental tool call data for streaming responses.
type ToolCallDelta struct {
	Index    int    `json:"index"`
	ID       string `json:"id,omitempty"`
	Function struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function,omitempty"`
}

// MockLLM is an echo implementation of LLM that returns the user's input back.
// It is useful for unit testing code that depends on the LLM interface.
//
// ChunkDelay, if non-zero, causes MockChatStream to sleep for that duration
// before yielding each chunk, simulating real streaming latency.
type MockLLM struct {
	ChunkDelay time.Duration
}

var _ LLM = (*MockLLM)(nil)

func (m *MockLLM) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("chat request cannot be nil")
	}
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

// MockChatStream is an iterator over pre-built ChatChunks for testing.
type MockChatStream struct {
	chunks     []ChatChunk
	pos        int
	closed     bool
	chunkDelay time.Duration
}

var _ ChatStream = (*MockChatStream)(nil)

func (s *MockChatStream) Next() bool {
	if s.closed {
		return false
	}
	if s.pos < len(s.chunks) {
		s.pos++
		if s.chunkDelay > 0 {
			time.Sleep(s.chunkDelay)
		}
		return true
	}
	return false
}

func (s *MockChatStream) Current() ChatChunk {
	if s.pos == 0 || s.pos > len(s.chunks) {
		return ChatChunk{}
	}
	return s.chunks[s.pos-1]
}

func (s *MockChatStream) Err() error { return nil }

func (s *MockChatStream) Close() error { s.closed = true; return nil }

// StreamChat returns a ChatStream that echoes the last user message content
// as a single content chunk followed by a final done chunk.
func (m *MockLLM) StreamChat(ctx context.Context, req *ChatRequest) (ChatStream, error) {
	if req == nil {
		return nil, fmt.Errorf("chat request cannot be nil")
	}
	content := ""
	if len(req.Messages) > 0 {
		content = req.Messages[len(req.Messages)-1].Content
	}
	chunks := []ChatChunk{
		{
			Content: content,
			Role:    RoleAssistant,
		},
		{
			FinishReason: FinishReasonStop,
			Usage: &UsageStats{
				PromptTokens:     len(content),
				CompletionTokens: len(content),
				TotalTokens:      len(content) * 2,
			},
		},
	}
	return &MockChatStream{chunks: chunks, chunkDelay: m.ChunkDelay}, nil
}
