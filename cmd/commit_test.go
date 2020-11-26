package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gong/git"
	"io/ioutil"
	"os"
	"path"
	"testing"

	lib "github.com/libgit2/git2go/v30"
)

func createTestRepo() (*git.Repository, error) {
	path, err := ioutil.TempDir("", "gong")
	if err != nil {
		return nil, err
	}

	return git.Init(path, false, "")
}

func seedRepo(repo *git.Repository) (commidID *lib.Oid, err error) {
	files := []string{"README.md", "gongo-bongo.go"}

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", repo.Core.Workdir(), f)
		if err = ioutil.WriteFile(path, []byte("temp\n"), 0644); err != nil {
			return
		}
	}

	treeID, err := repo.AddToIndex(files)
	if err != nil {
		return
	}

	return repo.Commit(treeID, commitMsg)
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
			files: []string{"README.md"},
			args:  []string{""},
		},
		{
			name: `Should stage and record changes to repository.
			Use <pathspec> or <pathpattern> for files that has been changed
			and are ready to be staged and recorded to repository.
			If no arguments were given add all changed files to stage.`,
			files: []string{"gong_test.go", "another_gong_file_in_a_new_directory.go", ".env"},
			args:  []string{""},
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
			args := []string{commitCmd.Name(), "-m", "testmsg"}
			args = append(args, tt.args...)

			rootCmd.SetArgs(args)

			newFiles := make(map[string]struct{})

			for i, f := range tt.files {
				var err error
				var file string
				if i == 1 {
					dirname, err := ioutil.TempDir(workdir, "tmp")
					if err != nil {
						t.Fatal(err)
					}
					base := path.Base(dirname)
					file = fmt.Sprintf("%s/%s", dirname, f)
					newFiles[fmt.Sprintf("%s/%s", base, f)] = struct{}{}
					err = ioutil.WriteFile(file, []byte{}, os.ModePerm)
				} else {
					file = fmt.Sprintf("%s/%s", workdir, f)
					newFiles[f] = struct{}{}
					err = ioutil.WriteFile(file, []byte{}, os.ModePerm)
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

			commit, err := repo.Core.LookupCommit(head.Target())
			if err != nil {
				t.Fatal(err)
			}

			tree, err := repo.Core.LookupTree(commit.TreeId())
			if err != nil {
				t.Fatal(err)
			}

			commits, err := repo.Commits()
			if err != nil {
				t.Fatal(err)
			}

			if len(commits) == 0 {
				t.Fatal(errors.New("commit was not recorded to repository"))
			}

			if commit.Id().String() != commits[0].Id().String() {
				t.Fatal(errors.New("commit head does not match commits"))
			}

			if len(commits) > 1 {
				parentTree, err := commits[1].Tree()
				if err != nil {
					t.Fatal(err)
				}
				diff, err := repo.Core.DiffTreeToTree(parentTree, tree, nil)
				if err != nil {
					t.Fatal(err)
				}

				raw, err := diff.ToBuf(lib.DiffFormatNameOnly)
				if err != nil {
					t.Fatal(err)
				}

				scanner := bufio.NewScanner(bytes.NewReader(raw))
				var includedFiles []string
				for scanner.Scan() {
					includedFiles = append(includedFiles, scanner.Text())
				}

				if len(includedFiles) != len(newFiles) {
					t.Fatal(errors.New("new files do not match diff of commit"))
				}

				for _, f := range includedFiles {
					if _, ok := newFiles[f]; !ok {
						t.Fatal(errors.New("diff file does not match new file"))
					}
				}
			}

			if tree.EntryCount() < 1 {
				t.Fatal(errors.New("commit not found in tree"))
			}
		})
	}
}
