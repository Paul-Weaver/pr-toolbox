package gitutils

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupRepo(t *testing.T) (func(), string) {
	dir, err := os.MkdirTemp("", "testrepo")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	commands := [][]string{
		{"git", "init"},
		{"sh", "-c", "echo 'hello world' > file.txt"},
		{"git", "add", "file.txt"},
		{"git", "commit", "-m", "initial commit"},
		{"git", "checkout", "-b", "main"}, // Ensure we create a "main" branch
		{"sh", "-c", "echo 'new change' >> file.txt"},
		{"git", "commit", "-am", "second commit"},
	}

	for _, cmd := range commands {
		command := exec.Command(cmd[0], cmd[1:]...)
		command.Dir = dir
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("Failed to run %v: %v, output: %s", cmd, err, output)
		}
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return cleanup, dir
}

func TestGetBaseBranch(t *testing.T) {
	cleanup, dir := setupRepo(t)
	defer cleanup()

	os.Chdir(dir)
	defer os.Chdir("..")

	baseBranch, err := GetBaseBranch()
	assert.NoError(t, err)
	assert.Contains(t, []string{"main", "master"}, baseBranch)
}

func TestGetGitDiff(t *testing.T) {
	cleanup, dir := setupRepo(t)
	defer cleanup()

	os.Chdir(dir)
	defer os.Chdir("..")

	baseBranch, err := GetBaseBranch()
	assert.NoError(t, err)

	diff, err := GetGitDiff(baseBranch)
	assert.NoError(t, err)
	assert.Contains(t, diff, "new change")
}
