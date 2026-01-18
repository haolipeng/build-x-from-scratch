# 第三天：基础命令行选项

## 学习目标

在这个阶段，你将学习如何：
1. 实现命令行选项解析
2. 支持常用的 grep 选项
3. 理解选项如何影响匹配和输出
4. 掌握 flag 包的使用

## grep 常用选项

grep 提供了丰富的命令行选项来控制搜索行为和输出格式。以下是最常用的几个：

### 基础选项

```bash
-n, --line-number       # 显示行号
-i, --ignore-case       # 忽略大小写
-v, --invert-match      # 反向匹配（显示不匹配的行）
-c, --count             # 只显示匹配的行数
```

### 使用示例

```bash
# 显示行号
grep -n "pattern" file.txt
# 输出: 3:matching line

# 忽略大小写
grep -i "HELLO" file.txt
# 匹配: hello, Hello, HELLO, HeLLo

# 反向匹配
grep -v "error" log.txt
# 显示所有不包含 "error" 的行

# 统计匹配数
grep -c "TODO" src/*.go
# 输出: 42
```

### 组合使用

```bash
# 忽略大小写 + 显示行号
grep -in "error" app.log

# 统计不匹配的行数
grep -cv "success" results.txt
```

## 本课实现内容

### 1. 命令行选项结构

我们将扩展 `Args` 结构体来支持选项：

```go
type Args struct {
    Pattern     string   // 搜索模式
    Files       []string // 文件列表

    // 选项
    LineNumber  bool     // -n: 显示行号
    IgnoreCase  bool     // -i: 忽��大小写
    InvertMatch bool     // -v: 反向匹配
    Count       bool     // -c: 统计数量
}
```

### 2. 选项解析

使用 Go 的 `flag` 包来解析命令行选项：

```go
import "flag"

func ParseArgs(args []string) (*Args, error) {
    fs := flag.NewFlagSet("grep", flag.ContinueOnError)

    // 定义选项
    lineNumber := fs.Bool("n", false, "显示行号")
    ignoreCase := fs.Bool("i", false, "忽略大小写")
    invertMatch := fs.Bool("v", false, "反向匹配")
    count := fs.Bool("c", false, "统计数量")

    // 解析
    fs.Parse(args)

    // ...
}
```

### 3. 选项功能实现

#### -n (显示行号)
```
输入:
  1: hello world
  2: foo bar
  3: hello again

输出 (grep -n "hello"):
  1:hello world
  3:hello again
```

#### -i (忽略大小写)
```
输入: Hello, HELLO, hello

模式: "hello"
  不使用 -i: 只匹配 "hello"
  使用 -i:   匹配所有三个
```

#### -v (反向匹配)
```
输入:
  error: file not found
  success: operation completed
  error: timeout

输出 (grep -v "error"):
  success: operation completed
```

#### -c (统计数量)
```
输入:
  TODO: implement feature A
  FIXME: bug in module B
  TODO: add tests

输出 (grep -c "TODO"):
  2
```

## 实现架构

```
┌─────────────────────┐
│  解析命令行参数      │
│  (选项 + 模式 + 文件)│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  编译正则表达式      │
│  (考虑 -i 选项)     │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  逐行读取输入       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  执行匹配           │
│  (考虑 -v 选项)     │
└──────────┬──────────┘
           │
           ▼
    ┌──────┴──────┐
    │  -c 选项？   │
    └──────┬──────┘
           │
    ┌──────┴──────┐
    │             │
   是            否
    │             │
    ▼             ▼
 统计数量      ┌─────────────┐
    │         │ -n 选项？    │
    │         └──────┬───────┘
    │                │
    │         ┌──────┴──────┐
    │         │             │
    │        是            否
    │         │             │
    ▼         ▼             ▼
  输出数量  行号:内容      内容
```

## 代码组织

### options.go (新增)
定义选项相关的结构和方法：
```go
type Options struct {
    LineNumber  bool
    IgnoreCase  bool
    InvertMatch bool
    Count       bool
}

func (o *Options) ShouldShowLine(matched bool) bool {
    if o.InvertMatch {
        return !matched
    }
    return matched
}
```

### args.go (修改)
扩展参数解析以支持选项：
- 使用 `flag` 包解析选项
- 验证选项组合的有效性

### matcher.go (修改)
支持忽略大小写匹配：
```go
func (m *RegexMatcher) SetIgnoreCase(ignore bool) {
    m.ignoreCase = ignore
}
```

### output.go (新增)
处理不同的输出格式：
- `PrintLine()`: 普通输出
- `PrintLineWithNumber()`: 带行号输出
- `PrintCount()`: 输出统计数量

## 测试用例

### 测试 1: -n 显示行号
```bash
echo -e "apple\nbanana\napricot" | go run . -n "ap"

# 输出:
# 1:apple
# 3:apricot
```

### 测试 2: -i 忽略大小写
```bash
echo -e "Hello\nWORLD\nhello" | go run . -i "hello"

# 输出:
# Hello
# hello
```

### 测试 3: -v 反向匹配
```bash
echo -e "error\nsuccess\nerror" | go run . -v "error"

# 输出:
# success
```

### 测试 4: -c 统计数量
```bash
echo -e "cat\ndog\ncat\nbird" | go run . -c "cat"

# 输出:
# 2
```

### 测试 5: 组合选项
```bash
# -i + -n: 忽略大小写并显示行号
echo -e "Hello\nworld\nHELLO" | go run . -in "hello"

# 输出:
# 1:Hello
# 3:HELLO

# -c + -v: 统计不匹配的行数
echo -e "a\nb\nc\na" | go run . -cv "a"

# 输出:
# 2
```

### 测试 6: 多文件 + 选项
```bash
echo "hello" > file1.txt
echo "world" > file2.txt
echo "hello" > file3.txt

go run . -c "hello" file1.txt file2.txt file3.txt

# 输出:
# file1.txt:1
# file2.txt:0
# file3.txt:1
```

## 关键知识点

### 1. flag 包的使用

```go
// 创建 FlagSet
fs := flag.NewFlagSet("grep", flag.ContinueOnError)

// 定义布尔选项
lineNum := fs.Bool("n", false, "显示行号")
lineNum2 := fs.Bool("line-number", false, "显示行号")

// 解析参数
err := fs.Parse(args)

// 获取非选项参数（模式和文件）
remaining := fs.Args()
```

### 2. 大小写转换

```go
import "strings"

// 转为小写进行匹配
func toLowerCase(s string) string {
    return strings.ToLower(s)
}

// 在匹配前转换
if ignoreCase {
    pattern = toLowerCase(pattern)
    text = toLowerCase(text)
}
```

### 3. 行号跟踪

```go
lineNumber := 0
for scanner.Scan() {
    lineNumber++
    line := scanner.Text()

    if matched {
        if showLineNumber {
            fmt.Printf("%d:%s\n", lineNumber, line)
        } else {
            fmt.Println(line)
        }
    }
}
```

### 4. 计数器

```go
matchCount := 0
for scanner.Scan() {
    if matched {
        matchCount++
    }
}

if countOnly {
    fmt.Println(matchCount)
}
```

## 选项优先级和冲突

某些选项组合需要特殊处理：

### 1. -c 优先级最高
当使用 `-c` 时，忽略 `-n`（因为只输出数字）：
```bash
grep -cn "pattern" file.txt
# 输出: 5  (不显示行号，只显示数量)
```

### 2. -v 改变匹配逻辑
```go
func shouldPrint(matched bool, invertMatch bool) bool {
    if invertMatch {
        return !matched
    }
    return matched
}
```

### 3. 多文件时的行为
使用 `-c` 时，每个文件单独统计：
```
file1.txt:3
file2.txt:0
file3.txt:7
```

## GNU grep 兼容性

我们的实现遵循 GNU grep 的行为：

| 选项 | 短选项 | 长选项 | 功能 |
|------|--------|--------|------|
| 行号 | `-n` | `--line-number` | 显示行号 |
| 忽略大小写 | `-i` | `--ignore-case` | 忽略大小写 |
| 反向匹配 | `-v` | `--invert-match` | 显示不匹配的行 |
| 统计 | `-c` | `--count` | 只显示匹配数量 |

## 性能考虑

### 1. 忽略大小写的性能影响

```go
// 方法 1: 每次匹配时转换（慢）
if ignoreCase {
    if strings.ToLower(text) contains strings.ToLower(pattern) {
        // ...
    }
}

// 方法 2: 预处理模式（快）
compiledPattern := compilePattern(pattern, ignoreCase)
// 只在匹配时转换文本
```

### 2. 计数模式优化

使用 `-c` 时不需要存储行内容：
```go
if countOnly {
    // 只增加计数器，不存储行
    count++
} else {
    // 需要存储行以便输出
    matchedLines = append(matchedLines, line)
}
```

## 常见错误和调试

### 1. 选项解析错误
```bash
# 错误：选项必须在模式之前
grep "pattern" -n file.txt  # ✗

# 正确
grep -n "pattern" file.txt  # ✓
```

### 2. 引号问题
```bash
# 模式包含空格时必须使用引号
grep -n hello world  # ✗ (world 被当作文件名)
grep -n "hello world"  # ✓
```

### 3. 选项组合
```bash
# 多个短选项可以组合
grep -in "pattern"  # 等同于 grep -i -n "pattern"
grep -niv "pattern" # 等同于 grep -n -i -v "pattern"
```

## 扩展思考

1. **更多选项**：如何添加 `-o`（只显示匹配部分）、`-A`（显示匹配后的行）？
2. **长选项**：如何同时支持 `-n` 和 `--line-number`？
3. **选项验证**：哪些选项组合是无效的？
4. **配置文件**：如何支持从配置文件读取默认选项？

## 下一步

完成本课后，你应该能够：
- ✓ 使用 flag 包解析命令行选项
- ✓ 实现常用的 grep 选项
- ✓ 处理选项组合和冲突
- ✓ 理解选项如何影响输出格式

**下一课预告**：第四天我们将实现字符类，支持 `\d`（数字）、`\w`（单词字符）、`\s`（空白字符）等正则表达式特性。

## 参考资料

- [GNU Grep Manual - Command-line Options](https://www.gnu.org/software/grep/manual/grep.html#Command_002dline-Options)
- [Go flag 包文档](https://pkg.go.dev/flag)
- [POSIX grep 标准](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/grep.html)
