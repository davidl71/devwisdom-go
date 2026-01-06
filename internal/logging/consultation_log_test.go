package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

func TestNewConsultationLogger(t *testing.T) {
	// Create temporary directory for tests
	tmpDir := t.TempDir()

	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	if logger == nil {
		t.Fatal("NewConsultationLogger returned nil")
	}

	// Verify log file was created
	logPath := filepath.Join(tmpDir, "consultations.jsonl")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestConsultationLogger_Log(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Create test consultation
	consultation := &wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "pistis_sophia",
		AdvisorIcon:      "📜",
		AdvisorName:      "Pistis Sophia",
		Rationale:        "Test rationale",
		ScoreAtTime:      75.0,
		ConsultationMode: "maturing",
		ModeIcon:         "🌱",
		ModeFrequency:    "milestones",
		Quote:            "Test quote",
		QuoteSource:      "test_source",
		Encouragement:    "Test encouragement",
	}

	// Log consultation
	if err := logger.Log(consultation); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Verify file was written
	logPath := filepath.Join(tmpDir, "consultations.jsonl")
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Log file is empty")
	}
}

func TestConsultationLogger_GetLogs(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Log multiple consultations with different timestamps
	now := time.Now()
	for i := 0; i < 5; i++ {
		consultation := &wisdom.Consultation{
			Timestamp:        now.AddDate(0, 0, -i).Format(time.RFC3339),
			ConsultationType: "advisor",
			Advisor:          "pistis_sophia",
			ScoreAtTime:      75.0,
			Quote:            "Test quote",
		}
		if err := logger.Log(consultation); err != nil {
			t.Fatalf("Log failed: %v", err)
		}
	}

	// Retrieve logs for last 3 days
	logs, err := logger.GetLogs(3)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	// Should get consultations from last 3 days (indices 0, 1, 2)
	if len(logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}
}

func TestConsultationLogger_GetLogs_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Get logs from empty file
	logs, err := logger.GetLogs(7)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	if len(logs) != 0 {
		t.Errorf("Expected empty logs, got %d", len(logs))
	}
}

func TestConsultationLogger_GetLogs_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create logger but don't write anything (file won't exist until first write)
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	logger.Close() // Close to release file handle

	// Delete the file
	logPath := filepath.Join(tmpDir, "consultations.jsonl")
	os.Remove(logPath)

	// Create new logger instance
	logger2, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger2.Close()

	// Get logs from non-existent file should return empty array
	logs, err := logger2.GetLogs(7)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	if len(logs) != 0 {
		t.Errorf("Expected empty logs for non-existent file, got %d", len(logs))
	}
}

func TestConsultationLogger_ThreadSafety(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			consultation := &wisdom.Consultation{
				Timestamp:        time.Now().Format(time.RFC3339),
				ConsultationType: "advisor",
				Advisor:          "pistis_sophia",
				ScoreAtTime:      float64(id),
				Quote:            "Test quote",
			}
			if err := logger.Log(consultation); err != nil {
				t.Errorf("Concurrent Log failed: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all consultations were logged
	logs, err := logger.GetLogs(1)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	if len(logs) != 10 {
		t.Errorf("Expected 10 logs from concurrent writes, got %d", len(logs))
	}
}

func TestConsultationLogger_LogRotation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a log file with yesterday's date
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayPath := filepath.Join(tmpDir, "consultations.jsonl")

	// Create file and set its modification time to yesterday
	file, err := os.Create(yesterdayPath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Set modification time to yesterday
	if err := os.Chtimes(yesterdayPath, yesterday, yesterday); err != nil {
		t.Fatalf("Failed to set file time: %v", err)
	}

	// Create logger - should rotate the file
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Check that rotated file was created
	rotatedPath := filepath.Join(tmpDir, fmt.Sprintf("consultations-%s.jsonl", yesterday.Format("2006-01-02")))
	if _, err := os.Stat(rotatedPath); os.IsNotExist(err) {
		t.Error("Rotated log file was not created")
	}

	// Check that current file exists
	currentPath := filepath.Join(tmpDir, "consultations.jsonl")
	if _, err := os.Stat(currentPath); os.IsNotExist(err) {
		t.Error("Current log file was not created after rotation")
	}
}

func TestConsultationLogger_GetLogs_WithRotatedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a rotated log file manually
	yesterday := time.Now().AddDate(0, 0, -1)
	rotatedPath := filepath.Join(tmpDir, fmt.Sprintf("consultations-%s.jsonl", yesterday.Format("2006-01-02")))

	// Write some test data to rotated file
	rotatedFile, err := os.Create(rotatedPath)
	if err != nil {
		t.Fatalf("Failed to create rotated file: %v", err)
	}

	consultation := &wisdom.Consultation{
		Timestamp:        yesterday.Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "pistis_sophia",
		ScoreAtTime:      75.0,
		Quote:            "Test quote from yesterday",
	}

	encoder := json.NewEncoder(rotatedFile)
	if err := encoder.Encode(consultation); err != nil {
		rotatedFile.Close()
		t.Fatalf("Failed to write to rotated file: %v", err)
	}
	rotatedFile.Close()

	// Create logger and write to current file
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Write to current file
	currentConsultation := &wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "stoic",
		ScoreAtTime:      80.0,
		Quote:            "Test quote from today",
	}
	if err := logger.Log(currentConsultation); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Get logs for last 2 days - should get both
	logs, err := logger.GetLogs(2)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	// Should have 2 logs (one from yesterday, one from today)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs (from rotated and current), got %d", len(logs))
	}

	// Verify we got both quotes
	quotes := make(map[string]bool)
	for _, log := range logs {
		quotes[log.Quote] = true
	}
	if !quotes["Test quote from yesterday"] {
		t.Error("Missing quote from rotated file")
	}
	if !quotes["Test quote from today"] {
		t.Error("Missing quote from current file")
	}
}

func TestConsultationLogger_RotationDuringActiveLogging(t *testing.T) {
	tmpDir := t.TempDir()

	// Create logger
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Write initial consultation
	consultation1 := &wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "pistis_sophia",
		ScoreAtTime:      75.0,
		Quote:            "First quote",
	}
	if err := logger.Log(consultation1); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Manually trigger rotation by manipulating the currentDate field
	// This simulates a date change
	logger.mu.Lock()
	logger.currentDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	logger.mu.Unlock()

	// Write another consultation - should trigger rotation
	consultation2 := &wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "stoic",
		ScoreAtTime:      80.0,
		Quote:            "Second quote after rotation",
	}
	if err := logger.Log(consultation2); err != nil {
		t.Fatalf("Log failed after rotation: %v", err)
	}

	// Check that rotated file exists
	rotatedPath := filepath.Join(tmpDir, fmt.Sprintf("consultations-%s.jsonl", time.Now().AddDate(0, 0, -1).Format("2006-01-02")))
	if _, err := os.Stat(rotatedPath); os.IsNotExist(err) {
		t.Error("Rotated file was not created during active logging")
	}

	// Verify both logs are retrievable
	logs, err := logger.GetLogs(2)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	if len(logs) < 2 {
		t.Errorf("Expected at least 2 logs, got %d", len(logs))
	}
}

func TestConsultationLogger_RotationIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create logger
	logger, err := NewConsultationLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewConsultationLogger failed: %v", err)
	}
	defer logger.Close()

	// Manually set currentDate to yesterday to simulate date change
	logger.mu.Lock()
	logger.currentDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	logger.mu.Unlock()

	// First rotation
	consultation1 := &wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "pistis_sophia",
		ScoreAtTime:      75.0,
		Quote:            "First quote",
	}
	if err := logger.Log(consultation1); err != nil {
		t.Fatalf("First Log failed: %v", err)
	}

	// Try to rotate again (should be idempotent - no error, no duplicate rotation)
	consultation2 := &wisdom.Consultation{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConsultationType: "advisor",
		Advisor:          "stoic",
		ScoreAtTime:      80.0,
		Quote:            "Second quote",
	}
	if err := logger.Log(consultation2); err != nil {
		t.Fatalf("Second Log failed: %v", err)
	}

	// Should only have one rotated file
	rotatedPath := filepath.Join(tmpDir, fmt.Sprintf("consultations-%s.jsonl", time.Now().AddDate(0, 0, -1).Format("2006-01-02")))
	if _, err := os.Stat(rotatedPath); os.IsNotExist(err) {
		t.Error("Rotated file was not created")
	}

	// Verify logs are correct
	logs, err := logger.GetLogs(2)
	if err != nil {
		t.Fatalf("GetLogs failed: %v", err)
	}

	if len(logs) < 2 {
		t.Errorf("Expected at least 2 logs, got %d", len(logs))
	}
}
