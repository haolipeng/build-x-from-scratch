package main

import (
	"fmt"
	"os"

	"geegit/beginner/day2-read-blob/blob"
	"geegit/beginner/day2-read-blob/hash"
)

func main() {
	fmt.Println("=== Day 2: Read a blob object ===")
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run *.go <git-directory> <object-hash>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run *.go /path/to/repo/.git e69de29bb2d1d6434b8b29ae775ad8c2e48c5391")
		fmt.Println()
		return
	}

	gitDir := os.Args[1]
	hashStr := os.Args[2]

	fmt.Println("✓ Step 1: Parsing hash")
	h, err := hash.NewHash(hashStr)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("  Hash: %s\n\n", h.String())

	fmt.Println("✓ Step 2: Reading object from .git/objects/")
	blobObj, err := blob.ReadBlob(gitDir, h)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("  Found blob object\n")
	fmt.Printf("  Size: %d bytes\n\n", len(blobObj.Data))

	fmt.Println("✓ Step 3: Displaying content")
	fmt.Println()
	fmt.Println("Result:")
	fmt.Println("─────────────────────────────")
	fmt.Print(string(blobObj.Data))
	if len(blobObj.Data) > 0 && blobObj.Data[len(blobObj.Data)-1] != '\n' {
		fmt.Println()
	}
	fmt.Println("─────────────────────────────")
	fmt.Println()

	fmt.Println("=== Day 2 Complete! ===")
}
