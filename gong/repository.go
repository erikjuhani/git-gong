package gong

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/erikjuhani/git-gong/cli"
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

type Repository struct {
	Head    *Head
	Path    string
	GitPath string
	core    *git.Repository
	index   *git.Index
}

func NewRepository(gitRepo *git.Repository, index *git.Index) *Repository {
	return &Repository{
		Head:    &Head{repository: gitRepo},
		Path:    gitRepo.Workdir(),
		GitPath: gitRepo.Path(),
		core:    gitRepo,
		index:   index,
	}
}

func NewTag(gitTag *git.Tag) *Tag {
	return &Tag{
		ID:   gitTag.Id(),
		Name: gitTag.Name(),
		core: gitTag,
	}
}

type Tag struct {
	ID   *git.Oid
	Name string
	core *git.Tag
}

type FileEntry struct{}

type Info struct{}

type Branch struct {
	Name      string
	Shorthand string
	core      *git.Branch
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

	idx, err := gitRepo.Index()
	if err != nil {
		return
	}

	return NewRepository(gitRepo, idx), nil
}

func (repo *Repository) DiffTreeToTree(oldTree *git.Tree, newTree *git.Tree) (*git.Diff, error) {
	return repo.core.DiffTreeToTree(oldTree, newTree, nil)
}

func (repo *Repository) StashAmount() (amount uint) {
	repo.core.Stashes.Foreach(func(index int, message string, id *git.Oid) error {
		amount++
		return nil
	})

	return
}

func (repo *Repository) FindBranch(branchName string, branchType git.BranchType) (*git.Branch, error) {
	return repo.core.LookupBranch(branchName, branchType)
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

	statusList, err := repo.core.StatusList(opts)
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

func (repo *Repository) CreateStash() (stash *git.Oid, err error) {
	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return
	}

	stashName := fmt.Sprintf(stashPattern, currentBranch.Name)

	repo.core.Stashes.Save(signature(), stashName, git.StashIncludeUntracked)

	return
}

func (repo *Repository) PopStash(branchName string) error {
	re := regexp.MustCompile(`@[\w-]+`)

	return repo.core.Stashes.Foreach(func(index int, message string, id *git.Oid) error {
		if branchName == strings.Trim(re.FindString(message), "@") {
			opts, err := git.DefaultStashApplyOptions()
			if err != nil {
				return err
			}

			return repo.core.Stashes.Pop(index, opts)
		}

		return nil
	})
}

func (repo *Repository) Tags() (tags []*Tag, err error) {
	err = repo.core.Tags.Foreach(func(name string, id *git.Oid) error {
		ref, err := repo.core.References.Lookup(name)
		if err != nil {
			return err
		}

		if ref.IsTag() {
			tagObj, err := ref.Peel(git.ObjectTag)
			if err != nil {
				return err
			}

			tag, err := tagObj.AsTag()
			if err != nil {
				return err
			}

			tags = append(tags, NewTag(tag))
		}

		return nil
	})

	return
}

func (repo *Repository) FindTag(tagName string) (tag *Tag, err error) {
	tags, err := repo.Tags()
	if err != nil {
		return
	}

	for _, t := range tags {
		if t.Name == tagName {
			tag = t
			break
		}
	}

	if tag == nil {
		return nil, fmt.Errorf("no tag found by tag name %s", tagName)
	}

	return
}

func (repo *Repository) FindCommit(id *git.Oid) (*Commit, error) {
	gitCommit, err := repo.core.LookupCommit(id)
	if err != nil {
		return nil, err
	}

	return NewCommit(gitCommit), nil
}

func (repo *Repository) CheckoutTree(tree *git.Tree, opts *git.CheckoutOptions) error {
	return repo.core.CheckoutTree(tree, opts)
}

func (repo *Repository) CheckoutTag(tagName string) (tag *Tag, err error) {
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	tags, err := repo.Tags()
	if err != nil {
		return
	}

	for _, t := range tags {
		if t.Name == tagName {
			tag = t
			break
		}
	}

	if tag == nil {
		return nil, fmt.Errorf("no tag found by tag name %s", tagName)
	}

	commit, err := repo.FindCommit(tag.core.TargetId())
	if err != nil {
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		return
	}

	defer tree.Free()

	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		return
	}

	err = repo.Head.Detach(commit.ID)
	return
}

func (repo *Repository) CheckoutCommit(hash string) (commit *Commit, err error) {
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	commits, err := repo.Commits()
	if err != nil {
		return
	}

	for _, commit = range commits {
		if commit.ID.String() == hash {
			break
		}
	}

	if commit == nil {
		return nil, fmt.Errorf("no commit found by hash %s", hash)
	}

	tree, err := commit.Tree()
	if err != nil {
		return
	}

	defer tree.Free()

	err = repo.core.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		return
	}

	err = repo.Head.Detach(commit.ID)

	return
}

func (repo *Repository) CheckoutBranch(branchName string) (branch *git.Branch, err error) {
	detached, err := repo.core.IsHeadDetached()
	if err != nil {
		return
	}

	if detached {
		ref, err := repo.core.References.Lookup(fmt.Sprintf("%s%s", headRef, branchName))
		if err != nil {
			return nil, err
		}

		repo.core.SetHead(ref.Name())

		checkoutOpts := &git.CheckoutOpts{
			Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
		}

		repo.core.CheckoutHead(checkoutOpts)
	}

	branch, err = repo.core.LookupBranch(branchName, git.BranchLocal)

	// Branch does not exist, create it first
	if branch == nil || err != nil {
		branch, err = repo.CreateLocalBranch(branchName)
	}

	defer branch.Free()

	_, err = repo.CreateStash()
	if err != nil {
		return
	}

	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}

	localCommit, err := repo.core.LookupCommit(branch.Target())
	if err != nil {
		return
	}

	defer localCommit.Free()

	tree, err := repo.core.LookupTree(localCommit.TreeId())
	if err != nil {
		return
	}

	defer tree.Free()

	err = repo.core.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		return
	}

	err = repo.core.SetHead(fmt.Sprintf("%s%s", headRef, branchName))
	if err != nil {
		return
	}

	err = repo.PopStash(branchName)
	return
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

	idx, err := gitRepo.Index()
	if err != nil {
		return
	}

	return NewRepository(gitRepo, idx), nil
}

func (repo *Repository) Changed() (err error) {
	diff, err := repo.core.DiffIndexToWorkdir(
		repo.index,
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

	status, err := repo.core.StatusList(&git.StatusOptions{})
	if err != nil {
		return
	}

	entries, err := status.EntryCount()
	if err != nil {
		return
	}

	defer status.Free()

	if changedFiles == 0 && entries == 0 {
		err = errors.New("no files changed, nothing to commit, working tree clean")
		return
	}

	return
}

func (repo *Repository) AddToIndex(pathspec []string) (treeID *git.Oid, err error) {
	if err = repo.Changed(); err != nil {
		return
	}

	idx := repo.index

	if err = idx.AddAll(pathspec, git.IndexAddDefault, nil); err != nil {
		return
	}

	treeID, err = idx.WriteTree()
	if err != nil {
		return
	}

	if err = idx.Write(); err != nil {
		return
	}

	return
}

func (repo *Repository) createCommit(treeID *git.Oid, commit *Commit, msg string) (id *git.Oid, err error) {
	tree, err := repo.core.LookupTree(treeID)
	if err != nil {
		return
	}

	sig := signature()

	if checkEmptyString(msg) {
		input, cliErr := cli.CaptureInput()
		if cliErr != nil {
			return nil, cliErr
		}

		msg = string(input)
	}

	if checkEmptyString(msg) {
		err = errors.New("aborting due to empty commit message")
		return
	}

	if commit != nil {
		id, err = repo.core.CreateCommit("HEAD", sig, sig, string(msg), tree, commit.core)
		if err != nil {
			return
		}

		err = repo.core.CheckoutHead(&git.CheckoutOpts{
			Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
		})
		return
	}

	// Initial commit
	id, err = repo.core.CreateCommit("HEAD", sig, sig, string(msg), tree)
	if err != nil {
		return
	}

	err = repo.core.CheckoutHead(&git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
	})

	return
}

func (repo *Repository) Commit(treeID *git.Oid, msg string) (commitID *git.Oid, err error) {
	unborn, err := repo.core.IsHeadUnborn()
	if err != nil {
		return
	}

	if unborn {
		commitID, err = repo.createCommit(treeID, nil, msg)
		if err != nil {
			return
		}

		return
	}

	currentTip, err := repo.Head.Commit()
	if err != nil {
		return
	}

	return repo.createCommit(treeID, currentTip, msg)
}

func (repo *Repository) References() ([]string, error) {
	var list []string
	iter, err := repo.core.NewReferenceIterator()
	if err != nil {
		return list, err
	}

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

	commits = append(commits, headCommit)

	if headCommit.core.ParentCount() != 0 {
		parent := headCommit.core.Parent(0)
		commits = append(commits, NewCommit(parent))

		for parent.ParentCount() != 0 {
			parent = parent.Parent(0)
			commits = append(commits, NewCommit(parent))
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

	gitTag, err := repo.core.Tags.Create(tagname, headCommit.core, signature(), message)
	if err != nil {
		return
	}

	return &Tag{ID: gitTag, Name: tagname}, nil
}

// CreateLocalBranch creates a local branch to repository.
func (repo *Repository) CreateLocalBranch(branchName string) (branch *git.Branch, err error) {
	// Check if branch already exists
	localBranch, err := repo.core.LookupBranch(branchName, git.BranchLocal)
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

	return repo.core.CreateBranch(branchName, headCommit.core, false)
}

// TODO get signature from git configuration
func signature() *git.Signature {
	return &git.Signature{
		Name:  "gong tester",
		Email: "gong@tester.com",
		When:  time.Now(),
	}
}
