package blob

import "geegit/beginner/day3-write-blob/hash"

// Blob 表示一个 blob 对象
type Blob struct {
	Hash hash.Hash
	Data []byte
}
