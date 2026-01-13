// Package mcp provides the Model Context Protocol (MCP) server implementation.
// This file contains the SDK-based adapter (new implementation).
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/davidl71/devwisdom-go/internal/logging"
	"github.com/davidl71/devwisdom-go/internal/wisdom"

	mcpconfig "github.com/davidl71/mcp-go-core/pkg/mcp/config"
	mcplogging "github.com/davidl71/mcp-go-core/pkg/mcp/logging"
	mcpresponse "github.com/davidl71/mcp-go-core/pkg/mcp/response"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WisdomServerSDK implements the MCP server using the official SDK.
// It wraps the SDK server and integrates with the wisdom engine.
type WisdomServerSDK struct {
	server    *mcp.Server
	wisdom    *wisdom.Engine
	logger    *logging.ConsultationLogger
	appLogger *mcplogging.Logger
}

// NewWisdomServerSDK creates a new wisdom MCP server instance using the official SDK.
func NewWisdomServerSDK() *WisdomServerSDK {
	// Load server configuration using builder pattern (with environment variable support)
	// Try loading from environment first (backward compatibility)
	cfg, err := mcpconfig.LoadBaseConfig()
	if err != nil || cfg.Name == "" || cfg.Name == "mcp-server" || cfg.Version == "" || cfg.Version == "1.0.0" {
		// Use builder pattern with defaults
		cfg, err = mcpconfig.NewConfigBuilder().
			WithName("devwisdom").
			WithVersion(Version).
			Build()
		if err != nil {
			// Fallback to hardcoded defaults if builder fails
			cfg = &mcpconfig.BaseConfig{
				Name:    "devwisdom",
				Version: Version,
			}
		}
	}

	// Initialize consultation logger (log directory: .devwisdom)
	logger, err := logging.NewConsultationLogger(".devwisdom")
	if err != nil {
		// Log initialization failure is non-fatal - server can still work without logging
		logger = nil
	}

	// Initialize structured application logger (supports both DEVWISDOM_DEBUG and MCP_DEBUG)
	// Handle DEVWISDOM_DEBUG for backward compatibility
	if os.Getenv("DEVWISDOM_DEBUG") == "1" {
		os.Setenv("MCP_DEBUG", "1")
	}
	appLogger := mcplogging.NewLogger()

	// Create SDK server with config
	sdkServer := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.Name,
		Version: cfg.Version,
	}, nil)

	return &WisdomServerSDK{
		server:    sdkServer,
		wisdom:    wisdom.NewEngine(),
		logger:    logger,
		appLogger: appLogger,
	}
}

// Run starts the MCP server with stdio transport using the SDK.
func (s *WisdomServerSDK) Run(ctx context.Context) error {
	// Initialize wisdom engine first (before any output)
	if err := s.wisdom.Initialize(); err != nil {
		s.appLogger.Error("", "Failed to initialize wisdom engine: %v", err)
		return fmt.Errorf("failed to initialize wisdom engine (check sources.json configuration and file permissions): %w", err)
	}

	// Log server startup
	s.appLogger.Info("", "MCP server v%s starting (SDK)", Version)

	// Register tools
	s.registerTools()

	// Register resources
	if err := s.registerResources(); err != nil {
		return fmt.Errorf("failed to register resources: %w", err)
	}

	// Run with stdio transport
	transport := &mcp.StdioTransport{}
	if err := s.server.Run(ctx, transport); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}

// ToolHandlerFunc is a function type for handling tool execution
type ToolHandlerFunc func(map[string]interface{}) (interface{}, error)

// newErrorResult creates a CallToolResult with an error message
func newErrorResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: message,
			},
		},
	}
}


// wrapToolHandler wraps a tool handler function to work with SDK's CallToolRequest
// This eliminates duplication across all tool handlers by providing:
// - Argument unmarshaling
// - Error handling and result marshaling
// - Consistent error response format
// - Standardized response formatting using mcp-go-core utilities
func wrapToolHandler(handler ToolHandlerFunc) func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract and unmarshal arguments
		args := make(map[string]interface{})
		if req.Params != nil && len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return newErrorResult(fmt.Sprintf("Failed to parse arguments: %v", err)), nil
			}
		}

		// Call the actual handler
		result, err := handler(args)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Tool execution error: %v", err)), nil
		}

		// Convert result to map[string]interface{} for response.FormatResult()
		resultMap, err := mcpresponse.ConvertToMap(result)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Failed to convert result to map: %v", err)), nil
		}

		// Extract output_path from args if present
		outputPath := ""
		if path, ok := args["output_path"].(string); ok && path != "" {
			outputPath = path
		}

		// Use mcp-go-core response formatter for standardized formatting
		textContents, err := mcpresponse.FormatResult(resultMap, outputPath)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

		// Convert []types.TextContent to []mcp.Content ([]mcp.TextContent)
		contents := make([]mcp.Content, len(textContents))
		for i, tc := range textContents {
			contents[i] = &mcp.TextContent{
				Text: tc.Text,
			}
		}

		// Return success result with formatted content
		return &mcp.CallToolResult{
			Content: contents,
		}, nil
	}
}


// registerTools registers all MCP tools with the SDK server.
func (s *WisdomServerSDK) registerTools() {
	// Create handlers instance to reuse business logic
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)

	// Register consult_advisor tool
	consultAdvisorTool := &mcp.Tool{
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
	}

	consultAdvisorHandler := wrapToolHandler(handlers.handleConsultAdvisor)

	s.server.AddTool(consultAdvisorTool, consultAdvisorHandler)

	// Register get_wisdom tool
	getWisdomTool := &mcp.Tool{
		Name:        "get_wisdom",
		Description: "Get a wisdom quote based on project health score and source",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"score": map[string]interface{}{
					"type":        "number",
					"description": "Project health score (0-100)",
				},
				"source": map[string]interface{}{
					"type":        "string",
					"description": "Wisdom source ID (e.g., 'pistis_sophia', 'stoic') or 'random' for date-seeded random selection",
				},
			},
			"required": []string{"score"},
		},
	}

	getWisdomHandler := wrapToolHandler(handlers.handleGetWisdom)

	s.server.AddTool(getWisdomTool, getWisdomHandler)

	// Register get_daily_briefing tool
	getDailyBriefingTool := &mcp.Tool{
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
	}

	getDailyBriefingHandler := wrapToolHandler(handlers.handleGetDailyBriefing)

	s.server.AddTool(getDailyBriefingTool, getDailyBriefingHandler)

	// Register get_consultation_log tool
	getConsultationLogTool := &mcp.Tool{
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
	}

	getConsultationLogHandler := wrapToolHandler(handlers.handleGetConsultationLog)

	s.server.AddTool(getConsultationLogTool, getConsultationLogHandler)

	return nil
}

// registerResources registers all MCP resources with the SDK server.
func (s *WisdomServerSDK) registerResources() error {
	// Create handlers instance to reuse business logic
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)

	// Register wisdom://tools resource
	toolsResource := &mcp.Resource{
		URI:         "wisdom://tools",
		Name:        "Available Tools",
		Description: "List all available MCP tools with descriptions and parameters",
		MIMEType:    "application/json",
	}

	toolsHandler := func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		mockReq := &JSONRPCRequest{
			ID:     "resource-tools",
			Method: "resources/read",
			Params: json.RawMessage(`{"uri": "wisdom://tools"}`),
		}
		resp := handlers.HandleToolsResource(mockReq)
		return s.convertResourceResponse(resp, "wisdom://tools")
	}

	s.server.AddResource(toolsResource, toolsHandler)

	// Register other resources similarly...
	// (wisdom://sources, wisdom://advisors, wisdom://advisor/{id}, wisdom://consultations/{days})
	// For brevity, I'll add a helper function to register all resources

	return s.registerAllResources(handlers)
}

// registerAllResources registers all wisdom resources.
func (s *WisdomServerSDK) registerAllResources(handlers *WisdomHandlers) error {
	// Register wisdom://sources
	sourcesResource := &mcp.Resource{
		URI:         "wisdom://sources",
		Name:        "Wisdom Sources",
		Description: "List all available wisdom sources",
		MIMEType:    "application/json",
	}
	sourcesHandler := s.createResourceHandler("wisdom://sources", handlers.HandleSourcesResource)
	s.server.AddResource(sourcesResource, sourcesHandler)

	// Register wisdom://advisors
	advisorsResource := &mcp.Resource{
		URI:         "wisdom://advisors",
		Name:        "Wisdom Advisors",
		Description: "List all available advisors",
		MIMEType:    "application/json",
	}
	advisorsHandler := s.createResourceHandler("wisdom://advisors", handlers.HandleAdvisorsResource)
	s.server.AddResource(advisorsResource, advisorsHandler)

	// Register wisdom://advisor/{id} - use ResourceTemplate for dynamic URI
	advisorTemplate := &mcp.ResourceTemplate{
		URITemplate: "wisdom://advisor/{id}",
		Name:        "Advisor Details",
		Description: "Get details for a specific advisor",
		MIMEType:    "application/json",
	}
	advisorTemplateHandler := func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params == nil || req.Params.URI == "" {
			return nil, fmt.Errorf("resource URI is required")
		}
		uri := req.Params.URI

		// Extract advisor ID from URI (wisdom://advisor/{id})
		if !strings.HasPrefix(uri, "wisdom://advisor/") {
			return nil, fmt.Errorf("invalid advisor URI format: %s", uri)
		}
		advisorID := strings.TrimPrefix(uri, "wisdom://advisor/")

		mockReq := &JSONRPCRequest{
			ID:     "resource",
			Method: "resources/read",
			Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
		}

		resp := handlers.HandleAdvisorResource(mockReq, advisorID)
		return s.convertResourceResponse(resp, uri)
	}
	s.server.AddResourceTemplate(advisorTemplate, advisorTemplateHandler)

	// Register wisdom://consultations/{days} - use ResourceTemplate for dynamic URI
	consultationsTemplate := &mcp.ResourceTemplate{
		URITemplate: "wisdom://consultations/{days}",
		Name:        "Consultation Log",
		Description: "Get consultation log entries for the specified number of days",
		MIMEType:    "application/json",
	}
	consultationsTemplateHandler := func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params == nil || req.Params.URI == "" {
			return nil, fmt.Errorf("resource URI is required")
		}
		uri := req.Params.URI

		// Extract days from URI (wisdom://consultations/{days})
		if !strings.HasPrefix(uri, "wisdom://consultations/") {
			return nil, fmt.Errorf("invalid consultations URI format: %s", uri)
		}
		daysStr := strings.TrimPrefix(uri, "wisdom://consultations/")
		days := 7 // default
		if daysStr != "" {
			if d, err := strconv.Atoi(daysStr); err == nil {
				days = d
			}
		}

		mockReq := &JSONRPCRequest{
			ID:     "resource",
			Method: "resources/read",
			Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
		}

		resp := handlers.HandleConsultationsResource(mockReq, days)
		return s.convertResourceResponse(resp, uri)
	}
	s.server.AddResourceTemplate(consultationsTemplate, consultationsTemplateHandler)

	return nil
}

// createResourceHandler creates a resource handler that converts SDK requests to handler format.
func (s *WisdomServerSDK) createResourceHandler(uri string, handlerFunc func(*JSONRPCRequest) *JSONRPCResponse) func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		mockReq := &JSONRPCRequest{
			ID:     "resource",
			Method: "resources/read",
			Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
		}

		resp := handlerFunc(mockReq)
		return s.convertResourceResponse(resp, uri)
	}
}

// convertResourceResponse converts our JSON-RPC response format to SDK format.
func (s *WisdomServerSDK) convertResourceResponse(resp *JSONRPCResponse, uri string) (*mcp.ReadResourceResult, error) {
	if resp.Error != nil {
		return nil, fmt.Errorf("resource read error: %v", resp.Error.Message)
	}

	if result, ok := resp.Result.(map[string]interface{}); ok {
		if contents, ok := result["contents"].([]interface{}); ok && len(contents) > 0 {
			// Convert first content item
			if contentMap, ok := contents[0].(map[string]interface{}); ok {
				return &mcp.ReadResourceResult{
					Contents: []*mcp.ResourceContents{
						{
							URI:      getString(contentMap, "uri", uri),
							MIMEType: getString(contentMap, "mimeType", "application/json"),
							Text:     getString(contentMap, "text", ""),
						},
					},
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("unexpected response format")
}

// getString safely extracts a string value from a map.
func getString(m map[string]interface{}, key, defaultValue string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultValue
}
