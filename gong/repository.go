package gong

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v31"
)

const (
	stashPattern = "@%s"
	headRef      = "refs/heads/"
)

// TODO: Move this somewhere more appropriate
func checkEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

type freer interface {
	Free()
}

func Free(f freer) {
	f.Free()
}

type Repository struct {
	Head    *Head
	Path    string
	GitPath string
	Index   *git.Index
	Tree    *git.Tree
	essence *git.Repository
	Stashes *StashCollection
}

func (repo *Repository) Free() {
	repo.Essence().Free()
}

func (repo *Repository) Essence() *git.Repository {
	return repo.essence
}

func (repo *Repository) FindTree(treeID *git.Oid) (*git.Tree, error) {
	return repo.Essence().LookupTree(treeID)
}

func NewRepository(gitRepo *git.Repository, index *git.Index) *Repository {
	return &Repository{
		Head:    NewHead(gitRepo),
		Path:    gitRepo.Workdir(),
		GitPath: gitRepo.Path(),
		Index:   index,
		Stashes: NewStashCollection(&gitRepo.Stashes),
		essence: gitRepo,
	}
}

type FileEntry struct{}

type Info struct{}

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

func (collection *StashCollection) Create(head *Head) (*Stash, error) {
	currentBranch, err := head.Branch()
	if err != nil {
		return nil, err
	}

	stashName := fmt.Sprintf(stashPattern, currentBranch.Name)

	stashID, err := collection.essence.Save(signature(), stashName, git.StashIncludeUntracked)

	stash := &Stash{ID: stashID, Message: stashName, Index: 0}
	collection.stashes[stashName] = stash

	return stash, nil
}

type Stash struct {
	ID      *git.Oid
	Message string
	Index   int
}

func (collection *StashCollection) Find(name string) (*Stash, error) {
	if stash, ok := collection.stashes[name]; ok {
		return stash, nil
	}

	return nil, fmt.Errorf("stash with name %s was not found", name)
}

func (collection *StashCollection) Pop(name string) error {
	stash, err := collection.Find(name)
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

	delete(collection.stashes, name)

	return nil
}

func (collection *StashCollection) Stashes() map[string]*Stash {
	return collection.stashes
}

var (
	DefaultReference = "main"
)

// Init initializes the repository.
// TODO: use the git config to set the initial default branch.
// more info: https://github.blog/2020-07-27-highlights-from-git-2-28/#introducing-init-defaultbranch
func Init(path string, bare bool, initialReference string) (repo *Repository, err error) {
	gitRepo, err := git.InitRepository(path, bare)
	if err != nil {
		return
	}

	if checkEmptyString(initialReference) {
		initialReference = DefaultReference
	}

	initRef := fmt.Sprintf("%s%s", headRef, initialReference)
	err = ioutil.WriteFile(fmt.Sprintf("%s/HEAD", gitRepo.Path()), []byte("ref: "+initRef), 0644)
	if err != nil {
		return
	}

	index, err := gitRepo.Index()
	if err != nil {
		return
	}

	defer Free(index)

	return NewRepository(gitRepo, index), nil
}

func (repo *Repository) DiffTreeToTree(oldTree *git.Tree, newTree *git.Tree) (*git.Diff, error) {
	return repo.Essence().DiffTreeToTree(oldTree, newTree, nil)
}

func (repo *Repository) FindBranch(branchName string, branchType git.BranchType) (*Branch, error) {
	gitBranch, err := repo.Essence().LookupBranch(branchName, branchType)
	if err != nil {
		return nil, err
	}

	return NewBranch(branchName, gitBranch), nil
}

func (repo *Repository) Info() (output string, err error) {
	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return
	}

	currentTip, err := repo.Head.Commit()
	if err != nil {
		return
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Branch %s\n", currentBranch.Name))
	sb.WriteString(fmt.Sprintf("Commit %s\n\n", currentTip.ID.String()))

	entries, err := repo.StatusEntries()
	if err != nil {
		return
	}

	if len(entries) == 0 {
		sb.WriteString("No changes.")
		return strings.TrimSuffix(sb.String(), "\n"), nil
	}

	sb.WriteString(fmt.Sprintf("Changes (%d):\n", len(entries)))
	for _, e := range entries {
		switch e.IndexToWorkdir.Status {
		case git.DeltaAdded:
			sb.WriteString(" A ")
		case git.DeltaModified:
			sb.WriteString(" M ")
		case git.DeltaRenamed:
			sb.WriteString(" R ")
		case git.DeltaDeleted:
			sb.WriteString(" D ")
		case git.DeltaUntracked:
			sb.WriteString("?? ")
		}

		sb.WriteString(fmt.Sprintf("%s\n", e.IndexToWorkdir.NewFile.Path))
	}

	return strings.TrimSuffix(sb.String(), "\n"), nil
}

func (repo *Repository) StatusEntries() (entries []git.StatusEntry, err error) {
	opts := &git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked | git.StatusOptRenamesHeadToIndex | git.StatusOptSortCaseSensitively,
	}

	statusList, err := repo.Essence().StatusList(opts)
	if err != nil {
		return
	}

	amount, err := statusList.EntryCount()
	if err != nil {
		return
	}

	for i := 0; i < amount; i++ {
		entry, err := statusList.ByIndex(i)
		if err != nil {
			return entries, err
		}

		entries = append(entries, entry)
	}

	return
}

func (repo *Repository) CurrentBranch() (*Branch, error) {
	return repo.Head.Branch()
}

func (repo *Repository) Tags() ([]*Tag, error) {
	var tags []*Tag
	err := repo.Essence().Tags.Foreach(func(name string, id *git.Oid) error {
		ref, err := repo.Essence().References.Lookup(name)
		if err != nil {
			return err
		}
		defer Free(ref)

		if ref.IsTag() {
			tagObj, err := ref.Peel(git.ObjectTag)
			if err != nil {
				return err
			}

			tag, err := tagObj.AsTag()
			if err != nil {
				return err
			}
			defer Free(tag)

			tags = append(tags, NewTag(tag))
		}

		return nil
	})

	if err != nil {
		return tags, err
	}

	return tags, nil
}

func (repo *Repository) FindTag(tagName string) (*Tag, error) {
	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var tag *Tag

	for _, tag = range tags {
		if tag.Name == tagName {
			break
		}
	}

	if tag == nil {
		return nil, fmt.Errorf("no tag found by tag name %s", tagName)
	}
	defer Free(tag)

	return tag, nil
}

func (repo *Repository) FindCommit(commitID *git.Oid) (*Commit, error) {
	gitCommit, err := repo.Essence().LookupCommit(commitID)
	if err != nil {
		return nil, err
	}

	return NewCommit(gitCommit), nil
}

func (repo *Repository) CheckoutTree(tree *git.Tree, opts *git.CheckoutOptions) error {
	return repo.Essence().CheckoutTree(tree, opts)
}

func (repo *Repository) CheckoutTag(tagName string) (*Tag, error) {
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var tag *Tag

	for _, tag = range tags {
		if tag.Name == tagName {
			break
		}
	}

	if tag == nil {
		return nil, fmt.Errorf("no tag found by tag name %s", tagName)
	}
	defer Free(tag)

	commit, err := repo.FindCommit(tag.Essence().TargetId())
	if err != nil {
		return nil, err
	}
	defer Free(commit)

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	defer Free(tree)

	if err := repo.CheckoutTree(tree, checkoutOpts); err != nil {
		return nil, err
	}

	if err := repo.Head.Detach(commit.ID); err != nil {
		return nil, err
	}

	return tag, nil
}

func (repo *Repository) CheckoutCommit(hash string) (*Commit, error) {
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	commits, err := repo.Commits()
	if err != nil {
		return nil, err
	}

	var commit *Commit
	for _, commit = range commits {
		if commit.ID.String() == hash {
			break
		}
	}

	if commit == nil {
		return nil, fmt.Errorf("no commit found by hash %s", hash)
	}
	defer Free(commit)

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	defer Free(tree)

	if err := repo.Essence().CheckoutTree(tree, checkoutOpts); err != nil {
		return nil, err
	}

	if err := repo.Head.Detach(commit.ID); err != nil {
		return nil, err
	}

	return commit, err
}

func (repo *Repository) CheckoutBranch(branchName string) (*Branch, error) {
	detached, err := repo.Head.IsDetached()
	if err != nil {
		return nil, err
	}

	if detached {
		ref, err := repo.Essence().References.Lookup(fmt.Sprintf("%s%s", headRef, branchName))
		if err != nil {
			return nil, err
		}
		defer Free(ref)

		err = repo.Head.SetReference(ref.Name())
		if err != nil {
			return nil, err
		}

		err = repo.Head.Checkout()
		if err != nil {
			return nil, err
		}
	}

	branch, err := repo.FindBranch(branchName, git.BranchLocal)
	// Branch does not exist, create it first
	if branch == nil || err != nil {
		branch, err = repo.CreateLocalBranch(branchName)
		if err != nil {
			return nil, err
		}
	}
	defer Free(branch)

	_, err = repo.Stashes.Create(repo.Head)
	if err != nil {
		return nil, err
	}

	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	localCommit, err := repo.FindCommit(branch.ReferenceID)
	if err != nil {
		return nil, err
	}
	defer Free(localCommit)

	tree, err := localCommit.Tree()
	if err != nil {
		return nil, err
	}
	defer Free(tree)

	if err := repo.Essence().CheckoutTree(tree, checkoutOpts); err != nil {
		return nil, err
	}

	if err := repo.Head.SetReference(headRef + branchName); err != nil {
		return nil, err
	}

	if err := repo.Stashes.Pop(branchName); err != nil {
		return nil, err
	}

	return branch, nil
}

func Open() (repo *Repository, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	gitRepo, err := git.OpenRepository(wd)
	if err != nil {
		return
	}

	index, err := gitRepo.Index()
	if err != nil {
		return
	}
	defer Free(index)

	return NewRepository(gitRepo, index), nil
}

func (repo *Repository) Changed() (err error) {
	diff, err := repo.Essence().DiffIndexToWorkdir(
		repo.Index,
		&git.DiffOptions{Flags: git.DiffIncludeUntracked},
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

	status, err := repo.Essence().StatusList(&git.StatusOptions{})
	if err != nil {
		return
	}

	entries, err := status.EntryCount()
	if err != nil {
		return
	}
	defer Free(status)

	if changedFiles == 0 && entries == 0 {
		err = errors.New("no files changed, nothing to commit, working tree clean")
		return
	}

	return
}

func (repo *Repository) AddToIndex(pathspec []string) (*git.Tree, error) {
	if err := repo.Changed(); err != nil {
		return nil, err
	}

	if err := repo.Index.AddAll(pathspec, git.IndexAddDefault, nil); err != nil {
		return nil, err
	}

	treeID, err := repo.Index.WriteTree()
	if err != nil {
		return nil, err
	}

	if err = repo.Index.Write(); err != nil {
		return nil, err
	}

	return repo.FindTree(treeID)
}

func (repo *Repository) CreateCommit(tree *git.Tree, message string) (*Commit, error) {
	// Implement later.
	/*
		if checkEmptyString(msg) {
			input, cliErr := cli.CaptureInput()
			if cliErr != nil {
				return nil, cliErr
			}

			msg = string(input)
		}
	*/

	if checkEmptyString(message) {
		return nil, ErrEmptyCommitMsg
	}

	head := repo.Head

	exists, err := head.Exists()
	if err != nil {
		return nil, err
	}

	var commitID *git.Oid

	if exists {
		headCommit, err := repo.Head.Commit()
		defer Free(headCommit)

		if err != nil {
			return nil, err
		}

		commitID, err = repo.Essence().CreateCommit(head.RefName, signature(), signature(), message, tree, headCommit.Essence())
		if err != nil {
			return nil, err
		}
	} else {
		// Initial commit.
		commitID, err = repo.Essence().CreateCommit(head.RefName, signature(), signature(), message, tree)
		if err != nil {
			return nil, err
		}
	}

	if err := head.Checkout(); err != nil {
		return nil, err
	}

	return repo.FindCommit(commitID)
}

func (repo *Repository) References() ([]string, error) {
	var list []string
	iter, err := repo.Essence().NewReferenceIterator()
	if err != nil {
		return list, err
	}
	defer Free(iter)

	nameIter := iter.Names()
	name, err := nameIter.Next()
	for err == nil {
		list = append(list, name)
		name, err = nameIter.Next()
	}

	return list, err
}

func (repo *Repository) Commits() (commits []*Commit, err error) {
	headCommit, err := repo.Head.Commit()
	if err != nil {
		return
	}
	defer Free(headCommit)

	commits = append(commits, headCommit)

	if !headCommit.IsRoot() {
		parent := headCommit.Parent()
		commits = append(commits, parent)

		for !parent.IsRoot() {
			parent = parent.Parent()
			commits = append(commits, parent)
		}
	}
	return
}

// CreateTag creates a git tag.
func (repo *Repository) CreateTag(tagname string, message string) (tag *Tag, err error) {
	headCommit, err := repo.Head.Commit()
	if err != nil {
		return
	}
	defer Free(headCommit)

	gitTag, err := repo.Essence().Tags.Create(tagname, headCommit.Essence(), signature(), message)
	if err != nil {
		return
	}

	return &Tag{ID: gitTag, Name: tagname}, nil
}

// CreateLocalBranch creates a local branch to repository.
func (repo *Repository) CreateLocalBranch(branchName string) (branch *Branch, err error) {
	// Check if branch already exists
	localBranch, err := repo.FindBranch(branchName, git.BranchLocal)
	if localBranch != nil && err != nil {
		return
	}

	// Branch already exists return existing branch and an error stating branch already exists.
	if localBranch != nil {
		return localBranch, fmt.Errorf("branch %s already exists", branchName)
	}

	headCommit, err := repo.Head.Commit()
	if err != nil {
		return
	}
	defer Free(headCommit)

	return repo.CreateBranch(branchName, headCommit, false)
}

func (repo *Repository) CreateBranch(branchName string, commit *Commit, force bool) (*Branch, error) {
	gitBranch, err := repo.Essence().CreateBranch(branchName, commit.Essence(), force)
	if err != nil {
		return nil, err
	}

	return NewBranch(branchName, gitBranch), nil
}

// TODO get signature from git configuration
func signature() *git.Signature {
	return &git.Signature{
		Name:  "gong tester",
		Email: "gong@tester.com",
		When:  time.Now(),
	}
}
