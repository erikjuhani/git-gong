package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestInfoCmd(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		args     []string
	}{
		{
			name:     `Should display difference between the index file and the current HEAD in short format.`,
			expected: "write later",
			args:     []string{""},
		},
		{
			name:     `Should show file changes in changed files.`,
			expected: "write later",
			args:     []string{"--changes"},
		},
	}

	repo, err := createTestRepo()
	if err != nil {
		t.Fatal(err)
	}

	_, err = seedRepo(repo)
	if err != nil {
		t.Fatal(err)
	}

	defer cleanupTestRepo(repo)

	workdir := repo.Core.Workdir()

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{infoCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			b := bytes.NewBufferString("")

			rootCmd.SetOut(b)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}

			if string(out) != tt.expected {
				t.Fatal(fmt.Errorf("actual info output %s did not match the expected %s output", out, tt.expected))
			}
		})
	}
}
