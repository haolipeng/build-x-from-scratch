package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=== Day 1: Init the .git directory ===")

	// 创建临时目录作为测试仓库
	testDir := "./test-repo"

	// 清理可能存在的旧目录
	os.RemoveAll(testDir)

	fmt.Println("\nStep 1: Creating test repository directory")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Created: %s\n", testDir)

	fmt.Println("\nStep 2: Initializing Git repository structure")
	if err := InitRepository(testDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Initialized .git directory")

	// 显示创建的目录结构
	fmt.Println("\nGit directory structure:")

	gitDir := filepath.Join(testDir, ".git")

	// 检查并显示每个目录/文件
	dirs := []string{
		".git",
		".git/objects",
		".git/refs",
		".git/refs/heads",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(testDir, dir)
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("  %s/\n", dir)
		}
	}

	// 显示 HEAD 文件内容
	headPath := filepath.Join(gitDir, "HEAD")
	if content, err := os.ReadFile(headPath); err == nil {
		fmt.Printf("  .git/HEAD\n")
		fmt.Printf("  Content: %s", content)
	}

	fmt.Printf("\n可以通过以下命令验证: tree %s/.git 或 ls -R %s/.git\n", testDir, testDir)

	fmt.Println("\n=== Day 1 Complete! ===")
	fmt.Println("Next: Day 2 - Read a blob object")
}
