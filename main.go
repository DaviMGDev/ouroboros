package main

import (
	"context"
	"fmt"
)

func main() {
	mock := &MockLLM{}
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

	fmt.Println("Response:", resp.Message.Content)

}
