package gong

import (
	git "github.com/libgit2/git2go/v31"
)

type Commit struct {
	ID      *git.Oid
	Message string
	essence *git.Commit
}

func (commit *Commit) Essence() *git.Commit {
	return commit.essence
}

func NewCommit(commit *git.Commit) *Commit {
	id := commit.Id()
	msg := commit.Message()

	return &Commit{
		ID:      id,
		Message: msg,
		essence: commit,
	}
}

func (commit *Commit) Parent() *Commit {
	return NewCommit(commit.essence.Parent(0))
}

func (commit *Commit) HasChildren() bool {
	return commit.essence.ParentCount() != 0
}

func (commit *Commit) Tree() (*git.Tree, error) {
	return commit.essence.Tree()
}

func (commit *Commit) Free() {
	commit.Essence().Free()
}
