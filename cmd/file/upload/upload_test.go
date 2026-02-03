package upload

import (
	"testing"
)

func TestUploadCommand(t *testing.T) {
	cmd := NewUploadCmd()

	// Check flags
	spaceIdFlag := cmd.Flag("space-id")
	if spaceIdFlag == nil {
		t.Fatal("space-id flag not found")
	}

	// Check arguments requirement
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error for missing arguments, got nil")
	}
}
