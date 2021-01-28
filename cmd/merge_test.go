package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/erikjuhani/git-gong/gong"
)

func TestMergeCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong merge <branchname>. Merge to <currentbranch>.
			e.g. gong merge -> Should merge the given <branch> to <currentbranch>`,
			args: []string{"gong-branch"},
		},
	}

	repo, clean, err := gong.TestRepo()
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	_, err = repo.Seed("default-commit")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.CheckoutBranch("gong-branch")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := repo.Seed("gong-branch-commit", "a.file")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.CheckoutBranch("main")
	if err != nil {
		t.Fatal(err)
	}

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{mergeCmd.Name()}
			args = append(args, tt.args...)
			rootCmd.SetArgs(args)

			if err := rootCmd.Execute(); err != nil {
				t.Fatal(err)
			}

			actual, err := repo.Head.Commit()
			if err != nil {
				t.Fatal(err)
			}

			if expected.ID.String() != actual.ID.String() {
				t.Fatal(fmt.Errorf("expected head state: %s did not match the actual state: %s", expected.ID.String(), actual.ID.String()))
			}
		})
	}
}
