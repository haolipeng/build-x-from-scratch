package main

import (
	"fmt"
	"os"

	"geegit/beginner/day4-read-tree/blob"
	"geegit/beginner/day4-read-tree/repository"
	"geegit/beginner/day4-read-tree/tree"
)

func main() {
	fmt.Println("=== Day 4: Read a tree object ===\n")

	testDir := "./test-repo-day4"
	fmt.Printf("✓ Step 1: Initialize test repository\n")
	os.RemoveAll(testDir)
	repository.InitRepository(testDir)

	gitDir := testDir + "/.git"

	fmt.Printf("\n✓ Step 2: Create blob objects\n")
	readmeHash, _ := blob.WriteBlob(gitDir, []byte("# Test Project\n"))
	mainHash, _ := blob.WriteBlob(gitDir, []byte("package main\n\nfunc main() {}\n"))

	fmt.Printf("\n✓ Step 3: Create tree object\n")
	treeContent := tree.BuildTreeContent([]tree.TreeEntry{
		{Mode: "100644", Name: "README.md", Hash: readmeHash},
		{Mode: "100644", Name: "main.go", Hash: mainHash},
	})
	treeHash, _ := tree.WriteRawTree(gitDir, treeContent)

	fmt.Printf("\n✓ Step 4: Read the tree object\n")
	treeObj, _ := tree.ReadTree(gitDir, treeHash)

	fmt.Printf("\nResult:\n")
	fmt.Printf("  Tree Hash: %s\n", treeObj.Hash.String())
	fmt.Printf("  Entries (%d):\n", len(treeObj.Entries))
	for i, entry := range treeObj.Entries {
		fmt.Printf("    [%d] %s %s %s\n", i+1, entry.Mode, entry.Hash.String()[:8], entry.Name)
	}

	fmt.Printf("\n=== Day 4 Complete! ===\n")
}
