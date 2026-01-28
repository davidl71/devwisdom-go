package mcp

import (
	"github.com/davidl71/mcp-go-core/pkg/mcp/protocol"
)

// This maintains backward compatibility while using standardized types.

// JSON-RPC types.
type (
	JSONRPCRequest  = protocol.JSONRPCRequest
	JSONRPCResponse = protocol.JSONRPCResponse
	JSONRPCError    = protocol.JSONRPCError
)

// Error codes.
const (
	ErrCodeParseError     = protocol.ErrCodeParseError
	ErrCodeInvalidRequest = protocol.ErrCodeInvalidRequest
	ErrCodeMethodNotFound = protocol.ErrCodeMethodNotFound
	ErrCodeInvalidParams  = protocol.ErrCodeInvalidParams
	ErrCodeInternalError  = protocol.ErrCodeInternalError
)

// Protocol types.
type (
	InitializeParams    = protocol.InitializeParams
	ClientCapabilities  = protocol.ClientCapabilities
	ClientInfo          = protocol.ClientInfo
	InitializeResult    = protocol.InitializeResult
	ServerCapabilities  = protocol.ServerCapabilities
	ToolsCapability     = protocol.ToolsCapability
	ResourcesCapability = protocol.ResourcesCapability
	ServerInfo          = protocol.ServerInfo
	Tool                = protocol.Tool
	Resource            = protocol.Resource
)

// Helper functions.
var (
	NewErrorResponse       = protocol.NewErrorResponse
	NewSuccessResponse     = protocol.NewSuccessResponse
	NewMethodNotFoundError = protocol.NewMethodNotFoundError
	NewInvalidParamsError  = protocol.NewInvalidParamsError
	NewInternalError       = protocol.NewInternalError
)

// This type is specific to devwisdom-go and not in mcp-go-core.
type ToolCallParams struct {
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Name      string                 `json:"name"`
}

// This type is specific to devwisdom-go and not in mcp-go-core.
type ResourceReadParams struct {
	URI string `json:"uri"`
}
