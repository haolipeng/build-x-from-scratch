# 第七天：克隆仓库

## 目标

在这个阶段，你将实现从 GitHub 克隆公共仓库的功能。

这是挑战的最后一个阶段，也可能是所有 CodeCrafters 中最难的一个！

我们未来可能会将其拆分为包含多个阶段的扩展，但目前它只是一个大阶段。

我们没有为这个阶段提供详细的说明，所以你需要自己探索。以下是一些帮助你入门的提示：

- [这篇论坛帖子](https://forum.codecrafters.io/t/step-for-git-clone-implementing-the-git-protocol/4407) 提供了一些关于如何逐步实现的建议。
- 你需要使用 Git 的 [Smart HTTP 传输协议](https://www.git-scm.com/docs/http-protocol)。
- 要了解更多关于协议格式的信息，我们推荐阅读：
  - [gitprotocol-pack](https://git-scm.com/docs/gitprotocol-pack)
  - [gitformat-pack](https://git-scm.com/docs/gitformat-pack)
  - [Unpacking Git packfiles](https://codewords.recurse.com/issues/three/unpacking-git-packfiles)
  - [Sneaky git number encoding](https://medium.com/@concertdaw/sneaky-git-number-encoding-ddcc5db5329f)

## 测试

测试器将这样运行你的程序：

```bash
$ /path/to/your_program.sh clone https://github.com/blah/blah <some_dir>
```

你的程序必须创建 `<some_dir>` 目录并将给定的仓库克隆到其中。

为了验证你的更改，测试器将：

- 检查随机文件的内容
- 从 `.git` 目录读取提交对象的属性