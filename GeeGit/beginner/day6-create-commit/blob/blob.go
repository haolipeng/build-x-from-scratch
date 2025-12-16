package blob

import "geegit/beginner/day6-create-commit/hash"

// Blob 表示一个 blob 对象
type Blob struct {
	Hash hash.Hash
	Data []byte
}
