package blob

import "geegit/beginner/day5-write-tree/hash"

// Blob 表示一个 blob 对象
type Blob struct {
	Hash hash.Hash
	Data []byte
}
