# Day 4: 读取 Tree 对象

在这个阶段，你将实现 `git ls-tree` 命令，用于检查 tree 对象的内容。

## 一、Tree 对象

在这个阶段，我们将处理下一个 Git 对象类型：[trees（树对象）](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects#_tree_objects)。

**Tree 对象用于存储目录结构。**

### Tree 对象的条目

一个 tree 对象包含多个"条目"（entries）。每个条目包括：

#### 1. SHA-1 哈希值
指向一个 blob 或 tree 对象：
- 如果条目是**文件**，则指向一个 blob 对象
- 如果条目是**目录**，则指向一个 tree 对象

#### 2. 文件/目录名称
条目对应的文件或目录的名称

#### 3. 模式（Mode）
文件/目录的模式，这是 Unix 文件系统权限的简化版本：

**文件的有效值：**
- `100644` - 普通文件
- `100755` - 可执行文件
- `120000` - 符号链接

**目录的值：**
- `40000` - 目录

> 注意：还有其他值用于子模块，但在本课程中不涉及。

### Tree 对象示例

假设你有如下的目录结构：

```
your_repo/
  - file1
  - dir1/
    - file_in_dir_1
    - file_in_dir_2
  - dir2/
    - file_in_dir_3
```

tree 对象中的条目会是这样：

```
40000 dir1 <tree_sha_1>
40000 dir2 <tree_sha_2>
100644 file1 <blob_sha_1>
```

**条目解释：**
- **第 1 行** (`40000 dir1 <tree_sha_1>`) - 表示 `dir1` 是一个目录，SHA 哈希为 `<tree_sha_1>`
- **第 2 行** (`40000 dir2 <tree_sha_2>`) - 表示 `dir2` 是一个目录，SHA 哈希为 `<tree_sha_2>`
- **第 3 行** (`100644 file1 <blob_sha_1>`) - 表示 `file1` 是一个普通文件，SHA 哈希为 `<blob_sha_1>`

> `dir1` 和 `dir2` 本身也是 tree 对象，它们的条目会包含其内部的文件/目录。

## 二、git ls-tree 命令

`git ls-tree` 命令用于检查 tree 对象的内容。

### 完整输出格式

对于上面的目录结构，`git ls-tree` 的输出如下：

```bash
$ git ls-tree <tree_sha>
040000 tree <tree_sha_1>    dir1
040000 tree <tree_sha_2>    dir2
100644 blob <blob_sha_1>    file1
```

> 注意：输出是按字母顺序排序的，这是 Git 内部存储 tree 对象条目的方式。

### --name-only 标志

在本阶段，你将实现带 `--name-only` 标志的 `git ls-tree` 命令。使用该标志时，输出只显示名称：

```bash
$ git ls-tree --name-only <tree_sha>
dir1
dir2
file1
```

**为什么使用 --name-only？**
- 测试器使用 `--name-only` 是因为这种输出格式更容易测试

**建议：**
我们建议你也实现完整的 `ls-tree` 输出，因为这需要你解析 tree 对象中的所有数据，而不仅仅是文件名。

## 三、Tree 对象存储

与 blob 对象一样，tree 对象也存储在 `.git/objects` 目录中。

### 存储路径

如果 tree 对象的哈希是 `e88f7a929cd70b0274c4ea33b209c97fa845fdbc`，则对象路径为：

```
.git/objects/e8/8f7a929cd70b0274c4ea33b209c97fa845fdbc
```

### 存储格式

tree 对象文件的格式如下（Zlib 解压后）：

```
tree <size>\0<mode> <name>\0<20_byte_sha><mode> <name>\0<20_byte_sha>...
```

> 注意：上面为了可读性使用了换行，但实际文件中没有换行符。

#### 格式说明

**1. 对象头部：**
```
tree <size>\0
```
- 文件以 `tree <size>\0` 开头
- 这是"对象头部"，类似于 blob 对象的头部
- `<size>` 是树对象内容的大小（字节）

**2. 条目格式：**
```
<mode> <name>\0<20_byte_sha>
```

每个条目的组成：
- `<mode>` - 文件/目录的模式（参见前面章节的有效值）
- `<name>` - 文件/目录的名称
- `\0` - 空字节
- `<20_byte_sha>` - 20 字节的 SHA-1 哈希值（**不是**十六进制格式，是二进制格式）

> 深入了解：你可以在[这里](https://stackoverflow.com/questions/14790681/what-is-the-internal-format-of-a-git-tree-object)阅读更多关于 tree 对象内部格式的信息。

### 重要提示

**SHA-1 格式的区别：**

| 位置 | 格式 | 长度 | 示例 |
|------|------|------|------|
| Tree 对象文件内 | 二进制（20 字节） | 20 字节 | `\x3b\x18\xe5\x12...` |
| 命令行输出 | 十六进制字符串 | 40 字符 | `3b18e512dba79e4c...` |

在解析时需要将 20 字节的二进制 SHA 转换为 40 字符的十六进制字符串。

## 四、实现步骤

要实现 `git ls-tree --name-only` 命令，你需要：

1. **读取 tree 对象文件**
   - 根据 SHA-1 哈希定位 `.git/objects` 中的文件

2. **解压缩**
   - 使用 Zlib 解压文件内容

3. **解析对象头部**
   - 读取 `tree <size>\0` 头部
   - 验证对象类型是 tree

4. **解析条目**
   - 读取 `<mode>`（以空格结束）
   - 读取 `<name>`（以 `\0` 结束）
   - 读取 `<20_byte_sha>`（固定 20 字节）
   - 重复直到读取所有条目

5. **处理输出**
   - 如果是 `--name-only`：只输出名称
   - 如果是完整格式：输出 `<mode> <type> <sha> <name>`
   - 确保输出按字母顺序排序

## 五、测试步骤

### 初始化仓库

测试器将使用你的程序初始化一个新的仓库：

```bash
$ mkdir test_dir && cd test_dir
$ /path/to/your_program.sh init
```

### 写入 tree 对象

测试器将向 `.git/objects` 目录写入一个 tree 对象。

### 运行你的程序

测试器将像这样运行你的程序：

```bash
$ /path/to/your_program.sh ls-tree --name-only <tree_sha>
```

### 验证输出

对于如下的目录结构：

```
your_repo/
  - file1
  - dir1/
    - file_in_dir_1
    - file_in_dir_2
  - dir2/
    - file_in_dir_3
```

预期的输出是：

```
dir1
dir2
file1
```

测试器将验证你的程序输出是否与 tree 对象的内容匹配。

## 六、学习要点

- ✅ 理解 Tree 对象的作用（存储目录结构）
- ✅ 掌握 Tree 对象的条目结构
- ✅ 了解文件模式（mode）的含义
- ✅ 学会解析 Tree 对象的二进制格式
- ✅ 理解 SHA-1 的二进制和十六进制表示
- ✅ 实现 `git ls-tree --name-only` 命令
- ✅ 理解 Git 如何存储目录层次结构
