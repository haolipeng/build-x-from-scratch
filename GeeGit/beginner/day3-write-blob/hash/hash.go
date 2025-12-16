package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// Hash 表示 Git 对象的 SHA-1 哈希值（20 字节）
type Hash [20]byte

// String 返回哈希的十六进制字符串表示
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// ObjectType 表示 Git 对象的类型
type ObjectType int

const (
	CommitObject ObjectType = iota
	TreeObject
	BlobObject
)

// String 返回对象类型的字符串表示
func (t ObjectType) String() string {
	switch t {
	case CommitObject:
		return "commit"
	case TreeObject:
		return "tree"
	case BlobObject:
		return "blob"
	default:
		return "unknown"
	}
}

// ComputeHash 计算给定对象类型和内容的哈希
// Git 对象的格式: <type> <size>\0<content>
func ComputeHash(objType ObjectType, content []byte) Hash {
	// 构建完整的对象内容
	header := []byte(objType.String())
	header = append(header, ' ')
	header = append(header, []byte(fmt.Sprintf("%d", len(content)))...)
	header = append(header, 0) // null byte

	// 计算整个对象的 SHA-1 哈希
	data := append(header, content...)
	return sha1.Sum(data)
}
