// Package logging provides consultation logging functionality with date-based rotation.
// It writes consultations to JSONL files and supports automatic log rotation by date.
package logging

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// ConsultationLogger handles consultation logging to JSONL file with date-based rotation.
// It is thread-safe and automatically rotates log files when the date changes.
type ConsultationLogger struct {
	mu          sync.Mutex
	logDir      string
	filePath    string
	file        *os.File
	encoder     *json.Encoder
	currentDate string // Track current date for rotation (YYYY-MM-DD format)
}

// NewConsultationLogger creates a new consultation logger.
// logDir is the directory where log files will be stored (e.g., ".devwisdom").
// If a log file exists from a previous date, it will be automatically rotated.
func NewConsultationLogger(logDir string) (*ConsultationLogger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create consultation log directory %q: %w", logDir, err)
	}

	// Log file path
	filePath := filepath.Join(logDir, "consultations.jsonl")

	// Get current date for rotation tracking
	currentDate := time.Now().Format("2006-01-02")

	// Check if file exists and get its modification date
	var file *os.File
	var err error
	if info, err := os.Stat(filePath); err == nil {
		// File exists - check if it needs rotation
		fileDate := info.ModTime().Format("2006-01-02")
		if fileDate != currentDate {
			// File is from a different date - rotate it
			rotatedPath := filepath.Join(logDir, fmt.Sprintf("consultations-%s.jsonl", fileDate))
			if err := os.Rename(filePath, rotatedPath); err != nil {
				return nil, fmt.Errorf("failed to rotate log file from %q to %q: %w", filePath, rotatedPath, err)
			}
		}
	}

	// Open or create log file in append mode
	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %q: %w", filePath, err)
	}

	logger := &ConsultationLogger{
		logDir:      logDir,
		filePath:    filePath,
		file:        file,
		encoder:     json.NewEncoder(file),
		currentDate: currentDate,
	}

	return logger, nil
}

// rotateIfNeeded checks if log rotation is needed based on date change
// Must be called with mutex held
func (l *ConsultationLogger) rotateIfNeeded() error {
	currentDate := time.Now().Format("2006-01-02")

	// If date hasn't changed, no rotation needed
	if l.currentDate == currentDate {
		return nil
	}

	// Date changed - need to rotate
	// Close current file
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("failed to close log file for rotation: %w", err)
		}
	}

	// Rename current file with date suffix (atomic operation)
	rotatedPath := filepath.Join(l.logDir, fmt.Sprintf("consultations-%s.jsonl", l.currentDate))
	if err := os.Rename(l.filePath, rotatedPath); err != nil {
		// If rename fails, try to reopen the file
		file, err2 := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			return fmt.Errorf("failed to rotate log file and failed to reopen: %w (original: %v)", err2, err)
		}
		l.file = file
		l.encoder = json.NewEncoder(file)
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Create new log file for current date
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file after rotation: %w", err)
	}

	l.file = file
	l.encoder = json.NewEncoder(file)
	l.currentDate = currentDate

	return nil
}

// Log writes a consultation to the JSONL log file
// Thread-safe: uses mutex to protect concurrent writes
// Automatically rotates log file if date has changed
func (l *ConsultationLogger) Log(consultation *wisdom.Consultation) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if rotation is needed (date-based)
	if err := l.rotateIfNeeded(); err != nil {
		return fmt.Errorf("log rotation failed: %w", err)
	}

	// Encode consultation as JSON and write as single line
	if err := l.encoder.Encode(consultation); err != nil {
		return fmt.Errorf("failed to encode consultation: %w", err)
	}

	// Flush to ensure data is written
	if err := l.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync log file: %w", err)
	}

	return nil
}

// readLogFile reads consultations from a single log file
func (l *ConsultationLogger) readLogFile(filePath string, cutoffTime time.Time) ([]*wisdom.Consultation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty array
			return []*wisdom.Consultation{}, nil
		}
		return nil, fmt.Errorf("failed to open log file %s: %w", filePath, err)
	}
	defer file.Close()

	var consultations []*wisdom.Consultation

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue // Skip empty lines
		}

		// Parse JSON line into Consultation
		var consultation wisdom.Consultation
		if err := json.Unmarshal(line, &consultation); err != nil {
			// Skip malformed JSON lines (log error but continue)
			continue
		}

		// Parse timestamp and filter by date
		timestamp, err := time.Parse(time.RFC3339, consultation.Timestamp)
		if err != nil {
			// Skip entries with invalid timestamps
			continue
		}

		// Include consultations from cutoff time onwards
		if timestamp.After(cutoffTime) || timestamp.Equal(cutoffTime) {
			consultations = append(consultations, &consultation)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file %s: %w", filePath, err)
	}

	return consultations, nil
}

// GetLogs retrieves consultations from log files filtered by days
// Returns consultations from the last N days, reading from both current and rotated files
func (l *ConsultationLogger) GetLogs(days int) ([]*wisdom.Consultation, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Calculate cutoff time
	cutoffTime := time.Now().AddDate(0, 0, -days)

	var allConsultations []*wisdom.Consultation

	// Read from current log file
	consultations, err := l.readLogFile(l.filePath, cutoffTime)
	if err != nil {
		return nil, err
	}
	allConsultations = append(allConsultations, consultations...)

	// Read from rotated log files (consultations-YYYY-MM-DD.jsonl)
	entries, err := os.ReadDir(l.logDir)
	if err != nil {
		// If we can't read directory, just return what we got from current file
		return allConsultations, nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Check if it's a rotated log file (starts with "consultations-" and ends with ".jsonl")
		if strings.HasPrefix(name, "consultations-") && strings.HasSuffix(name, ".jsonl") && name != "consultations.jsonl" {
			// Extract date from filename (consultations-YYYY-MM-DD.jsonl)
			dateStr := strings.TrimPrefix(name, "consultations-")
			dateStr = strings.TrimSuffix(dateStr, ".jsonl")

			// Parse date to check if it's within range
			fileDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				// Skip files with invalid date format
				continue
			}

			// Only read files that might contain logs within our date range
			// (file date should be >= cutoff date)
			if fileDate.After(cutoffTime) || fileDate.Equal(cutoffTime) {
				rotatedPath := filepath.Join(l.logDir, name)
				consultations, err := l.readLogFile(rotatedPath, cutoffTime)
				if err != nil {
					// Log error but continue with other files
					continue
				}
				allConsultations = append(allConsultations, consultations...)
			}
		}
	}

	return allConsultations, nil
}

// Close closes the log file
func (l *ConsultationLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
