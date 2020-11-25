package cmd

import (
	"os"
	"testing"

	lib "github.com/libgit2/git2go/v30"
)

func TestCreateBranchCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong create branch <branchname>.
			Should create a new branch with <branchname>`,
			args: []string{"gong-branch"},
		},
	}
	repo, err := createTestRepo()
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
			args := []string{createCmd.Name(), createBranchCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Core.LookupBranch(tt.args[0], lib.BranchLocal)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
