package gong

import (
	"fmt"

	git "github.com/libgit2/git2go/v31"
)

type StashCollection struct {
	essence *git.StashCollection
	stashes map[string]*Stash
}

func NewStashCollection(gitStash *git.StashCollection) *StashCollection {
	stashes := make(map[string]*Stash)

	gitStash.Foreach(func(index int, message string, id *git.Oid) error {
		stashes[message] = &Stash{ID: id, Message: message, Index: index}
		return nil
	})

	return &StashCollection{
		essence: gitStash,
		stashes: stashes,
	}
}

func (collection *StashCollection) Essence() *git.StashCollection {
	return collection.essence
}

func (collection *StashCollection) Create(currentBranch *Branch) (*Stash, error) {
	branchID := currentBranch.ReferenceID.String()
	stashID, err := collection.essence.Save(signature(), branchID, git.StashIncludeUntracked)
	if err != nil {
		return nil, err
	}

	stash := &Stash{ID: stashID, Message: branchID, Index: 0}
	collection.stashes[branchID] = stash

	return stash, nil
}

func (collection *StashCollection) Find(branch *Branch) (*Stash, error) {
	if stash, ok := collection.stashes[branch.ReferenceID.String()]; ok {
		return stash, nil
	}

	return nil, fmt.Errorf("stash with name %s was not found", branch.Name)
}

func (collection *StashCollection) Has(branch *Branch) bool {
	_, ok := collection.stashes[branch.ReferenceID.String()]
	return ok
}

func (collection *StashCollection) Pop(branch *Branch) error {
	stash, err := collection.Find(branch)
	if err != nil {
		return err
	}

	opts, err := git.DefaultStashApplyOptions()
	if err != nil {
		return err
	}

	if err := collection.Essence().Pop(stash.Index, opts); err != nil {
		return err
	}

	delete(collection.stashes, branch.ReferenceID.String())

	return nil
}

func (collection *StashCollection) Stashes() map[string]*Stash {
	return collection.stashes
}
