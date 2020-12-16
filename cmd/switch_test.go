package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	lib "github.com/libgit2/git2go/v31"
)

func TestSwitchBranchCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong switch branch <branchname>.
				Should change to a branch with <branchname>.
				If branch does not exist create it. e.g.
				gong switch branch gong-branch ->
				Creates a new branch with name gong-branch, and switches to it.`,
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

	path := fmt.Sprintf("%s/%s", repo.Path, "stash.me")
	if err = ioutil.WriteFile(path, []byte("---i-am-untracked-and-i-shall-be-stashed---\n"), 0644); err != nil {
		t.Fatal(err)
	}

	defer cleanupTestRepo(repo)

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{switchCmd.Name(), switchBranchCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()

			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.FindBranch(tt.args[0], lib.BranchLocal)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSwitchCommitCmd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: `Command gong switch commit <commithash>.
				Should change to a commit with <commithash>. e.g.
				gong switch commit abc3 -> Should change to commit abc3 in the
				currently checked out branch.`,
		},
	}

	repo, err := createTestRepo()
	if err != nil {
		t.Fatal(err)
	}

	firstCommit, err := seedRepo(repo)
	if err != nil {
		t.Fatal(err)
	}

	expectedID := firstCommit.String()

	_, err = seedRepo(repo, "commit.me")
	if err != nil {
		t.Fatal(err)
	}

	defer cleanupTestRepo(repo)

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{switchCmd.Name(), switchCommitCmd.Name()}
			args = append(args, expectedID)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			currentTip, err := repo.Head.Commit()
			if err != nil {
				return
			}

			currentTipID := currentTip.ID.String()

			if currentTipID != expectedID {
				t.Fatal(fmt.Errorf("current tip %s does not equal to expected tip %s", currentTipID, expectedID))
			}
		})
	}
}

func TestSwitchTagCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong switch tag <tagname>. Should change to a tag <tagname>`,
			args: []string{"v0.1.0"},
		},
	}

	repo, err := createTestRepo()
	if err != nil {
		t.Fatal(err)
	}

	firstCommit, err := seedRepo(repo)
	if err != nil {
		t.Fatal(err)
	}

	expectedID := firstCommit.String()

	_, err = repo.CreateTag("v0.1.0", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = seedRepo(repo, "commit.me")
	if err != nil {
		t.Fatal(err)
	}

	defer cleanupTestRepo(repo)

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{switchCmd.Name(), switchTagCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			currentTip, err := repo.Head.Commit()
			if err != nil {
				return
			}

			currentTipID := currentTip.ID.String()

			if currentTipID != expectedID {
				t.Fatal(fmt.Errorf("current tip %s does not equal to expected tip %s", currentTipID, expectedID))
			}
		})
	}
}

func TestSwitchReleaseCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: `Command gong switch release <releasename>. Should change to a release <releasename>`,
			args: []string{"v0.1.0"},
		},
	}

	repo, err := createTestRepo()
	if err != nil {
		t.Fatal(err)
	}

	firstCommit, err := seedRepo(repo)
	if err != nil {
		t.Fatal(err)
	}

	expectedID := firstCommit.String()

	_, err = repo.CreateTag("v0.1.0", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = seedRepo(repo, "commit.me")
	if err != nil {
		t.Fatal(err)
	}

	defer cleanupTestRepo(repo)

	workdir := repo.Path

	if err := os.Chdir(workdir); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{switchCmd.Name(), switchReleaseCmd.Name()}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			currentTip, err := repo.Head.Commit()
			if err != nil {
				return
			}

			currentTipID := currentTip.ID.String()

			if currentTipID != expectedID {
				t.Fatal(fmt.Errorf("current tip %s does not equal to expected tip %s", currentTipID, expectedID))
			}
		})
	}
}
