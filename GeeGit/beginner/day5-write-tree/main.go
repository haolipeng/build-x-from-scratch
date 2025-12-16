package main

import (
	"fmt"
	"os"

	"geegit/beginner/day5-write-tree/blob"
	"geegit/beginner/day5-write-tree/repository"
	"geegit/beginner/day5-write-tree/tree"
)

func main() {
	fmt.Println("=== Day 5: Write a tree object ===\n")

	testDir := "./test-repo-day5"
	fmt.Printf("✓ Step 1: Initialize test repository\n")
	os.RemoveAll(testDir)
	repository.InitRepository(testDir)

	gitDir := testDir + "/.git"

	fmt.Printf("\n✓ Step 2: Create blob objects\n")
	readmeHash, _ := blob.WriteBlob(gitDir, []byte("# My Project\n"))
	mainHash, _ := blob.WriteBlob(gitDir, []byte("package main\n\nfunc main() {}\n"))

	fmt.Printf("\n✓ Step 3: Create tree object\n")
	entries := []tree.TreeEntry{
		{Mode: "100644", Name: "main.go", Hash: mainHash},
		{Mode: "100644", Name: "README.md", Hash: readmeHash},
	}

	treeHash, _ := tree.WriteTree(gitDir, entries)
	fmt.Printf("  Tree created: %s\n", treeHash.String()[:8])

	fmt.Printf("\n=== Day 5 Complete! ===\n")
}
