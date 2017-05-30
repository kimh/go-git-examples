package main

// Usage: go run git.go "github.com:circleci/circle.git" /tmp/circle ~/.ssh/git-dec

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
)

// Basic example of how to clone a repository using clone options.
func main() {
	CheckArgs("<url>", "<directory>", "<key>", "<branch>")
	url := os.Args[1]
	directory := os.Args[2]
	key := os.Args[3]
	branch := os.Args[4]

	if branch == "" {
		branch = "master"
	}

	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))

	// Clone the given repository to the given directory
	Info("git clone %s %s", url, directory)

	buffer, err := ioutil.ReadFile(key)
	CheckIfError(err)

	signer, err := ssh.ParsePrivateKey(buffer)
	CheckIfError(err)

	var r *git.Repository

	_, err = os.Stat(directory)

	if _, err = os.Stat(directory); err == nil {
		r, err = git.PlainOpen(directory)
	} else {
		r, err = git.PlainClone(directory, false, &git.CloneOptions{
			URL: url,
			// This will automatically checkout to the branch
			ReferenceName: branchRef,
			Auth: &gitssh.PublicKeys{
				User:   "git",
				Signer: signer,
			},
			Progress: os.Stdout,
		})
	}

	CheckIfError(err)

	// Fetch Example
	// ------------------------------
	fmt.Println("doing fetch...")
	err = r.Fetch(&git.FetchOptions{
		Progress: os.Stdout,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/pull/*/head:refs/remotes/origin/pr/*"),
		},
		Auth: &gitssh.PublicKeys{
			User:   "git",
			Signer: signer,
		},
	})
	if err != nil && err.Error() != "already up-to-date" {
		CheckIfError(err)
	}

	// Checkout Example
	// ------------------------------
	fmt.Println("doing checkout...")
	// branch doesn't probably exist in refs/heads. So I need to check remotes refs
	// in this example, we checkout 20170530094317 branch
	ref, err := r.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", "20170530094317")), false)
	fmt.Println(ref)
	w, err := r.Worktree()
	CheckIfError(err)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
	})
	CheckIfError(err)

	// Reset Example
	// ------------------------------
	fmt.Println("doing reset...")
	// we move HEAD to 382aca0d2dbc8320562fe067b5bf9bfd8a7df598 as an example
	commit := plumbing.NewHash("382aca0d2dbc8320562fe067b5bf9bfd8a7df598")
	err = w.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: commit,
	})
	CheckIfError(err)
}
