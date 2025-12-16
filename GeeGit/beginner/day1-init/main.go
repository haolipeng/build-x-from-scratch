package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("=== Day 1: Init the .git directory ===")
	fmt.Println()

	// 创建临时目录作为测试仓库
	testDir := "./test-repo"

	// 清理可能存在的旧目录
	os.RemoveAll(testDir)

	fmt.Println("✓ Step 1: Creating test repository directory")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("  Created: %s\n\n", testDir)

	fmt.Println("✓ Step 2: Initializing Git repository structure")
	if err := InitRepository(testDir); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("  Initialized .git directory")
	fmt.Println()

	fmt.Println("✓ Step 3: Verifying directory structure")
	fmt.Println()

	// 显示创建的目录结构
	fmt.Println("Result - Git directory structure:")
	fmt.Println()

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
			fmt.Printf("  📁 %s/\n", dir)
		}
	}

	// 显示 HEAD 文件内容
	headPath := filepath.Join(gitDir, "HEAD")
	if content, err := os.ReadFile(headPath); err == nil {
		fmt.Printf("  📄 .git/HEAD\n")
		fmt.Printf("     Content: %s", content)
	}

	fmt.Println()
	fmt.Println("Verification:")
	fmt.Println()

	// 使用 tree 命令显示目录结构（如果可用）
	if _, err := exec.LookPath("tree"); err == nil {
		fmt.Println("  You can verify the structure with:")
		fmt.Printf("  $ tree %s/.git\n", testDir)
		fmt.Println()

		cmd := exec.Command("tree", "-L", "2", filepath.Join(testDir, ".git"))
		output, _ := cmd.CombinedOutput()
		fmt.Println(string(output))
	} else {
		fmt.Println("  You can verify the structure with:")
		fmt.Printf("  $ ls -la %s/.git\n", testDir)
		fmt.Println()
	}

	fmt.Println("=== Day 1 Complete! ===")
	fmt.Println()
	fmt.Println("What you learned:")
	fmt.Println("  • Git repository structure (.git directory)")
	fmt.Println("  • Where objects are stored (.git/objects)")
	fmt.Println("  • Where branches are stored (.git/refs/heads)")
	fmt.Println("  • What HEAD points to (current branch)")
	fmt.Println()
	fmt.Println("Next: Day 2 - Read a blob object")
}
