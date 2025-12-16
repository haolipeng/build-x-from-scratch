package tree

import "geegit/beginner/day5-write-tree/hash"

// TreeEntry 表示 tree 对象中的一个条目
type TreeEntry struct {
	Mode string      // 文件模式
	Name string      // 文件或目录名
	Hash hash.Hash   // 指向的对象哈希
}

// Tree 表示一个 tree 对象（目录结构）
type Tree struct {
	Hash    hash.Hash
	Entries []TreeEntry
}
