# 第一天：项目初始化

## 学习目标

在这个阶段，你将学习如何：
1. 搭建 grep 工具的基本项目结构
2. 实现命令行参数解析
3. 处理基本的输入输出
4. 理解 grep 的核心工作流程

## grep 是什么？

`grep` (Global Regular Expression Print) 是 Unix/Linux 系统中最常用的文本搜索工具之一。它可以在文件或标准输入中搜索匹配指定模式的行，并将结果输出。

基本用法示例：
```bash
# 在文件中搜索包含 "hello" 的行
grep "hello" file.txt

# 在标准输入中搜索
echo "hello world" | grep "hello"

# 递归搜索目录
grep -r "pattern" /path/to/dir
```

## 命令行参数结构

grep 的基本命令行格式为：
```
grep [OPTIONS] PATTERN [FILE...]
```

- `PATTERN`: 要搜索的模式（字符串或正则表达式）
- `FILE`: 一个或多个要搜索的文件（可选，不提供则从标准输入读取）
- `OPTIONS`: 可选的命令行选项（如 `-i`, `-n`, `-v` 等）

## 本课实现内容

在第一天，我们将实现：

### 1. 项目结构

```
GeeGrep/
├── beginner/
│   └── day1-project-init/
│       ├── main.go          # 主程序入口
│       ├── args.go          # 命令行参数解析
│       ├── matcher.go       # 匹配逻辑（基础版）
│       └── README.md        # 本课说明文档
```

### 2. 命令行参数解析

实现一个简单的命令行参数解析器，能够：
- 解析 PATTERN（搜索模式）
- 解析 FILE 列表
- 为后续课程预留选项解析的接口

### 3. 基本输入输出处理

实现两种输入模式：
- **从文件读取**：逐行读取文件内容
- **从标准输入读取**：从 stdin 读取内容（通过管道）

### 4. 简单字符串匹配

在这一课中，我们先实现最简单的字面字符串匹配（不涉及正则表达式）：
- 检查每一行是否包含指定的搜索模式
- 如果包含，则输出该行

## 工作流程

```
┌─────────────────┐
│  解析命令行参数  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  打开输入源     │
│ (文件/stdin)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  逐行读取输入   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  字符串匹配     │
│  (包含检查)     │
└────────┬────────┘
         │
         ▼
    匹配？
    ┌─┴─┐
   否│   │是
     │   ▼
     │ ┌─────────────┐
     │ │  输出该行   │
     │ └─────────────┘
     │
     ▼
  下一行/结束
```

## 测试用例

### 测试 1：从文件读取
```bash
# 创建测试文件
echo -e "hello world\nfoo bar\nhello again" > test.txt

# 运行程序
go run . "hello" test.txt

# 期望输出：
# hello world
# hello again
```

### 测试 2：从标准输入读取
```bash
echo -e "apple\nbanana\napricot" | go run . "ap"

# 期望输出：
# apple
# apricot
```

### 测试 3：无匹配
```bash
echo -e "foo\nbar\nbaz" | go run . "xyz"

# 期望输出：(无输出)
# 退出码：1
```

### 测试 4：多文件搜索
```bash
# 创建多个测试文件
echo "first file match" > file1.txt
echo "second file match" > file2.txt

go run . "match" file1.txt file2.txt

# 期望输出：
# first file match
# second file match
```

## 关键知识点

### 1. Go 命令行参数

Go 程序的命令行参数通过 `os.Args` 获取：
```go
// os.Args[0] 是程序名
// os.Args[1:] 是实际参数
args := os.Args[1:]
```

### 2. 文件 vs 标准输入

判断是从文件还是标准输入读取：
```go
if len(files) == 0 {
    // 从标准输入读取
    reader := bufio.NewReader(os.Stdin)
} else {
    // 从文件读取
    file, err := os.Open(filename)
    reader := bufio.NewReader(file)
}
```

### 3. 逐行读取

使用 `bufio.Scanner` 高效地逐行读取：
```go
scanner := bufio.NewScanner(reader)
for scanner.Scan() {
    line := scanner.Text()
    // 处理每一行
}
```

### 4. 字符串包含检查

使用 `strings.Contains` 进行简单匹配：
```go
if strings.Contains(line, pattern) {
    fmt.Println(line)
}
```

## 代码结构说明

### main.go
- 程序入口
- 调用参数解析
- 协调整体流程
- 处理错误和退出码

### args.go
- 定义 `Args` 结构体存储解析结果
- 实现 `ParseArgs()` 函数解析命令行参数
- 验证参数有效性

### matcher.go
- 定义 `Matcher` 接口（为后续扩展做准备）
- 实现 `LiteralMatcher`（字面字符串匹配器）
- 实现 `Search()` 函数执行搜索逻辑

## 扩展思考

1. **错误处理**：如何优雅地处理文件不存在、权限不足等错误？
2. **性能优化**：大文件如何高效处理？需要缓冲吗？
3. **标准兼容**：GNU grep 的退出码规范是什么？
   - 0: 找到匹配
   - 1: 未找到匹配
   - 2: 发生错误

## 下一步

完成本课后，你应该能够：
- ✓ 理解 grep 的基本工作原理
- ✓ 解析简单的命令行参数
- ✓ 从文件或标准输入读取数据
- ✓ 实现基本的字符串匹配

**下一课预告**：第二天我们将学习如何实现简单的正则表达式匹配，支持 `.`（任意字符）和字面字符。

## 参考资料

- [GNU Grep 官方文档](https://www.gnu.org/software/grep/manual/grep.html)
- [Go bufio 包文档](https://pkg.go.dev/bufio)
- [Go flag 包文档](https://pkg.go.dev/flag)
