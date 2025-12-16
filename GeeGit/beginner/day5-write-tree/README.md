# Day 5: Write a tree object

## 学习目标

在 Day 5，你将学会：
- 如何创建 Git tree 对象
- 如何正确格式化 tree 的二进制内容
- 理解 tree 条目的排序规则
- 构建嵌套的目录结构

## 关键概念

### 1. Tree 对象的写入流程

创建 tree 对象的步骤：
1. **准备条目**：收集所有文件和子目录的信息
2. **排序**：按名称对条目排序
3. **序列化**：转换为二进制格式
4. **计算哈希**：对整个对象计算 SHA-1
5. **压缩**：使用 zlib 压缩
6. **写入**：保存到 `.git/objects/`

### 2. Tree 条目的排序

**关键规则**：Tree 中的条目**必须**按名称字典序排序。

```go
// 正确的排序
entries := []TreeEntry{
	{Name: "LICENSE", ...},
	{Name: "README.md", ...},
	{Name: "main.go", ...},
	{Name: "src", ...},  // 目录不带 / 后缀
}
```

**为什么要排序？**
1. **规范化**：相同内容总是产生相同哈希
2. **高效查找**：可以使用二分搜索
3. **Delta 压缩**：有序数据更容易压缩

### 3. 文件模式的使用

创建 tree 时需要指定正确的模式：

| 文件类型 | Mode | 示例 |
|---------|------|------|
| 普通文件 | `100644` | README.md, LICENSE |
| 可执行文件 | `100755` | build.sh, run |
| 符号链接 | `120000` | link-to-file |
| 子目录 | `040000` | src/, docs/ |

### 4. 嵌套 Tree 结构

Tree 可以包含其他 tree：

```
Root Tree (hash: aaa...)
├── 100644 blob bbb... README.md
├── 100644 blob ccc... main.go
└── 040000 tree ddd... src/
    └── Src Tree (hash: ddd...)
        └── 100644 blob eee... util.go
```

创建嵌套结构的顺序：**从内到外**
1. 先创建最深层的 tree（src/）
2. 再创建父 tree，引用子 tree 的哈希

## 代码说明

### WriteTree() 函数

`write.go` 中的 `WriteTree()` 函数：

```go
func WriteTree(gitDir string, entries []TreeEntry) (Hash, error)
```

**参数**：
- `gitDir`: .git 目录路径
- `entries`: tree 的条目列表（无需预先排序）

**返回值**：
- 成功：tree 对象的哈希
- 失败：错误信息

**执行步骤**：
1. 复制并排序条目
2. 调用 `buildTreeContent()` 生成二进制内容
3. 计算哈希
4. 构建完整对象：`tree <size>\0<content>`
5. zlib 压缩
6. 写入文件系统

### buildTreeContent() 函数

构建 tree 的二进制内容：

```go
func buildTreeContent(entries []TreeEntry) []byte
```

**格式**：`<mode> <name>\0<20-byte-hash><mode> <name>\0<20-byte-hash>...`

**示例**：
```
输入：
  {Mode: "100644", Name: "README.md", Hash: [20]byte{...}}
  {Mode: "100644", Name: "main.go", Hash: [20]byte{...}}

输出（二进制）：
  100644 README.md\0<20字节>100644 main.go\0<20字节>
```

### 文件结构

```
day5-write-tree/
├── hash.go       - Hash 类型（来自 Day 1）
├── object.go     - TreeEntry, Tree（来自 Day 4）
├── init.go       - InitRepository()（来自 Day 1）
├── read.go       - ReadBlob(), ReadTree()（来自 Day 4）
├── write.go      - **扩展** 添加 WriteTree() 和 buildTreeContent()
└── main.go       - 演示程序
```

## 运行演示

```bash
cd beginner/day5-write-tree
go run *.go
```

**预期输出**：
```
=== Day 5: Write a tree object ===

✓ Step 1: Initialize test repository at ./test-repo-day5

✓ Step 2: Create blob objects for files
  Created README.md: <hash>...
  Created main.go: <hash>...
  Created LICENSE: <hash>...

✓ Step 3: Create tree object with these files
  Tree written successfully!
  Tree Hash: <40-character-hex>

✓ Step 4: Read back the tree to verify

Result:
  Tree contains 3 entries:
    [1] 100644 LICENSE <hash>...
    [2] 100644 README.md <hash>...
    [3] 100644 main.go <hash>...

✓ Step 5: Verify entries are sorted correctly
  ✓ Entries are correctly sorted by name

✓ Step 6: Create a nested tree structure
  Created src/ tree: <hash>...
  Created root tree: <hash>...

  Root tree structure:
    100644 blob LICENSE
    100644 blob README.md
    100644 blob main.go
    040000 tree src

Verification with real Git:
  You can verify these trees using real git commands:
  $ cd ./test-repo-day5
  $ git ls-tree <hash>
  $ git cat-file -p <hash>

=== Summary ===
✓ Created multiple blob objects
✓ Created flat tree with 3 files
✓ Created nested tree with subdirectory
✓ Verified tree reading works correctly
✓ All entries are properly sorted

You have successfully implemented Git's tree writing!
```

## 与真实 Git 对比

| 我们实现的 | Git 命令 | 说明 |
|-----------|---------|------|
| `WriteTree(gitDir, entries)` | `git write-tree` | 创建 tree 对象 |
| `entries` 参数 | 工作目录状态 | Git 从索引读取 |
| `treeHash` | `git write-tree` 输出 | 返回的哈希值 |

### 手动验证示例

```bash
# 1. 运行我们的程序
cd beginner/day5-write-tree
go run *.go

# 2. 进入测试仓库
cd test-repo-day5

# 3. 查看 tree
git ls-tree <tree-hash>
# 输出：
# 100644 blob <hash>	LICENSE
# 100644 blob <hash>	README.md
# 100644 blob <hash>	main.go

# 4. 查看嵌套 tree
git ls-tree <root-tree-hash>
# 输出：
# 100644 blob <hash>	LICENSE
# 100644 blob <hash>	README.md
# 100644 blob <hash>	main.go
# 040000 tree <hash>	src

# 5. 递归查看
git ls-tree -r <root-tree-hash>
# 输出：
# 100644 blob <hash>	LICENSE
# 100644 blob <hash>	README.md
# 100644 blob <hash>	main.go
# 100644 blob <hash>	src/util.go
```

## 深入理解

### 为什么排序如此重要？

不排序的后果：

```go
// 情况 1：
entries1 := []TreeEntry{
	{Name: "a.txt", ...},
	{Name: "b.txt", ...},
}
hash1 := WriteTree(gitDir, entries1)

// 情况 2：相同内容，不同顺序
entries2 := []TreeEntry{
	{Name: "b.txt", ...},
	{Name: "a.txt", ...},
}
hash2 := WriteTree(gitDir, entries2)

// 如果不排序：hash1 != hash2（错误！）
// 正确实现：hash1 == hash2（因为内容相同）
```

### Tree 和文件系统的对应关系

```
文件系统：
my-project/
├── README.md (文件)
├── LICENSE (文件)
├── main.go (文件)
└── src/ (目录)
    └── util.go (文件)

Git Tree：
Root Tree:
  100644 blob <hash> LICENSE
  100644 blob <hash> README.md
  100644 blob <hash> main.go
  040000 tree <hash> src
    ↓
Src Tree:
  100644 blob <hash> util.go
```

### Git 如何处理空目录？

Git **不存储**空目录！

```bash
mkdir empty-dir
git add empty-dir/
# Git 不会 add 这个目录

# 常见解决方案：
touch empty-dir/.gitkeep
git add empty-dir/.gitkeep
```

### 目录模式的特殊性

目录的模式是 `040000`，但在某些上下文中，Git 会将其视为 `40000`：

```bash
# Git 内部存储
040000

# Git 命令输出
$ git ls-tree <hash>
040000 tree <hash>	src
```

我们的实现使用 `"040000"`（字符串），与 Git 内部格式一致。

## 练习题

### 练习 1：验证哈希计算
创建一个简单的 tree，手动计算其哈希：

```
Tree:
  100644 blob e69de29bb2d1d6434b8b29ae775ad8c2e48c5391 file.txt
```

步骤：
1. 构建内容：`100644 file.txt\0<20字节>`
2. 添加头部：`tree <size>\0<content>`
3. 计算 SHA-1

用我们的程序验证结果。

### 练习 2：实现 write-tree 命令
编写工具模拟 `git write-tree`：

```bash
# 从工作目录创建 tree
./geegit write-tree /path/to/directory
```

需要：
- 遍历目录
- 为每个文件创建 blob
- 递归处理子目录
- 返回 root tree 哈希

### 练习 3：比较两个 Tree
实现函数比较两个 tree 的差异：

```go
func DiffTrees(gitDir string, hash1, hash2 Hash) {
	// 输出添加、删除、修改的文件
}
```

### 练习 4：优化排序
当前实现总是排序。如果条目已排序，可以跳过：

```go
func WriteTree(gitDir string, entries []TreeEntry) (Hash, error) {
	if !isSorted(entries) {
		sort.Slice(entries, ...)
	}
	// ...
}
```

实现 `isSorted()` 函数。

## 常见问题

### Q1: 为什么不能直接传入目录路径？
**A**: 为了保持函数简单。真实的 `git write-tree` 从**索引**（staging area）读取，不是直接从目录。我们的实现要求调用者提供条目列表。

### Q2: 如果条目重名怎么办？
**A**: Git 不允许同一 tree 中有重名条目。我们的实现没有检查，依赖调用者保证唯一性。生产代码应该添加验证。

### Q3: 模式可以是其他值吗？
**A**: Git 支持的模式有限：
- `100644` - 普通文件
- `100755` - 可执行文件
- `120000` - 符号链接
- `040000` - 目录
- `160000` - Git 子模块

其他值会被 Git 拒绝。

### Q4: Tree 对象可以为空吗？
**A**: 可以！空 tree 是合法的（包含 0 个条目）。内容只有头部：`tree 0\0`。

## 与 Day 4 的联系

现在我们完整实现了 tree 的读写：

- **Day 4**: `ReadTree()` - 读取 tree 对象
- **Day 5**: `WriteTree()` - 创建 tree 对象

完整流程：
```
1. 创建文件 blob (WriteBlob)
2. 创建 tree 引用这些 blob (WriteTree)
3. 读取 tree 验证 (ReadTree)
4. 读取每个 blob (ReadBlob)
```

## 下一步

在 Day 6，我们将学习：
- 如何创建 **commit** 对象
- Commit 的格式和元数据
- 如何链接 tree、author、message

Commit 是 Git 的核心！它将 tree（快照）、author（谁）、时间（when）和 message（为什么）组合在一起。

---

**恭喜你完成 Day 5！** 🎉

你现在已经掌握了 Git 目录结构的完整生命周期。继续前进！
