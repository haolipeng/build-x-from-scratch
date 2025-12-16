# Day 6: Create a commit

## 学习目标

在 Day 6，你将学会：
- 理解 Git commit 对象的结构
- 学习 commit 的元数据（author、committer、timestamp）
- 创建 commit 对象并写入仓库
- 构建 commit 历史链（parent 关系）
- 理解 commit 的完整生命周期

## 关键概念

### 1. Commit 是什么？

Commit 是 Git 的核心概念，它记录了：
- **快照**：tree 对象（项目在某个时刻的完整状态）
- **历史**：parent commit（从哪个 commit 演化而来）
- **作者**：谁创建了这些更改
- **时间**：何时创建的
- **原因**：为什么做这些更改（commit message）

### 2. Commit 对象的格式

Commit 是**纯文本**格式（不像 tree 是二进制）：

```
tree <tree-hash>
parent <parent-hash> (可选，可以有多个)
author Name <email> timestamp timezone
committer Name <email> timestamp timezone

Commit message goes here.

Can have multiple lines.
```

**示例**：
```
tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904
parent a1b2c3d4e5f6789012345678901234567890abcd
author Alice <alice@example.com> 1609459200 +0800
committer Alice <alice@example.com> 1609459200 +0800

Initial commit

This is the first commit in the repository.
```

### 3. Author vs Committer

- **Author**：创建这些更改的人
- **Committer**：将更改提交到仓库的人

**何时不同？**
- 通过 email patch 应用更改
- 使用 `git commit --amend`
- Cherry-pick 其他人的 commit

大多数情况下，author 和 committer 相同。

### 4. Timestamp 格式

Git 使用 Unix 时间戳 + 时区：

```
1609459200 +0800
│          │
│          └─ 时区（UTC+8）
└─ Unix 时间戳（秒）
```

**时区格式**：
- `+0800`：UTC+8（北京时间）
- `-0500`：UTC-5（美国东部时间）
- `+0000`：UTC

### 5. Parent 关系

- **Root commit**：没有 parent（第一个 commit）
- **普通 commit**：一个 parent
- **Merge commit**：两个或多个 parent

```
C3 (merge)
├─ parent: C1
└─ parent: C2

C2
└─ parent: C1

C1 (root)
└─ (no parent)
```

## 代码说明

### Signature 结构体

```go
type Signature struct {
	Name  string    // 姓名
	Email string    // 邮箱
	When  time.Time // 时间
}
```

### Commit 结构体

```go
type Commit struct {
	Hash      Hash      // commit 的哈希
	Tree      Hash      // 指向的 tree 对象
	Parents   []Hash    // 父 commit（可以有多个）
	Author    Signature // 作者信息
	Committer Signature // 提交者信息
	Message   string    // 提交信息
}
```

### WriteCommit() 函数

`write.go` 中的 `WriteCommit()` 函数：

```go
func WriteCommit(gitDir string, commit *Commit) (Hash, error)
```

**执行步骤**：
1. 调用 `buildCommitContent()` 构建文本内容
2. 计算哈希
3. 添加头部：`commit <size>\0<content>`
4. zlib 压缩
5. 写入 `.git/objects/`

### buildCommitContent() 函数

构建 commit 的文本内容：

```go
func buildCommitContent(commit *Commit) []byte
```

**格式化规则**：
- `tree` 行：必须
- `parent` 行：可选，可多个
- `author` 行：必须
- `committer` 行：必须
- 空行：分隔元数据和消息
- message：可多行

### ReadCommit() 函数

`read.go` 中新增的 `ReadCommit()` 函数：

```go
func ReadCommit(gitDir string, hash Hash) (*Commit, error)
```

解析 commit 对象的文本内容，提取所有字段。

### 文件结构

```
day6-create-commit/
├── hash.go       - Hash 类型（来自 Day 1）
├── object.go     - **扩展** 添加 Signature 和 Commit 结构体
├── init.go       - InitRepository()（来自 Day 1）
├── read.go       - **扩展** 添加 ReadCommit() 和 parseCommit()
├── write.go      - **扩展** 添加 WriteCommit() 和 buildCommitContent()
└── main.go       - 演示程序
```

## 运行演示

```bash
cd beginner/day6-create-commit
go run *.go
```

**预期输出**：
```
=== Day 6: Create a commit ===

✓ Step 1: Initialize test repository at ./test-repo-day6

✓ Step 2: Create blob objects
  Created README.md blob: <hash>...
  Created main.go blob: <hash>...

✓ Step 3: Create tree object
  Created tree: <hash>...

✓ Step 4: Create initial commit (no parent)
  Commit Hash: <40-character-hex>
  Author: Alice Developer <alice@example.com>
  Message: Initial commit...

✓ Step 5: Read back the commit to verify

Commit Details:
  Tree: <hash>...
  Parents: 0
  Author: Alice Developer <alice@example.com>
  Committer: Alice Developer <alice@example.com>
  Message:
    Initial commit...

✓ Step 6: Create second commit with parent
  Commit Hash: <hash>
  Parent: <hash>...
  Author: Bob Contributor

✓ Step 7: Display commit history chain

  Commit Chain:
    [<hash>] (HEAD)
    │ Update README...
    │ by Bob Contributor
    ↓
    [<hash>] (root)
      Initial commit...
      by Alice Developer

Verification with real Git:
  You can verify these commits using real git commands:
  $ cd ./test-repo-day6
  $ git cat-file -t <hash>  # should output: commit
  $ git cat-file -p <hash>  # show commit content
  $ git log --oneline <hash>  # show commit history

=== Summary ===
✓ Created 2 blob objects
✓ Created 2 tree objects
✓ Created initial commit (no parent)
✓ Created second commit (with parent)
✓ Verified commit reading works correctly
✓ Built a commit history chain

You have successfully implemented Git's commit creation!
```

## 与真实 Git 对比

| 我们实现的 | Git 命令 | 说明 |
|-----------|---------|------|
| `WriteCommit(gitDir, commit)` | `git commit` | 创建 commit |
| `commit.Tree` | `git write-tree` | Tree 哈希 |
| `commit.Parents` | `git commit-tree -p` | Parent 指定 |
| `commit.Message` | `git commit -m` | Commit message |

### 手动验证示例

```bash
# 1. 运行我们的程序
cd beginner/day6-create-commit
go run *.go

# 2. 进入测试仓库
cd test-repo-day6

# 3. 查看 commit
git cat-file -p <commit-hash>
# 输出：
# tree <tree-hash>
# author Alice Developer <alice@example.com> 1234567890 +0800
# committer Alice Developer <alice@example.com> 1234567890 +0800
#
# Initial commit
#
# Add README and main.go files.

# 4. 查看 commit 历史
git log --oneline <commit2-hash>
# 输出：
# <hash2> Update README
# <hash1> Initial commit

# 5. 查看差异
git diff <commit1-hash> <commit2-hash>
```

## 深入理解

### Commit 的不可变性

Commit 创建后**永不改变**：
- 修改任何字段（tree、parent、message）→ 新哈希 → 新 commit
- `git commit --amend` 实际创建**新 commit**，旧 commit 仍存在

### Commit 和 Tree 的关系

```
Commit A                    Commit B
├─ tree: abc123            ├─ tree: def456
│  ├─ README.md (v1)       │  ├─ README.md (v2)  ← 变化
│  └─ main.go (v1)         │  └─ main.go (v1)    ← 未变
└─ parent: (none)          └─ parent: A
```

**关键点**：
- Commit 指向完整的 tree（整个项目快照）
- 即使只改一个文件，也会创建新 tree
- 但 tree 中未变的文件仍引用旧 blob

### 为什么需要 Author 和 Committer？

**场景 1：Patch workflow**
```bash
# Alice 创建 patch
git format-patch -1

# Bob 应用 patch
git am < 0001-fix-bug.patch
# Author: Alice
# Committer: Bob
```

**场景 2：Rebase**
```bash
git rebase master
# Author: 原作者
# Committer: 执行 rebase 的人（时间也更新）
```

### 时间戳的意义

```go
Author Date:    Wed Jan 1 00:00:00 2020 +0800  // 何时写的代码
Commit Date:    Thu Jan 2 10:00:00 2020 +0800  // 何时提交的
```

这对于理解开发历史很重要！

### Merge Commit 的特殊性

Merge commit 有**多个 parent**：

```
tree <tree-hash>
parent <parent1-hash>
parent <parent2-hash>
author ...
committer ...

Merge branch 'feature' into master
```

我们的实现已支持多个 parent（`Parents []Hash`）。

## 练习题

### 练习 1：计算 Commit 哈希
手动计算以下 commit 的哈希：

```
tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904
author Test <test@example.com> 0 +0000
committer Test <test@example.com> 0 +0000

test
```

提示：
1. 构建完整内容（注意换行）
2. 添加头部：`commit <size>\0<content>`
3. 计算 SHA-1

### 练习 2：实现 commit-tree 命令
编写工具模拟 `git commit-tree`：

```bash
./geegit commit-tree <tree-hash> -p <parent-hash> -m "message"
```

### 练习 3：显示 Commit 历史
实现函数显示 commit 链：

```go
func PrintHistory(gitDir string, commitHash Hash) {
	// 递归读取 parent，显示历史
}
```

输出格式：
```
<hash> Initial commit (Alice, 2024-01-01)
```

### 练习 4：实现简单的 log
实现 `git log` 的简化版：

```go
func Log(gitDir string, commitHash Hash, maxCount int) {
	// 显示最近 maxCount 个 commit
}
```

## 常见问题

### Q1: 为什么 message 可以有多行？
**A**: Commit message 通常包含：
- 第一行：简短总结（<50 字符）
- 空行
- 详细描述

这是 Git 的最佳实践。

### Q2: Parent 可以为空吗？
**A**: 可以！这叫 **root commit**（初始 commit）。每个仓库至少有一个 root commit。

### Q3: 可以有多个 root commit 吗？
**A**: 可以（罕见）！例如：
- 合并两个独立的仓库
- `git checkout --orphan`

### Q4: Commit 对象可以删除吗？
**A**: 技术上可以（删除文件），但：
- 如果有引用指向它（branch、tag），Git 不会删除
- `git gc` 会清理无引用的 commit

### Q5: 时区为什么重要？
**A**: 跨时区协作时，时区信息帮助理解实际时间：
```
Author Date: 2024-01-01 08:00:00 +0800  // 北京时间早上8点
等价于:      2024-01-01 00:00:00 +0000  // UTC时间午夜
```

## 与 Day 5 的联系

现在我们有了完整的对象模型：

- **Day 2-3**: Blob（文件内容）
- **Day 4-5**: Tree（目录结构）
- **Day 6**: Commit（快照 + 历史 + 元数据）

完整流程：
```
1. 创建文件 blob (WriteBlob)
2. 创建 tree 引用 blob (WriteTree)
3. 创建 commit 引用 tree (WriteCommit)
4. 更新 branch 指向 commit
```

## 下一步

在 Day 7，我们将学习：
- 如何**克隆**远程仓库（简化版）
- Git 网络协议（HTTP Smart Protocol）
- Pkt-line 编解码
- Packfile 基础（无 Delta）

这是最后一天的初级课程！完成后你将理解 Git 的核心原理。

---

**恭喜你完成 Day 6！** 🎉

你现在已经掌握了 Git 的三大核心对象（blob、tree、commit）的完整生命周期。只差最后一步：从远程获取这些对象！
