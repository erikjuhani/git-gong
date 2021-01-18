package gong

import (
	"errors"

	git "github.com/libgit2/git2go/v31"
)

const RefName = "HEAD"

var (
	ErrHeadNotExists  = errors.New("head does not exist")
	ErrEmptyCommitMsg = errors.New("aborting due to empty commit message")
)

type Head struct {
	RefName    string
	repository *git.Repository
}

func NewHead(gitRepo *git.Repository) *Head {
	return &Head{RefName: RefName, repository: gitRepo}
}

func (head *Head) SetReference(refName string) error {
	return head.repository.SetHead(refName)
}

func (head *Head) Branch() (branch *Branch, err error) {
	ref, err := head.Reference()
	if err != nil {
		return
	}

	branchName, err := ref.Branch().Name()
	if err != nil {
		return
	}

	return NewBranch(branchName, ref.Branch()), nil
}

func (head *Head) Commit() (commit *Commit, err error) {
	ref, err := head.Reference()
	if err != nil {
		return
	}

	gitCommit, err := head.repository.LookupCommit(ref.Target())
	if err != nil {
		return
	}

	return NewCommit(gitCommit), nil
}

func (head *Head) Reference() (*git.Reference, error) {
	exists, err := head.Exists()
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrHeadNotExists
	}

	return head.repository.Head()
}

func (head *Head) Detach(id *git.Oid) error {
	return head.repository.SetHeadDetached(id)
}

func (head *Head) IsDetached() (bool, error) {
	return head.repository.IsHeadDetached()
}

func (head *Head) Exists() (bool, error) {
	unborn, err := head.repository.IsHeadUnborn()
	return !unborn, err
}

func (head *Head) Checkout() error {
	opts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}
	return head.repository.CheckoutHead(opts)
}
