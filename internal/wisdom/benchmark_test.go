package wisdom

import (
	"testing"
)

// BenchmarkEngine_GetWisdom benchmarks wisdom quote retrieval
func BenchmarkEngine_GetWisdom(b *testing.B) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		b.Fatalf("Failed to initialize engine: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.GetWisdom(75.0, "stoic")
		if err != nil {
			b.Fatalf("GetWisdom failed: %v", err)
		}
	}
}

// BenchmarkEngine_GetWisdom_Random benchmarks random source selection
func BenchmarkEngine_GetWisdom_Random(b *testing.B) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		b.Fatalf("Failed to initialize engine: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.GetWisdom(75.0, "random")
		if err != nil {
			b.Fatalf("GetWisdom failed: %v", err)
		}
	}
}

// BenchmarkEngine_GetRandomSource benchmarks random source ID retrieval
func BenchmarkEngine_GetRandomSource(b *testing.B) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		b.Fatalf("Failed to initialize engine: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.GetRandomSource(true) // date-seeded
		if err != nil {
			b.Fatalf("GetRandomSource failed: %v", err)
		}
	}
}

// BenchmarkEngine_ListSources benchmarks source listing
func BenchmarkEngine_ListSources(b *testing.B) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		b.Fatalf("Failed to initialize engine: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = engine.ListSources()
	}
}

// BenchmarkEngine_GetSource benchmarks individual source retrieval
func BenchmarkEngine_GetSource(b *testing.B) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		b.Fatalf("Failed to initialize engine: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.GetSource("stoic")
	}
}

// BenchmarkAdvisorRegistry_GetAdvisorForMetric benchmarks advisor lookup by metric
func BenchmarkAdvisorRegistry_GetAdvisorForMetric(b *testing.B) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.GetAdvisorForMetric("security")
	}
}

// BenchmarkAdvisorRegistry_GetAdvisorForTool benchmarks advisor lookup by tool
func BenchmarkAdvisorRegistry_GetAdvisorForTool(b *testing.B) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.GetAdvisorForTool("project_scorecard")
	}
}

// BenchmarkAdvisorRegistry_GetAdvisorForStage benchmarks advisor lookup by stage
func BenchmarkAdvisorRegistry_GetAdvisorForStage(b *testing.B) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = registry.GetAdvisorForStage("daily_checkin")
	}
}

// BenchmarkGetAeonLevel benchmarks aeon level calculation
func BenchmarkGetAeonLevel(b *testing.B) {
	scores := []float64{25.0, 40.0, 60.0, 75.0, 90.0}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		score := scores[i%len(scores)]
		_ = GetAeonLevel(score)
	}
}

// BenchmarkSource_GetQuote benchmarks quote retrieval from source
func BenchmarkSource_GetQuote(b *testing.B) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		b.Fatalf("Failed to initialize engine: %v", err)
	}

	source, exists := engine.GetSource("stoic")
	if !exists {
		b.Fatal("Source 'stoic' not found")
	}

	aeonLevels := []string{"chaos", "lower_aeons", "middle_aeons", "upper_aeons", "treasury"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		level := aeonLevels[i%len(aeonLevels)]
		_ = source.GetQuote(level)
	}
}

// BenchmarkEngine_Initialize benchmarks engine initialization
func BenchmarkEngine_Initialize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		engine := NewEngine()
		_ = engine.Initialize()
	}
}

