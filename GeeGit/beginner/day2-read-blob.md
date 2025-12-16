# Day 2: 读取 Blob 对象

在这个阶段，你将添加使用 `git cat-file` 命令读取 blob 对象的支持。

## 一、Git 对象

在这个课程中，我们将处理三种 Git 对象：

### 1. Blobs（本阶段）

- 用于存储文件数据
- Blob 只存储文件的内容，不包含文件名或权限信息

### 2. Trees（后续阶段）

- 用于存储目录结构
- 存储的信息包括树中有哪些文件/目录、它们的名称和权限等

### 3. Commits（后续阶段）

- 用于存储提交数据
- 存储的信息包括提交消息、作者、提交者、父提交等

---

**重要概念：** 所有 Git 对象都由一个 40 字符的 SHA-1 哈希值标识，也称为"对象哈希"。

**示例：** `e88f7a929cd70b0274c4ea33b209c97fa845fdbc`



## 二、Git 对象存储

Git 对象存储在 `.git/objects` 目录中。对象的路径是根据其哈希值派生的。

### 存储路径规则

哈希值为 `e88f7a929cd70b0274c4ea33b209c97fa845fdbc` 的对象，其路径为：

```
.git/objects/e8/8f7a929cd70b0274c4ea33b209c97fa845fdbc
```

**说明：**
- 文件并不是直接放在 `.git/objects` 目录中
- 而是放在一个以对象哈希**前两个字符**命名的目录中
- 剩余的 **38 个字符**用作文件名

每个 Git 对象都有自己的存储格式。我们将在本阶段学习 Blob 的存储方式，并在后续阶段介绍其他对象。



## 三、Blob 对象存储

每个 Git Blob 在 `.git/objects` 目录中作为单独的文件存储。该文件包含一个头部和 blob 对象的内容，使用 **Zlib 压缩**。

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



## 四、cat-file 命令

在这个阶段，你将通过从 `.git/objects` 目录读取内容来读取 git 仓库中的 blob。

你将使用我们在这个课程中遇到的第一个"**底层**"命令：`git cat-file`。

### 命令说明

`git cat-file` 用于查看对象的类型、大小和内容。

**使用示例：**

```bash
$ git cat-file -p <blob_sha>
hello world  # 这是 blob 的内容
```

### 实现步骤

要实现这个功能，你需要：

1. 从 `.git/objects` 目录读取 blob 对象文件的内容
2. 使用 Zlib 解压缩内容
3. 从解压后的数据中提取实际的"内容"
4. 将内容打印到标准输出



## 五、测试步骤

测试器将首先使用你的程序初始化一个新的 git 仓库，然后向 `.git/objects` 目录插入一个包含随机内容的 blob：

```bash
$ mkdir /tmp/test_dir && cd /tmp/test_dir
$ /path/to/your_program.sh init
$ echo "hello world" > test.txt  # 测试器将使用随机字符串，而不是 "hello world"
$ git hash-object -w test.txt
3b18e512dba79e4c8300dd08aeb37f8e728b8dad
```

之后，它将像这样运行你的程序：

```bash
$ /path/to/your_program.sh cat-file -p 3b18e512dba79e4c8300dd08aeb37f8e728b8dad
hello world
```

测试器将验证你的程序输出是否与 blob 的内容匹配。

## 学习要点

- ✅ 理解 Git 对象的三种类型
- ✅ 掌握 Git 对象的存储路径规则
- ✅ 了解 Blob 对象的存储格式
- ✅ 学会使用 Zlib 解压缩
- ✅ 实现 `git cat-file -p` 命令
