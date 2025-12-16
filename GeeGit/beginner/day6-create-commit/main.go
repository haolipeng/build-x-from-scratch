package main

import (
	"fmt"
	"os"
	"time"

	"geegit/beginner/day6-create-commit/blob"
	"geegit/beginner/day6-create-commit/commit"
	"geegit/beginner/day6-create-commit/repository"
	"geegit/beginner/day6-create-commit/tree"
)

func main() {
	fmt.Println("=== Day 6: Create a commit object ===\n")

	testDir := "./test-repo-day6"
	fmt.Printf("✓ Step 1: Initialize test repository\n")
	os.RemoveAll(testDir)
	repository.InitRepository(testDir)

	gitDir := testDir + "/.git"

	fmt.Printf("\n✓ Step 2: Create blob and tree objects\n")
	readmeHash, _ := blob.WriteBlob(gitDir, []byte("# Test Project\n"))
	mainHash, _ := blob.WriteBlob(gitDir, []byte("package main\n\nfunc main() {}\n"))
	
	treeHash, _ := tree.WriteTree(gitDir, []tree.TreeEntry{
		{Mode: "100644", Name: "README.md", Hash: readmeHash},
		{Mode: "100644", Name: "main.go", Hash: mainHash},
	})

	fmt.Printf("\n✓ Step 3: Create commit\n")
	c := &commit.Commit{
		Tree: treeHash,
		Author: commit.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Committer: commit.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Message: "Initial commit\n",
	}

	commitHash, _ := commit.WriteCommit(gitDir, c)
	fmt.Printf("  Commit created: %s\n", commitHash.String()[:8])

	fmt.Printf("\n=== Day 6 Complete! ===\n")
}
