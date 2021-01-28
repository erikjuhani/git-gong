package gong

import git "github.com/libgit2/git2go/v31"

type Branch struct {
	ReferenceID *git.Oid
	RefName     string
	Name        string
	Shorthand   string
	essence     *git.Branch
}

func NewBranch(branchName string, gitBranch *git.Branch) *Branch {
	return &Branch{
		ReferenceID: gitBranch.Target(),
		RefName:     gitBranch.Reference.Name(),
		Name:        branchName,
		Shorthand:   gitBranch.Shorthand(),
		essence:     gitBranch,
	}
}

func (branch *Branch) AnnotatedCommit() (*git.AnnotatedCommit, error) {
	owner := branch.Essence().Owner()
	defer Free(owner)

	return owner.AnnotatedCommitFromRef(branch.Essence().Reference)
}

func (branch *Branch) Essence() *git.Branch {
	return branch.essence
}

func (branch *Branch) Free() {
	branch.Essence().Free()
}
