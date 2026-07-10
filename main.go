package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func main() {
	mock := &MockLLM{
		ChunkDelay: 160 * time.Millisecond,
	}

	// Blocking chat example
	{
		req := &ChatRequest{
			Messages: []Message{
				{Role: RoleUser, Content: "Hello, how are you?"},
			},
			Model:       "mock-model",
			Temperature: 0.7,
			MaxTokens:   100,
		}

		resp, err := mock.Chat(context.Background(), req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Chat Response:", resp.Message.Content)
	}

	// Streaming chat example
	{
		req := &ChatRequest{
			Messages: []Message{
				{Role: RoleUser, Content: "Hello from streaming!"},
			},
			Model: "mock-model",
		}

		stream, err := mock.StreamChat(context.Background(), req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer stream.Close()

		var full strings.Builder
		for stream.Next() {
			chunk := stream.Current()
			fmt.Print(chunk.Content)
			full.WriteString(chunk.Content)
			if chunk.FinishReason != "" {
				fmt.Printf("\n[%s] tokens: %d\n", chunk.FinishReason, chunk.Usage.TotalTokens)
			}
		}
		if err := stream.Err(); err != nil {
			fmt.Println("\nStream error:", err)
			return
		}
		fmt.Println("\nStreamed complete:", full.String())
	}
}
