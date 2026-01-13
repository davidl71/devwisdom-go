package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ConvertToolsToSDK converts protocol.Tool slice to SDK mcp.Tool slice.
// This allows getToolDefinitions() to be used by both server.go and sdk_adapter.go.
func ConvertToolsToSDK(tools []Tool) []*mcp.Tool {
	sdkTools := make([]*mcp.Tool, len(tools))
	for i, tool := range tools {
		sdkTools[i] = &mcp.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		}
	}
	return sdkTools
}
