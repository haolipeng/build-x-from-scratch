# Day 3: 基础命令行选项 - 示例代码

## 项目结构

```
day3-command-options/
├── main.go       # 主程序入口
├── args.go       # 命令行参数和选项解析
├── regex.go      # 正则表达式解析器（与 Day 2 相同）
├── engine.go     # 匹配引擎（增加忽略大小写支持）
├── output.go     # 输出格式化器
├── test.sh       # 测试脚本
├── go.mod        # Go 模块文件
└── README.md     # 本文档
```

## 新增功能

### 1. 命令行选项支持

Day 3 添加了四个常用的 grep 选项：

| 选项 | 长选项 | 功能 | 示例 |
|------|--------|------|------|
| `-n` | `--line-number` | 显示行号 | `grep -n "pattern" file.txt` |
| `-i` | `--ignore-case` | 忽略大小写 | `grep -i "HELLO" file.txt` |
| `-v` | `--invert-match` | 反向匹配 | `grep -v "error" log.txt` |
| `-c` | `--count` | 统计数量 | `grep -c "TODO" *.go` |

### 2. 核心模块

#### args.go - 选项解析
- 使用 Go 的 `flag` 包解析命令行选项
- `Options` 结构体存储所有选项
- 支持选项组合

#### engine.go - 忽略大小写匹配
- 添加 `SetIgnoreCase()` 方法
- 在匹配时转换文本为小写
- 优化性能：只转换必要的字符

#### output.go - 输出格式化
- `OutputFormatter` 统一处理输出格式
- 根据选项决定输出格式
- 支持行号、统计等多种模式

## 使用示例

### 基础用法

```bash
# -n: 显示行号
echo -e "apple\nbanana\napricot" | go run . -n "ap"
# 输出:
# 1:apple
# 3:apricot

# -i: 忽略大小写
echo -e "Hello\nWORLD\nhello" | go run . -i "hello"
# 输出:
# Hello
# hello

# -v: 反向匹配（显示不匹配的行）
echo -e "error\nsuccess\nerror" | go run . -v "error"
# 输出:
# success

# -c: 统计匹配数量
echo -e "cat\ndog\ncat\nbird" | go run . -c "cat"
# 输出:
# 2
```

### 选项组合

```bash
# -i + -n: 忽略大小写并显示行号
echo -e "Hello\nworld\nHELLO" | go run . -i -n "hello"
# 输出:
# 1:Hello
# 3:HELLO

# -c + -v: 统计不匹配的行数
echo -e "a\nb\nc\na" | go run . -c -v "a"
# 输出:
# 2

# -v + -n: 反向匹配并显示行号
echo -e "error\ninfo\nwarning" | go run . -v -n "error"
# 输出:
# 2:info
# 3:warning
```

### 正则表达式 + 选项

```bash
# 通配符 + 忽略大小写
echo -e "Cat\ncot\nCUT" | go run . -i "c.t"
# 输出:
# Cat
# cot
# CUT

# 通配符 + 行号
echo -e "hello\nhallo\nhxllo" | go run . -n "h.llo"
# 输出:
# 1:hello
# 2:hallo
# 3:hxllo
```

### 文件操作

```bash
# 从文件搜索
echo -e "line 1: hello\nline 2: world\nline 3: hello" > test.txt
go run . -n "hello" test.txt
# 输出:
# 1:line 1: hello
# 3:line 3: hello

# 多文件统计
echo "hello" > file1.txt
echo "world" > file2.txt
go run . -c "hello" file1.txt file2.txt
# 输出:
# file1.txt:1
# file2.txt:0
```

## 运行测试

```bash
# 运行完整测试套件
bash test.sh

# 手动测试
echo -e "apple\nbanana\napricot" | go run . -n "ap"
```

## 实现细节

### 1. 选项解析

使用 Go 的 `flag` 包：

```go
fs := flag.NewFlagSet("grep", flag.ContinueOnError)

// 定义选项
opts := &Options{}
fs.BoolVar(&opts.LineNumber, "n", false, "显示行号")
fs.BoolVar(&opts.IgnoreCase, "i", false, "忽略大小写")
fs.BoolVar(&opts.InvertMatch, "v", false, "反向匹配")
fs.BoolVar(&opts.Count, "c", false, "统计数量")

// 解析
fs.Parse(args)
```

### 2. 忽略大小写实现

```go
func (m *RegexMatcher) Match(text string) bool {
    searchText := text
    if m.ignoreCase {
        searchText = strings.ToLower(text)
    }
    // ... 匹配逻辑
}

func (m *RegexMatcher) matchTokens(...) bool {
    // 字面字符匹配时也要考虑大小写
    if m.ignoreCase {
        patternChar = toLowerByte(patternChar)
    }
    // ...
}
```

### 3. 反向匹配逻辑

```go
func (f *OutputFormatter) ShouldPrint(matched bool) bool {
    if f.opts.InvertMatch {
        return !matched  // 反向：不匹配的才打印
    }
    return matched       // 正常：匹配的才打印
}
```

### 4. 输出格式控制

```go
func (f *OutputFormatter) PrintLine(result *MatchResult) {
    if !f.ShouldPrint(result.Matched) {
        return
    }

    // -c 选项：不打印行，只统计
    if f.opts.Count {
        return
    }

    // -n 选项：带行号
    if f.opts.LineNumber {
        fmt.Printf("%d:%s\n", result.LineNumber, result.Line)
    } else {
        fmt.Println(result.Line)
    }
}
```

## 选项行为说明

### -n (行号)
- 格式：`行号:内容`
- 行号从 1 开始
- 与 `-c` 冲突时，`-c` 优先

### -i (忽略大小写)
- 匹配时不区分大小写
- 输出时保持原始大小写
- 对正则表达式也有效

### -v (反向匹配)
- 显示**不**匹配模式的行
- 可与其他选项组合
- 影响退出码：有不匹配的行则退出 0

### -c (统计)
- 只输出匹配的行数
- 多文件时，每个文件单独统计：`filename:count`
- 覆盖 `-n` 的效果

## 与 Day 2 的区别

| 特性 | Day 2 | Day 3 |
|------|-------|-------|
| 命令行选项 | 不支持 | 支持 -n, -i, -v, -c |
| 大小写敏感 | 是 | 可选（-i） |
| 输出格式 | 固定 | 可配置 |
| 统计功能 | 无 | 有（-c） |
| 反向匹配 | 无 | 有（-v） |

### 示例对比

```bash
# Day 2: 只能精确匹配大小写
echo "Hello" | day2/go run . "hello"  # 无输出

# Day 3: 可以忽略大小写
echo "Hello" | day3/go run . -i "hello"  # 输出: Hello
```

## 退出码

| 情况 | 退出码 | 说明 |
|------|--------|------|
| 找到匹配 | 0 | 成功 |
| 未找到匹配 | 1 | 正常，但无匹配 |
| 参数错误 | 2 | 命令行错误 |
| 文件错误 | 2 | 无法读取文件 |

**注意**: 使用 `-v` 时，"有不匹配的行" 视为成功（退出 0）。

## 限制和注意事项

1. **选项组合**：Go 的 flag 包不支持 `-in` 这样的组合，必须写成 `-i -n`
2. **选项位置**：选项必须在模式之前：`grep -n "pattern" file`（不能是 `grep "pattern" -n file`）
3. **长选项**：当前只支持短选项（`-n`），不支持长选项（`--line-number`）
4. **Unicode**：`-i` 选项只处理 ASCII 大小写，不支持 Unicode 字符

这些限制将在后续课程中改进。

## 调试技巧

### 查看解析的参数

修改 `main.go` 临时添加调试输出：

```go
args, err := ParseArgs(os.Args[1:])
fmt.Fprintf(os.Stderr, "DEBUG: %s\n", args.String())
```

### 测试单个选项

```bash
# 测试 -n
echo -e "a\nb\nc" | go run . -n "b"

# 测试 -i
echo "Hello" | go run . -i "hello"

# 测试 -v
echo -e "a\nb" | go run . -v "a"

# 测试 -c
echo -e "a\na\nb" | go run . -c "a"
```

## 性能考虑

### 忽略大小写的开销

```go
// 优化：只转换一次模式
if ignoreCase {
    // 模式在编译时已转换
    // 只需在匹配时转换文本
}
```

### 统计模式优化

```go
// -c 模式下不需要存储行内容
if opts.Count {
    count++  // 只计数
} else {
    lines = append(lines, line)  // 存储用于输出
}
```

## 下一步

完成本课后，你应该掌握了：
- ✓ 使用 flag 包解析命令行选项
- ✓ 实现常用 grep 选项的功能
- ✓ 处理选项组合和优先级
- ✓ 理解选项如何影响匹配和输出

**下一课预告**: Day 4 将实现字符类，支持：
- `\d`: 匹配数字 `[0-9]`
- `\w`: 匹配单词字符 `[a-zA-Z0-9_]`
- `\s`: 匹配空白字符
- 转义字符支持

## 常见问题

### Q: 为什么 `-in` 不能组合使用？
A: Go 的 flag 包不支持短选项组合。必须写成 `-i -n`。如果需要支持，需要自己实现选项解析器。

### Q: `-c` 和 `-n` 一起用会怎样？
A: `-c` 优先级更高，只输出数字，不显示行号。

### Q: `-v` 如何影响退出码？
A: 如果有不匹配的行（即有输出），退出码为 0。这符合 grep 的行为。

### Q: 为什么多文件时 `-c` 输出格式不同？
A: 为了区分不同文件的统计结果，输出格式为 `filename:count`。
