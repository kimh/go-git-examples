package main

// Usage: go run git.go "github.com:circleci/circle.git" /tmp/circle ~/.ssh/git-dec

import (
	"fmt"
	"os"

	"srcd.works/go-git.v4"
	//"srcd.works/go-git.v4/config"
	. "srcd.works/go-git.v4/examples"
	"srcd.works/go-git.v4/plumbing"
)

// Basic example of how to clone a repository using clone options.
func main() {
	CheckArgs("<directory>", "<branch>")
	directory := os.Args[1]
	branch := os.Args[2]

	_, err := os.Stat(directory)
	CheckIfError(err)

	r, err := git.PlainOpen(directory)
	CheckIfError(err)

	//refspec := config.RefSpec("+refs/heads/*:refs/remotes/origin/*")
	branchRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
	ref, err := r.Reference(branchRef, false)
	head := plumbing.NewHashReference(plumbing.HEAD, ref.Hash())

	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	fmt.Println(r.Head())
	err = w.Checkout(head.Hash())
	CheckIfError(err)
	fmt.Println(r.Head())

	/// ... retrieving the commit object
	//commit, err := r.Commit(ref.Hash())
	//CheckIfError(err)
	//fmt.Println(commit)
}
