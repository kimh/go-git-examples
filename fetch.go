package main

// Usage: go run git.go "github.com:circleci/circle.git" /tmp/circle ~/.ssh/git-dec

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"srcd.works/go-git.v4"
	//"srcd.works/go-git.v4/config"
	. "srcd.works/go-git.v4/examples"
	"srcd.works/go-git.v4/plumbing"
	gitssh "srcd.works/go-git.v4/plumbing/transport/ssh"
)

// Basic example of how to clone a repository using clone options.
func main() {
	CheckArgs("<directory>", "<key>", "<branch>")
	directory := os.Args[1]
	key := os.Args[2]
	branch := os.Args[3]

	Info("git fetch %s", directory)

	buffer, err := ioutil.ReadFile(key)
	CheckIfError(err)

	signer, err := ssh.ParsePrivateKey(buffer)
	CheckIfError(err)

	_, err = os.Stat(directory)

	CheckIfError(err)

	r, err := git.PlainOpen(directory)
	CheckIfError(err)

	//refspec := config.RefSpec("+refs/heads/*:refs/remotes/origin/*")
	fmt.Println("doing fetching")
	err = r.Fetch(&git.FetchOptions{
		//RefSpecs: []config.RefSpec{refspec},
		Auth: &gitssh.PublicKeys{
			User:   "git",
			Signer: signer,
		},
		Progress: os.Stdout,
	})

	CheckIfError(err)

	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
	ref, err := r.Reference(branchRef, false)

	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)
	err = w.Checkout(ref.Hash())
	CheckIfError(err)

	/// ... retrieving the commit object
	commit, err := r.Commit(ref.Hash())
	CheckIfError(err)
	fmt.Println(commit)
}
