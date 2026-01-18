# Day 1: 项目初始化 - 示例代码

## 项目结构

```
day1-project-init/
├── main.go       # 主程序入口
├── args.go       # 命令行参数解析
├── matcher.go    # 匹配逻辑
└── README.md     # 本文档
```

## 代码说明

### args.go
定义了命令行参数的数据结构和解析逻辑：
- `Args` 结构体：存储解析后的参数
- `ParseArgs()`：解析命令行参数
- `Validate()`：验证参数有效性

### matcher.go
实现了字符串匹配的核心逻辑：
- `Matcher` 接口：定义匹配器接口，为后续扩展做准备
- `LiteralMatcher`：实现字面字符串匹配
- `Search()`：在输入流中搜索匹配的行
- `SearchFiles()`：在多个文件中搜索

### main.go
程序入口，协调各模块：
1. 解析命令行参数
2. 创建匹配器
3. 执行搜索
4. 输出结果并设置退出码

## 运行方式

### 1. 从文件搜索

```bash
# 创建测试文件
echo -e "hello world\nfoo bar\nhello again" > test.txt

# 搜索包含 "hello" 的行
go run . "hello" test.txt
```

期望输出：
```
hello world
hello again
```

### 2. 从标准输入搜索

```bash
# 使用管道传入数据
echo -e "apple\nbanana\napricot" | go run . "ap"
```

期望输出：
```
apple
apricot
```

### 3. 多文件搜索

```bash
# 创建多个测试文件
echo "first file match" > file1.txt
echo "second file match" > file2.txt
echo "no pattern here" > file3.txt

# 在多个文件中搜索
go run . "match" file1.txt file2.txt file3.txt
```

期望输出：
```
first file match
second file match
```

### 4. 无匹配情况

```bash
# 搜索不存在的模式
echo -e "foo\nbar\nbaz" | go run . "xyz"

# 检查退出码
echo $?  # 应该输出 1
```

## 退出码说明

- `0`：找到匹配
- `1`：未找到匹配
- `2`：发生错误（参数错误、文件无法打开等）

## 测试脚本

创建 `test.sh` 文件来自动化测试：

```bash
#!/bin/bash

echo "=== Day 1 Tests ==="
echo

# 测试 1: 基本匹配
echo "Test 1: Basic matching from stdin"
result=$(echo -e "hello\nworld\nhello again" | go run . "hello")
expected="hello
hello again"
if [ "$result" = "$expected" ]; then
    echo "✓ PASSED"
else
    echo "✗ FAILED"
    echo "Expected: $expected"
    echo "Got: $result"
fi
echo

# 测试 2: 从文件读取
echo "Test 2: Matching from file"
echo -e "apple pie\nbanana\napricot" > /tmp/test.txt
result=$(go run . "ap" /tmp/test.txt)
expected="apple pie
apricot"
if [ "$result" = "$expected" ]; then
    echo "✓ PASSED"
else
    echo "✗ FAILED"
fi
rm /tmp/test.txt
echo

# 测试 3: 无匹配
echo "Test 3: No match"
echo "foo" | go run . "bar" > /dev/null
exitcode=$?
if [ $exitcode -eq 1 ]; then
    echo "✓ PASSED (exit code 1)"
else
    echo "✗ FAILED (expected exit code 1, got $exitcode)"
fi
echo

# 测试 4: 多文件搜索
echo "Test 4: Multiple files"
echo "match1" > /tmp/file1.txt
echo "match2" > /tmp/file2.txt
result=$(go run . "match" /tmp/file1.txt /tmp/file2.txt)
expected="match1
match2"
if [ "$result" = "$expected" ]; then
    echo "✓ PASSED"
else
    echo "✗ FAILED"
fi
rm /tmp/file1.txt /tmp/file2.txt
echo

echo "=== All tests completed ==="
```

运行测试：
```bash
chmod +x test.sh
./test.sh
```

## 学习要点

### 1. 命令行参数解析
```go
// 获取程序参数（不包括程序名）
args := os.Args[1:]

// 第一个参数是模式，其余是文件列表
pattern := args[0]
files := args[1:]
```

### 2. 从不同输入源读取
```go
var reader io.Reader

if len(files) == 0 {
    // 从标准输入读取
    reader = os.Stdin
} else {
    // 从文件读取
    file, _ := os.Open(filename)
    reader = file
}
```

### 3. 逐行处理
```go
scanner := bufio.NewScanner(reader)
for scanner.Scan() {
    line := scanner.Text()
    // 处理每一行
}
```

### 4. 字符串包含检查
```go
if strings.Contains(line, pattern) {
    fmt.Println(line)
}
```

### 5. 退出码
```go
os.Exit(0)  // 成功
os.Exit(1)  // 未找到匹配
os.Exit(2)  // 错误
```

## 下一步

完成本课后，你应该掌握了：
- ✓ 基本的 Go 项目结构
- ✓ 命令行参数解析
- ✓ 文件和标准输入的处理
- ✓ 简单的字符串匹配
- ✓ 退出码的正确使用

**下一课**：我们将实现正则表达式的基础功能，支持通配符 `.` 和字面字符匹配。

## 常见问题

### Q: 为什么要定义 Matcher 接口？
A: 虽然现在只有字面字符串匹配，但后续会添加正则表达式匹配。使用接口可以轻松扩展新的匹配器类型。

### Q: SearchResult 的作用是什么？
A: 它存储搜索结果的元信息（是否有匹配、匹配数量等），虽然现在只用于判断退出码，但后续会用于 `-c` 选项（统计匹配数）。

### Q: 为什么使用 bufio.Scanner 而不是直接读取？
A: Scanner 提供了方便的逐行读取功能，并且性能更好。它内部使用缓冲，适合处理大文件。

### Q: 错误处理为什么用 defer file.Close()？
A: defer 确保函数返回时文件会被关闭，即使发生错误也能正确释放资源。这是 Go 的最佳实践。
