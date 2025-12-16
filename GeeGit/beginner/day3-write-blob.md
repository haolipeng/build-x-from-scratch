# Day 3: 写入 Blob 对象

在这个阶段，你将实现使用 [`git hash-object`](https://git-scm.com/docs/git-hash-object) 命令创建 blob 对象的支持。

## 一、git hash-object 命令

`git hash-object` 用于计算 Git 对象的 SHA-1 哈希值。当使用 `-w` 标志时，它还会将对象写入 `.git/objects` 目录。

### 命令使用示例

```bash
# 创建一个包含内容的文件
$ echo -n "hello world" > test.txt

# 计算文件的 SHA-1 哈希值并写入 .git/objects
$ git hash-object -w test.txt
95d09f2b10159347eece71399a7e2e907ea3df4f

# 验证文件已写入 .git/objects
$ file .git/objects/95/d09f2b10159347eece71399a7e2e907ea3df4f
.git/objects/95/d09f2b10159347eece71399a7e2e907ea3df4f: zlib compressed data
```

**关键参数：**
- `-w` 标志：将对象写入对象数据库（`.git/objects` 目录）

## 二、Blob 对象存储（回顾）

如前一阶段所述，每个 Git Blob 在 `.git/objects` 目录中作为单独的文件存储。该文件包含一个头部和 blob 对象的内容，使用 **Zlib 压缩**。

### 存储格式

blob 对象文件的格式如下（Zlib 解压后）：

```
blob <size>\0<content>
```

**格式说明：**
- `<size>` - 内容的大小（以字节为单位）
- `\0` - 空字节（null byte）
- `<content>` - 文件的实际内容

### 示例

如果文件的内容是 `hello world`，那么 blob 对象文件的内容如下（Zlib 解压后）：

```
blob 11\0hello world
```

> 注意：`hello world` 的长度是 11 字节

## 三、实现步骤

要实现 `git hash-object -w` 命令，你需要：

1. **读取文件内容**
   - 从指定的文件路径读取文件数据

2. **构建 blob 对象**
   - 按照格式 `blob <size>\0<content>` 构建对象内容

3. **计算 SHA-1 哈希值**
   - 对构建的对象内容计算 SHA-1 哈希

4. **压缩对象**
   - 使用 Zlib 压缩对象内容

5. **写入对象文件**
   - 根据哈希值确定存储路径
   - 创建必要的子目录（哈希前 2 个字符）
   - 将压缩后的内容写入文件

6. **输出哈希值**
   - 将 40 字符的 SHA-1 哈希值打印到标准输出

## 四、测试步骤

### 初始化仓库

测试器将首先使用你的程序初始化一个新的 git 仓库：

```bash
$ mkdir test_dir && cd test_dir
$ /path/to/your_program.sh init
```

### 创建测试文件

向文件中写入一些随机数据：

```bash
$ echo "hello world" > test.txt
```

### 运行你的程序

测试器将像这样运行你的程序：

```bash
$ ./your_program.sh hash-object -w test.txt
3b18e512dba79e4c8300dd08aeb37f8e728b8dad
```

### 验证要求

测试器将验证：

- ✅ 你的程序向标准输出打印一个 40 字符的 SHA-1 哈希值
- ✅ 写入 `.git/objects` 的文件内容与官方 `git` 实现写入的内容完全匹配

## 五、学习要点

- ✅ 理解 `git hash-object -w` 命令的作用
- ✅ 掌握 blob 对象的构建过程
- ✅ 学会计算 SHA-1 哈希值
- ✅ 实现 Zlib 压缩
- ✅ 理解 Git 对象的存储路径规则
- ✅ 完成从读取文件到写入对象数据库的完整流程

## 六、与 Day 2 的对比

| 操作 | Day 2 | Day 3 |
|------|-------|-------|
| 命令 | `git cat-file -p` | `git hash-object -w` |
| 方向 | 读取（Read） | 写入（Write） |
| 主要步骤 | 解压 → 提取内容 → 输出 | 读取 → 构建 → 压缩 → 写入 |
| 输入 | 对象哈希值 | 文件路径 |
| 输出 | 文件内容 | 对象哈希值 |

现在你已经掌握了 Git 对象的读写操作，这是理解 Git 内部工作原理的重要基础！
