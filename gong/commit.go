package gong

import (
	git "github.com/libgit2/git2go/v31"
)

type Commit struct {
	ID      *git.Oid
	Message string
	core    *git.Commit
}

func NewCommit(commit *git.Commit) *Commit {
	id := commit.Id()
	msg := commit.Message()

	return &Commit{
		ID:      id,
		Message: msg,
		core:    commit,
	}
}

func (commit *Commit) Tree() (*git.Tree, error) {
	return commit.core.Tree()
}
