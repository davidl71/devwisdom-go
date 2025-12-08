package wisdom

import "fmt"

// AdvisorRegistry manages advisor mappings for metrics, tools, and stages
type AdvisorRegistry struct {
	metricAdvisors map[string]*AdvisorInfo
	toolAdvisors   map[string]*AdvisorInfo
	stageAdvisors  map[string]*AdvisorInfo
	initialized    bool
}

// NewAdvisorRegistry creates a new advisor registry
func NewAdvisorRegistry() *AdvisorRegistry {
	return &AdvisorRegistry{
		metricAdvisors: make(map[string]*AdvisorInfo),
		toolAdvisors:   make(map[string]*AdvisorInfo),
		stageAdvisors:  make(map[string]*AdvisorInfo),
	}
}

// Initialize loads advisor mappings
func (r *AdvisorRegistry) Initialize() {
	if r.initialized {
		return
	}

	// Load metric advisors
	r.loadMetricAdvisors()

	// Load tool advisors
	r.loadToolAdvisors()

	// Load stage advisors
	r.loadStageAdvisors()

	r.initialized = true
}

// loadMetricAdvisors populates metric → advisor mappings
func (r *AdvisorRegistry) loadMetricAdvisors() {
	r.metricAdvisors["security"] = &AdvisorInfo{
		Advisor:   "bofh",
		Icon:      "😈",
		Rationale: "BOFH is paranoid about security, expects users to break everything",
		HelpsWith: "Finding vulnerabilities, defensive thinking, access control",
	}

	r.metricAdvisors["testing"] = &AdvisorInfo{
		Advisor:   "stoic",
		Icon:      "🏛️",
		Rationale: "Stoics teach discipline through adversity - tests reveal truth",
		HelpsWith: "Persistence through failures, accepting harsh feedback",
	}

	// TODO: Port all metric advisors from Python version
	// See: project_management_automation/tools/wisdom/advisors.py METRIC_ADVISORS
}

// loadToolAdvisors populates tool → advisor mappings
func (r *AdvisorRegistry) loadToolAdvisors() {
	r.toolAdvisors["project_scorecard"] = &AdvisorInfo{
		Advisor:   "pistis_sophia",
		Rationale: "Journey through aeons mirrors project health stages",
	}

	// TODO: Port all tool advisors from Python version
	// See: project_management_automation/tools/wisdom/advisors.py TOOL_ADVISORS
}

// loadStageAdvisors populates stage → advisor mappings
func (r *AdvisorRegistry) loadStageAdvisors() {
	r.stageAdvisors["daily_checkin"] = &AdvisorInfo{
		Advisor:   "pistis_sophia",
		Icon:      "📜",
		Rationale: "Start each day with enlightenment journey wisdom",
	}

	// TODO: Port all stage advisors from Python version
	// See: project_management_automation/tools/wisdom/advisors.py STAGE_ADVISORS
}

// GetAdvisorForMetric returns the advisor for a given metric
func (r *AdvisorRegistry) GetAdvisorForMetric(metric string) (*AdvisorInfo, error) {
	advisor, exists := r.metricAdvisors[metric]
	if !exists {
		return nil, fmt.Errorf("no advisor for metric: %s", metric)
	}
	return advisor, nil
}

// GetAdvisorForTool returns the advisor for a given tool
func (r *AdvisorRegistry) GetAdvisorForTool(tool string) (*AdvisorInfo, error) {
	advisor, exists := r.toolAdvisors[tool]
	if !exists {
		return nil, fmt.Errorf("no advisor for tool: %s", tool)
	}
	return advisor, nil
}

// GetAdvisorForStage returns the advisor for a given stage
func (r *AdvisorRegistry) GetAdvisorForStage(stage string) (*AdvisorInfo, error) {
	advisor, exists := r.stageAdvisors[stage]
	if !exists {
		return nil, fmt.Errorf("no advisor for stage: %s", stage)
	}
	return advisor, nil
}
