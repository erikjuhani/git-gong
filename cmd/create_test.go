package cmd

import (
	"fmt"
	"os"
	"testing"

	lib "github.com/libgit2/git2go/v30"
)

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}

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

func TestCreateFileCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong create file <pathspec>.
				Should create a new regular file with <pathspec>.
				e.g. "gong create file module/gongo-bongo.go"
				Creates a file gongo-bongo.go to a directory module.
				If directory does not exists create directory.`,
			args: []string{"module/gongo-bongo.go"},
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
			args := []string{createCmd.Name(), createFileCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if err != nil {
				t.Fatal(err)
			}

			filepath := fmt.Sprintf("%s%s", workdir, tt.args[0])

			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateDirectoryCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong create directory <dirname>. Should create a new directory <dirname>`,
			args: []string{"gong-folder"},
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
			args := []string{createCmd.Name(), createDirectoryCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if err != nil {
				t.Fatal(err)
			}

			dirpath := fmt.Sprintf("%s%s", workdir, tt.args[0])

			if _, err := os.Stat(dirpath); os.IsNotExist(err) {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateTagCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong create tag <tagname>. Should create a new tag <tagname>`,
			args: []string{"v0.1.0"},
		},
	}
	repo, err := createTestRepo()
	if err != nil {
		t.Fatal(err)
	}

	defer cleanupTestRepo(repo)

	_, err = seedRepo(repo)
	if err != nil {
		t.Fatal(err)
	}

	workdir := repo.Core.Workdir()

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{createCmd.Name(), createTagCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if err != nil {
				t.Fatal(err)
			}

			tags, err := repo.Core.Tags.List()
			if err != nil {
				t.Fatal(err)
			}

			if !contains(tags, tt.args[0]) {
				t.Fatal(fmt.Errorf("tag %s does not exist in tags list %q", tt.args[0], tags))
			}
		})
	}
}

func TestCreateReleaseCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong create release <releasename>. Should create a new tag <releasename>`,
			args: []string{"v0.1.0"},
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
			args := []string{createCmd.Name(), createReleaseCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if err != nil {
				t.Fatal(err)
			}

			tags, err := repo.Core.Tags.List()
			if err != nil {
				t.Fatal(err)
			}

			if !contains(tags, tt.args[0]) {
				t.Fatal(fmt.Errorf("release %s does not exist in tags list %q", tt.args[0], tags))
			}
		})
	}
}
