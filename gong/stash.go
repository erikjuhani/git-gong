package gong

import git "github.com/libgit2/git2go/v31"

type Stash struct {
	ID      *git.Oid
	Message string
	Index   int
}
