package blob

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"geegit/beginner/day6-create-commit/hash"
)

// ReadBlob 从 .git/objects 读取一个 blob 对象
func ReadBlob(gitDir string, h hash.Hash) (*Blob, error) {
	hashStr := h.String()
	if len(hashStr) < 2 {
		return nil, fmt.Errorf("invalid hash")
	}

	objPath := filepath.Join(gitDir, "objects", hashStr[:2], hashStr[2:])

	file, err := os.Open(objPath)
	if err != nil {
		return nil, fmt.Errorf("object not found: %v", err)
	}
	defer file.Close()

	zr, err := zlib.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("zlib decompress failed: %v", err)
	}
	defer zr.Close()

	data, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	nullIdx := bytes.IndexByte(data, 0)
	if nullIdx < 0 {
		return nil, fmt.Errorf("invalid object format")
	}

	header := string(data[:nullIdx])
	content := data[nullIdx+1:]

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
