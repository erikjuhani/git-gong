package gong

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/erikjuhani/git-gong/config"
)

func commandExists(command string) bool {
	if _, err := exec.LookPath(command); err == nil {
		return true
	}

	return false
}

func parseArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}

	repo, err := Open()
	if err != nil {
		return err
	}
	defer Free(repo)

	branch, err := repo.Head.Branch()
	if err != nil {
		return err
	}

	gitCmd := args[0]

	switch gitCmd {
	case "add":
		fallthrough
	case "commit":
		if config.ProtectedBranchPatterns.Match(branch.Name) {
			return errors.New("trying to commit on a protected branch, operation aborted")
		}
		return nil
	case "branch":
		if args[1][0:1] == "-" {
			return nil
		}

		if !config.AllowedBranchPatterns.Match(args[1]) {
			return errors.New("error branch name did not match allowed template patterns")
		}

		return nil
	case "checkout":
		if strings.ToLower(args[1]) != "-b" || len(args) < 3 {
			return nil
		}

		if !config.AllowedBranchPatterns.Match(args[2]) {
			return errors.New("error branch name did not match allowed template patterns")
		}
		return nil
	}

	return nil
}

func RunGitCommand(args []string) error {
	if !commandExists("git") {
		return errors.New("git executable not found in path")
	}

	if err := parseArgs(args); err != nil {
		return err
	}

	git := exec.Command("git", args...)
	git.Stdin = os.Stdin
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr

	git.Run()

	return nil
}
