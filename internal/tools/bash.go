package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/DaviMGDev/thoth-agent/internal/llm"
)

var _ llm.Tool = (*BashTool)(nil)

// BashTool executes a bash command and returns its stdout and stderr.
type BashTool struct {
	// MaxOutput limits the combined stdout+stderr output size.
	// Defaults to 10,000 bytes if zero.
	MaxOutput int
}

func (t *BashTool) Name() string { return "bash" }

func (t *BashTool) Description() string {
	return "Execute a bash command and return its output. Use this for running shell commands, scripts, or checking system state. Returns stdout and stderr combined."
}

func (t *BashTool) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The bash command to execute.",
			},
		},
		"required": []string{"command"},
	}
}

func (t *BashTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	cmd, ok := args["command"]
	if !ok {
		return "", fmt.Errorf("missing required argument: command")
	}
	cmdStr, ok := cmd.(string)
	if !ok {
		return "", fmt.Errorf("command must be a string, got %T", cmd)
	}

	// Use context-aware command execution
	c := exec.CommandContext(ctx, "bash", "-c", cmdStr)

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()

	var output strings.Builder
	if stdout.Len() > 0 {
		output.WriteString(stdout.String())
	}
	if stderr.Len() > 0 {
		if output.Len() > 0 {
			output.WriteString("\n")
		}
		output.WriteString(stderr.String())
	}

	// Truncate if needed
	maxLen := t.MaxOutput
	if maxLen <= 0 {
		maxLen = 10_000
	}

	result := output.String()
	if len(result) > maxLen {
		result = result[:maxLen] + fmt.Sprintf("\n\n... [truncated, %d more bytes]", len(result)-maxLen)
	}

	if err != nil {
		// Include exit error info but still return the output
		if exitErr, ok := err.(*exec.ExitError); ok {
			return result + fmt.Sprintf("\n[exit code: %d]", exitErr.ExitCode()), nil
		}
		return result + fmt.Sprintf("\n[error: %v]", err), nil
	}

	if result == "" {
		return "(no output)", nil
	}

	return result, nil
}
