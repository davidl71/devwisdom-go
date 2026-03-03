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
		// Try loading from environment first (backward compatibility).
	cfg, err := mcpconfig.LoadBaseConfig()
	if err != nil || cfg.Name == "" || cfg.Name == "mcp-server" || cfg.Version == "" || cfg.Version == "1.0.0" {
				// Use builder pattern with defaults.
		cfg, err = mcpconfig.NewConfigBuilder().
			WithName("devwisdom").
			WithVersion(Version).
			Build()
		if err != nil {
						// Fallback to hardcoded defaults if builder fails.
			cfg = &mcpconfig.BaseConfig{
				Name:    "devwisdom",
				Version: Version,
			}
		}
	}

		// Initialize consultation logger (log directory: .devwisdom).
	logger, err := logging.NewConsultationLogger(".devwisdom")
	if err != nil {
				// Log initialization failure is non-fatal - server can still work without logging.
		logger = nil
	}

		// Handle DEVWISDOM_DEBUG for backward compatibility.
	if os.Getenv("DEVWISDOM_DEBUG") == "1" {
		os.Setenv("MCP_DEBUG", "1")
	}
	appLogger := mcplogging.NewLogger()

		// Create SDK server with config.
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
		// Initialize wisdom engine first (before any output).
	if err := s.wisdom.Initialize(); err != nil {
		s.appLogger.Error("", "Failed to initialize wisdom engine: %v", err)
		return fmt.Errorf("failed to initialize wisdom engine (check sources.json configuration and file permissions): %w", err)
	}

		// Log server startup.
	s.appLogger.Info("", "MCP server v%s starting (SDK)", Version)

		// Register tools.
	s.registerTools()

		// Register resources.
	if err := s.registerResources(); err != nil {
		return fmt.Errorf("failed to register resources: %w", err)
	}

		// Run with stdio transport.
	transport := &mcp.StdioTransport{}
	if err := s.server.Run(ctx, transport); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}

// ToolHandlerFunc is a function type for handling tool execution.
type ToolHandlerFunc func(map[string]interface{}) (interface{}, error)

// newErrorResult creates a CallToolResult with an error message.
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


// - Standardized response formatting using mcp-go-core utilities.
func wrapToolHandler(handler ToolHandlerFunc) func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				// Extract and unmarshal arguments.
		args := make(map[string]interface{})
		if req.Params != nil && len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return newErrorResult(fmt.Sprintf("Failed to parse arguments: %v", err)), nil
			}
		}

				// Call the actual handler.
		result, err := handler(args)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Tool execution error: %v", err)), nil
		}

				// Convert result to map[string]interface{} for response.FormatResult().
		resultMap, err := mcpresponse.ConvertToMap(result)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Failed to convert result to map: %v", err)), nil
		}

				// Extract output_path from args if present.
		outputPath := ""
		if path, ok := args["output_path"].(string); ok && path != "" {
			outputPath = path
		}

				// Use mcp-go-core response formatter for standardized formatting.
		textContents, err := mcpresponse.FormatResult(resultMap, outputPath)
		if err != nil {
			return newErrorResult(fmt.Sprintf("Failed to format result: %v", err)), nil
		}

				// Convert []types.TextContent to []mcp.Content ([]mcp.TextContent).
		contents := make([]mcp.Content, len(textContents))
		for i, tc := range textContents {
			contents[i] = &mcp.TextContent{
				Text: tc.Text,
			}
		}

				// Return success result with formatted content.
		return &mcp.CallToolResult{
			Content: contents,
		}, nil
	}
}


// registerTools registers all MCP tools with the SDK server.
func (s *WisdomServerSDK) registerTools() {
		// Create handlers instance to reuse business logic.
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)

		// Get tool definitions from shared function and convert to SDK format.
	toolDefs := getToolDefinitions()
	sdkTools := ConvertToolsToSDK(toolDefs)

		// Map tool names to handler functions.
	handlerMap := map[string]ToolHandlerFunc{
		"consult_advisor":     handlers.handleConsultAdvisor,
		"get_wisdom":          handlers.handleGetWisdom,
		"get_daily_briefing": handlers.handleGetDailyBriefing,
		"get_consultation_log": handlers.handleGetConsultationLog,
	}

		// Register all tools.
	for _, tool := range sdkTools {
		handler, ok := handlerMap[tool.Name]
		if !ok {
			continue // Skip if handler not found (shouldn't happen).
		}
		wrappedHandler := wrapToolHandler(handler)
		s.server.AddTool(tool, wrappedHandler)
	}
}

// registerResources registers all MCP resources with the SDK server.
func (s *WisdomServerSDK) registerResources() error {
		// Create handlers instance to reuse business logic.
	handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)

		// Register wisdom://tools resource.
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

		// For brevity, I'll add a helper function to register all resources.

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

		// Register wisdom://advisor/{id} - use ResourceTemplate for dynamic URI.
	advisorTemplate := &mcp.ResourceTemplate{
		URITemplate: "wisdom://advisor/{id}",
		Name:        "Advisor Details",
		Description: "Get details for a specific advisor",
		MIMEType:    "application/json",
	}
		// Extract advisor ID (string parameter).
	extractAdvisorID := func(suffix string) (string, error) {
		if suffix == "" {
			return "", fmt.Errorf("advisor ID is required")
		}
		return suffix, nil
	}
	advisorTemplateHandler := s.createResourceTemplateHandler(
		"wisdom://advisor/",
		extractAdvisorID,
		handlers.HandleAdvisorResource,
	)
	s.server.AddResourceTemplate(advisorTemplate, advisorTemplateHandler)

		// Register wisdom://consultations/{days} - use ResourceTemplate for dynamic URI.
	consultationsTemplate := &mcp.ResourceTemplate{
		URITemplate: "wisdom://consultations/{days}",
		Name:        "Consultation Log",
		Description: "Get consultation log entries for the specified number of days",
		MIMEType:    "application/json",
	}
		// Extract days (int parameter with default value 7).
	extractDays := func(suffix string) (int, error) {
		if suffix == "" {
			return 7, nil // Use default if empty.
		}
		days, err := strconv.Atoi(suffix)
		if err != nil {
			return 0, fmt.Errorf("invalid days value: %q", suffix)
		}
		return days, nil
	}
		consultationsTemplateHandler := s.createResourceTemplateHandlerInt(
		"wisdom://consultations/",
		extractDays,
		7, // default value.
		handlers.HandleConsultationsResource,
	)
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

// createResourceTemplateHandler creates a resource template handler factory.
// This eliminates duplication between advisor and consultations template handlers.
// Parameters:
//   - uriPrefix: The URI prefix to match (e.g., "wisdom://advisor/")
//   - extractParam: Function to extract parameter from URI suffix (returns string, error)
//   - handlerFunc: Handler function that takes mockReq and extracted parameter
//
// Returns a handler function compatible with SDK's ResourceTemplate handler signature.
func (s *WisdomServerSDK) createResourceTemplateHandler(
	uriPrefix string,
	extractParam func(string) (string, error),
	handlerFunc func(*JSONRPCRequest, string) *JSONRPCResponse,
) func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params == nil || req.Params.URI == "" {
			return nil, fmt.Errorf("resource URI is required")
		}
		uri := req.Params.URI

				// Validate URI prefix.
		if !strings.HasPrefix(uri, uriPrefix) {
			return nil, fmt.Errorf("invalid URI format: expected prefix %q, got %q", uriPrefix, uri)
		}

				// Extract parameter from URI suffix.
		param, err := extractParam(strings.TrimPrefix(uri, uriPrefix))
		if err != nil {
			return nil, fmt.Errorf("failed to extract parameter from URI %q: %w", uri, err)
		}

				// Create mock JSON-RPC request.
		mockReq := &JSONRPCRequest{
			ID:     "resource",
			Method: "resources/read",
			Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
		}

				// Call handler with extracted parameter.
		resp := handlerFunc(mockReq, param)
		return s.convertResourceResponse(resp, uri)
	}
}

// createResourceTemplateHandlerInt creates a resource template handler factory for int parameters.
// Similar to createResourceTemplateHandler but handles int parameters with optional default value.
func (s *WisdomServerSDK) createResourceTemplateHandlerInt(
	uriPrefix string,
	extractParam func(string) (int, error),
	defaultValue int,
	handlerFunc func(*JSONRPCRequest, int) *JSONRPCResponse,
) func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		if req.Params == nil || req.Params.URI == "" {
			return nil, fmt.Errorf("resource URI is required")
		}
		uri := req.Params.URI

				// Validate URI prefix.
		if !strings.HasPrefix(uri, uriPrefix) {
			return nil, fmt.Errorf("invalid URI format: expected prefix %q, got %q", uriPrefix, uri)
		}

				// Extract parameter from URI suffix.
		param, err := extractParam(strings.TrimPrefix(uri, uriPrefix))
		if err != nil {
						// Use default value if extraction fails.
			param = defaultValue
		}

				// Create mock JSON-RPC request.
		mockReq := &JSONRPCRequest{
			ID:     "resource",
			Method: "resources/read",
			Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
		}

				// Call handler with extracted parameter.
		resp := handlerFunc(mockReq, param)
		return s.convertResourceResponse(resp, uri)
	}
}

// convertResourceResponse converts our JSON-RPC response format to SDK format.
func (s *WisdomServerSDK) convertResourceResponse(resp *JSONRPCResponse, uri string) (*mcp.ReadResourceResult, error) {
	if resp.Error != nil {
		return nil, fmt.Errorf("resource read error: %v", resp.Error.Message)
	}

	if result, ok := resp.Result.(map[string]interface{}); ok {
		contentsValue := result["contents"]
		var contentMap map[string]interface{}

				// Try []map[string]interface{} first (actual type from newResourceResponse).
		if contents, ok := contentsValue.([]map[string]interface{}); ok && len(contents) > 0 {
			contentMap = contents[0]
		} else if contents, ok := contentsValue.([]interface{}); ok && len(contents) > 0 {
						// Fallback to []interface{}.
			var ok2 bool
			contentMap, ok2 = contents[0].(map[string]interface{})
			if !ok2 {
				return nil, fmt.Errorf("unexpected response format: content item is not map[string]interface{}")
			}
		} else {
			return nil, fmt.Errorf("unexpected response format: contents is not []map[string]interface{} or []interface{}")
		}

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

	return nil, fmt.Errorf("unexpected response format: result is not map[string]interface{}")
}

// getString safely extracts a string value from a map.
func getString(m map[string]interface{}, key, defaultValue string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultValue
}
