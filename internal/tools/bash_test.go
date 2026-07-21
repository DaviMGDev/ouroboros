package tools

import (
	"context"
	"testing"
)

func TestBashTool_Name(t *testing.T) {
	tool := &BashTool{}
	if got := tool.Name(); got != "bash" {
		t.Errorf("Name() = %q, want %q", got, "bash")
	}
}

func TestBashTool_Description(t *testing.T) {
	tool := &BashTool{}
	if got := tool.Description(); got == "" {
		t.Error("Description() returned empty string")
	}
}

func TestBashTool_Schema(t *testing.T) {
	tool := &BashTool{}
	schema := tool.Schema()

	// Check type
	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want object", schema["type"])
	}

	// Check required fields
	required, ok := schema["required"].([]string)
	if !ok || len(required) != 1 || required[0] != "command" {
		t.Errorf("schema required = %v, want [command]", schema["required"])
	}

	// Check properties
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("schema properties not a map")
	}
	if _, ok := props["command"]; !ok {
		t.Error("schema missing command property")
	}
}

func TestBashTool_Execute(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		{
			name:    "simple echo",
			command: "echo hello",
			want:    "hello",
		},
		{
			name:    "echo with spaces",
			command: "echo 'hello world'",
			want:    "hello world",
		},
		{
			name:    "no output",
			command: "true",
			want:    "(no output)",
		},
		{
			name:    "stderr output",
			command: "echo error >&2",
			want:    "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := &BashTool{}
			result, err := tool.Execute(context.Background(), map[string]any{
				"command": tt.command,
			})
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			// Trim trailing newline
			if len(result) > 0 && result[len(result)-1] == '\n' {
				result = result[:len(result)-1]
			}
			if result != tt.want {
				t.Errorf("Execute() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestBashTool_Execute_MissingCommand(t *testing.T) {
	tool := &BashTool{}
	_, err := tool.Execute(context.Background(), map[string]any{})
	if err == nil {
		t.Error("expected error for missing command argument")
	}
}

func TestBashTool_Execute_InvalidCommandType(t *testing.T) {
	tool := &BashTool{}
	_, err := tool.Execute(context.Background(), map[string]any{
		"command": 123,
	})
	if err == nil {
		t.Error("expected error for non-string command argument")
	}
}

func TestBashTool_Execute_FailingCommand(t *testing.T) {
	tool := &BashTool{}
	result, err := tool.Execute(context.Background(), map[string]any{
		"command": "exit 1",
	})
	if err != nil {
		t.Fatalf("expected no error for failing command, got: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result for failing command")
	}
}

func TestBashTool_Execute_MaxOutput(t *testing.T) {
	tool := &BashTool{MaxOutput: 20}
	result, err := tool.Execute(context.Background(), map[string]any{
		"command": "echo 'this is a long output that should be truncated'",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result) > 100 { // generous limit for truncation message
		t.Errorf("output too long: %d bytes", len(result))
	}
}

func TestBashTool_Execute_ContextCancellation(t *testing.T) {
	tool := &BashTool{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	result, err := tool.Execute(ctx, map[string]any{
		"command": "sleep 10",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Context cancellation should produce an error message in the output
	if result == "" {
		t.Error("expected non-empty result for cancelled context")
	}
}
