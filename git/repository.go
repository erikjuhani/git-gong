package git

import (
	"errors"
	"fmt"
	"gong/cli"
	"io/ioutil"
	"os"
	"strings"
	"time"

	lib "github.com/libgit2/git2go/v30"
)

// TODO: Move this somewhere more appropriate
func emptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

type Repository struct {
	Core   *lib.Repository
	index  *lib.Index
	unborn bool
}

var (
	DefaultReference = "main"
)

// Free from memory
func (r *Repository) Free() {
	if r.Core != nil {
		r.Core.Free()
	}
	if r.index != nil {
		r.index.Free()
	}
}

func Init(path string, bare bool, initialReference string) (repo *Repository, err error) {
	gitRepo, err := lib.InitRepository(path, bare)
	if err != nil {
		return
	}

	if emptyString(initialReference) {
		initialReference = DefaultReference
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

	unborn, err := gitRepo.IsHeadUnborn()
	if err != nil {
		return
	}

	repo = &Repository{
		Core:   gitRepo,
		index:  idx,
		unborn: unborn,
	}

	return
}

func Open() (repo *Repository, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	gitRepo, err := lib.OpenRepository(wd)
	if err != nil {
		return
	}

	idx, err := gitRepo.Index()
	if err != nil {
		return
	}

	unborn, err := gitRepo.IsHeadUnborn()
	if err != nil {
		return
	}

	repo = &Repository{
		Core:   gitRepo,
		index:  idx,
		unborn: unborn,
	}

	return
}

func (r *Repository) Head() (head *lib.Reference, err error) {
	if r.unborn {
		return nil, errors.New("head is unborn")
	}

	head, err = r.Core.Head()
	if err != nil {
		return
	}

	return
}

func (r *Repository) Changed() (err error) {
	diff, err := r.Core.DiffIndexToWorkdir(
		r.index,
		&lib.DiffOptions{Flags: lib.DiffIncludeUntracked},
	)
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

	status, err := r.Core.StatusList(&lib.StatusOptions{})
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

func (r *Repository) AddToIndex(pathspec []string) (treeID *lib.Oid, err error) {
	if err = r.Changed(); err != nil {
		return
	}

	idx := r.index

	if err = idx.AddAll(pathspec, lib.IndexAddDefault, nil); err != nil {
		return
	}

	treeID, err = idx.WriteTree()
	if err != nil {
		return
	}

	if err = idx.Write(); err != nil {
		return
	}

	return
}

func (r *Repository) createCommit(treeID *lib.Oid, commit *lib.Commit, msg string) (id *lib.Oid, err error) {
	tree, err := r.Core.LookupTree(treeID)
	if err != nil {
		return
	}

	sig := signature()

	if emptyString(msg) {
		input, cliErr := cli.CaptureInput()
		if cliErr != nil {
			return nil, cliErr
		}

		msg = string(input)
	}

	if emptyString(msg) {
		err = errors.New("Aborting due to empty commit message")
		return
	}

	if commit != nil {
		id, err = r.Core.CreateCommit("HEAD", sig, sig, string(msg), tree, commit)
		if err != nil {
			return
		}
		err = r.Core.CheckoutHead(&lib.CheckoutOpts{
			Strategy: lib.CheckoutSafe | lib.CheckoutRecreateMissing,
		})
		return
	}

	// Initial commit
	id, err = r.Core.CreateCommit("HEAD", sig, sig, string(msg), tree)
	if err != nil {
		return
	}

	err = r.Core.CheckoutHead(&lib.CheckoutOpts{
		Strategy: lib.CheckoutSafe | lib.CheckoutRecreateMissing,
	})

	return
}

func (r *Repository) Commit(treeID *lib.Oid, msg string) (commitID *lib.Oid, err error) {
	if r.unborn {
		commitID, err = r.createCommit(treeID, nil, msg)
		if err != nil {
			return
		}

		return
	}

	head, err := r.Head()
	if err != nil {
		return
	}

	currentTip, err := r.Core.LookupCommit(head.Target())
	if err != nil {
		return
	}

	return r.createCommit(treeID, currentTip, msg)
}

func (r *Repository) References() ([]string, error) {
	var list []string
	iter, err := r.Core.NewReferenceIterator()
	if err != nil {
		return list, err
	}

	nameIter := iter.Names()
	name, err := nameIter.Next()
	for err == nil {
		list = append(list, name)
		name, err = nameIter.Next()
	}

	return list, err
}

func (r *Repository) Commits() (commits []*lib.Commit, err error) {
	unborn, err := r.Core.IsHeadUnborn()
	if err != nil {
		return
	}

	if unborn {
		err = errors.New("no existing commits")
		return
	}

	head, err := r.Head()
	if err != nil {
		return
	}
	defer head.Free()

	headCommit, err := r.Core.LookupCommit(head.Target())
	if err != nil {
		return
	}
	defer headCommit.Free()

	commits = append(commits, headCommit)

	if headCommit.ParentCount() != 0 {
		parent := headCommit.Parent(0)
		defer parent.Free()
		commits = append(commits, parent)

		for parent.ParentCount() != 0 {
			parent = parent.Parent(0)
			defer parent.Free()
			commits = append(commits, parent)
		}
	}
	return
}

func (r *Repository) CreateTag(tagname string, message string) (tag *lib.Oid, err error) {
	head, err := r.Head()
	if err != nil {
		return
	}
	defer head.Free()

	headCommit, err := r.Core.LookupCommit(head.Target())
	if err != nil {
		return
	}
	defer headCommit.Free()

	return r.Core.Tags.Create(tagname, headCommit, signature(), message)
}

func (r *Repository) CreateLocalBranch(branchName string) (branch *lib.Branch, err error) {
	head, err := r.Head()
	if err != nil {
		return
	}

	// Check if branch already exists
	localBranch, err := r.Core.LookupBranch(branchName, lib.BranchLocal)
	if localBranch != nil && err != nil {
		return
	}

	// Branch already exists return existing branch and an error stating branch already exists.
	if localBranch != nil {
		return localBranch, fmt.Errorf("branch %s already exists", branchName)
	}

	commit, err := r.Core.LookupCommit(head.Target())
	if err != nil {
		return
	}

	return r.Core.CreateBranch(branchName, commit, false)
}

// TODO get signature from git configuration
func signature() *lib.Signature {
	return &lib.Signature{
		Name:  "gong tester",
		Email: "gong@tester.com",
		When:  time.Now(),
	}
}
