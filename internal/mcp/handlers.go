// Package mcp provides handler methods for wisdom tools and resources.
// These handlers are used by the SDK adapter (sdk_adapter.go).
package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davidl71/devwisdom-go/internal/logging"
	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// WisdomHandlers provides handler methods for tools and resources.
// This struct contains the business logic that is reused by the SDK adapter.
type WisdomHandlers struct {
	wisdom    *wisdom.Engine
	logger    *logging.ConsultationLogger
	appLogger *logging.Logger
}

// NewWisdomHandlers creates a new handlers instance.
func NewWisdomHandlers(wisdomEngine *wisdom.Engine, consultationLogger *logging.ConsultationLogger, appLogger *logging.Logger) *WisdomHandlers {
	return &WisdomHandlers{
		wisdom:    wisdomEngine,
		logger:    consultationLogger,
		appLogger: appLogger,
	}
}

// HandleToolCall processes MCP tool calls
func (h *WisdomHandlers) HandleToolCall(name string, params map[string]interface{}) (interface{}, error) {
	switch name {
	case "consult_advisor":
		return h.handleConsultAdvisor(params)
	case "get_wisdom":
		return h.handleGetWisdom(params)
	case "get_daily_briefing":
		return h.handleGetDailyBriefing(params)
	case "get_consultation_log":
		return h.handleGetConsultationLog(params)
	default:
		availableTools := []string{"consult_advisor", "get_wisdom", "get_daily_briefing", "get_consultation_log"}
		return nil, fmt.Errorf("unknown tool %q (available tools: %v). Check tool name spelling", name, availableTools)
	}
}

// handleConsultAdvisor implements consult_advisor tool
func (h *WisdomHandlers) handleConsultAdvisor(params map[string]interface{}) (interface{}, error) {
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
		advisorInfo, err = h.wisdom.GetAdvisors().GetAdvisorForMetric(metric)
	} else if tool != "" {
		advisorInfo, err = h.wisdom.GetAdvisors().GetAdvisorForTool(tool)
	} else if stage != "" {
		advisorInfo, err = h.wisdom.GetAdvisors().GetAdvisorForStage(stage)
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
	quote, err := h.wisdom.GetWisdom(score, advisorInfo.Advisor)
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

	// Log consultation if logger is available
	if h.logger != nil {
		if err := h.logger.Log(&consultation); err != nil {
			// Logging failure is non-fatal
			h.appLogger.Warn("", "Failed to log consultation: %v", err)
		}
	}

	return consultation, nil
}

// handleGetWisdom implements get_wisdom tool
func (h *WisdomHandlers) handleGetWisdom(params map[string]interface{}) (interface{}, error) {
	var score float64
	var source string

	// Score is required
	if sc, ok := params["score"].(float64); ok {
		score = sc
	} else if sc, ok := params["score"].(int); ok {
		score = float64(sc)
	} else {
		return nil, fmt.Errorf("score parameter is required and must be a number between 0-100")
	}

	// Validate and clamp score
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Source is optional
	if s, ok := params["source"].(string); ok {
		source = s
	}

	// Get wisdom quote
	quote, err := h.wisdom.GetWisdom(score, source)
	if err != nil {
		return nil, fmt.Errorf("failed to get wisdom quote: %w", err)
	}

	return quote, nil
}

// handleGetDailyBriefing implements get_daily_briefing tool
func (h *WisdomHandlers) handleGetDailyBriefing(params map[string]interface{}) (interface{}, error) {
	var score float64

	if sc, ok := params["score"].(float64); ok {
		score = sc
	} else if sc, ok := params["score"].(int); ok {
		score = float64(sc)
	}

	// Validate and clamp score
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Get multiple quotes from different sources
	sources := h.wisdom.ListSources()
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
		quote, err := h.wisdom.GetWisdom(score, src)
		if err == nil {
			quotes = append(quotes, quote)
		}
	}

	briefing["quotes"] = quotes
	return briefing, nil
}

// handleGetConsultationLog implements get_consultation_log tool
func (h *WisdomHandlers) handleGetConsultationLog(params map[string]interface{}) (interface{}, error) {
	days := 7 // default

	if d, ok := params["days"].(float64); ok {
		days = int(d)
	} else if d, ok := params["days"].(int); ok {
		days = d
	}

	// Retrieve consultations from logger
	var consultations []interface{}
	if h.logger != nil {
		logs, err := h.logger.GetLogs(days)
		if err == nil {
			// Convert to interface{} slice for JSON serialization
			consultations = make([]interface{}, len(logs))
			for i, consultation := range logs {
				consultations[i] = consultation
			}
		}
		// If logger is nil or error occurs, consultations remains empty array
	}

	return consultations, nil
}

// HandleToolsResource handles wisdom://tools resource
func (h *WisdomHandlers) HandleToolsResource(req *JSONRPCRequest) *JSONRPCResponse {
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
		"contents": []map[string]interface{}{
			{
				"uri":      "wisdom://tools",
				"mimeType": "application/json",
				"text":     string(mustMarshalJSONCompact(tools)),
			},
		},
	})
}

// HandleSourcesResource handles wisdom://sources resource
func (h *WisdomHandlers) HandleSourcesResource(req *JSONRPCRequest) *JSONRPCResponse {
	sourceIDs := h.wisdom.ListSources()
	sourceDetails := make([]map[string]interface{}, 0, len(sourceIDs))

	for _, id := range sourceIDs {
		source, found := h.wisdom.GetSource(id)
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

// HandleAdvisorsResource handles wisdom://advisors resource
func (h *WisdomHandlers) HandleAdvisorsResource(req *JSONRPCRequest) *JSONRPCResponse {
	advisorRegistry := h.wisdom.GetAdvisors()

	// Get all advisors
	metricAdvisors := advisorRegistry.GetAllMetricAdvisors()
	toolAdvisors := advisorRegistry.GetAllToolAdvisors()
	stageAdvisors := advisorRegistry.GetAllStageAdvisors()

	// Build response
	advisors := map[string]interface{}{
		"metric_advisors": metricAdvisors,
		"tool_advisors":   toolAdvisors,
		"stage_advisors":  stageAdvisors,
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

// HandleAdvisorResource handles wisdom://advisor/{id} resource
func (h *WisdomHandlers) HandleAdvisorResource(req *JSONRPCRequest, advisorID string) *JSONRPCResponse {
	advisorRegistry := h.wisdom.GetAdvisors()

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

// HandleConsultationsResource handles wisdom://consultations/{days} resource
func (h *WisdomHandlers) HandleConsultationsResource(req *JSONRPCRequest, days int) *JSONRPCResponse {
	// Retrieve consultations from logger
	var consultations []interface{}
	if h.logger != nil {
		logs, err := h.logger.GetLogs(days)
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

