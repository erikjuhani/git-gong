package cmd

import (
	"fmt"
	"gong/git"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func createTestRepo() (*git.Repository, error) {
	path, err := ioutil.TempDir("", "gong")
	if err != nil {
		return nil, err
	}

	return git.Init(path, false, "")
}

func cleanupTestRepo(r *git.Repository) {
	if err := os.RemoveAll(r.Core.Workdir()); err != nil {
		panic(err)
	}
	r.Core.Free()
}

func TestCommitCmd(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		args  []string
	}{
		{
			name: `Should stage and record changes to repository.
			Use <pathspec> or <pathpattern> for files that has been changed
			and are ready to be staged and recorded to repository.
			If no arguments were given add all changed files to stage,
			commit and record to repository.`,
			files: []string{"gong_file"},
			args:  []string{""},
		},
		{
			name: `Should stage and record changes to repository.
			Use <pathspec> or <pathpattern> for files that has been changed
			and are ready to be staged and recorded to repository.
			If no arguments were given add all changed files to stage.`,
			files: []string{"gong_file", "another_gong_file_in_a_new_directory"},
			args:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := createTestRepo()
			if err != nil {
				t.Fatal(err)
			}

			defer cleanupTestRepo(repo)

			if err := os.Chdir(repo.Core.Workdir()); err != nil {
				t.Fatal(err)
			}

			args := []string{commitCmd.Name(), "-m", "testmsg"}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			for i, f := range tt.files {
				var err error
				if i == 1 {
					dirName, err := ioutil.TempDir(repo.Core.Workdir(), "tmp")
					if err != nil {
						t.Fatal(err)
					}
					log.Println(dirName)
					err = ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", repo.Core.Workdir(), dirName, f), []byte{}, 0644)
				} else {
					err = ioutil.WriteFile(fmt.Sprintf("%s/%s", repo.Core.Workdir(), f), []byte{}, 0644)
				}

				if err != nil {
					t.Fatal(err)
				}
			}

			err = rootCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			head, err := repo.Core.Head()
			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Core.LookupCommit(head.Target())
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
