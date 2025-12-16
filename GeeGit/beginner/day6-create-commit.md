# 第六天：创建 Commit 对象

在这个阶段，你将实现 `git commit-tree` 命令，用于创建一个 commit 对象。

## Commit 对象

在本挑战中我们要处理的最后一个 git 对象是 [commit 对象](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects#_git_commit_objects)。

一个 commit 对象文件在 Zlib 压缩之前的格式如下：

```
commit <size>\0tree <tree_sha>
parent <parent_sha>
author <name> <<email>> <timestamp> <timezone>
committer <name> <<email>> <timestamp> <timezone>

<commit message>
```

格式说明：
- 文件以 `commit <size>\0` 开头（与 blob 和 tree 对象的头部格式相同）
- 在头部之后，内容是纯文本（非二进制）
- 每一行以换行符（`\n`）结尾
- tree 和 parent 的 SHA 值使用十六进制格式（40 个字符），而不像 tree 对象中那样使用 20 字节
- 元数据和提交消息之间有一个空行
- 时间戳格式为：`<自纪元以来的秒数> <时区偏移量>`（例如：`1234567890 +0000`）

下面是一个 commit 对象的示例：

```
commit 177\0tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904
parent 3b18e512dba79e4c8300dd08aeb37f8e728b8dad
author John Doe <john@example.com> 1234567890 +0000
committer John Doe <john@example.com> 1234567890 +0000

Initial commit
```

你可以在[这里](https://stackoverflow.com/questions/22968856/what-is-the-file-format-of-a-git-commit-object-data-structure)阅读更多关于 commit 对象格式的内容。

## `git commit-tree` 命令

`git commit-tree` 命令用于创建一个 commit 对象。

使用示例：

```bash
$ mkdir test_dir && cd test_dir
$ git init
Initialized empty Git repository in /path/to/test_dir/.git/

# 创建一个 tree，获取其 SHA
$ echo "hello world" > test.txt
$ git add test.txt
$ git write-tree
4b825dc642cb6eb9a060e54bf8d69288fbee4904

# 创建初始 commit
$ git commit-tree 4b825dc642cb6eb9a060e54bf8d69288fbee4904 -m "Initial commit"
3b18e512dba79e4c8300dd08aeb37f8e728b8dad

# 写入一些更改，获取另一个 tree SHA
$ echo "hello world 2" > test.txt
$ git add test.txt
$ git write-tree
5b825dc642cb6eb9a060e54bf8d69288fbee4904

# 使用新的 tree SHA 和 parent 创建一个新的 commit
$ git commit-tree 5b825dc642cb6eb9a060e54bf8d69288fbee4904 -p 3b18e512dba79e4c8300dd08aeb37f8e728b8dad -m "Second commit"
6c18e512dba79e4c8300dd08aeb37f8e728b8dad
```

`git commit-tree` 的输出是写入到 `.git/objects` 目录的 commit 对象的 40 字符 SHA-1 哈希值。

## 测试要求

你的程序将这样被调用：

```bash
$ ./your_program.sh commit-tree <tree_sha> -p <commit_sha> -m <message>
```

你的程序必须创建一个 commit 对象并将其 40 字符的 SHA-1 哈希值打印到标准输出。

为了简化实现：

- 你将收到恰好一个父 commit
- 你将收到恰好一行提交消息
- 你可以为 author/committer 字段硬编码任何有效的姓名和邮箱
- 你可以使用任何有效的时间戳（例如：当前时间或硬编码的值）

测试器将通过使用 `git show` 命令从 `.git` 目录读取 commit 对象来验证你的更改。