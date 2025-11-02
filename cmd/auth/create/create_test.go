package create

import (
	"strings"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	cmd := NewCreateCmd()

	if cmd.Use != "create" {
		t.Errorf("Use = %v, want create", cmd.Use)
	}

	if cmd.Flag("name") == nil {
		t.Error("name flag not found")
	}
	if cmd.Flag("root-path") == nil {
		t.Error("root-path flag not found")
	}
	if cmd.Flag("listen-address") == nil {
		t.Error("listen-address flag not found")
	}
}

func TestCreateCommandInteractiveInput(t *testing.T) {
	cmd := NewCreateCmd()

	cmd.SetIn(strings.NewReader("TestAccount\n"))
	cmd.SetArgs([]string{})

	if cmd.RunE == nil {
		t.Error("RunE function is not set")
	}
}
