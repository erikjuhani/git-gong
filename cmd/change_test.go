package cmd

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/erikjuhani/git-gong/gong"
)

func TestChangeCmd(t *testing.T) {
	repo, clean, err := gong.TestRepo()
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	tests := []struct {
		name     string
		args     []string
		fn       func(t *testing.T) string
		expected string
	}{
		{
			name: `Command gong change branch [<branchname>] --rename <newbranchname>. Should change branch <branchname> name as <newbranchname>`,
			args: []string{"branch", "main", "--rename", "dev"},
			fn: func(t *testing.T) string {
				branch, err := repo.CurrentBranch()
				if err != nil {
					t.Fatal(err)
				}
				return branch.Name
			},
			expected: "dev",
		},
		{
			name: `Command gong change file [<filename>] --rename <newfilename>. Should change filename <filename> to <newfilename>.`,
			args: []string{"file", "README.md", "--rename", "NOTREADME.md"},
			fn: func(t *testing.T) string {
				filepath := fmt.Sprintf("%s%s", repo.Path, "README.md")

				if _, err := os.Stat(filepath); !os.IsNotExist(err) {
					t.Fatal(errors.New("file should not exist"))
				}

				filepath = fmt.Sprintf("%s%s", repo.Path, "NOTREADME.md")
				if _, err := os.Stat(filepath); os.IsNotExist(err) {
					t.Fatal(err)
				}

				return "NOTREADME"
			},
			expected: "NOTREADME.md",
		},
		{
			name: `Command gong change commit --message [<newmessage>]. Should change the last commit message with <newmessage>.`,
			args: []string{"commit", "--message", "Changed commit message"},
			fn: func(t *testing.T) string {
				commit, err := repo.Head.Commit()
				if err != nil {
					t.Fatal(err)
				}
				return commit.Message
			},
			expected: "Changed commit message",
		},
	}

	_, err = repo.Seed("commitmsg")
	if err != nil {
		t.Fatal(err)
	}

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{changeCmd.Name()}
			args = append(args, tt.args...)
			rootCmd.SetArgs(args)

			if err := rootCmd.Execute(); err != nil {
				t.Fatal(err)
			}

			actual := tt.fn(t)

			if actual != tt.expected {
				t.Fatal(fmt.Errorf("expected state: %s did not match the actual state: %s", tt.expected, actual))
			}
		})
	}
}
