package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig() returned nil")
	}
	if cfg.Source != "pistis_sophia" {
		t.Errorf("NewConfig().Source = %q, want %q", cfg.Source, "pistis_sophia")
	}
	if cfg.HebrewEnabled {
		t.Error("NewConfig().HebrewEnabled = true, want false")
	}
	if cfg.HebrewOnly {
		t.Error("NewConfig().HebrewOnly = true, want false")
	}
	if cfg.Disabled {
		t.Error("NewConfig().Disabled = true, want false")
	}
	if cfg.configPath == "" {
		t.Error("NewConfig().configPath is empty")
	}
}

func TestConfig_Load_EnvironmentVariables(t *testing.T) {
	// Save original values.
	originalSource := os.Getenv("EXARP_WISDOM_SOURCE")
	originalHebrew := os.Getenv("EXARP_WISDOM_HEBREW")
	originalHebrewOnly := os.Getenv("EXARP_WISDOM_HEBREW_ONLY")
	originalDisabled := os.Getenv("EXARP_DISABLE_WISDOM")

	// Clean up after test.
	defer func() {
		if originalSource != "" {
			os.Setenv("EXARP_WISDOM_SOURCE", originalSource)
		} else {
			os.Unsetenv("EXARP_WISDOM_SOURCE")
		}
		if originalHebrew != "" {
			os.Setenv("EXARP_WISDOM_HEBREW", originalHebrew)
		} else {
			os.Unsetenv("EXARP_WISDOM_HEBREW")
		}
		if originalHebrewOnly != "" {
			os.Setenv("EXARP_WISDOM_HEBREW_ONLY", originalHebrewOnly)
		} else {
			os.Unsetenv("EXARP_WISDOM_HEBREW_ONLY")
		}
		if originalDisabled != "" {
			os.Setenv("EXARP_DISABLE_WISDOM", originalDisabled)
		} else {
			os.Unsetenv("EXARP_DISABLE_WISDOM")
		}
	}()

	tests := []struct {
		name           string
		setEnv         func()
		wantSource     string
		wantHebrew     bool
		wantHebrewOnly bool
		wantDisabled   bool
	}{
		{
			name: "no environment variables",
			setEnv: func() {
				os.Unsetenv("EXARP_WISDOM_SOURCE")
				os.Unsetenv("EXARP_WISDOM_HEBREW")
				os.Unsetenv("EXARP_WISDOM_HEBREW_ONLY")
				os.Unsetenv("EXARP_DISABLE_WISDOM")
			},
			wantSource:     "pistis_sophia",
			wantHebrew:     false,
			wantHebrewOnly: false,
			wantDisabled:   false,
		},
		{
			name: "set source",
			setEnv: func() {
				os.Setenv("EXARP_WISDOM_SOURCE", "stoic")
				os.Unsetenv("EXARP_WISDOM_HEBREW")
				os.Unsetenv("EXARP_WISDOM_HEBREW_ONLY")
				os.Unsetenv("EXARP_DISABLE_WISDOM")
			},
			wantSource:     "stoic",
			wantHebrew:     false,
			wantHebrewOnly: false,
			wantDisabled:   false,
		},
		{
			name: "set hebrew enabled",
			setEnv: func() {
				os.Unsetenv("EXARP_WISDOM_SOURCE")
				os.Setenv("EXARP_WISDOM_HEBREW", "1")
				os.Unsetenv("EXARP_WISDOM_HEBREW_ONLY")
				os.Unsetenv("EXARP_DISABLE_WISDOM")
			},
			wantSource:     "pistis_sophia",
			wantHebrew:     true,
			wantHebrewOnly: false,
			wantDisabled:   false,
		},
		{
			name: "set hebrew only",
			setEnv: func() {
				os.Unsetenv("EXARP_WISDOM_SOURCE")
				os.Unsetenv("EXARP_WISDOM_HEBREW")
				os.Setenv("EXARP_WISDOM_HEBREW_ONLY", "1")
				os.Unsetenv("EXARP_DISABLE_WISDOM")
			},
			wantSource:     "pistis_sophia",
			wantHebrew:     false,
			wantHebrewOnly: true,
			wantDisabled:   false,
		},
		{
			name: "set disabled",
			setEnv: func() {
				os.Unsetenv("EXARP_WISDOM_SOURCE")
				os.Unsetenv("EXARP_WISDOM_HEBREW")
				os.Unsetenv("EXARP_WISDOM_HEBREW_ONLY")
				os.Setenv("EXARP_DISABLE_WISDOM", "1")
			},
			wantSource:     "pistis_sophia",
			wantHebrew:     false,
			wantHebrewOnly: false,
			wantDisabled:   true,
		},
		{
			name: "all environment variables",
			setEnv: func() {
				os.Setenv("EXARP_WISDOM_SOURCE", "tao")
				os.Setenv("EXARP_WISDOM_HEBREW", "1")
				os.Setenv("EXARP_WISDOM_HEBREW_ONLY", "1")
				os.Setenv("EXARP_DISABLE_WISDOM", "1")
			},
			wantSource:     "tao",
			wantHebrew:     true,
			wantHebrewOnly: true,
			wantDisabled:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setEnv()
			cfg := NewConfig()
			if err := cfg.Load(); err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.Source != tt.wantSource {
				t.Errorf("Source = %q, want %q", cfg.Source, tt.wantSource)
			}
			if cfg.HebrewEnabled != tt.wantHebrew {
				t.Errorf("HebrewEnabled = %v, want %v", cfg.HebrewEnabled, tt.wantHebrew)
			}
			if cfg.HebrewOnly != tt.wantHebrewOnly {
				t.Errorf("HebrewOnly = %v, want %v", cfg.HebrewOnly, tt.wantHebrewOnly)
			}
			if cfg.Disabled != tt.wantDisabled {
				t.Errorf("Disabled = %v, want %v", cfg.Disabled, tt.wantDisabled)
			}
		})
	}
}

func TestConfig_Load_ConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".exarp_wisdom_config")

	cfg := NewConfig()
	cfg.configPath = configPath

	// Test with valid config file.
	configData := map[string]interface{}{
		"source":         "stoic",
		"hebrew_enabled": true,
		"hebrew_only":    false,
		"disabled":       false,
	}
	data, err := json.Marshal(configData)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	if err := cfg.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Source != "stoic" {
		t.Errorf("Source = %q, want %q", cfg.Source, "stoic")
	}
	if !cfg.HebrewEnabled {
		t.Error("HebrewEnabled = false, want true")
	}
}

func TestConfig_Load_MarkerFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create marker file.
	markerFile := ".exarp_no_wisdom"
	if err := os.WriteFile(markerFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create marker file: %v", err)
	}
	defer os.Remove(markerFile)

	cfg := NewConfig()
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Disabled {
		t.Error("Disabled = false, want true (marker file should disable)")
	}
}

func TestConfig_Load_MissingConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "nonexistent", ".exarp_wisdom_config")

	cfg := NewConfig()
	cfg.configPath = configPath

	// Should not error if config file doesn't exist.
	if err := cfg.Load(); err != nil {
		t.Errorf("Load() error = %v, want nil (config file is optional)", err)
	}
}

func TestConfig_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".exarp_wisdom_config")

	cfg := NewConfig()
	cfg.configPath = configPath
	cfg.Source = "test_source"
	cfg.HebrewEnabled = true

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created.
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("Config file was not created: %v", err)
	}

	// Verify file contents.
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal config file: %v", err)
	}

	if loaded.Source != "test_source" {
		t.Errorf("Loaded Source = %q, want %q", loaded.Source, "test_source")
	}
	if !loaded.HebrewEnabled {
		t.Error("Loaded HebrewEnabled = false, want true")
	}
}

func TestConfig_Save_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", ".exarp_wisdom_config")

	cfg := NewConfig()
	cfg.configPath = configPath

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify directory was created.
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); err != nil {
		t.Fatalf("Config directory was not created: %v", err)
	}
}

func TestGetConfigPath(t *testing.T) {
	// Test that GetConfigPath returns a non-empty path.
	path := GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath() returned empty path")
	}

	// Should end with .exarp_wisdom_config.
	if filepath.Base(path) != ".exarp_wisdom_config" {
		t.Errorf("GetConfigPath() = %q, want path ending with .exarp_wisdom_config", path)
	}
}

func TestConfig_Load_EnvironmentOverridesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".exarp_wisdom_config")

	// Create config file with one value.
	configData := map[string]interface{}{
		"source": "file_source",
	}
	data, err := json.Marshal(configData)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Set environment variable (should override file).
	originalSource := os.Getenv("EXARP_WISDOM_SOURCE")
	defer func() {
		if originalSource != "" {
			os.Setenv("EXARP_WISDOM_SOURCE", originalSource)
		} else {
			os.Unsetenv("EXARP_WISDOM_SOURCE")
		}
	}()

	os.Setenv("EXARP_WISDOM_SOURCE", "env_source")

	cfg := NewConfig()
	cfg.configPath = configPath
	if err := cfg.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Environment should override file.
	if cfg.Source != "env_source" {
		t.Errorf("Source = %q, want %q (environment should override file)", cfg.Source, "env_source")
	}
}
