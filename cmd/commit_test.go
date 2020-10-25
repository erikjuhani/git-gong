package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	git "github.com/libgit2/git2go/v30"
)

func createTestRepo() (*git.Repository, error) {
	path, err := ioutil.TempDir("", "gong-clone")
	if err != nil {
		return nil, err
	}

	repo, err := git.InitRepository(path, false)
	if err != nil {
		return nil, err
	}

	sig := &git.Signature{
		Name:  "gong tester",
		Email: "gong@tester.com",
		When:  time.Now(),
	}

	tmpfile := "TMPGONG"
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, tmpfile), []byte{}, 0644)

	idx, err := repo.Index()
	if err != nil {
		return nil, err
	}

	err = idx.AddByPath("TMPGONG")
	if err != nil {
		return nil, err
	}

	err = idx.Write()
	if err != nil {
		return nil, err
	}

	treeID, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return nil, err
	}

	_, err = repo.CreateCommit("HEAD", sig, sig, "test", tree)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func cleanupTestRepo(r *git.Repository) {
	if err := os.RemoveAll(r.Workdir()); err != nil {
		panic(err)
	}
	r.Free()
}

func TestCommitCmd(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		stageOnly bool
	}{
		{
			name: `Should stage and record changes to repository.
			Use <pathspec> or <pathpattern> for files that has been changed
			and are ready to be staged and recorded to repository.
			If no arguments were given add all changed files to stage,
			commit and record to repository.`,
			files:     []string{"gong_file"},
			stageOnly: false,
		},
		{
			name: `Should stage and record changes to repository.
			Use <pathspec> or <pathpattern> for files that has been changed
			and are ready to be staged and recorded to repository.
			If no arguments were given add all changed files to stage.`,
			files:     []string{"gong_file", "another_gong_file"},
			stageOnly: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := createTestRepo()
			if err != nil {
				t.Fatal(err)
			}

			defer cleanupTestRepo(repo)

			if err := os.Chdir(repo.Workdir()); err != nil {
				t.Fatal(err)
			}

			args := []string{commitCmd.Name()}

			rootCmd.SetArgs(args)

			for _, f := range tt.files {
				err := ioutil.WriteFile(fmt.Sprintf("%s/%s", repo.Workdir(), f), []byte{}, 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			head, err := repo.Head()
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.LookupCommit(head.Target())
			if err != nil {
				t.Fatal(err)
			}

			// TODO implement testing phase fully after code is done.
		})
	}
}
