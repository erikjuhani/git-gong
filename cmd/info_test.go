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
		files    []string
	}{
		{
			name:     `Should display no changes message if no files have changed.`,
			expected: "Branch %s\nCommit %s\n\nNo changes.\n",
		},
		{
			name:     `Should display difference between the index file and the current HEAD in short format.`,
			expected: "Branch %s\nCommit %s\n\nChanges (2):\n M README.md\n?? beam-me-up.scotty\n",
			files:    []string{"beam-me-up.scotty", "README.md"},
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

			for _, f := range tt.files {
				_, err := os.Create(f)
				if err != nil {
					t.Fatal(err)
				}
			}

			rootCmd.SetArgs(args)

			outBuff := bytes.NewBuffer(nil)

			rootCmd.SetOut(outBuff)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			out, err := ioutil.ReadAll(outBuff)
			if err != nil {
				t.Fatal(err)
			}

			cb, err := repo.CurrentBranch()
			if err != nil {
				t.Fatal(err)
			}

			hc, err := repo.HeadCommit()
			if err != nil {
				t.Fatal(err)
			}

			expected := fmt.Sprintf(tt.expected, cb.Shorthand(), hc.Id().String())

			if !bytes.Equal(out, []byte(expected)) {
				t.Fatal(fmt.Errorf("actual info output %s did not match the expected %s output", out, expected))
			}
		})
	}
}