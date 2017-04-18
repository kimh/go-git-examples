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

	fmt.Println(fmt.Sprintf("refs/heads/%s", branch))

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
	})

	//if err.Error() != "already up-to-date" {
	//	CheckIfError(err)
	//}

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Reference("refs/remotes/origin/pr/62", false)

	fmt.Println(ref)

	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)
	err = w.Checkout(ref.Hash())
	CheckIfError(err)

	/// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)
	fmt.Println(commit)
}

//git clone $CIRCLE_REPOSITORY_URL .
//git fetch --force origin 20170310120922:remotes/origin/20170310120922
//git reset --hard $CIRCLE_SHA1
//git checkout -q -B $CIRCLE_BRANCH
//git reset --hard $CIRCLE_SHA1
