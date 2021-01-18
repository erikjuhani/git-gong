package gong

import git "github.com/libgit2/git2go/v31"

type Branch struct {
	ReferenceID *git.Oid
	Name        string
	Shorthand   string
	essence     *git.Branch
}

func NewBranch(branchName string, gitBranch *git.Branch) *Branch {
	return &Branch{
		ReferenceID: gitBranch.Target(),
		Name:        branchName,
		Shorthand:   gitBranch.Shorthand(),
		essence:     gitBranch,
	}
}

func (branch *Branch) Essence() *git.Branch {
	return branch.essence
}

func (branch *Branch) Free() {
	branch.Essence().Free()
}
