package tree

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"geegit/beginner/day4-read-tree/hash"
)

// ReadTree 从 .git/objects 读取一个 tree 对象
func ReadTree(gitDir string, hash hash.Hash) (*Tree, error) {
	hashStr := hash.String()
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
	if objTypeStr != "tree" {
		return nil, fmt.Errorf("expected tree, got %s", objTypeStr)
	}

	// 解析 tree 内容
	entries, err := parseTreeEntries(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tree entries: %v", err)
	}

	return &Tree{
		Hash:    hash,
		Entries: entries,
	}, nil
}

// parseTreeEntries 解析 tree 对象的二进制内容
// 格式: <mode> <name>\0<20-byte-hash> ...
func parseTreeEntries(data []byte) ([]TreeEntry, error) {
	var entries []TreeEntry
	offset := 0

	for offset < len(data) {
		// 1. 读取模式和名称（以空格分隔）
		spaceIdx := bytes.IndexByte(data[offset:], ' ')
		if spaceIdx < 0 {
			break
		}
		mode := string(data[offset : offset+spaceIdx])
		offset += spaceIdx + 1

		// 2. 读取名称（以 null 结束）
		nullIdx := bytes.IndexByte(data[offset:], 0)
		if nullIdx < 0 {
			return nil, fmt.Errorf("invalid tree format: missing null after name")
		}
		name := string(data[offset : offset+nullIdx])
		offset += nullIdx + 1

		// 3. 读取 20 字节哈希
		if offset+20 > len(data) {
			return nil, fmt.Errorf("invalid tree format: truncated hash")
		}
		var hash hash.Hash
		copy(hash[:], data[offset:offset+20])
		offset += 20

		entries = append(entries, TreeEntry{
			Mode: mode,
			Name: name,
			Hash: hash,
		})
	}

	return entries, nil
}
