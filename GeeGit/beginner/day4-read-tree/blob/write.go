package blob

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"

	"geegit/beginner/day4-read-tree/hash"
)

// WriteBlob 将 blob 对象写入 .git/objects 目录
func WriteBlob(gitDir string, content []byte) (hash.Hash, error) {
	// 1. 计算哈希
	h := hash.ComputeHash(hash.BlobObject, content)

	// 2. 构建对象内容: blob <size>\0<content>
	header := []byte(fmt.Sprintf("blob %d", len(content)))
	header = append(header, 0) // null byte
	data := append(header, content...)

	// 3. zlib 压缩
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		return hash.Hash{}, fmt.Errorf("zlib compress failed: %v", err)
	}
	if err := zw.Close(); err != nil {
		return hash.Hash{}, fmt.Errorf("zlib close failed: %v", err)
	}

	// 4. 写入文件系统
	hashStr := h.String()
	if len(hashStr) < 2 {
		return hash.Hash{}, fmt.Errorf("invalid hash")
	}

	// 创建目录: .git/objects/xx/
	objDir := filepath.Join(gitDir, "objects", hashStr[:2])
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return hash.Hash{}, fmt.Errorf("failed to create object directory: %v", err)
	}

	// 写入文件: .git/objects/xx/xxxxx...
	objPath := filepath.Join(objDir, hashStr[2:])
	if err := os.WriteFile(objPath, buf.Bytes(), 0444); err != nil {
		return hash.Hash{}, fmt.Errorf("failed to write object file: %v", err)
	}

	return h, nil
}
