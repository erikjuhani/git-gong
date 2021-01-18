package gong

import git "github.com/libgit2/git2go/v31"

type Tag struct {
	ID       *git.Oid
	CommitID *git.Oid
	Name     string
	essence  *git.Tag
}

func NewTag(gitTag *git.Tag) *Tag {
	return &Tag{
		ID:       gitTag.Id(),
		CommitID: gitTag.TargetId(),
		Name:     gitTag.Name(),
		essence:  gitTag,
	}
}

func (tag *Tag) Essence() *git.Tag {
	return tag.essence
}

func (tag *Tag) Free() {
	tag.Essence().Free()
}
