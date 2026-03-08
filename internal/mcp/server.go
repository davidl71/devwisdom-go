// Package mcp provides the Model Context Protocol (MCP) server implementation.
//
// DEPRECATED: This file contains legacy code from the custom JSON-RPC 2.0 implementation.
// Handler logic has been moved to handlers.go and is used by the SDK adapter (sdk_adapter.go).
//
// This file is kept for:
// - WisdomServer struct (used by tests and SDK adapter for backward compatibility)
// - NewWisdomServer() function (used by tests)
//
// All handler methods have been moved to handlers.go (WisdomHandlers).
// All JSON-RPC protocol code has been removed (now handled by official SDK).
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davidl71/devwisdom-go/internal/logging"
	"github.com/davidl71/devwisdom-go/internal/wisdom"

	mcpconfig "github.com/davidl71/mcp-go-core/pkg/mcp/config"
	mcplogging "github.com/davidl71/mcp-go-core/pkg/mcp/logging"
	"github.com/davidl71/mcp-go-core/pkg/mcp/protocol"
)

// Version is the devwisdom-go MCP server version (default).
const Version = "0.1.0"

// getServerName returns the server name from config or default.
func getServerName() string {
		// Try loading from environment first (backward compatibility).
	if cfg, err := mcpconfig.LoadBaseConfig(); err == nil && cfg.Name != "" && cfg.Name != "mcp-server" {
		return cfg.Name
	}

		// Use builder pattern with defaults.
	cfg, err := mcpconfig.NewConfigBuilder().
		WithName("devwisdom").
		Build()
	if err == nil {
		return cfg.Name
	}

	return "devwisdom"
}

// getServerVersion returns the server version from config or default.
func getServerVersion() string {
		// Try loading from environment first (backward compatibility).
	if cfg, err := mcpconfig.LoadBaseConfig(); err == nil && cfg.Version != "" && cfg.Version != "1.0.0" {
		return cfg.Version
	}

		// Use builder pattern with defaults.
	cfg, err := mcpconfig.NewConfigBuilder().
		WithVersion(Version).
		Build()
	if err == nil {
		return cfg.Version
	}

	return Version
}

// WisdomServer implements the MCP server for wisdom tools and resources.
// It handles JSON-RPC 2.0 requests and provides tools and resources for wisdom access.
type WisdomServer struct {
	wisdom      *wisdom.Engine
	logger      *logging.ConsultationLogger
	appLogger   *mcplogging.Logger // Structured logger for application logging.
	initialized bool
}

// NewWisdomServer creates a new wisdom MCP server instance.
// The server must be started with Run() to begin processing requests.
func NewWisdomServer() *WisdomServer {
		// Initialize consultation logger (log directory: .devwisdom).
	logger, err := logging.NewConsultationLogger(".devwisdom")
	if err != nil {
				// In production, you might want to log this to stderr or handle it differently.
		logger = nil
	}

		// Handle DEVWISDOM_DEBUG for backward compatibility.
	if os.Getenv("DEVWISDOM_DEBUG") == "1" {
		os.Setenv("MCP_DEBUG", "1")
	}
	appLogger := mcplogging.NewLogger()

	return &WisdomServer{
		wisdom:    wisdom.NewEngine(),
		logger:    logger,
		appLogger: appLogger,
	}
}

// Run starts the MCP server with stdio transport.
func (s *WisdomServer) Run(ctx context.Context, stdin io.Reader, stdout io.Writer) error {
		// Initialize wisdom engine first (before any output).
	if err := s.wisdom.Initialize(); err != nil {
		s.appLogger.Error("", "Failed to initialize wisdom engine: %v", err)
		return fmt.Errorf("failed to initialize wisdom engine (check sources.json configuration and file permissions): %w", err)
	}

		// Log server startup.
	s.appLogger.Info("", "MCP server v%s starting", Version)

		// Set up JSON-RPC 2.0 handlers.
	decoder := json.NewDecoder(stdin)
	encoder := json.NewEncoder(stdout)
		// Some clients have issues parsing indented JSON over stdio.
	encoder.SetIndent("", "") // Explicitly set to compact (no indentation).

		// Process messages.
	for {
		var req JSONRPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				s.appLogger.Info("", "EOF received, shutting down")
				break
			}
						// Send parse error (id must be null for parse errors per JSON-RPC 2.0 spec).
			parseErrMsg := fmt.Sprintf("JSON parse error: invalid JSON-RPC request format (%v). Ensure request is valid JSON and follows JSON-RPC 2.0 specification", err)
			s.appLogger.Error("", "JSON parse error: %v", err)
			resp := NewErrorResponse(nil, ErrCodeParseError, parseErrMsg, nil)
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("failed to send parse error response to client: %w", err)
			}
						// The decoder can't recover from parse errors, so we must exit.
			break
		}

				// Skip notifications (requests without id) - per JSON-RPC 2.0 spec.
		if req.ID == nil {
						// Notifications don't get responses, just continue.
			s.appLogger.Debug("", "Received notification (no ID): %s", req.Method)
			continue
		}

				// Log request start and measure duration.
		requestID := protocol.FormatRequestID(req.ID)
		startTime := time.Now()
		s.appLogger.LogRequest(requestID, req.Method)

		resp := s.handleRequest(&req)

				// Log request completion with duration.
		duration := time.Since(startTime)
		s.appLogger.LogRequestComplete(requestID, req.Method, duration)

		if resp != nil {
			if err := encoder.Encode(resp); err != nil {
				s.appLogger.Error(requestID, "Failed to encode response: %v", err)
				return fmt.Errorf("failed to encode JSON-RPC response (method: %q, id: %v): %w", req.Method, req.ID, err)
			}
		}
	}

	s.appLogger.Info("", "MCP server shutting down")
	return nil
}

// formatRequestID converts a JSON-RPC request ID to a string for logging.
// Deprecated: use protocol.FormatRequestID from mcp-go-core instead.
var formatRequestID = protocol.FormatRequestID

// handleRequest processes a JSON-RPC request.
func (s *WisdomServer) handleRequest(req *JSONRPCRequest) *JSONRPCResponse {
		// Validate JSON-RPC version.
	if req.JSONRPC != "2.0" {
		return NewErrorResponse(req.ID, ErrCodeInvalidRequest, fmt.Sprintf("Invalid JSON-RPC version: expected \"2.0\", got %q. Ensure client is using JSON-RPC 2.0 protocol", req.JSONRPC), nil)
	}

		// Handle different methods.
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolCall(req)
	case "resources/list":
		return s.handleResourcesList(req)
	case "resources/read":
		return s.handleResourceRead(req)
	default:
		return NewMethodNotFoundError(req.ID, req.Method)
	}
}

// handleInitialize handles the initialize request.
func (s *WisdomServer) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	var params InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsError(req.ID, fmt.Sprintf("Invalid initialize params: %v", err))
	}

	s.initialized = true

	result := InitializeResult{
		ProtocolVersion: "2024-11-05", // MCP protocol version.
		Capabilities: ServerCapabilities{
			Tools:     &ToolsCapability{},
			Resources: &ResourcesCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    getServerName(),
			Version: getServerVersion(),
		},
	}

	return NewSuccessResponse(req.ID, result)
}

// handleToolsList returns the list of available tools.
func (s *WisdomServer) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	tools := getToolDefinitions()

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"tools": tools,
	})
}

// handleToolCall processes a tool call request.
func (s *WisdomServer) handleToolCall(req *JSONRPCRequest) *JSONRPCResponse {
	var params ToolCallParams
	requestID := protocol.FormatRequestID(req.ID)

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.appLogger.LogError(requestID, "Tool call params parse", err)
		return NewInvalidParamsError(req.ID, fmt.Sprintf("invalid tool call params (failed to parse JSON from request): %v. Ensure params is valid JSON object", err))
	}

		// Log tool call start and measure duration.
	startTime := time.Now()
	s.appLogger.LogToolCall(requestID, params.Name, params.Arguments)

	result, err := s.HandleToolCall(params.Name, params.Arguments)

		// Log tool call completion.
	duration := time.Since(startTime)
	s.appLogger.LogToolCallComplete(requestID, params.Name, duration)

	if err != nil {
		s.appLogger.LogError(requestID, fmt.Sprintf("Tool call: %s", params.Name), err)
		return NewInternalError(req.ID, fmt.Sprintf("tool call %q failed with arguments %v: %v. Check tool parameters and ensure wisdom engine is initialized", params.Name, params.Arguments, err))
	}

	return NewSuccessResponse(req.ID, result)
}

// handleResourcesList returns the list of available resources.
func (s *WisdomServer) handleResourcesList(req *JSONRPCRequest) *JSONRPCResponse {
	resources := GetResourceList()

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"resources": resources,
	})
}

// handleResourceRead reads a resource.
func (s *WisdomServer) handleResourceRead(req *JSONRPCRequest) *JSONRPCResponse {
	var params ResourceReadParams
	requestID := protocol.FormatRequestID(req.ID)

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.appLogger.LogError(requestID, "Resource read params parse", err)
		return NewInvalidParamsError(req.ID, fmt.Sprintf("invalid resource read params (failed to parse JSON from request): %v. Ensure params contains valid 'uri' field", err))
	}

		// Log resource read.
	startTime := time.Now()
	s.appLogger.Debug(requestID, "Reading resource: %s", params.URI)

		// Parse resource URI.
	uri := params.URI
	var resp *JSONRPCResponse

	if uri == "wisdom://tools" {
		resp = s.handleToolsResource(req)
	} else if strings.HasPrefix(uri, "wisdom://sources") {
		resp = s.handleSourcesResource(req)
	} else if uri == "wisdom://advisors" {
		resp = s.handleAdvisorsResource(req)
	} else if strings.HasPrefix(uri, "wisdom://advisor/") {
		// Handle wisdom://advisor/{id}
		parts := strings.Split(uri, "/")
		if len(parts) >= 3 {
			advisorID := parts[len(parts)-1]
			resp = s.handleAdvisorResource(req, advisorID)
		} else {
			resp = NewInvalidParamsError(req.ID, fmt.Sprintf("invalid advisor resource URI: expected format 'wisdom://advisor/{id}', got %q", uri))
		}
	} else if strings.HasPrefix(uri, "wisdom://consultations/") {
		parts := strings.Split(uri, "/")
		if len(parts) >= 3 {
			daysStr := parts[len(parts)-1]
			days, err := strconv.Atoi(daysStr)
			if err != nil {
				resp = NewInvalidParamsError(req.ID, fmt.Sprintf("invalid days parameter %q in URI: must be a number (got %q)", daysStr, uri))
			} else {
				resp = s.handleConsultationsResource(req, days)
			}
		} else {
			resp = NewInvalidParamsError(req.ID, fmt.Sprintf("invalid consultations resource URI: expected format 'wisdom://consultations/{days}', got %q", uri))
		}
	} else {
		resp = NewErrorResponse(req.ID, -32602, fmt.Sprintf("unknown resource URI %q. Use 'wisdom://sources', 'wisdom://advisors', 'wisdom://advisor/{id}', or 'wisdom://consultations/{days}'", uri), nil)
	}

		// Log resource read completion.
	duration := time.Since(startTime)
	s.appLogger.LogPerformance(requestID, fmt.Sprintf("Resource read: %s", uri), duration)

	return resp
}

// HandleToolCall processes MCP tool calls
// DEPRECATED: Delegates to WisdomHandlers. Kept for backward compatibility with tests.
func (s *WisdomServer) HandleToolCall(name string, params map[string]interface{}) (interface{}, error) {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleToolCall(name, params)
}

// handleConsultAdvisor implements consult_advisor tool.
func (s *WisdomServer) handleConsultAdvisor(params map[string]interface{}) (interface{}, error) {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.handleConsultAdvisor(params)
}

// Resource handlers.

// handleToolsResource returns all available tools.
func (s *WisdomServer) handleToolsResource(req *JSONRPCRequest) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleToolsResource(req)
}

// handleSourcesResource returns all wisdom sources.
func (s *WisdomServer) handleSourcesResource(req *JSONRPCRequest) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleSourcesResource(req)
}

// handleAdvisorsResource returns all advisors.
func (s *WisdomServer) handleAdvisorsResource(req *JSONRPCRequest) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleAdvisorsResource(req)
}

// handleAdvisorResource returns a specific advisor.
func (s *WisdomServer) handleAdvisorResource(req *JSONRPCRequest, advisorID string) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleAdvisorResource(req, advisorID)
}

// handleConsultationsResource returns consultation log entries.
func (s *WisdomServer) handleConsultationsResource(req *JSONRPCRequest, days int) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleConsultationsResource(req, days)
}

// (formatRequestID is defined earlier in this file for backward compatibility).