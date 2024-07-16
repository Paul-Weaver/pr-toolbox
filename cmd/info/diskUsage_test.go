package info

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand executes a cobra.Command with given arguments, capturing its output.
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
	return buf.String(), err
}

// TestDiskUsageCmd tests the diskUsage command.
func TestDiskUsageCmd(t *testing.T) {
	output, err := executeCommand(InfoCmd, "diskUsage")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if output == "" {
		t.Error("Expected non-empty output")
	} else {
		fmt.Printf("Command output: %s\n", output)
	}
}
