package gong

import (
	"errors"

	git "github.com/libgit2/git2go/v31"
)

type Head struct {
	repository *git.Repository
}

func NewHead(gitRepo *git.Repository) *Head {
	return &Head{repository: gitRepo}
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

	return &Branch{Name: branchName, Shorthand: ref.Shorthand(), core: ref.Branch()}, nil
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

func (head *Head) Reference() (ref *git.Reference, err error) {
	exists, err := head.Exists()
	if err != nil {
		return
	}

	if !exists {
		return nil, errors.New("head does not exist")
	}

	ref, err = head.repository.Head()
	if err != nil {
		return
	}

	return
}

func (head *Head) Detach(id *git.Oid) error {
	return head.repository.SetHeadDetached(id)
}

func (head *Head) IsDetached() (detached bool, err error) {
	return head.repository.IsHeadDetached()
}

func (head *Head) Exists() (exists bool, err error) {
	unborn, err := head.repository.IsHeadUnborn()
	return !unborn, err
}
