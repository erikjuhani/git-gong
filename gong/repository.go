package gong

import (
	"errors"
	"fmt"
	"gong/cli"
	"io/ioutil"
	"os"
	"time"

	git "github.com/libgit2/git2go/v30"
)

type Repository struct {
	core  *git.Repository
	index *git.Index
}

func (r *Repository) clean() {
	r.index.Free()
}

func Init(path string, bare bool, initialReference string) (repo *Repository, err error) {
	gitRepo, err := git.InitRepository(path, bare)
	if err != nil {
		return
	}

	initRef := fmt.Sprintf("refs/heads/%s", initialReference)
	err = ioutil.WriteFile(fmt.Sprintf("%s/HEAD", gitRepo.Path()), []byte("ref: "+initRef), 0644)
	if err != nil {
		return
	}

	idx, err := gitRepo.Index()
	if err != nil {
		return
	}

	repo = &Repository{
		core:  gitRepo,
		index: idx,
	}

	return
}

func Open() (repo *Repository, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	gitRepo, err := git.OpenRepository(wd)
	if err != nil {
		return
	}

	idx, err := gitRepo.Index()
	if err != nil {
		return
	}

	repo = &Repository{
		core:  gitRepo,
		index: idx,
	}

	return
}

func (r *Repository) Changed() (err error) {
	diff, err := r.core.DiffIndexToWorkdir(r.index, nil)
	if err != nil {
		return
	}

	defer diff.Free()

	stats, err := diff.Stats()
	if err != nil {
		return
	}

	defer stats.Free()

	changedFiles := stats.FilesChanged()

	status, err := r.core.StatusList(&git.StatusOptions{})
	if err != nil {
		return
	}

	entries, err := status.EntryCount()
	if err != nil {
		return
	}

	defer status.Free()

	if changedFiles == 0 && entries == 0 {
		err = errors.New("no files changed, nothing to commit, working tree clean")
		return
	}

	return
}

func (r *Repository) AddToIndex(pathspec []string) (treeID *git.Oid, err error) {
	if err = r.Changed(); err != nil {
		return
	}

	idx := r.index

	if err = idx.AddAll(pathspec, git.IndexAddDefault, nil); err != nil {
		return
	}

	if err = idx.Write(); err != nil {
		return
	}

	treeID, err = idx.WriteTree()
	if err != nil {
		return
	}

	return
}

func (r *Repository) createCommit(treeID *git.Oid, commit *git.Commit) (id *git.Oid, err error) {
	tree, err := r.core.LookupTree(treeID)
	if err != nil {
		return
	}

	sig := signature()

	msg, err := cli.CaptureInput()
	if err != nil {
		return
	}

	if len(msg) == 0 {
		err = errors.New("Aborting due to empty commit message")
		return
	}

	if commit != nil {
		return r.core.CreateCommit("HEAD", sig, sig, string(msg), tree, commit)
	}

	// Initial commit
	return r.core.CreateCommit("HEAD", sig, sig, string(msg), tree)
}

func (r *Repository) Commit(treeID *git.Oid) (commitID *git.Oid, err error) {
	unborn, err := r.core.IsHeadUnborn()
	if err != nil {
		return
	}

	if unborn {
		_, err = r.createCommit(treeID, nil)
		return
	}

	head, err := r.core.Head()
	if err != nil {
		return
	}

	currentTip, err := r.core.LookupCommit(head.Target())
	if err != nil {
		return
	}

	return r.createCommit(treeID, currentTip)
}

// TODO get signature from git configuration
func signature() *git.Signature {
	return &git.Signature{
		Name:  "gong tester",
		Email: "gong@tester.com",
		When:  time.Now(),
	}
}
