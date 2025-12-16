package main

import (
	"fmt"
	"os"

	"geegit/beginner/day3-write-blob/blob"
	"geegit/beginner/day3-write-blob/repository"
)

func main() {
	fmt.Println("=== Day 3: Create a blob object ===\n")

	testDir := "./test-repo-day3"
	fmt.Printf("✓ Step 1: Initialize test repository at %s\n", testDir)
	os.RemoveAll(testDir)
	repository.InitRepository(testDir)

	gitDir := testDir + "/.git"

	content := []byte("Hello, GeeGit!\nThis is my first blob object.\n")
	fmt.Printf("\n✓ Step 2: Create blob with content\n")

	hash, _ := blob.WriteBlob(gitDir, content)
	fmt.Printf("\n✓ Step 3: Blob object written successfully\n")
	fmt.Printf("\nResult:\n")
	fmt.Printf("  Hash: %s\n", hash.String())

	blobObj, _ := blob.ReadBlob(gitDir, hash)
	if string(blobObj.Data) == string(content) {
		fmt.Printf("  ✓ Content matches! Verification successful.\n")
	}

	fmt.Printf("\n=== Day 3 Complete! ===\n")
}
