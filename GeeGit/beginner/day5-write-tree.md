# Day 5: 写入 Tree 对象

在这个阶段，你将实现将 tree 对象写入 `.git/objects` 目录的功能。

## 一、Tree 对象存储（回顾）

作为回顾，tree 对象用于存储目录结构，并保存在 `.git/objects` 目录中。

### 存储路径

如果 tree 对象的哈希是 `e88f7a929cd70b0274c4ea33b209c97fa845fdbc`，则对象路径为：

```
.git/objects/e8/8f7a929cd70b0274c4ea33b209c97fa845fdbc
```

### 存储格式

tree 对象文件的格式如下（Zlib 压缩前）：

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
- `<mode>` - 文件/目录的类型和权限
- `<name>` - 文件/目录的名称
- `\0` - 空字节
- `<20_byte_sha>` - 20 字节的 SHA-1 哈希值（二进制格式）



## 二、mode 字段详解

`<mode>` 字段表示每个条目的类型和权限。以下是一些有效值：

### 常用模式值

| 模式 | 类型 | 说明 |
|------|------|------|
| `100644` | 普通文件 | 不可执行的常规文件 |
| `100755` | 可执行文件 | 具有执行权限的文件 |
| `40000` | 目录 | Tree 对象 |

### 重要提示：目录模式

**存储格式 vs 显示格式：**

| 位置 | 格式 | 说明 |
|------|------|------|
| Tree 对象文件内 | `40000` | 实际存储的模式 |
| `git ls-tree` 输出 | `040000` | 为了可读性显示的格式 |

> **注意：** 虽然 Git 命令（如 `git ls-tree`）显示目录模式为 `040000` 以提高可读性，但在 tree 对象中实际存储的模式是 `40000`。

> 深入了解：你可以在[这里](https://stackoverflow.com/questions/14790681/what-is-the-internal-format-of-a-git-tree-object)阅读更多关于 tree 对象内部格式的信息。



## 三、git write-tree 命令

`git write-tree` 命令根据"暂存区"（staging area）的当前状态创建一个 tree 对象。暂存区是运行 `git add` 时更改存放的地方。

### 命令说明

**在本课程中的简化：**
- 你**不需要**实现暂存区
- 假设工作目录中的所有文件都已经暂存

### 命令使用示例

```bash
# 创建一个包含内容的文件
$ echo "hello world" > test.txt

# 将文件添加到暂存区（本课程中不实现暂存区）
$ git add test.txt

# 将 tree 写入 .git/objects
$ git write-tree
4b825dc642cb6eb9a060e54bf8d69288fbee4904
```

**输出：**
- `git write-tree` 的输出是写入 `.git/objects` 的 tree 对象的 40 字符 SHA-1 哈希值

### 与本地 git 测试

如果你在本地使用 `git` 测试，确保在运行 `git write-tree` 之前先运行 `git add .`，以便暂存工作目录中的所有文件。



## 四、实现步骤

要实现 `git write-tree` 命令，你需要：

### 步骤 1: 遍历工作目录
- 读取当前工作目录中的所有文件和子目录
- 跳过 `.git` 目录

### 步骤 2: 处理文件条目
对于每个文件：
1. 读取文件内容
2. 创建 blob 对象（使用 Day 3 实现的功能）
3. 记录其 SHA-1 哈希值
4. 确定文件模式（`100644` 或 `100755`）

### 步骤 3: 递归处理目录
对于每个子目录：
1. 递归调用 `write-tree` 处理子目录
2. 创建子目录的 tree 对象
3. 记录其 SHA-1 哈希值
4. 模式固定为 `40000`

### 步骤 4: 按名称排序
- 将所有条目**按名称字母顺序**排序
- 这是 Git 的标准行为，必须遵守

### 步骤 5: 构建 tree 对象
1. 创建对象头部：`tree <size>\0`
2. 按顺序添加所有条目：`<mode> <name>\0<20_byte_sha>`
3. 计算完整内容的 SHA-1 哈希值

### 步骤 6: 写入对象文件
1. 使用 Zlib 压缩 tree 对象内容
2. 根据 SHA-1 哈希确定存储路径
3. 创建必要的子目录（哈希前 2 个字符）
4. 将压缩内容写入文件
5. 将 40 字符的 SHA-1 哈希值打印到标准输出

### 关键细节

**排序规则：**
```go
// 按字母顺序排序条目名称
sort.Strings(entryNames)
```

**SHA-1 转换：**
```
十六进制字符串 (40 字符) → 二进制字节 (20 字节)
"3b18e512dba79e4c..." → []byte{0x3b, 0x18, 0xe5, ...}
```



## 五、测试步骤

### 初始化仓库

测试器将使用你的程序初始化一个新的 Git 仓库：

```bash
$ mkdir test_dir && cd test_dir
$ /path/to/your_program.sh init
```

### 创建测试文件和目录

测试器将创建一些随机文件和目录：

```bash
$ echo "hello world" > test_file_1.txt
$ mkdir test_dir_1
$ echo "hello world" > test_dir_1/test_file_2.txt
$ mkdir test_dir_2
$ echo "hello world" > test_dir_2/test_file_3.txt
```

目录结构：
```
test_dir/
├── test_file_1.txt
├── test_dir_1/
│   └── test_file_2.txt
└── test_dir_2/
    └── test_file_3.txt
```

### 运行你的程序

测试器将像这样运行你的程序：

```bash
$ /path/to/your_program.sh write-tree
4b825dc642cb6eb9a060e54bf8d69288fbee4904
```

### 验证要求

- ✅ 你的程序应该将整个工作目录作为 tree 对象写入
- ✅ 向标准输出打印 40 字符的 SHA-1 哈希值
- ✅ 写入的 tree 对象内容与官方 `git` 实现完全匹配
