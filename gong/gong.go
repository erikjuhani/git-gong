package gong

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrNothingToCommit = errors.New("nothing to commit")
)

var (
	DefaultReference = "main"
)

const (
	headRef = "refs/heads/"
)

func TestRepo() (*testRepository, func(), error) {
	path, err := ioutil.TempDir("", "gong")
	if err != nil {
		return nil, nil, err
	}

	repo, err := Init(path, false, "")
	if err != nil {
		return nil, nil, err
	}

	return &testRepository{repo}, cleanup(repo), nil
}

// TODO: Move this somewhere more appropriate
func checkEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

type freer interface {
	Free()
}

// Free function frees memory for any struct that implements a Free() function.
// This is an utility functionality to free pointers from memory from the underlying libgit.
func Free(f freer) {
	f.Free()
}

type testRepository struct {
	*Repository
}

func (repo *testRepository) Seed(commitMsg string, files ...string) (*Commit, error) {
	if len(files) == 0 {
		files = []string{"README.md", "gongo-bongo.go"}
	}

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", repo.Path, f)
		if err := ioutil.WriteFile(path, []byte("temp\n"), 0644); err != nil {
			return nil, err
		}
	}

	tree, err := repo.AddToIndex(files)
	if err != nil {
		return nil, err
	}
	defer Free(tree)

	return repo.CreateCommit(tree, commitMsg)
}

func cleanup(r *Repository) func() {
	return func() {
		if err := os.RemoveAll(r.Path); err != nil {
			panic(err)
		}
	}
}
