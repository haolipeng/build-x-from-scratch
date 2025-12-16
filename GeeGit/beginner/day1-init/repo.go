package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitRepository 初始化一个新的 Git 仓库
// 创建 .git 目录和必要的子目录结构
func InitRepository(path string) error {
	gitDir := filepath.Join(path, ".git")

	// 创建 .git 目录
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return fmt.Errorf("failed to create .git directory: %v", err)
	}

	// 创建 objects 目录（存储所有 Git 对象）
	objectsDir := filepath.Join(gitDir, "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create objects directory: %v", err)
	}

	// 创建 refs/heads 目录（存储分支引用）
	refsHeadsDir := filepath.Join(gitDir, "refs", "heads")
	if err := os.MkdirAll(refsHeadsDir, 0755); err != nil {
		return fmt.Errorf("failed to create refs/heads directory: %v", err)
	}

	// 创建 HEAD 文件，指向默认分支 main
	headPath := filepath.Join(gitDir, "HEAD")
	headContent := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(headPath, headContent, 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	return nil
}
