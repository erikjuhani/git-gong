package gong

import (
	"fmt"
	"io/ioutil"
	"os"
)

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

func cleanup(r *Repository) func() {
	return func() {
		if err := os.RemoveAll(r.Path); err != nil {
			panic(err)
		}
	}
}
