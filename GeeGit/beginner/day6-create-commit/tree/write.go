package tree

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"sort"

	"geegit/beginner/day6-create-commit/hash"
)

// buildTreeContent 构建 tree 对象的二进制内容
// 格式: <mode> <name>\0<20-byte-hash> ...
func BuildTreeContent(entries []TreeEntry) []byte {
	// Git 要求 tree 条目按名称排序
	sorted := make([]TreeEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	var buf []byte
	for _, entry := range sorted {
		// mode + space
		buf = append(buf, []byte(entry.Mode)...)
		buf = append(buf, ' ')
		// name + null
		buf = append(buf, []byte(entry.Name)...)
		buf = append(buf, 0)
		// 20-byte hash
		buf = append(buf, entry.Hash[:]...)
	}
	return buf
}

// WriteRawTree 写入原始 tree 内容（用于演示）
func WriteRawTree(gitDir string, content []byte) (hash.Hash, error) {
	h := hash.ComputeHash(hash.TreeObject, content)

	header := []byte(fmt.Sprintf("tree %d", len(content)))
	header = append(header, 0)
	data := append(header, content...)

	// zlib 压缩
	var compressed bytes.Buffer
	zw := zlib.NewWriter(&compressed)
	if _, err := zw.Write(data); err != nil {
		return hash.Hash{}, err
	}
	if err := zw.Close(); err != nil {
		return hash.Hash{}, err
	}

	// 写入文件
	hashStr := h.String()
	objDir := gitDir + "/objects/" + hashStr[:2]
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return hash.Hash{}, err
	}

	objPath := objDir + "/" + hashStr[2:]
	if err := os.WriteFile(objPath, compressed.Bytes(), 0444); err != nil {
		return hash.Hash{}, err
	}

	return h, nil
}

// WriteTree 将 tree 对象写入 .git/objects 目录
func WriteTree(gitDir string, entries []TreeEntry) (hash.Hash, error) {
	// 1. 对条目进行排序（Git 要求）
	sorted := make([]TreeEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	// 2. 构建 tree 的二进制内容
	content := BuildTreeContent(sorted)

	// 3. 计算哈希
	h := hash.ComputeHash(hash.TreeObject, content)

	// 4. 构建完整对象: tree <size>\0<content>
	header := []byte(fmt.Sprintf("tree %d", len(content)))
	header = append(header, 0) // null byte
	data := append(header, content...)

	// 5. zlib 压缩
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return hash.Hash{}, fmt.Errorf("zlib compress failed: %v", err)
	}
	if err := zw.Close(); err != nil {
		return hash.Hash{}, fmt.Errorf("zlib close failed: %v", err)
	}

	// 6. 写入文件系统
	hashStr := h.String()
	if len(hashStr) < 2 {
		return hash.Hash{}, fmt.Errorf("invalid hash")
	}

	objDir := gitDir + "/objects/" + hashStr[:2]
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return hash.Hash{}, fmt.Errorf("failed to create object directory: %v", err)
	}

	objPath := objDir + "/" + hashStr[2:]
	if err := os.WriteFile(objPath, buf.Bytes(), 0444); err != nil {
		return hash.Hash{}, fmt.Errorf("failed to write object file: %v", err)
	}

	return h, nil
}
