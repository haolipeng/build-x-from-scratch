package blob

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"geegit/beginner/day2-read-blob/hash"
)

// ReadBlob 从 .git/objects 读取一个 blob 对象
// gitDir: .git 目录路径
// h: 对象的 SHA-1 哈希值
func ReadBlob(gitDir string, h hash.Hash) (*Blob, error) {
	// 构建对象文件路径
	// 格式: .git/objects/e6/9de29bb2d1d6434b8b29ae775ad8c2e48c5391
	hashStr := h.String()
	if len(hashStr) < 2 {
		return nil, fmt.Errorf("invalid hash")
	}

	objPath := filepath.Join(gitDir, "objects", hashStr[:2], hashStr[2:])

	// 打开文件
	file, err := os.Open(objPath)
	if err != nil {
		return nil, fmt.Errorf("object not found: %v", err)
	}
	defer file.Close()

	// 使用 zlib 解压
	zr, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("zlib decompress failed: %v", err)
	}
	defer zr.Close()

	// 读取所有内容
	data, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	// 解析对象格式: <type> <size>\0<content>
	nullIdx := bytes.IndexByte(data, 0)
	if nullIdx < 0 {
		return nil, fmt.Errorf("invalid object format")
	}

	header := string(data[:nullIdx])
	content := data[nullIdx+1:]

	// 解析 header
	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid object header: %s", header)
	}

	objTypeStr := parts[0]
	if objTypeStr != "blob" {
		return nil, fmt.Errorf("expected blob, got %s", objTypeStr)
	}

	return &Blob{
		Hash: h,
		Data: content,
	}, nil
}
