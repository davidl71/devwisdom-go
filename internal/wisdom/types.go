package wisdom

// Quote represents a wisdom quote with metadata
type Quote struct {
	Quote         string `json:"quote"`
	Source        string `json:"source"`
	Encouragement string `json:"encouragement"`
	WisdomSource  string `json:"wisdom_source,omitempty"`
	WisdomIcon    string `json:"wisdom_icon,omitempty"`
}

// Source represents a wisdom source with quotes by aeon level
type Source struct {
	Name        string             `json:"name"`
	Icon        string             `json:"icon"`
	Quotes      map[string][]Quote `json:"quotes"` // Key: aeon level (chaos, lower_aeons, etc.)
	Description string             `json:"description,omitempty"`
}

// GetQuote retrieves a quote for the given aeon level
func (s *Source) GetQuote(aeonLevel string) *Quote {
	quotes, exists := s.Quotes[aeonLevel]
	if !exists || len(quotes) == 0 {
		// Fallback to any available quotes
		for _, qs := range s.Quotes {
			if len(qs) > 0 {
				quotes = qs
				break
			}
		}
	}

	if len(quotes) == 0 {
		return &Quote{
			Quote:         "Silence is also wisdom.",
			Source:        "Unknown",
			Encouragement: "Sometimes reflection is the answer.",
		}
	}

	// TODO: Implement daily consistent random selection
	// For now, return first quote
	return &quotes[0]
}

// Consultation represents an advisor consultation
type Consultation struct {
	Timestamp        string  `json:"timestamp"`
	ConsultationType string  `json:"consultation_type"`
	Advisor          string  `json:"advisor"`
	AdvisorIcon      string  `json:"advisor_icon"`
	AdvisorName      string  `json:"advisor_name"`
	Rationale        string  `json:"rationale"`
	Metric           string  `json:"metric,omitempty"`
	Tool             string  `json:"tool,omitempty"`
	Stage            string  `json:"stage,omitempty"`
	ScoreAtTime      float64 `json:"score_at_time"`
	ConsultationMode string  `json:"consultation_mode"`
	ModeIcon         string  `json:"mode_icon"`
	ModeFrequency    string  `json:"mode_frequency"`
	Quote            string  `json:"quote"`
	QuoteSource      string  `json:"quote_source"`
	Encouragement    string  `json:"encouragement"`
	Context          string  `json:"context,omitempty"`
	SessionMode      string  `json:"session_mode,omitempty"`
	ModeGuidance     string  `json:"mode_guidance,omitempty"`
}

// AdvisorInfo represents advisor metadata
type AdvisorInfo struct {
	Advisor   string `json:"advisor"`
	Icon      string `json:"icon"`
	Rationale string `json:"rationale"`
	HelpsWith string `json:"helps_with,omitempty"`
	Language  string `json:"language,omitempty"`
}

// AeonLevel represents project health stages
type AeonLevel string

const (
	AeonChaos    AeonLevel = "chaos"        // < 30%
	AeonLower    AeonLevel = "lower_aeons"  // 30-50%
	AeonMiddle   AeonLevel = "middle_aeons" // 50-70%
	AeonUpper    AeonLevel = "upper_aeons"  // 70-85%
	AeonTreasury AeonLevel = "treasury"     // > 85%
)

// GetAeonLevel returns the aeon level based on score
func GetAeonLevel(score float64) string {
	switch {
	case score < 30:
		return string(AeonChaos)
	case score < 50:
		return string(AeonLower)
	case score < 70:
		return string(AeonMiddle)
	case score < 85:
		return string(AeonUpper)
	default:
		return string(AeonTreasury)
	}
}

// ConsultationMode represents project health consultation modes
type ConsultationMode string

const (
	ModeChaos    ConsultationMode = "chaos"    // < 30%
	ModeBuilding ConsultationMode = "building" // 30-60%
	ModeMaturing ConsultationMode = "maturing" // 60-80%
	ModeMastery  ConsultationMode = "mastery"  // > 80%
)

// ConsultationModeConfig represents configuration for a consultation mode
type ConsultationModeConfig struct {
	Name        string  `json:"name"`
	MinScore    float64 `json:"min_score"`
	MaxScore    float64 `json:"max_score"`
	Frequency   string  `json:"frequency"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
}

// SessionMode represents different session interaction modes
type SessionMode string

const (
	SessionModeAgent  SessionMode = "AGENT"
	SessionModeAsk    SessionMode = "ASK"
	SessionModeManual SessionMode = "MANUAL"
)

// ModeConfig represents configuration for session mode-aware advisor selection
type ModeConfig struct {
	PreferredAdvisors []string `json:"preferred_advisors"`
	Tone              string   `json:"tone"`
	Focus             string   `json:"focus"`
}
