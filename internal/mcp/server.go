// Package mcp provides the Model Context Protocol (MCP) server implementation.
// It implements JSON-RPC 2.0 protocol over stdio transport for wisdom consultations.
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
func (s *WisdomServer) HandleToolCall(name string, params map[string]interface{}) (interface{}, error) {
	switch name {
	case "consult_advisor":
		return s.handleConsultAdvisor(params)
	case "get_wisdom":
		return s.handleGetWisdom(params)
	case "get_daily_briefing":
		return s.handleGetDailyBriefing(params)
	case "get_consultation_log":
		return s.handleGetConsultationLog(params)
	default:
		availableTools := []string{"consult_advisor", "get_wisdom", "get_daily_briefing", "get_consultation_log"}
		return nil, fmt.Errorf("unknown tool %q (available tools: %v). Check tool name spelling", name, availableTools)
	}
}

// handleConsultAdvisor implements consult_advisor tool
func (s *WisdomServer) handleConsultAdvisor(params map[string]interface{}) (interface{}, error) {
	// Extract parameters
	var metric, tool, stage, context string
	var score float64

	if m, ok := params["metric"].(string); ok {
		metric = m
	}
	if t, ok := params["tool"].(string); ok {
		tool = t
	}
	if st, ok := params["stage"].(string); ok {
		stage = st
	}
	if c, ok := params["context"].(string); ok {
		context = c
	}
	if sc, ok := params["score"].(float64); ok {
		score = sc
	} else if sc, ok := params["score"].(int); ok {
		score = float64(sc)
	}
	// Validate and clamp score to 0-100 range
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Determine advisor based on metric, tool, or stage
	var advisorInfo *wisdom.AdvisorInfo
	var err error

	if metric != "" {
		advisorInfo, err = s.wisdom.GetAdvisors().GetAdvisorForMetric(metric)
	} else if tool != "" {
		advisorInfo, err = s.wisdom.GetAdvisors().GetAdvisorForTool(tool)
	} else if stage != "" {
		advisorInfo, err = s.wisdom.GetAdvisors().GetAdvisorForStage(stage)
	} else {
		// Default advisor
		advisorInfo = &wisdom.AdvisorInfo{
			Advisor:   "pistis_sophia",
			Icon:      "📜",
			Rationale: "Default wisdom advisor",
		}
	}

	if err != nil {
		// Fallback to default
		advisorInfo = &wisdom.AdvisorInfo{
			Advisor:   "pistis_sophia",
			Icon:      "📜",
			Rationale: "Default wisdom advisor",
		}
	}

	// Get wisdom quote
	quote, err := s.wisdom.GetWisdom(score, advisorInfo.Advisor)
	if err != nil {
		// Fallback quote
		quote = &wisdom.Quote{
			Quote:         "Wisdom comes from experience.",
			Source:        "Unknown",
			Encouragement: "Keep learning and growing.",
		}
	}

	// Get consultation mode based on score
	modeConfig := wisdom.GetConsultationMode(score)

	// Create consultation
	consultation := wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          advisorInfo.Advisor,
		AdvisorIcon:      advisorInfo.Icon,
		AdvisorName:      advisorInfo.Advisor,
		Rationale:        advisorInfo.Rationale,
		ScoreAtTime:      score,
		ConsultationMode: modeConfig.Name,
		ModeIcon:         modeConfig.Icon,
		ModeFrequency:    modeConfig.Frequency,
		ModeGuidance:     modeConfig.Description,
		Quote:            quote.Quote,
		QuoteSource:      quote.Source,
		Encouragement:    quote.Encouragement,
		Context:          context,
	}

	if metric != "" {
		consultation.Metric = metric
	}
	if tool != "" {
		consultation.Tool = tool
	}
	if stage != "" {
		consultation.Stage = stage
	}

	// Log consultation (non-blocking - if logging fails, still return consultation)
	if s.logger != nil {
		if err := s.logger.Log(&consultation); err != nil {
			// Logging failure doesn't break the consultation response
			// In production, you might want to log this error
			_ = err // Explicitly ignore error for graceful degradation
		}
	}

	return consultation, nil
}

// handleGetWisdom implements get_wisdom tool
func (s *WisdomServer) handleGetWisdom(params map[string]interface{}) (interface{}, error) {
	// Extract score
	var score float64
	if sc, ok := params["score"].(float64); ok {
		score = sc
	} else if sc, ok := params["score"].(int); ok {
		score = float64(sc)
	} else {
		return nil, fmt.Errorf("score parameter is required and must be a number between 0-100")
	}
	// Validate and clamp score to 0-100 range
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Extract source (optional)
	source := "pistis_sophia" // Default
	if src, ok := params["source"].(string); ok && src != "" {
		source = src
	}

		// Get wisdom quote
		quote, err := s.wisdom.GetWisdom(score, source)
		if err != nil {
			return nil, fmt.Errorf("failed to get wisdom quote (source: %q, score: %.1f): %w", source, score, err)
		}

	return quote, nil
}

// handleGetDailyBriefing implements get_daily_briefing tool
func (s *WisdomServer) handleGetDailyBriefing(params map[string]interface{}) (interface{}, error) {
	var score float64
	if sc, ok := params["score"].(float64); ok {
		score = sc
	} else if sc, ok := params["score"].(int); ok {
		score = float64(sc)
	}
	// Default to 50 if score not provided or invalid type
	// Validate and clamp score to 0-100 range
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Get multiple quotes from different sources
	sources := s.wisdom.ListSources()
	briefing := map[string]interface{}{
		"date":    time.Now().Format("2006-01-02"),
		"score":   score,
		"quotes":  []interface{}{},
		"sources": sources,
	}

	// Get quotes from a few sources
	selectedSources := []string{"pistis_sophia", "stoic", "tao"}
	if len(sources) > 0 {
		selectedSources = sources
		if len(selectedSources) > 3 {
			selectedSources = selectedSources[:3]
		}
	}

	quotes := []interface{}{}
	for _, src := range selectedSources {
		quote, err := s.wisdom.GetWisdom(score, src)
		if err == nil {
			quotes = append(quotes, quote)
		}
	}

	briefing["quotes"] = quotes
	return briefing, nil
}

// handleGetConsultationLog implements get_consultation_log tool
func (s *WisdomServer) handleGetConsultationLog(params map[string]interface{}) (interface{}, error) {
	days := 7 // Default
	if d, ok := params["days"].(float64); ok {
		days = int(d)
	} else if d, ok := params["days"].(int); ok {
		days = d
	}

	// Retrieve consultations from logger
	if s.logger == nil {
		// Logger not initialized, return empty array
		return []interface{}{}, nil
	}

	consultations, err := s.logger.GetLogs(days)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve consultation log for last %d days (check log file permissions and .devwisdom directory): %w", days, err)
	}

	// Convert to interface{} slice for JSON serialization
	result := make([]interface{}, len(consultations))
	for i, consultation := range consultations {
		result[i] = consultation
	}

	return result, nil
}

// Resource handlers

// handleToolsResource returns all available tools
func (s *WisdomServer) handleToolsResource(req *JSONRPCRequest) *JSONRPCResponse {
	tools := []map[string]interface{}{
		{
			"name":        "consult_advisor",
			"description": "Consult a wisdom advisor based on metric, tool, or stage",
			"parameters": map[string]interface{}{
				"metric":  "Metric name (e.g., 'security', 'testing')",
				"tool":    "Tool name (e.g., 'project_scorecard')",
				"stage":   "Stage name (e.g., 'daily_checkin')",
				"score":   "Project health score (0-100)",
				"context": "Additional context for the consultation",
			},
		},
		{
			"name":        "get_wisdom",
			"description": "Get a wisdom quote based on project health score and source",
			"parameters": map[string]interface{}{
				"score":  "Project health score (0-100) - required",
				"source": "Wisdom source ID (e.g., 'pistis_sophia', 'stoic') or 'random' for date-seeded random selection - optional",
			},
		},
		{
			"name":        "get_daily_briefing",
			"description": "Get a daily wisdom briefing with quotes and guidance",
			"parameters": map[string]interface{}{
				"score": "Project health score (0-100) - optional",
			},
		},
		{
			"name":        "get_consultation_log",
			"description": "Retrieve consultation log entries",
			"parameters": map[string]interface{}{
				"days": "Number of days to retrieve (default: 7)",
			},
		},
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "wisdom://tools",
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(tools)),
			},
		},
	})
}

// handleSourcesResource returns all wisdom sources
func (s *WisdomServer) handleSourcesResource(req *JSONRPCRequest) *JSONRPCResponse {
	sources := s.wisdom.ListSources()
	sourceDetails := make([]map[string]interface{}, 0, len(sources))

	for _, id := range sources {
		source, found := s.wisdom.GetSource(id)
		if found {
			sourceDetails = append(sourceDetails, map[string]interface{}{
				"id":          id,
				"name":        source.Name,
				"icon":        source.Icon,
				"description": source.Description,
			})
		}
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "wisdom://sources",
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(sourceDetails)),
			},
		},
	})
}

// handleAdvisorsResource returns all advisors
func (s *WisdomServer) handleAdvisorsResource(req *JSONRPCRequest) *JSONRPCResponse {
	advisorRegistry := s.wisdom.GetAdvisors()

	// Build comprehensive advisor listing
	advisorList := map[string]interface{}{
		"metric_advisors": make([]map[string]interface{}, 0),
		"tool_advisors":   make([]map[string]interface{}, 0),
		"stage_advisors":  make([]map[string]interface{}, 0),
	}

	// Get all metric advisors
	metricAdvisors := advisorRegistry.GetAllMetricAdvisors()
	metricList := make([]map[string]interface{}, 0, len(metricAdvisors))
	for metric, info := range metricAdvisors {
		item := map[string]interface{}{
			"metric":    metric,
			"advisor":   info.Advisor,
			"icon":      info.Icon,
			"rationale": info.Rationale,
		}
		if info.HelpsWith != "" {
			item["helps_with"] = info.HelpsWith
		}
		if info.Language != "" {
			item["language"] = info.Language
		}
		metricList = append(metricList, item)
	}
	advisorList["metric_advisors"] = metricList

	// Get all tool advisors
	toolAdvisors := advisorRegistry.GetAllToolAdvisors()
	toolList := make([]map[string]interface{}, 0, len(toolAdvisors))
	for tool, info := range toolAdvisors {
		item := map[string]interface{}{
			"tool":      tool,
			"advisor":   info.Advisor,
			"rationale": info.Rationale,
		}
		if info.Language != "" {
			item["language"] = info.Language
		}
		toolList = append(toolList, item)
	}
	advisorList["tool_advisors"] = toolList

	// Get all stage advisors
	stageAdvisors := advisorRegistry.GetAllStageAdvisors()
	stageList := make([]map[string]interface{}, 0, len(stageAdvisors))
	for stage, info := range stageAdvisors {
		item := map[string]interface{}{
			"stage":     stage,
			"advisor":   info.Advisor,
			"icon":      info.Icon,
			"rationale": info.Rationale,
		}
		if info.Language != "" {
			item["language"] = info.Language
		}
		stageList = append(stageList, item)
	}
	advisorList["stage_advisors"] = stageList

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "wisdom://advisors",
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(advisorList)),
			},
		},
	})
}

// handleAdvisorResource returns a specific advisor
func (s *WisdomServer) handleAdvisorResource(req *JSONRPCRequest, advisorID string) *JSONRPCResponse {
	advisorRegistry := s.wisdom.GetAdvisors()

	// Try to find advisor in metric, tool, or stage advisors
	var advisorInfo *wisdom.AdvisorInfo
	var advisorType string

	// Try metric advisors first
	if info, err := advisorRegistry.GetAdvisorForMetric(advisorID); err == nil {
		advisorInfo = info
		advisorType = "metric"
	} else if info, err := advisorRegistry.GetAdvisorForTool(advisorID); err == nil {
		// Try tool advisors
		advisorInfo = info
		advisorType = "tool"
	} else if info, err := advisorRegistry.GetAdvisorForStage(advisorID); err == nil {
		// Try stage advisors
		advisorInfo = info
		advisorType = "stage"
	} else {
		// Advisor not found
		return NewErrorResponse(req.ID, ErrCodeInvalidParams, fmt.Sprintf("advisor not found: %q. Use 'wisdom://advisors' resource to list available advisors", advisorID), nil)
	}

	// Build advisor response
	advisor := map[string]interface{}{
		"id":        advisorID,
		"type":      advisorType,
		"advisor":   advisorInfo.Advisor,
		"rationale": advisorInfo.Rationale,
	}
	if advisorInfo.Icon != "" {
		advisor["icon"] = advisorInfo.Icon
	}
	if advisorInfo.HelpsWith != "" {
		advisor["helps_with"] = advisorInfo.HelpsWith
	}
	if advisorInfo.Language != "" {
		advisor["language"] = advisorInfo.Language
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "wisdom://advisor/" + advisorID,
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(advisor)),
			},
		},
	})
}

// handleConsultationsResource returns consultation log entries
func (s *WisdomServer) handleConsultationsResource(req *JSONRPCRequest, days int) *JSONRPCResponse {
	// Retrieve consultations from logger
	var consultations []interface{}
	if s.logger != nil {
		logs, err := s.logger.GetLogs(days)
		if err == nil {
			// Convert to interface{} slice for JSON serialization
			consultations = make([]interface{}, len(logs))
			for i, consultation := range logs {
				consultations[i] = consultation
			}
		}
		// If logger is nil or error occurs, consultations remains empty array
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      fmt.Sprintf("wisdom://consultations/%d", days),
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(consultations)),
			},
		},
	})
}

// Helper functions

// mustMarshalJSONCompact marshals to compact JSON (no indentation)
// Used for embedding JSON strings in resource responses
func mustMarshalJSONCompact(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		// Don't log to stderr in MCP server - it breaks stdio protocol
		// Return empty JSON object on error
		return []byte("{}")
	}
	return data
}
