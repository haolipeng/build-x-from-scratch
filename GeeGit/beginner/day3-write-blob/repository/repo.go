package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitRepository 初始化一个新的 Git 仓库
func InitRepository(path string) error {
	gitDir := filepath.Join(path, ".git")

	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return fmt.Errorf("failed to create .git directory: %v", err)
	}

	objectsDir := filepath.Join(gitDir, "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create objects directory: %v", err)
	}

	refsHeadsDir := filepath.Join(gitDir, "refs", "heads")
	if err := os.MkdirAll(refsHeadsDir, 0755); err != nil {
		return fmt.Errorf("failed to create refs/heads directory: %v", err)
	}

	headPath := filepath.Join(gitDir, "HEAD")
	headContent := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(headPath, headContent, 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	return nil
}
