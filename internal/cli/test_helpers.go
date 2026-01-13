package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

// validateJSONOutput validates that the output is valid JSON.
func validateJSONOutput(output string) bool {
	if len(output) == 0 {
		return false
	}
	var v interface{}
	return json.Unmarshal([]byte(output), &v) == nil
}

// captureOutput captures stdout from a function call.
// Returns the captured output and any error.
func captureOutput(fn func() error) (string, error) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	var buf bytes.Buffer
	done := make(chan bool)
	go func() {
		_, _ = buf.ReadFrom(r)
		done <- true
	}()

	err := fn()

	w.Close()
	os.Stdout = oldStdout
	<-done

	return buf.String(), err
}

// runCLITest runs a CLI test with output capture and validation.
func runCLITest(t *testing.T, name string, fn func() error, wantErr bool, check func(string) bool) {
	t.Run(name, func(t *testing.T) {
		output, err := captureOutput(fn)

		if (err != nil) != wantErr {
			t.Errorf("error = %v, wantErr %v", err, wantErr)
			return
		}

		if !wantErr && check != nil {
			if !check(output) {
				t.Errorf("output validation failed. Output: %s", output)
			}
		}
	})
}
