# Day 4: Read a tree object

## 学习目标

在 Day 4，你将学会：
- 理解 Git 的目录结构表示（tree 对象）
- 解析 tree 对象的二进制格式
- 理解文件模式（mode）的含义
- 构建目录树的层级结构

## 关键概念

### 1. Tree 对象是什么？

在 Git 中：
- **Blob** 存储文件内容
- **Tree** 存储目录结构

Tree 对象记录了：
- 每个文件/子目录的名称
- 每个文件/子目录的模式（权限）
- 每个文件/子目录指向的对象哈希

### 2. Tree 对象的二进制格式

Tree 对象的内容是**二进制格式**，不是纯文本：

```
<mode> <name>\0<20-byte-hash><mode> <name>\0<20-byte-hash>...
```

**示例**：
```
100644 README.md\0<20字节哈希>100644 main.go\0<20字节哈希>
```

**格式说明**：
- `<mode>`: ASCII 文本，如 "100644"
- ` `: 单个空格
- `<name>`: ASCII 文本文件名
- `\0`: 一个 null 字节
- `<20-byte-hash>`: 原始 20 字节哈希（**不是十六进制字符串**）

### 3. 文件模式（Mode）

Git 使用 Unix 风格的文件模式：

| Mode | 含义 | Git 类型 |
|------|------|---------|
| `100644` | 普通文件 | blob |
| `100755` | 可执行文件 | blob |
| `120000` | 符号链接 | blob |
| `040000` | 目录 | tree |
| `160000` | Gitlink（子模块） | commit |

**注意**：
- 普通文件是 `100644`，不是 Unix 的 `644`
- 目录是 `040000`，不是 Unix 的 `40000`

### 4. Tree 条目的排序

Git 要求 tree 中的条目**必须按名称排序**：
- 按字典序（lexicographic order）
- 目录名在排序时**不带** `/` 后缀

## 代码说明

### TreeEntry 结构体

```go
type TreeEntry struct {
	Mode string // "100644", "040000", etc.
	Name string // 文件或目录名
	Hash Hash   // 指向的对象哈希
}
```

### Tree 结构体

```go
type Tree struct {
	Hash    Hash
	Entries []TreeEntry
}
```

### ReadTree() 函数

`read.go` 中的 `ReadTree()` 函数：

```go
func ReadTree(gitDir string, hash Hash) (*Tree, error)
```

**执行步骤**：
1. 读取对象文件（与 ReadBlob 类似）
2. zlib 解压
3. 解析头部 `tree <size>\0`
4. 调用 `parseTreeEntries()` 解析二进制内容

### parseTreeEntries() 函数

解析 tree 的二进制内容：

```go
func parseTreeEntries(data []byte) ([]TreeEntry, error)
```

**解析流程**：
1. 读取到空格：得到 `mode`
2. 读取到 `\0`：得到 `name`
3. 读取 20 字节：得到 `hash`
4. 重复以上步骤，直到数据结束

### 文件结构

```
day4-read-tree/
├── hash.go       - Hash 类型（来自 Day 1）
├── object.go     - **扩展** 添加 TreeEntry 和 Tree
├── init.go       - InitRepository()（来自 Day 1）
├── read.go       - **扩展** 添加 ReadTree() 和 parseTreeEntries()
├── write.go      - WriteBlob()（来自 Day 3）
└── main.go       - 演示程序
```

## 运行演示

```bash
cd beginner/day4-read-tree
go run *.go
```

**预期输出**：
```
=== Day 4: Read a tree object ===

✓ Step 1: Initialize test repository at ./test-repo-day4

✓ Step 2: Create some blob objects
  Created README.md blob: <hash>
  Created main.go blob: <hash>

✓ Step 3: Create a tree object manually
  Created tree: <hash>

✓ Step 4: Read the tree object

Result:
  Tree Hash: <40-character-hex>
  Entries (2):
    [1] 100644 <hash> README.md
    [2] 100644 <hash> main.go

✓ Step 5: Verify each entry is readable
  ✓ README.md (27 bytes)
  ✓ main.go (47 bytes)

Verification with real Git:
  You can verify this tree using real git commands:
  $ cd ./test-repo-day4
  $ git ls-tree <hash>
  $ git cat-file -p <hash>

=== Summary ===
✓ Created a tree with 2 entries
✓ Successfully read and parsed tree object
✓ All entries are valid blob objects

You have successfully implemented Git's tree reading!
```

## 与真实 Git 对比

| 我们实现的 | Git 命令 | 说明 |
|-----------|---------|------|
| `ReadTree(gitDir, hash)` | `git ls-tree <hash>` | 读取 tree 对象 |
| `tree.Entries` | `git ls-tree <hash>` | 显示条目列表 |
| `entry.Hash` | `git ls-tree <hash>` | 每个条目的哈希 |

### 手动验证示例

```bash
# 1. 运行我们的程序
cd beginner/day4-read-tree
go run *.go

# 2. 进入测试仓库
cd test-repo-day4

# 3. 用 git 查看 tree
git ls-tree <tree-hash>
# 输出:
# 100644 blob <hash>	README.md
# 100644 blob <hash>	main.go

# 4. 用 cat-file 查看原始内容
git cat-file -p <tree-hash>
# 输出相同

# 5. 查看 tree 的二进制内容（十六进制）
git cat-file tree <tree-hash> | xxd
```

## 深入理解

### 为什么 Tree 用二进制格式？

1. **空间效率**：
   - 十六进制哈希：40 字节
   - 二进制哈希：20 字节
   - 节省 50% 空间

2. **解析效率**：
   - 固定 20 字节，无需解析
   - 直接内存拷贝

### Tree 如何表示嵌套目录？

Tree 可以嵌套：

```
root-tree:
  100644 blob <hash>  README.md
  040000 tree <hash>  src/
    ↓
src-tree:
  100644 blob <hash>  main.go
  040000 tree <hash>  utils/
    ↓
utils-tree:
  100644 blob <hash>  helper.go
```

每个 `040000` 模式的条目指向另一个 tree 对象。

### 为什么条目要排序？

1. **规范化**：相同目录结构→相同哈希
2. **高效查找**：可以用二分搜索
3. **Delta 压缩**：排序后的数据更容易压缩

### 空目录怎么办？

Git **不存储空目录**！

如果需要保留空目录，常见做法：
```bash
touch empty-dir/.gitkeep
```

## 练习题

### 练习 1：计算 Tree 哈希
手动计算以下目录结构的 tree 哈希：

```
README.md (blob: e69de29bb2d1d6434b8b29ae775ad8c2e48c5391)
```

提示：
1. 构建内容：`100644 README.md\0<20字节>`
2. 计算：`sha1("tree <size>\0<content>")`

### 练习 2：实现 ls-tree 命令
编写一个工具，模拟 `git ls-tree`：

```bash
./geegit ls-tree <tree-hash>
```

输出格式：
```
100644 blob <hash>    README.md
100644 blob <hash>    main.go
```

### 练习 3：递归遍历 Tree
修改 `main.go`，添加递归遍历：

```go
func PrintTreeRecursive(gitDir string, hash Hash, prefix string) {
	tree, _ := ReadTree(gitDir, hash)
	for _, entry := range tree.Entries {
		fmt.Printf("%s%s %s\n", prefix, entry.Mode, entry.Name)
		if entry.Mode == "040000" {
			PrintTreeRecursive(gitDir, entry.Hash, prefix+"  ")
		}
	}
}
```

### 练习 4：验证 Tree 排序
创建一个包含多个文件的 tree，验证：
- 条目是否按名称排序？
- 如果不排序，哈希会变化吗？

## 常见问题

### Q1: 为什么哈希是 20 字节而不是 40 字符？
**A**:
- 十六进制字符串（"a1b2..."）：40 字节
- 原始二进制（[0xa1, 0xb2, ...]）：20 字节
- Tree 存储原始字节以节省空间

### Q2: 如何区分文件和目录？
**A**: 通过 `mode` 字段：
- `100644` 或 `100755` → 文件（blob）
- `040000` → 目录（tree）

### Q3: 为什么有些文件是 100755？
**A**: `100755` 表示可执行文件：
```bash
chmod +x script.sh
git add script.sh
```

### Q4: Gitlink 是什么？
**A**: `160000` 模式表示 Git 子模块（submodule）：
```
160000 commit <hash>  lib/external
```
指向另一个仓库的 commit。

## 与 Day 3 的联系

现在我们有了：
- **Day 2**: 读取 blob（文件内容）
- **Day 3**: 写入 blob
- **Day 4**: 读取 tree（目录结构）

Tree 和 Blob 的关系：
```
Tree
├── entry[0] → Blob (README.md)
├── entry[1] → Blob (main.go)
└── entry[2] → Tree (子目录)
```

## 下一步

在 Day 5，我们将学习：
- 如何**创建** tree 对象
- 如何从工作目录构建 tree
- 如何正确排序和格式化 tree 条目

有了 write tree 的能力，我们就可以在 Day 6 创建完整的 commit 了！

---

**恭喜你完成 Day 4！** 🎉

你现在已经掌握了 Git 目录结构的核心原理。继续前进！
