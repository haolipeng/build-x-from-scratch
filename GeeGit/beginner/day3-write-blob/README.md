# Day 3: Create a blob object

## 学习目标

在 Day 3，你将学会：
- 如何计算 Git 对象的 SHA-1 哈希
- 如何使用 zlib 压缩对象内容
- 如何将 blob 对象写入 `.git/objects` 目录
- 理解 Git 对象存储的完整流程

## 关键概念

### 1. Git 对象的存储格式

Git 存储对象时使用以下格式：
```
<type> <size>\0<content>
```

例如，对于内容 "hello world"，存储格式为：
```
blob 11\0hello world
```

### 2. SHA-1 哈希计算

Git 对**整个对象内容**（包括头部）计算 SHA-1 哈希：

```go
hash := sha1.Sum([]byte("blob 11\0hello world"))
// 结果: 95d09f2b10159347eece71399a7e2e907ea3df4f
```

### 3. zlib 压缩

计算哈希后，Git 使用 zlib 压缩整个对象：

```go
zw := zlib.NewWriter(file)
zw.Write([]byte("blob 11\0hello world"))
```

### 4. 文件系统存储

压缩后的内容存储在 `.git/objects/xx/yyyyyy...`：
- 前 2 个字符作为目录名
- 剩余 38 个字符作为文件名

例如，哈希 `95d09f2b...` 存储在：
```
.git/objects/95/d09f2b10159347eece71399a7e2e907ea3df4f
```

## 代码说明

### WriteBlob() 函数

`write.go` 中的 `WriteBlob()` 函数实现了 blob 对象的写入：

```go
func WriteBlob(gitDir string, content []byte) (Hash, error)
```

**执行步骤**：
1. 计算哈希：`hash := ComputeHash(BlobObject, content)`
2. 构建对象内容：`blob <size>\0<content>`
3. zlib 压缩
4. 创建目录：`.git/objects/xx/`
5. 写入文件：`.git/objects/xx/yyyyyy...`

**返回值**：
- 成功：返回 Hash 和 nil
- 失败：返回零值 Hash 和错误信息

### 文件结构

```
day3-write-blob/
├── hash.go       - Hash 类型和计算（来自 Day 1）
├── object.go     - ObjectType 和 Blob（来自 Day 2）
├── init.go       - InitRepository()（来自 Day 1）
├── read.go       - ReadBlob()（来自 Day 2）
├── write.go      - **新增** WriteBlob()
└── main.go       - 演示程序
```

## 运行演示

```bash
cd beginner/day3-write-blob
go run *.go
```

**预期输出**：
```
=== Day 3: Create a blob object ===

✓ Step 1: Initialize test repository at ./test-repo-day3

✓ Step 2: Create blob with content:
   "Hello, GeeGit!
This is my first blob object.
"

✓ Step 3: Blob object written successfully

Result:
  Hash: <40-character-hex-string>
  Path: .git/objects/xx/xxxxx...

✓ Step 4: Read back the blob to verify
  ✓ Content matches! Verification successful.

Verification with real Git:
  You can verify this blob using real git commands:
  $ cd ./test-repo-day3
  $ git cat-file -t <hash>  # should output: blob
  $ git cat-file -p <hash>  # should output the content
  $ git cat-file -s <hash>  # should output: 46

✓ Step 5: Create another blob (shorter content)
  Hash: <another-hash>

=== Summary ===
✓ Created 2 blob objects
✓ Verified blob reading works correctly
✓ All objects stored in .git/objects/

You have successfully implemented Git's blob writing!
```

## 与真实 Git 对比

| 我们实现的 | Git 命令 | 说明 |
|-----------|---------|------|
| `WriteBlob(gitDir, content)` | `git hash-object -w <file>` | 创建 blob 对象 |
| `ReadBlob(gitDir, hash)` | `git cat-file -p <hash>` | 读取 blob 内容 |
| `hash.String()` | `git hash-object <file>` | 计算哈希（不写入） |

### 手动验证示例

```bash
# 1. 运行我们的程序
cd beginner/day3-write-blob
go run *.go

# 2. 进入测试仓库
cd test-repo-day3

# 3. 用 git 验证对象
git cat-file -t <hash>  # 输出: blob
git cat-file -p <hash>  # 输出: Hello, GeeGit!...

# 4. 自己创建一个 blob 对比
echo "test" | git hash-object -w --stdin
# 这个哈希应该和我们程序中第二个 blob 的哈希一致！
```

## 深入理解

### 为什么要分目录存储？

如果所有对象都放在 `.git/objects/` 下，会有几十万个文件，导致：
- 文件系统性能下降
- 目录列表缓慢

使用前 2 个字符作为目录名，将对象分散到 256 个子目录中：
- 每个目录平均只有几百个文件
- 文件系统性能更好

### 为什么要 zlib 压缩？

1. **节省磁盘空间**：文本文件压缩率可达 50-90%
2. **网络传输更快**：克隆仓库时传输的数据更少
3. **Git 标准格式**：所有 Git 实现都使用 zlib

### 哈希冲突怎么办？

SHA-1 生成 160 位哈希，冲突概率极低：
- 需要创建 2^80 个对象才有 50% 概率冲突
- 实际上几乎不可能发生

Git 也在逐步迁移到 SHA-256（256 位）。

## 练习题

### 练习 1：计算哈希
不运行程序，手动计算以下内容的 blob 哈希：
```
hello
```

提示：
1. 构建对象：`blob 6\0hello\n`（注意换行符）
2. 计算 SHA-1

可以用这个命令验证：
```bash
echo "hello" | git hash-object --stdin
```

### 练习 2：查看压缩效果
修改 `main.go`，添加代码显示：
- 原始内容大小
- 压缩后文件大小
- 压缩率

### 练习 3：处理二进制文件
尝试用 `WriteBlob()` 存储图片文件：
```go
data, _ := os.ReadFile("image.png")
hash, _ := WriteBlob(gitDir, data)
```

验证 Git 能否正确读取：
```bash
git cat-file blob <hash> > output.png
```

### 练习 4：实现 hash-object 命令
编写一个小工具，模拟 `git hash-object`：
```bash
# 只计算哈希
./geegit hash-object file.txt

# 计算并写入
./geegit hash-object -w file.txt
```

## 常见问题

### Q1: 为什么文件权限是 0444（只读）？
**A**: Git 对象是**不可变的**。一旦创建，内容永不改变。只读权限防止意外修改。

### Q2: 如果文件已存在怎么办？
**A**: 我们的实现会覆盖。真实 Git 会先检查文件是否存在，如果存在则跳过写入（内容相同，哈希必然相同）。

### Q3: WriteFile 失败怎么办？
**A**: 可能原因：
- 磁盘空间不足
- 权限问题
- 文件系统只读

生产代码应该添加重试和错误恢复逻辑。

### Q4: 能否并发写入？
**A**: 可以。不同哈希的对象存储在不同文件中，天然支持并发。相同内容的并发写入也是安全的（幂等操作）。

## 与 Day 2 的联系

Day 2 我们实现了 `ReadBlob()`，Day 3 实现了 `WriteBlob()`：

```
Day 2: .git/objects/xx/yy... → ReadBlob() → Blob 结构体
Day 3: Blob 数据 → WriteBlob() → .git/objects/xx/yy...
```

现在你可以：
1. 创建 blob 对象
2. 读取 blob 对象
3. 验证读写的正确性

## 下一步

在 Day 4，我们将学习：
- 如何读取 **tree** 对象（目录结构）
- Tree 对象的二进制格式
- 如何解析文件模式和名称

Tree 对象是 Git 的核心，它连接了 blob（文件内容）和 commit（版本历史）。

---

**恭喜你完成 Day 3！** 🎉

你现在已经掌握了 Git 对象存储的核心原理。继续前进！
