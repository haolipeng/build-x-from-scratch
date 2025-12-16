package commit

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"

	"geegit/beginner/day6-create-commit/hash"
)

// WriteCommit 将 commit 对象写入 .git/objects 目录
func WriteCommit(gitDir string, commit *Commit) (hash.Hash, error) {
	// 1. 构建 commit 的文本内容
	content := buildCommitContent(commit)

	// 2. 计算哈希
	h := hash.ComputeHash(hash.CommitObject, content)

	// 3. 构建完整对象: commit <size>\0<content>
	header := []byte(fmt.Sprintf("commit %d", len(content)))
	header = append(header, 0) // null byte
	data := append(header, content...)

	// 4. zlib 压缩
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return hash.Hash{}, fmt.Errorf("zlib compress failed: %v", err)
	}
	if err := zw.Close(); err != nil {
		return hash.Hash{}, fmt.Errorf("zlib close failed: %v", err)
	}

	// 5. 写入文件系统
	hashStr := h.String()
	if len(hashStr) < 2 {
		return hash.Hash{}, fmt.Errorf("invalid hash")
	}

	objDir := filepath.Join(gitDir, "objects", hashStr[:2])
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return hash.Hash{}, fmt.Errorf("failed to create object directory: %v", err)
	}

	objPath := filepath.Join(objDir, hashStr[2:])
	if err := os.WriteFile(objPath, buf.Bytes(), 0444); err != nil {
		return hash.Hash{}, fmt.Errorf("failed to write object file: %v", err)
	}

	return h, nil
}

// buildCommitContent 构建 commit 对象的文本内容
func buildCommitContent(commit *Commit) []byte {
	var buf bytes.Buffer

	// tree 行
	buf.WriteString(fmt.Sprintf("tree %s\n", commit.Tree.String()))

	// parent 行（可以有多个）
	for _, parent := range commit.Parents {
		buf.WriteString(fmt.Sprintf("parent %s\n", parent.String()))
	}

	// author 行
	buf.WriteString(fmt.Sprintf("author %s\n", formatSignature(commit.Author)))

	// committer 行
	buf.WriteString(fmt.Sprintf("committer %s\n", formatSignature(commit.Committer)))

	// 空行分隔
	buf.WriteString("\n")

	// commit message
	buf.WriteString(commit.Message)

	return buf.Bytes()
}

// formatSignature 格式化签名
// 格式: Name <email> timestamp timezone
func formatSignature(sig Signature) string {
	timestamp := sig.When.Unix()
	_, offset := sig.When.Zone()
	timezone := fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60)
	return fmt.Sprintf("%s <%s> %d %s", sig.Name, sig.Email, timestamp, timezone)
}
