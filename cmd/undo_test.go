package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/erikjuhani/git-gong/gong"
)

func TestUndoCmd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: `Command gong undo. Should undo the last command the user has done.
e.g. User commit file to wrong branch.`,
		},
	}

	repo, clean, err := gong.TestRepo()
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	expected, err := repo.Seed("a")
	if err != nil {
		t.Fatal(err)
	}

	// Mistake commit
	_, err = repo.Seed("b", "b.file")
	if err != nil {
		t.Fatal(err)
	}

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{undoCmd.Name()}
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
