package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xrpl-commons/bedrock/internal/cli"
)

// execute captures the output of a command for testing
func execute(args ...string) (string, error) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set the arguments for the root command
	cli.RootCmd.SetArgs(args)

	err := cli.RootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String(), err
}

func TestRootCommand(t *testing.T) {
	// Test the root command with no arguments
	output, err := execute()
	assert.NoError(t, err)
	assert.Contains(t, output, "BEDROCK - XRPL Smart Contract CLI")

	// Test the help command
	output, err = execute("help")
	assert.NoError(t, err)
	assert.Contains(t, output, "BEDROCK - XRPL Smart Contract CLI")

	// Test an unknown command
	_, err = execute("unknown")
	assert.Error(t, err)
}

func TestBuildCommand_NoConfig(t *testing.T) {
	_, err := execute("build")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}
