# Day 4: 字符类基础 - 示例代码

## 项目结构

```
day4-character-classes/
├── main.go       # 主程序入口
├── args.go       # 命令行参数解析
├── regex.go      # 正则表达式解析器（支持字符类和转义）
├── engine.go     # 匹配引擎（支持字符类匹配）
├── output.go     # 输出格式化器
├── test.sh       # 测试脚本
├── go.mod        # Go 模块文件
└── README.md     # 本文档
```

## 新增功能

### 1. 预定义字符类

Day 4 添加了对常用字符类的支持：

| 字符类 | 含义 | 等价于 | 示例 |
|--------|------|--------|------|
| `\d` | 数字 | `[0-9]` | 匹配 `0`, `5`, `9` |
| `\w` | 单词字符 | `[a-zA-Z0-9_]` | 匹配 `a`, `Z`, `5`, `_` |
| `\s` | 空白字符 | 空格、制表符、换行 | 匹配 ` `, `\t`, `\n` |
| `\D` | 非数字 | `[^0-9]` | 匹配 `a`, `!`, ` ` |
| `\W` | 非单词字符 | `[^a-zA-Z0-9_]` | 匹配 `!`, `@`, ` ` |
| `\S` | 非空白 | 非空白字符 | 匹配 `a`, `1`, `!` |

### 2. 字面转义

支持转义特殊字符以匹配字面值：

```bash
\.  # 匹配字面点号 "."
\\  # 匹配字面反斜杠 "\"
\*  # 匹配字面星号 "*"
\+  # 匹配字面加号 "+"
\?  # 匹配字面问号 "?"
```

## 使用示例

### 基础用法

```bash
# \d: 匹配数字
echo -e "abc\n123\na1b" | go run . "\\d"
# 输出:
# 123
# a1b

# \w: 匹配单词字符
echo -e "hello\nhello_world\nhello-world" | go run . "hello\\wworld"
# 输出:
# hello_world

# \s: 匹配空白
echo -e "hello world\nhelloworld" | go run . "hello\\sworld"
# 输出:
# hello world

# \.: 匹配字面点号
echo -e "3.14\n3a14" | go run . "3\\.14"
# 输出:
# 3.14
```

### 实际应用场景

#### 1. 匹配电话号码

```bash
echo "Call: 555-1234" | go run . "\\d\\d\\d-\\d\\d\\d\\d"
# 输出: Call: 555-1234
```

#### 2. 匹配变量名

```bash
echo "var_name_123" | go run . "\\w\\w\\w_\\w\\w\\w\\w_\\d\\d\\d"
# 输出: var_name_123
```

#### 3. 匹配日期格式

```bash
echo "Date: 12/31/2024" | go run . "\\d\\d/\\d\\d/\\d\\d\\d\\d"
# 输出: Date: 12/31/2024
```

#### 4. 匹配 IP 地址

```bash
echo "IP: 192.168.1.1" | go run . "\\d\\d\\d\\.\\d\\d\\d\\.\\d\\.\\d"
# 输出: IP: 192.168.1.1
```

#### 5. 匹配文件扩展名

```bash
echo -e "file.txt\nfile.doc\nfiletxt" | go run . "file\\.\\w\\w\\w"
# 输出:
# file.txt
# file.doc
```

### 否定字符类

```bash
# \D: 匹配非数字
echo -e "123\nabc\n1a2" | go run . "\\D"
# 输出:
# abc
# 1a2

# \W: 匹配非单词字符
echo -e "hello\nhello@world\nabc_123" | go run . "\\W"
# 输出:
# hello@world

# \S: 匹配非空白
echo -e "   \nhello\n\t" | go run . "\\S"
# 输出:
# hello
```

### 与命令行选项组合

```bash
# \d + -n: 显示包含数字的行号
echo -e "no digits\nline 123\nmore text" | go run . -n "\\d"
# 输出:
# 2:line 123

# \d + -v: 显示不包含数字的行
echo -e "abc\n123\nxyz" | go run . -v "\\d"
# 输出:
# abc
# xyz

# \d + -c: 统计包含数字的行数
echo -e "a1\nb2\nc3\nde" | go run . -c "\\d"
# 输出:
# 3

# \w + -i: 忽略大小写匹配单词字符
echo -e "ABC\n123\nxyz" | go run . -i "\\w\\w\\w"
# 输出:
# ABC
# 123
# xyz
```

## 运行测试

```bash
# 运行完整测试套件
bash test.sh

# 手动测试
echo -e "abc\n123\na1b" | go run . "\\d"
```

## 实现细节

### 1. Token 类型扩展

新增了多种字符类 token：

```go
const (
    TokenLiteral  TokenType = iota
    TokenWildcard
    TokenDigit      // \d
    TokenWord       // \w
    TokenSpace      // \s
    TokenNotDigit   // \D
    TokenNotWord    // \W
    TokenNotSpace   // \S
)
```

### 2. 转义序列解析

```go
func parseEscape(pattern string, pos int) (Token, int, error) {
    next := pattern[pos+1]

    switch next {
    case 'd':
        return Token{Type: TokenDigit}, pos + 2, nil
    case 'w':
        return Token{Type: TokenWord}, pos + 2, nil
    case 's':
        return Token{Type: TokenSpace}, pos + 2, nil
    case '.', '*', '+', '?', '\\':
        // 字面转义
        return Token{Type: TokenLiteral, Value: next}, pos + 2, nil
    // ...
    }
}
```

### 3. 字符类匹配

```go
func (m *RegexMatcher) matchChar(token Token, char byte) bool {
    switch token.Type {
    case TokenDigit:
        return isDigit(char)
    case TokenWord:
        return isWordChar(char)
    case TokenSpace:
        return isSpace(char)
    // ...
    }
}
```

### 4. 辅助函数

```go
func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

func isWordChar(c byte) bool {
    return (c >= 'a' && c <= 'z') ||
           (c >= 'A' && c <= 'Z') ||
           (c >= '0' && c <= '9') ||
           (c == '_')
}

func isSpace(c byte) bool {
    return c == ' ' || c == '\t' ||
           c == '\n' || c == '\r'
}
```

## 字符类定义

### \d - 数字

```
匹配: 0 1 2 3 4 5 6 7 8 9
不匹配: a-z A-Z 特殊字符
```

### \w - 单词字符

```
匹配: a-z A-Z 0-9 _
不匹配: 空格 @ ! # 等特殊字符
```

### \s - 空白字符

```
匹配: 空格 ' ', 制表符 '\t', 换行 '\n', 回车 '\r'
不匹配: 字母 数字 特殊字符
```

## Shell 中的转义

在 shell 中使用字符类时，需要注意转义：

```bash
# 方法 1: 双引号 + 双反斜杠
grep "\\d" file.txt

# 方法 2: 单引号 + 单反斜杠
grep '\d' file.txt

# 错误示例
grep "\d" file.txt  # 某些 shell 可能将 \d 解析为 d
```

## 与 Day 3 的区别

| 特性 | Day 3 | Day 4 |
|------|-------|-------|
| 字符类 | 不支持 | 支持 \d, \w, \s 等 |
| 转义字符 | 不支持 | 支持 \., \\, \*, 等 |
| 匹配能力 | 基础 | 增强 |
| 实际应用 | 有限 | 广泛 |

### 示例对比

```bash
# Day 3: 无法直接匹配数字
echo "line 123" | day3/go run . "123"  # 需要知道确切数字

# Day 4: 可以匹配任意数字
echo "line 123" | day4/go run . "\\d\\d\\d"  # 匹配任意3位数字

# Day 3: 无法匹配字面点号
echo "3.14" | day3/go run . "3.14"  # . 被当作通配符

# Day 4: 可以精确匹配
echo "3.14" | day4/go run . "3\\.14"  # \. 匹配字面点号
```

## 调试技巧

### 1. 查看 Token 解析

修改 `regex.go` 添加调试输出：

```go
tokens := ParsePattern("\\d\\.\\d")
for _, token := range tokens {
    fmt.Println(token)
}
// 输出:
// Digit(\d)
// Literal('.')
// Digit(\d)
```

### 2. 测试单个字符类

```bash
# 测试 \d
echo "1" | go run . "\\d"  # 应该匹配
echo "a" | go run . "\\d"  # 不应该匹配

# 测试 \w
echo "_" | go run . "\\w"  # 应该匹配
echo "@" | go run . "\\w"  # 不应该匹配

# 测试 \s
echo " " | go run . "\\s"  # 应该匹配
echo "a" | go run . "\\s"  # 不应该匹配
```

## 限制和注意事项

1. **ASCII 限制**: 当前只支持 ASCII 字符，不支持 Unicode
2. **空白字符**: 只支持常见的空白字符（空格、制表符、换行、回车）
3. **单字节**: 假设所有字符都是单字节
4. **Shell 转义**: 在命令行使用时需要正确转义反斜杠

这些限制将在后续课程中改进。

## 性能考虑

### 字符类匹配的性能

```go
// 直接比较（快）
func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

// 可选优化：查表法
var digitTable [256]bool
func init() {
    for c := '0'; c <= '9'; c++ {
        digitTable[c] = true
    }
}
func isDigit(c byte) bool {
    return digitTable[c]
}
```

## 下一步

完成本课后，你应该掌握了：
- ✓ 字符类的概念和用途
- ✓ 转义字符的解析和处理
- ✓ `\d`, `\w`, `\s` 及其否定形式
- ✓ 字面转义的实现
- ✓ 字符类在实际场景中的应用

**下一课预告**: Day 5 将实现字符组 `[abc]` 和否定字符组 `[^abc]`，支持：
- 字符枚举：`[aeiou]`
- 字符范围：`[a-z]`, `[0-9]`
- 否定字符组：`[^0-9]`
- 组合使用：`[a-zA-Z0-9_]`

## 常见问题

### Q: 为什么 `\d` 不支持 Unicode 数字？
A: 为了简化实现，当前只支持 ASCII。后续课程会扩展 Unicode 支持。

### Q: 如何匹配字面的反斜杠？
A: 使用 `\\\\`（shell 中需要4个反斜杠，或在单引号中用 `'\\\\'`）。

### Q: `\w` 为什么包含下划线？
A: 这是正则表达式的标准定义，来源于编程语言中标识符的定义。

### Q: 能否添加自定义字符类？
A: Day 5 会实现字符组 `[...]`，可以自定义字符集合。
