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

	r.metricAdvisors["documentation"] = &AdvisorInfo{
		Advisor:   "confucius",
		Icon:      "🎓",
		Rationale: "Confucius emphasized teaching and transmitting wisdom",
		HelpsWith: "Clear explanations, teaching future maintainers",
	}

	r.metricAdvisors["completion"] = &AdvisorInfo{
		Advisor:   "art_of_war",
		Icon:      "⚔️",
		Rationale: "Sun Tzu teaches strategy and decisive execution",
		HelpsWith: "Prioritization, knowing when to attack vs wait",
	}

	r.metricAdvisors["alignment"] = &AdvisorInfo{
		Advisor:   "tao",
		Icon:      "☯️",
		Rationale: "Tao emphasizes balance, flow, and purpose",
		HelpsWith: "Ensuring work serves project goals, finding harmony",
	}

	r.metricAdvisors["clarity"] = &AdvisorInfo{
		Advisor:   "gracian",
		Icon:      "🎭",
		Rationale: "Gracián's maxims are models of clarity and pragmatism",
		HelpsWith: "Simplifying complexity, clear communication",
	}

	r.metricAdvisors["ci_cd"] = &AdvisorInfo{
		Advisor:   "kybalion",
		Icon:      "⚗️",
		Rationale: "Kybalion teaches cause and effect - CI/CD is pure causation",
		HelpsWith: "Understanding pipelines, automation philosophy",
	}

	r.metricAdvisors["dogfooding"] = &AdvisorInfo{
		Advisor:   "murphy",
		Icon:      "🔧",
		Rationale: "Murphy's Law: if it can break, it will - use your own tools!",
		HelpsWith: "Finding edge cases, eating your own cooking",
	}

	r.metricAdvisors["uniqueness"] = &AdvisorInfo{
		Advisor:   "shakespeare",
		Icon:      "🎭",
		Rationale: "Shakespeare created unique works that transcended his time",
		HelpsWith: "Creative differentiation, memorable design",
	}

	r.metricAdvisors["codebase"] = &AdvisorInfo{
		Advisor:   "enochian",
		Icon:      "🔮",
		Rationale: "Enochian mysticism reveals hidden structure and patterns",
		HelpsWith: "Architecture, finding hidden connections",
	}

	r.metricAdvisors["parallelizable"] = &AdvisorInfo{
		Advisor:   "tao_of_programming",
		Icon:      "💻",
		Rationale: "The Tao of Programming teaches elegant parallel design",
		HelpsWith: "Decomposition, independent task design",
	}

	// Hebrew Advisors - Jewish wisdom traditions
	r.metricAdvisors["ethics"] = &AdvisorInfo{
		Advisor:   "rebbe",
		Icon:      "🕎",
		Rationale: "The Rebbe teaches ethical conduct and righteous behavior (מוסר)",
		HelpsWith: "Code ethics, proper conduct, doing the right thing",
		Language:  "hebrew",
	}

	r.metricAdvisors["perseverance"] = &AdvisorInfo{
		Advisor:   "tzaddik",
		Icon:      "✡️",
		Rationale: "The Tzaddik (righteous one) demonstrates steadfast commitment",
		HelpsWith: "Persistence, staying on the righteous path, not giving up",
		Language:  "hebrew",
	}

	r.metricAdvisors["wisdom"] = &AdvisorInfo{
		Advisor:   "chacham",
		Icon:      "📜",
		Rationale: "The Chacham (sage) seeks deep understanding through Torah",
		HelpsWith: "Deep analysis, seeking understanding, learning from tradition",
		Language:  "hebrew",
	}
}

// loadToolAdvisors populates tool → advisor mappings
func (r *AdvisorRegistry) loadToolAdvisors() {
	r.toolAdvisors["project_scorecard"] = &AdvisorInfo{
		Advisor:   "pistis_sophia",
		Rationale: "Journey through aeons mirrors project health stages",
	}

	r.toolAdvisors["project_overview"] = &AdvisorInfo{
		Advisor:   "kybalion",
		Rationale: "Hermetic principles for holistic understanding",
	}

	r.toolAdvisors["sprint_automation"] = &AdvisorInfo{
		Advisor:   "art_of_war",
		Rationale: "Sprint is a campaign requiring strategy",
	}

	r.toolAdvisors["check_documentation_health"] = &AdvisorInfo{
		Advisor:   "confucius",
		Rationale: "Teaching requires good documentation",
	}

	r.toolAdvisors["analyze_todo2_alignment"] = &AdvisorInfo{
		Advisor:   "tao",
		Rationale: "Alignment is balance and flow",
	}

	r.toolAdvisors["detect_duplicate_tasks"] = &AdvisorInfo{
		Advisor:   "bofh",
		Rationale: "Duplicates are user error manifested",
	}

	r.toolAdvisors["scan_dependency_security"] = &AdvisorInfo{
		Advisor:   "bofh",
		Rationale: "Security paranoia is a feature",
	}

	r.toolAdvisors["run_tests"] = &AdvisorInfo{
		Advisor:   "stoic",
		Rationale: "Tests teach through failure",
	}

	r.toolAdvisors["validate_ci_cd_workflow"] = &AdvisorInfo{
		Advisor:   "kybalion",
		Rationale: "CI/CD is cause and effect",
	}

	r.toolAdvisors["dev_reload"] = &AdvisorInfo{
		Advisor:   "murphy",
		Rationale: "Hot reload because Murphy says restarts will fail at the worst time",
	}

	// Hebrew advisor tools - for ethical and wisdom-focused operations
	r.toolAdvisors["ethics_check"] = &AdvisorInfo{
		Advisor:   "rebbe",
		Rationale: "Rebbe guides ethical code review and conduct",
		Language:  "hebrew",
	}

	r.toolAdvisors["wisdom_reflection"] = &AdvisorInfo{
		Advisor:   "chacham",
		Rationale: "Chacham provides deep wisdom for retrospectives",
		Language:  "hebrew",
	}
}

// loadStageAdvisors populates stage → advisor mappings
func (r *AdvisorRegistry) loadStageAdvisors() {
	r.stageAdvisors["daily_checkin"] = &AdvisorInfo{
		Advisor:   "pistis_sophia",
		Icon:      "📜",
		Rationale: "Start each day with enlightenment journey wisdom",
	}

	r.stageAdvisors["planning"] = &AdvisorInfo{
		Advisor:   "art_of_war",
		Icon:      "⚔️",
		Rationale: "Planning is strategy - Sun Tzu is the master",
	}

	r.stageAdvisors["implementation"] = &AdvisorInfo{
		Advisor:   "tao_of_programming",
		Icon:      "💻",
		Rationale: "During coding, let the code flow naturally",
	}

	r.stageAdvisors["debugging"] = &AdvisorInfo{
		Advisor:   "bofh",
		Icon:      "😈",
		Rationale: "BOFH knows all the ways things break",
	}

	r.stageAdvisors["review"] = &AdvisorInfo{
		Advisor:   "stoic",
		Icon:      "🏛️",
		Rationale: "Review requires accepting harsh truths with equanimity",
	}

	r.stageAdvisors["retrospective"] = &AdvisorInfo{
		Advisor:   "confucius",
		Icon:      "🎓",
		Rationale: "Retrospectives are about learning and teaching",
	}

	r.stageAdvisors["celebration"] = &AdvisorInfo{
		Advisor:   "shakespeare",
		Icon:      "🎭",
		Rationale: "Celebrate with drama and poetry!",
	}

	// Hebrew advisor stages
	r.stageAdvisors["shabbat"] = &AdvisorInfo{
		Advisor:   "rebbe",
		Icon:      "🕎",
		Rationale: "Shabbat is for reflection and spiritual renewal (מנוחה)",
		Language:  "hebrew",
	}

	r.stageAdvisors["teshuvah"] = &AdvisorInfo{
		Advisor:   "tzaddik",
		Icon:      "✡️",
		Rationale: "Teshuvah (repentance) is for fixing past mistakes and returning to the right path",
		Language:  "hebrew",
	}

	r.stageAdvisors["learning"] = &AdvisorInfo{
		Advisor:   "chacham",
		Icon:      "📜",
		Rationale: "Torah study and continuous learning (לימוד)",
		Language:  "hebrew",
	}
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
