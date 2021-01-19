package gong

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v31"
)

// Repository represents is an abstraction of the underlying *git.Repository.
// To access the essence (*git.Repository) call Essence() function.
type Repository struct {
	Head    *Head
	Path    string
	GitPath string
	Index   *git.Index
	Tree    *git.Tree
	essence *git.Repository
	Stashes *StashCollection
}

// Free frees git repository pointer.
func (repo *Repository) Free() {
	repo.Essence().Free()
}

// Essence returns *git.Repository.
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

// Init initializes the repository.
// TODO: use the git config to set the initial default branch.
// more info: https://github.blog/2020-07-27-highlights-from-git-2-28/#introducing-init-defaultbranch
func Init(path string, bare bool, initialReference string) (*Repository, error) {
	gitRepo, err := git.InitRepository(path, bare)
	if err != nil {
		return nil, err
	}

	if checkEmptyString(initialReference) {
		initialReference = DefaultReference
	}

	initRef := fmt.Sprintf("%s%s", headRef, initialReference)
	if err := ioutil.WriteFile(fmt.Sprintf("%s/HEAD", gitRepo.Path()), []byte("ref: "+initRef), 0644); err != nil {
		return nil, err
	}

	index, err := gitRepo.Index()
	if err != nil {
		return nil, err
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

func (repo *Repository) Info() (string, error) {
	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return "", err
	}

	currentTip, err := repo.Head.Commit()
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Branch %s\n", currentBranch.Name))
	sb.WriteString(fmt.Sprintf("Commit %s\n\n", currentTip.ID.String()))

	entries, err := repo.StatusEntries()
	if err != nil {
		return "", err
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

func (repo *Repository) StatusEntries() ([]git.StatusEntry, error) {
	opts := &git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked | git.StatusOptRenamesHeadToIndex | git.StatusOptSortCaseSensitively,
	}

	var entries []git.StatusEntry

	statusList, err := repo.Essence().StatusList(opts)
	if err != nil {
		return entries, err
	}
	defer Free(statusList)

	entryCount, err := statusList.EntryCount()
	if err != nil {
		return entries, err
	}

	for i := 0; i < entryCount; i++ {
		entry, err := statusList.ByIndex(i)
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
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

		if err := repo.Head.SetReference(ref.Name()); err != nil {
			return nil, err
		}

		if err := repo.Head.Checkout(); err != nil {
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

	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return nil, err
	}

	_, err = repo.Stashes.Create(currentBranch)
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

	// No existing stash.
	if !repo.Stashes.Has(branch) {
		return branch, nil
	}

	if err := repo.Stashes.Pop(branch); err != nil {
		return nil, err
	}

	return branch, nil
}

// Clone clones a git repository from source location to a target location.
// If target location is an empty string clone to a directory named after source.
func Clone(source string, target string) (*Repository, error) {
	opts := git.CloneOptions{}

	// Check that the source is a valid url.
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	src := strings.TrimSuffix(u.String(), ".git")

	gitRepo, err := git.Clone(src, target, &opts)
	if err != nil {
		return nil, err
	}
	defer Free(gitRepo)

	return NewRepository(gitRepo, nil), nil

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

func (repo *Repository) Changed() error {
	diff, err := repo.Essence().DiffIndexToWorkdir(
		repo.Index,
		&git.DiffOptions{Flags: git.DiffIncludeUntracked},
	)
	if err != nil {
		return err
	}
	defer diff.Free()

	stats, err := diff.Stats()
	if err != nil {
		return err
	}
	defer stats.Free()

	changeCount := stats.FilesChanged()

	status, err := repo.Essence().StatusList(&git.StatusOptions{})
	if err != nil {
		return err
	}
	defer Free(status)

	entryCount, err := status.EntryCount()
	if err != nil {
		return err
	}

	if changeCount == 0 && entryCount == 0 {
		return fmt.Errorf("no files changed, %w", ErrNothingToCommit)
	}

	return nil
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
	iter, err := repo.Essence().NewReferenceIterator()
	if err != nil {
		return nil, err
	}
	defer Free(iter)

	var list []string

	nameIter := iter.Names()
	name, err := nameIter.Next()
	for err == nil {
		list = append(list, name)
		name, err = nameIter.Next()
	}

	return list, err
}

func (repo *Repository) Commits() ([]*Commit, error) {
	currentTip, err := repo.Head.Commit()
	if err != nil {
		return nil, err
	}
	defer Free(currentTip)

	commits := []*Commit{currentTip}

	parent := currentTip
	for parent.HasChildren() {
		parent = parent.Parent()
		commits = append(commits, parent)
	}

	return commits, nil
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

	return repo.createBranch(branchName, headCommit, false)
}

func (repo *Repository) createBranch(branchName string, commit *Commit, force bool) (*Branch, error) {
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
