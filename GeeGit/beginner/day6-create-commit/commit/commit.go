package commit

import (
	"time"

	"geegit/beginner/day6-create-commit/hash"
)

// Signature 表示 Git 的签名信息（author 或 committer）
type Signature struct {
	Name  string
	Email string
	When  time.Time
}

// Commit 表示一个 commit 对象
type Commit struct {
	Hash      hash.Hash
	Tree      hash.Hash
	Parents   []hash.Hash
	Author    Signature
	Committer Signature
	Message   string
}
