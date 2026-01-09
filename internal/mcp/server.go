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
	"strconv"
	"strings"
	"time"

	"github.com/davidl71/devwisdom-go/internal/logging"
	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// Version is the devwisdom-go MCP server version.
const Version = "0.1.0"

// WisdomServer implements the MCP server for wisdom tools and resources.
// It handles JSON-RPC 2.0 requests and provides tools and resources for wisdom access.
type WisdomServer struct {
	wisdom      *wisdom.Engine
	logger      *logging.ConsultationLogger
	appLogger   *logging.Logger // Structured logger for application logging
	initialized bool
}

// NewWisdomServer creates a new wisdom MCP server instance.
// The server must be started with Run() to begin processing requests.
func NewWisdomServer() *WisdomServer {
	// Initialize consultation logger (log directory: .devwisdom)
	logger, err := logging.NewConsultationLogger(".devwisdom")
	if err != nil {
		// Log initialization failure is non-fatal - server can still work without logging
		// In production, you might want to log this to stderr or handle it differently
		logger = nil
	}

	// Initialize structured application logger
	appLogger := logging.NewLogger()

	return &WisdomServer{
		wisdom:    wisdom.NewEngine(),
		logger:    logger,
		appLogger: appLogger,
	}
}

// Run starts the MCP server with stdio transport
func (s *WisdomServer) Run(ctx context.Context, stdin io.Reader, stdout io.Writer) error {
	// Initialize wisdom engine first (before any output)
	if err := s.wisdom.Initialize(); err != nil {
		s.appLogger.Error("", "Failed to initialize wisdom engine: %v", err)
		return fmt.Errorf("failed to initialize wisdom engine (check sources.json configuration and file permissions): %w", err)
	}

	// Log server startup
	s.appLogger.Info("", "MCP server v%s starting", Version)

	// Set up JSON-RPC 2.0 handlers
	decoder := json.NewDecoder(stdin)
	encoder := json.NewEncoder(stdout)
	// Use compact JSON (no indentation) for better compatibility with MCP clients
	// Some clients have issues parsing indented JSON over stdio
	encoder.SetIndent("", "") // Explicitly set to compact (no indentation)

	// Process messages
	for {
		var req JSONRPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				s.appLogger.Info("", "EOF received, shutting down")
				break
			}
			// Send parse error (id must be null for parse errors per JSON-RPC 2.0 spec)
			parseErrMsg := fmt.Sprintf("JSON parse error: invalid JSON-RPC request format (%v). Ensure request is valid JSON and follows JSON-RPC 2.0 specification", err)
			s.appLogger.Error("", "JSON parse error: %v", err)
			resp := NewErrorResponse(nil, ErrCodeParseError, parseErrMsg, nil)
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("failed to send parse error response to client: %w", err)
			}
			// After sending parse error, break to avoid infinite loop on invalid input
			// The decoder can't recover from parse errors, so we must exit
			break
		}

		// Handle request
		// Skip notifications (requests without id) - per JSON-RPC 2.0 spec
		if req.ID == nil {
			// Notifications don't get responses, just continue
			s.appLogger.Debug("", "Received notification (no ID): %s", req.Method)
			continue
		}

		// Log request start and measure duration
		requestID := formatRequestID(req.ID)
		startTime := time.Now()
		s.appLogger.LogRequest(requestID, req.Method)

		resp := s.handleRequest(&req)

		// Log request completion with duration
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

// formatRequestID converts a JSON-RPC request ID to a string for logging
func formatRequestID(id interface{}) string {
	if id == nil {
		return "null"
	}
	switch v := id.(type) {
	case string:
		return v
	case float64:
		// JSON numbers are decoded as float64
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", id)
	}
}

// handleRequest processes a JSON-RPC request
func (s *WisdomServer) handleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		return NewErrorResponse(req.ID, ErrCodeInvalidRequest, fmt.Sprintf("Invalid JSON-RPC version: expected \"2.0\", got %q. Ensure client is using JSON-RPC 2.0 protocol", req.JSONRPC), nil)
	}

	// Handle different methods
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

// handleInitialize handles the initialize request
func (s *WisdomServer) handleInitialize(req *JSONRPCRequest) *JSONRPCResponse {
	var params InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsError(req.ID, fmt.Sprintf("Invalid initialize params: %v", err))
	}

	s.initialized = true

	result := InitializeResult{
		ProtocolVersion: "2024-11-05", // MCP protocol version
		Capabilities: ServerCapabilities{
			Tools:     &ToolsCapability{},
			Resources: &ResourcesCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "devwisdom",
			Version: Version,
		},
	}

	return NewSuccessResponse(req.ID, result)
}

// handleToolsList returns the list of available tools
func (s *WisdomServer) handleToolsList(req *JSONRPCRequest) *JSONRPCResponse {
	tools := []Tool{
		{
			Name:        "consult_advisor",
			Description: "Consult a wisdom advisor based on metric, tool, or stage",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"metric": map[string]interface{}{
						"type":        "string",
						"description": "Metric name (e.g., 'security', 'testing')",
					},
					"tool": map[string]interface{}{
						"type":        "string",
						"description": "Tool name (e.g., 'project_scorecard')",
					},
					"stage": map[string]interface{}{
						"type":        "string",
						"description": "Stage name (e.g., 'daily_checkin')",
					},
					"score": map[string]interface{}{
						"type":        "number",
						"description": "Project health score (0-100)",
					},
					"context": map[string]interface{}{
						"type":        "string",
						"description": "Additional context for the consultation",
					},
				},
			},
		},
		{
			Name:        "get_wisdom",
			Description: "Get a wisdom quote based on project health score and source",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"score": map[string]interface{}{
						"type":        "number",
						"description": "Project health score (0-100)",
						"required":    true,
					},
					"source": map[string]interface{}{
						"type":        "string",
						"description": "Wisdom source ID (e.g., 'pistis_sophia', 'stoic') or 'random' for date-seeded random selection",
					},
				},
				"required": []string{"score"},
			},
		},
		{
			Name:        "get_daily_briefing",
			Description: "Get a daily wisdom briefing with quotes and guidance",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"score": map[string]interface{}{
						"type":        "number",
						"description": "Project health score (0-100)",
					},
				},
			},
		},
		{
			Name:        "get_consultation_log",
			Description: "Retrieve consultation log entries",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"days": map[string]interface{}{
						"type":        "number",
						"description": "Number of days to retrieve (default: 7)",
					},
				},
			},
		},
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"tools": tools,
	})
}

// handleToolCall processes a tool call request
func (s *WisdomServer) handleToolCall(req *JSONRPCRequest) *JSONRPCResponse {
	var params ToolCallParams
	requestID := formatRequestID(req.ID)

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.appLogger.LogError(requestID, "Tool call params parse", err)
		return NewInvalidParamsError(req.ID, fmt.Sprintf("invalid tool call params (failed to parse JSON from request): %v. Ensure params is valid JSON object", err))
	}

	// Log tool call start and measure duration
	startTime := time.Now()
	s.appLogger.LogToolCall(requestID, params.Name, params.Arguments)

	result, err := s.HandleToolCall(params.Name, params.Arguments)

	// Log tool call completion
	duration := time.Since(startTime)
	s.appLogger.LogToolCallComplete(requestID, params.Name, duration)

	if err != nil {
		s.appLogger.LogError(requestID, fmt.Sprintf("Tool call: %s", params.Name), err)
		return NewInternalError(req.ID, fmt.Sprintf("tool call %q failed with arguments %v: %v. Check tool parameters and ensure wisdom engine is initialized", params.Name, params.Arguments, err))
	}

	return NewSuccessResponse(req.ID, result)
}

// handleResourcesList returns the list of available resources
func (s *WisdomServer) handleResourcesList(req *JSONRPCRequest) *JSONRPCResponse {
	resources := []Resource{
		{
			URI:         "wisdom://tools",
			Name:        "Available Tools",
			Description: "List all available MCP tools with descriptions and parameters",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://sources",
			Name:        "Wisdom Sources",
			Description: "List all available wisdom sources",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://advisors",
			Name:        "Wisdom Advisors",
			Description: "List all available advisors",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://advisor/{id}",
			Name:        "Advisor Details",
			Description: "Get details for a specific advisor",
			MimeType:    "application/json",
		},
		{
			URI:         "wisdom://consultations/{days}",
			Name:        "Consultation Log",
			Description: "Get consultation log entries for the specified number of days",
			MimeType:    "application/json",
		},
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"resources": resources,
	})
}

// handleResourceRead reads a resource
func (s *WisdomServer) handleResourceRead(req *JSONRPCRequest) *JSONRPCResponse {
	var params ResourceReadParams
	requestID := formatRequestID(req.ID)

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.appLogger.LogError(requestID, "Resource read params parse", err)
		return NewInvalidParamsError(req.ID, fmt.Sprintf("invalid resource read params (failed to parse JSON from request): %v. Ensure params contains valid 'uri' field", err))
	}

	// Log resource read
	startTime := time.Now()
	s.appLogger.Debug(requestID, "Reading resource: %s", params.URI)

	// Parse resource URI
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

	// Log resource read completion
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

// DEPRECATED: Handler methods moved to handlers.go. Keeping for reference only.
// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleConsultAdvisor implements consult_advisor tool
func (s *WisdomServer) handleConsultAdvisor(params map[string]interface{}) (interface{}, error) {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.handleConsultAdvisor(params)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleGetWisdom implements get_wisdom tool
func (s *WisdomServer) handleGetWisdom(params map[string]interface{}) (interface{}, error) {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.handleGetWisdom(params)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleGetDailyBriefing implements get_daily_briefing tool
func (s *WisdomServer) handleGetDailyBriefing(params map[string]interface{}) (interface{}, error) {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.handleGetDailyBriefing(params)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleGetConsultationLog implements get_consultation_log tool
func (s *WisdomServer) handleGetConsultationLog(params map[string]interface{}) (interface{}, error) {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.handleGetConsultationLog(params)
}

// Resource handlers

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleToolsResource returns all available tools
func (s *WisdomServer) handleToolsResource(req *JSONRPCRequest) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleToolsResource(req)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleSourcesResource returns all wisdom sources
func (s *WisdomServer) handleSourcesResource(req *JSONRPCRequest) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleSourcesResource(req)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleAdvisorsResource returns all advisors
func (s *WisdomServer) handleAdvisorsResource(req *JSONRPCRequest) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleAdvisorsResource(req)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleAdvisorResource returns a specific advisor
func (s *WisdomServer) handleAdvisorResource(req *JSONRPCRequest, advisorID string) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleAdvisorResource(req, advisorID)
}

// DEPRECATED: Handler methods moved to handlers.go. This delegates to handlers.go.
// handleConsultationsResource returns consultation log entries
func (s *WisdomServer) handleConsultationsResource(req *JSONRPCRequest, days int) *JSONRPCResponse {
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
	return handlers.HandleConsultationsResource(req, days)
}

// Helper functions
// (formatRequestID is defined earlier in this file for backward compatibility)
