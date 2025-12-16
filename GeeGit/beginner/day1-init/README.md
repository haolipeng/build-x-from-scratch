# Day 1: Init the .git directory

## 学习目标

通过这一天的学习，你将：
- ✅ 理解 Git 仓库的基本目录结构
- ✅ 知道 `.git` 目录下各个子目录的作用
- ✅ 学会创建一个空的 Git 仓库
- ✅ 了解 HEAD 文件的作用

## 关键概念

### Git 仓库结构

当你运行 `git init` 时，Git 会创建以下目录结构：

```
.git/
├── HEAD                 # 指向当前分支
├── objects/             # 存储所有 Git 对象（blob, tree, commit）
│   ├── info/
│   └── pack/
└── refs/                # 存储引用（分支、标签）
    ├── heads/           # 本地分支
    └── tags/            # 标签
```

### HEAD 文件

`HEAD` 文件指向当前所在的分支，内容通常是：
```
ref: refs/heads/main
```

这表示当前在 `main` 分支上。

### objects 目录

这是 Git 最核心的目录，存储所有的：
- **Blob**: 文件内容
- **Tree**: 目录结构
- **Commit**: 提交信息

对象以 SHA-1 哈希命名，例如：
```
objects/e6/9de29bb2d1d6434b8b29ae775ad8c2e48c5391
```

## 代码说明

### hash.go
定义了 `Hash` 类型和相关函数：
- `Hash`: 20字节的 SHA-1 哈希值
- `ComputeHash()`: 计算对象的 SHA-1 哈希

### object.go
定义了对象类型：
- `BlobObject`: 文件对象
- `TreeObject`: 目录对象
- `CommitObject`: 提交对象

### init.go
实现了仓库初始化函数：
- `InitRepository()`: 创建 .git 目录结构

## 运行演示

```bash
cd beginner/day1-init
go run *.go
```

预期输出：
```
=== Day 1: Init the .git directory ===

✓ Step 1: Creating test repository directory
  Created: ./test-repo

✓ Step 2: Initializing Git repository structure
  Initialized .git directory

✓ Step 3: Verifying directory structure

Result - Git directory structure:

  📁 .git/
  📁 .git/objects/
  📁 .git/refs/
  📁 .git/refs/heads/
  📄 .git/HEAD
     Content: ref: refs/heads/main

...
```

## 与真实 Git 对比

| 我们的实现 | Git 命令 | 说明 |
|-----------|---------|------|
| `InitRepository(path)` | `git init` | 初始化仓库 |
| 创建 `.git/objects/` | Git 自动创建 | 对象存储目录 |
| 创建 `.git/refs/heads/` | Git 自动创建 | 分支目录 |
| `HEAD` -> `refs/heads/main` | Git 默认行为 | 指向默认分支 |

### 验证我们的实现

```bash
# 运行我们的程序
go run *.go

# 用真实的 git 检查目录
cd test-repo
git status

# 应该看到：
# On branch main
# No commits yet
```

## 练习

1. **修改默认分支名**
   - 将 `init.go` 中的 `main` 改为 `master`
   - 运行后查看 `HEAD` 文件内容

2. **添加更多目录**
   - 尝试创建 `.git/info/` 目录
   - 创建 `.git/hooks/` 目录

3. **探索真实仓库**
   ```bash
   cd /path/to/your/git/repo
   tree .git
   cat .git/HEAD
   ```

## 常见问题

**Q: 为什么 objects 目录是空的？**
A: 因为我们还没有创建任何对象。在 Day 2-3 会学习如何创建对象。

**Q: refs/heads 下为什么没有文件？**
A: 因为还没有任何提交。第一次提交后会创建分支文件。

**Q: 可以用真实的 git 命令操作这个仓库吗？**
A: 可以！我们创建的结构和真实 Git 完全兼容。

## 下一步

在 **Day 2**，我们将学习：
- 如何读取 blob 对象
- zlib 压缩/解压
- Git 对象的存储格式

➡️ [前往 Day 2: Read a blob object](../day2-read-blob/)
