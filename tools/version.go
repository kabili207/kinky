package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5"

	//nolint:stylecheck // looks better with dot imports
	. "github.com/dave/jennifer/jen"
)

const GitHashLength = 6

var errMissingArgs = errors.New("missing arguments")

func parseArgs() (string, string) {
	args := os.Args[1:]
	if args[0] == "--" {
		args = args[1:]
	}

	if len(args) != 2 {
		panic(errMissingArgs)
	}

	file := path.Clean(args[0])

	return file, args[1]
}

func parseGit() (git.Status, *plumbing.Reference, *object.Commit) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		panic(fmt.Errorf("failed to open git repository: %w", err))
	}
	repoHead, err := repo.Head()
	if err != nil {
		panic(fmt.Errorf("failed to get git head: %w", err))
	}
	repoCommit, err := repo.CommitObject(repoHead.Hash())
	if err != nil {
		panic(fmt.Errorf("failed to get git commit: %w", err))
	}
	worktree, err := repo.Worktree()
	if err != nil {
		panic(fmt.Errorf("failed to get git worktree: %w", err))
	}
	worktreeStatus, err := worktree.Status()
	if err != nil {
		panic(fmt.Errorf("failed to get git worktree status: %w", err))
	}

	return worktreeStatus, repoHead, repoCommit
}

func main() {
	file, pkg := parseArgs()

	fmt.Printf("creating parent directory: %v\n", path.Dir(file))
	if err := os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create parent folder: %w", err))
	}

	now := time.Now()
	buffer, err := os.ReadFile("VERSION")
	if err != nil {
		panic(fmt.Errorf("failed to read version file: %w", err))
	}
	version := strings.TrimSpace(string(buffer))
	worktreeStatus, repoHead, repoCommit := parseGit()

	fmt.Println("generating file")
	f := NewFilePath(pkg)

	f.HeaderComment("Code generated by tools/version.go - DO NOT EDIT.")

	f.Const().Id("BuildTimestamp").Op("=").Lit(now.Unix())
	f.Const().Id("BuildDate").Op("=").Lit(now.Format(time.RFC3339))
	f.Const().Id("Version").Op("=").Lit(version)

	f.Const().Id("GitDirty").Op("=").Lit(!worktreeStatus.IsClean())
	f.Const().Id("GitSha").Op("=").Lit(repoHead.Hash().String())
	f.Const().Id("GitCommitTimestamp").Op("=").Lit(repoCommit.Author.When.Unix())
	f.Const().Id("GitCommitDate").Op("=").Lit(repoCommit.Author.When.Format(time.RFC3339))
	f.Const().Id("GitBranch").Op("=").Lit(repoHead.Name().Short())

	f.Func().Id("BuildTime").Params().Qual("time", "Time").Block(
		Return(Qual("time", "Unix").Call(Id("BuildTimestamp"), Lit(0))),
	)

	f.Func().Id("GitCommitTime").Params().Qual("time", "Time").Block(
		Return(Qual("time", "Unix").Call(Id("GitCommitTimestamp"), Lit(0))),
	)

	f.Func().Id("GitShaShort").Params().String().Block(
		Return(Id("GitSha").Index(Empty(), Lit(GitHashLength))),
	)

	fmt.Printf("saving file %v\n", file)
	if err := f.Save(file); err != nil {
		panic(fmt.Errorf("failed to save file: %w", err))
	}
}
