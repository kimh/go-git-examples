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

	fmt.Println("doing fetching")
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

	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)), false)

	w, err := r.Worktree()
	CheckIfError(err)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: ref.Hash(),
	})
	CheckIfError(err)
}
