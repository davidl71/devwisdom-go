package mcp

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"
)

// BenchmarkServer_Initialize benchmarks server initialization
func BenchmarkServer_Initialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		server := NewWisdomServer()
		ctx := context.Background()
		stdin := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`)
		stdout := &strings.Builder{}

		// Initialize in background (non-blocking)
		go func() {
			_ = server.Run(ctx, stdin, stdout)
		}()

		// Wait a bit for initialization
		time.Sleep(10 * time.Millisecond)
	}
}

// BenchmarkServer_ToolCall_GetWisdom benchmarks get_wisdom tool call
func BenchmarkServer_ToolCall_GetWisdom(b *testing.B) {
	server := NewWisdomServer()
	ctx := context.Background()

	// Initialize request
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`

	// Tool call request
	toolReq := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(initReq + "\n" + toolReq)
		output := &strings.Builder{}

		// Run in background
		done := make(chan error, 1)
		go func() {
			done <- server.Run(ctx, input, output)
		}()

		// Wait for completion or timeout
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			b.Fatal("Benchmark timeout")
		}
	}
}

// BenchmarkServer_ToolCall_ConsultAdvisor benchmarks consult_advisor tool call
func BenchmarkServer_ToolCall_ConsultAdvisor(b *testing.B) {
	server := NewWisdomServer()
	ctx := context.Background()

	// Initialize request
	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`

	// Tool call request
	toolReq := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"consult_advisor","arguments":{"metric":"security","score":40.0}}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(initReq + "\n" + toolReq)
		output := &strings.Builder{}

		// Run in background
		done := make(chan error, 1)
		go func() {
			done <- server.Run(ctx, input, output)
		}()

		// Wait for completion or timeout
		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			b.Fatal("Benchmark timeout")
		}
	}
}

// BenchmarkHandleToolCall_GetWisdom benchmarks direct tool call handling (without JSON-RPC overhead)
func BenchmarkHandleToolCall_GetWisdom(b *testing.B) {
	server := NewWisdomServer()
	// Initialize engine
	ctx := context.Background()
	stdin := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`)
	stdout := io.Discard
	go func() {
		_ = server.Run(ctx, stdin, stdout)
	}()
	time.Sleep(50 * time.Millisecond) // Wait for initialization

	params := map[string]interface{}{
		"score":  75.0,
		"source": "stoic",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = server.HandleToolCall("get_wisdom", params)
	}
}

// BenchmarkHandleToolCall_ConsultAdvisor benchmarks direct tool call handling (without JSON-RPC overhead)
func BenchmarkHandleToolCall_ConsultAdvisor(b *testing.B) {
	server := NewWisdomServer()
	// Initialize engine
	ctx := context.Background()
	stdin := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`)
	stdout := io.Discard
	go func() {
		_ = server.Run(ctx, stdin, stdout)
	}()
	time.Sleep(50 * time.Millisecond) // Wait for initialization

	params := map[string]interface{}{
		"metric": "security",
		"score":  40.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = server.HandleToolCall("consult_advisor", params)
	}
}

// BenchmarkJSONRPC_Parse benchmarks JSON-RPC request parsing
func BenchmarkJSONRPC_Parse(b *testing.B) {
	jsonStr := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_wisdom","arguments":{"score":75.0,"source":"stoic"}}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var req JSONRPCRequest
		_ = json.Unmarshal([]byte(jsonStr), &req)
	}
}
