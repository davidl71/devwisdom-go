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

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// Version is the devwisdom-go MCP server version
const Version = "0.1.0"

// WisdomServer implements the MCP server for wisdom tools and resources
type WisdomServer struct {
	wisdom      *wisdom.Engine
	initialized bool
}

// NewWisdomServer creates a new wisdom MCP server instance
func NewWisdomServer() *WisdomServer {
	return &WisdomServer{
		wisdom: wisdom.NewEngine(),
	}
}

// Run starts the MCP server with stdio transport
func (s *WisdomServer) Run(ctx context.Context, stdin io.Reader, stdout io.Writer) error {
	// Initialize wisdom engine first (before any output)
	if err := s.wisdom.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine: %w", err)
	}

	// Print version to stderr for debugging (after initialization, before JSON-RPC loop)
	// Use fmt.Fprintf to stderr explicitly - this should not interfere with stdout JSON-RPC
	// However, some MCP clients merge stderr with stdout, so we make this conditional
	// Only print if DEBUG environment variable is set to avoid breaking clients that merge streams
	if os.Getenv("DEVWISDOM_DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, "devwisdom-go MCP server v%s starting...\n", Version)
	}

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
				break
			}
			// Send parse error (id must be null for parse errors per JSON-RPC 2.0 spec)
			resp := NewErrorResponse(nil, ErrCodeParseError, "Parse error", nil)
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("failed to send parse error: %w", err)
			}
			// After sending parse error, break to avoid infinite loop on invalid input
			// The decoder can't recover from parse errors, so we must exit
			break
		}

		// Handle request
		// Skip notifications (requests without id) - per JSON-RPC 2.0 spec
		if req.ID == nil {
			// Notifications don't get responses, just continue
			continue
		}

		resp := s.handleRequest(&req)
		if resp != nil {
			if err := encoder.Encode(resp); err != nil {
				return fmt.Errorf("failed to encode response: %w", err)
			}
		}
	}

	return nil
}

// handleRequest processes a JSON-RPC request
func (s *WisdomServer) handleRequest(req *JSONRPCRequest) *JSONRPCResponse {
	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		return NewErrorResponse(req.ID, ErrCodeInvalidRequest, "Invalid JSON-RPC version", nil)
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
						"description": "Wisdom source ID (e.g., 'pistis_sophia', 'stoic')",
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
		{
			Name:        "export_for_podcast",
			Description: "Export consultations as podcast episodes",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"days": map[string]interface{}{
						"type":        "number",
						"description": "Number of days to export (default: 7)",
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
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsError(req.ID, fmt.Sprintf("Invalid tool call params: %v", err))
	}

	result, err := s.HandleToolCall(params.Name, params.Arguments)
	if err != nil {
		return NewInternalError(req.ID, fmt.Sprintf("Tool call failed: %v", err))
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
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsError(req.ID, fmt.Sprintf("Invalid resource read params: %v", err))
	}

	// Parse resource URI
	uri := params.URI
	if uri == "wisdom://tools" {
		return s.handleToolsResource(req)
	} else if strings.HasPrefix(uri, "wisdom://sources") {
		return s.handleSourcesResource(req)
	} else if strings.HasPrefix(uri, "wisdom://advisors") {
		if uri == "wisdom://advisors" {
			return s.handleAdvisorsResource(req)
		}
		// Handle wisdom://advisor/{id}
		parts := strings.Split(uri, "/")
		if len(parts) >= 3 {
			advisorID := parts[len(parts)-1]
			return s.handleAdvisorResource(req, advisorID)
		}
		return NewInvalidParamsError(req.ID, "Invalid advisor resource URI")
	} else if strings.HasPrefix(uri, "wisdom://consultations/") {
		parts := strings.Split(uri, "/")
		if len(parts) >= 3 {
			daysStr := parts[len(parts)-1]
			days, err := strconv.Atoi(daysStr)
			if err != nil {
				return NewInvalidParamsError(req.ID, fmt.Sprintf("Invalid days parameter: %s", daysStr))
			}
			return s.handleConsultationsResource(req, days)
		}
		return NewInvalidParamsError(req.ID, "Invalid consultations resource URI")
	}

	return NewErrorResponse(req.ID, -32602, "Unknown resource URI", nil)
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
	case "export_for_podcast":
		return s.handleExportForPodcast(params)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
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

	// Create consultation
	consultation := wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          advisorInfo.Advisor,
		AdvisorIcon:      advisorInfo.Icon,
		AdvisorName:      advisorInfo.Advisor,
		Rationale:        advisorInfo.Rationale,
		ScoreAtTime:      score,
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
		return nil, fmt.Errorf("score parameter is required")
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
		return nil, fmt.Errorf("failed to get wisdom: %w", err)
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

	// TODO: Implement actual consultation log retrieval
	// days parameter will be used when log retrieval is implemented
	_ = days // Suppress unused variable warning until implementation
	return []interface{}{}, nil
}

// handleExportForPodcast implements export_for_podcast tool
func (s *WisdomServer) handleExportForPodcast(params map[string]interface{}) (interface{}, error) {
	days := 7 // Default
	if d, ok := params["days"].(float64); ok {
		days = int(d)
	} else if d, ok := params["days"].(int); ok {
		days = d
	}

	// TODO: Implement actual podcast export
	// For now, return empty episodes
	return map[string]interface{}{
		"episodes": []interface{}{},
		"days":     days,
	}, nil
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
				"source": "Wisdom source ID (e.g., 'pistis_sophia', 'stoic') - optional",
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
		{
			"name":        "export_for_podcast",
			"description": "Export consultations as podcast episodes",
			"parameters": map[string]interface{}{
				"days": "Number of days to export (default: 7)",
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
	// TODO: Implement advisor listing
	// For now, return placeholder
	advisors := []map[string]interface{}{
		{
			"id":        "pistis_sophia",
			"name":      "Pistis Sophia",
			"icon":      "📜",
			"rationale": "Journey through aeons mirrors project health stages",
		},
		{
			"id":        "stoic",
			"name":      "Stoic Philosophers",
			"icon":      "🏛️",
			"rationale": "Stoics teach discipline through adversity",
		},
	}

	return NewSuccessResponse(req.ID, map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      "wisdom://advisors",
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(advisors)),
			},
		},
	})
}

// handleAdvisorResource returns a specific advisor
func (s *WisdomServer) handleAdvisorResource(req *JSONRPCRequest, advisorID string) *JSONRPCResponse {
	// TODO: Implement advisor retrieval
	advisor := map[string]interface{}{
		"id":        advisorID,
		"name":      advisorID,
		"icon":      "📜",
		"rationale": "Wisdom advisor",
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
	// TODO: Implement consultation log retrieval
	consultations := []interface{}{}

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
