# GeeGrep - 从零实现 grep

## 项目简介

GeeGrep 是一个从零开始实现 grep 命令行工具的教程项目。通过循序渐进的方式，你将学习到：

- 命令行工具开发
- 正则表达式引擎实现
- 文件系统操作
- 文本处理算法
- 性能优化技巧

## 学习路线

### 第一部分：基础准备（3课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 1 | [项目初始化](beginner/day1-project-init.md) | ✅ 已完成 |
| Day 2 | [正则表达式基础](beginner/day2-regex-basics.md) | ✅ 已完成 |
| Day 3 | [基础命令行选项](beginner/day3-command-options.md) | ✅ 已完成 |

### 第二部分：正则表达式基础（8课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 4 | [字符类基础](beginner/day4-character-classes.md) (`\d`, `\w`, `\s`) | ✅ 已完成 |
| Day 5 | [字符组](beginner/day5-character-groups.md) (`[abc]`, `[^abc]`) | ✅ 已完成 |
| Day 6 | [组合字符类](beginner/day6-combined-classes.md) | ✅ 已完成 |
| Day 7 | 锚点 (`^`, `$`, `\b`) | 🚧 规划中 |
| Day 8 | 量词（一）(`.`, `?`, `+`, `*`) | 🚧 规划中 |
| Day 9 | 量词（二）(`{n}`, `{n,}`, `{n,m}`) | 🚧 规划中 |
| Day 10 | 贪婪与非贪婪匹配 | 🚧 规划中 |
| Day 11 | 分支与分组 (`|`, `()`) | 🚧 规划中 |

### 第三部分：匹配输出控制（5课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 12 | 匹配行处理 | 🚧 规划中 |
| Day 13 | 匹配内容提取 (`-o`) | 🚧 规划中 |
| Day 14 | 上下文显示 (`-A`, `-B`, `-C`) | 🚧 规划中 |
| Day 15 | 高亮显示（一） | 🚧 规划中 |
| Day 16 | 高亮显示（二）(`--color`) | 🚧 规划中 |

### 第四部分：文件处理（4课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 17 | 单文件搜索 | 🚧 规划中 |
| Day 18 | 多文件搜索 | 🚧 规划中 |
| Day 19 | 递归搜索 (`-r`, `-R`) | 🚧 规划中 |
| Day 20 | 文件类型过滤 | 🚧 规划中 |

### 第五部分：高级正则表达式（4课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 21 | 捕获组 | 🚧 规划中 |
| Day 22 | 反向引用（一） | 🚧 规划中 |
| Day 23 | 反向引用（二） | 🚧 规划中 |
| Day 24 | 断言 (`(?=)`, `(?!)`, `(?<=)`, `(?<!)`) | 🚧 规划中 |

### 第六部分：性能与优化（3课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 25 | 正则表达式优化 | 🚧 规划中 |
| Day 26 | 缓冲与流处理 | 🚧 规划中 |
| Day 27 | 并发处理 | 🚧 规划中 |

### 第七部分：完善与测试（3课）

| 课程 | 主题 | 状态 |
|------|------|------|
| Day 28 | 边界情况处理 | 🚧 规划中 |
| Day 29 | 兼容性与标准 | 🚧 规划中 |
| Day 30 | 综合测试与发布 | 🚧 规划中 |

## 快速开始

### 环境要求

- Go 1.19 或更高版本
- Linux/macOS/Windows
- 基本的命令行使用经验

### 开始学习

1. 克隆本仓库
```bash
git clone https://github.com/yourusername/build-x-from-scratch
cd build-x-from-scratch/GeeGrep
```

2. 从第一课开始
```bash
cd beginner/day1-project-init
go run . --help
```

3. 按照每一课的 README 和课程文档学习

## 项目结构

```
GeeGrep/
├── README.md              # 本文件
├── beginner/              # 基础课程
│   ├── day1-project-init/
│   │   ├── main.go
│   │   ├── args.go
│   │   ├── matcher.go
│   │   └── README.md
│   ├── day1-project-init.md
│   └── ...
└── advanced/              # 高级课程（待添加）
```

## 学习建议

1. **循序渐进**：按照顺序完成每一课，不要跳过
2. **动手实践**：不要只看代码，要自己动手写
3. **理解原理**：不要只是复制代码，要理解为什么这样做
4. **多做测试**：每完成一个功能，写测试用例验证
5. **参考标准**：对比 GNU grep 的行为，确保实现正确

## 相关资源

### 官方文档
- [GNU Grep 官方文档](https://www.gnu.org/software/grep/manual/grep.html)
- [POSIX 正则表达式标准](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap09.html)

### 学习资料
- [Regular-Expressions.info](https://www.regular-expressions.info/)
- [Regex101 - 正则表达式测试工具](https://regex101.com/)
- [Go 标准库 regexp 包](https://pkg.go.dev/regexp)

### 相关项目
- [ripgrep](https://github.com/BurntSushi/ripgrep) - 用 Rust 实现的高性能 grep
- [ag (The Silver Searcher)](https://github.com/ggreer/the_silver_searcher) - 速度极快的代码搜索工具

## 贡献

欢迎提交 Issue 和 Pull Request！

如果你发现了错误或有改进建议，请：
1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的修改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 常见问题

### Q: 需要多长时间完成全部课程？
A: 根据你的经验水平，预计需要 2-4 周时间。每天 1-2 小时的学习时间是比较理想的。

### Q: 我是 Go 新手，能学习这个课程吗？
A: 可以，但建议先学习 Go 的基础知识。推荐先完成 [A Tour of Go](https://go.dev/tour/)。

### Q: 和真正的 grep 比，这个实现有什么区别？
A: 教学版本注重原理讲解和代码清晰度，真正的 grep 经过了几十年的优化，在性能和功能上更加完善。

### Q: 完成后能达到什么水平？
A: 你将能够：
- 理解正则表达式的底层实现原理
- 开发命令行工具
- 处理复杂的文本搜索需求
- 理解性能优化的基本方法

## 致谢

本项目灵感来源于：
- [CodeCrafters - Build your own grep](https://app.codecrafters.io/courses/grep)
- [Build Your Own Text Editor](https://viewsourcecode.org/snaptoken/kilo/)
- [Crafting Interpreters](https://craftinginterpreters.com/)

## 许可证

MIT License - 详见 [LICENSE](../LICENSE) 文件

## 联系方式

如有问题或建议，欢迎通过以下方式联系：
- 提交 Issue
- 发送 Pull Request
- Email: your-email@example.com

---

开始你的 grep 实现之旅吧！🚀
