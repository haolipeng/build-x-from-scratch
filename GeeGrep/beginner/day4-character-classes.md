# 第四天：字符类基础

## 学习目标

在这个阶段，你将学习如何：
1. 理解正则表达式中的字符类概念
2. 实现转义字符解析
3. 支持预定义字符类 `\d`、`\w`、`\s`
4. 实现字面转义（如 `\.` 匹配字面点号）
5. 理解字符类的匹配逻辑

## 什么是字符类？

字符类（Character Classes）是正则表达式中用于匹配特定字符集合的简写方式。

### 预定义字符类

| 字符类 | 含义 | 等价于 | 示例 |
|--------|------|--------|------|
| `\d` | 数字 | `[0-9]` | 匹配 `0`, `5`, `9` |
| `\w` | 单词字符 | `[a-zA-Z0-9_]` | 匹配 `a`, `Z`, `5`, `_` |
| `\s` | 空白字符 | 空格、制表符、换行等 | 匹配 ` `, `\t`, `\n` |
| `\D` | 非数字 | `[^0-9]` | 匹配除数字外的字符 |
| `\W` | 非单词字符 | `[^a-zA-Z0-9_]` | 匹配 `!`, `@`, ` ` |
| `\S` | 非空白 | 非空白字符 | 匹配 `a`, `1`, `!` |

### 字面转义

某些字符在正则表达式中有特殊含义，需要转义才能匹配字面值：

```
\.  匹配字面点号 "."
\*  匹配字面星号 "*"
\+  匹配字面加号 "+"
\?  匹配字面问号 "?"
\\  匹配字面反斜杠 "\"
```

## 使用示例

### 示例 1: 匹配数字 `\d`

```bash
# 匹配包含数字的行
echo -e "abc\n123\na1b" | grep "\d"
# 输出:
# 123
# a1b

# 匹配电话号码模式
echo "Call: 555-1234" | grep "\d\d\d-\d\d\d\d"
# 输出: Call: 555-1234
```

### 示例 2: 匹配单词字符 `\w`

```bash
# 匹配标识符（3个单词字符）
echo -e "foo\nbar_123\n@#$" | grep "\w\w\w"
# 输出:
# foo
# bar_123

# 匹配变量名模式
echo "var_name_123" | grep "\w\w\w_\w\w\w\w"
# 输出: var_name_123
```

### 示例 3: 匹配空白 `\s`

```bash
# 匹配包含空格的行
echo -e "hello world\nhelloworld" | grep "hello\sworld"
# 输出: hello world

# 匹配缩进的代码
echo -e "\tfoo\nbar" | grep "\s"
# 输出:     foo
```

### 示例 4: 字面转义

```bash
# 匹配字面点号
echo -e "3.14\n3a14" | grep "3\.14"
# 输出: 3.14

# 匹配文件扩展名
echo -e "file.txt\nfileXtxt" | grep "file\.txt"
# 输出: file.txt
```

## 实现架构

### 1. Token 类型扩展

```go
const (
    TokenLiteral     TokenType = iota  // 字面字符
    TokenWildcard                      // . 通配符
    TokenDigit                         // \d 数字
    TokenWord                          // \w 单词字符
    TokenSpace                         // \s 空白字符
    TokenNotDigit                      // \D 非数字
    TokenNotWord                       // \W 非单词字符
    TokenNotSpace                      // \S 非空白
)
```

### 2. 转义字符解析

解析过程需要识别反斜杠 `\` 并处理后续字符：

```
输入: "a\d\s\."
       ↓
解析为 Token 序列:
[
  Token{Type: Literal, Value: 'a'},
  Token{Type: Digit},
  Token{Type: Space},
  Token{Type: Literal, Value: '.'}
]
```

### 3. 匹配逻辑

```go
func matchCharClass(token Token, char byte) bool {
    switch token.Type {
    case TokenDigit:
        return char >= '0' && char <= '9'

    case TokenWord:
        return (char >= 'a' && char <= 'z') ||
               (char >= 'A' && char <= 'Z') ||
               (char >= '0' && char <= '9') ||
               (char == '_')

    case TokenSpace:
        return char == ' ' || char == '\t' ||
               char == '\n' || char == '\r'

    // ...
    }
}
```

## 解析流程

```
┌─────────────────────┐
│  读取字符           │
└──────────┬──────────┘
           │
           ▼
    是反斜杠 \ ？
    ┌──────┴──────┐
    │             │
   是            否
    │             │
    ▼             ▼
┌─────────────┐  直接作为
│ 读取下一字符 │  字面字符
└──────┬──────┘     │
       │            │
       ▼            │
   是转义序列？      │
   ┌────┴────┐      │
   │         │      │
  是        否      │
   │         │      │
   ▼         ▼      │
字符类    字面转义   │
(\d,\w)   (\.,\\)   │
   │         │      │
   └─────┬───┴──────┘
         │
         ▼
   创建 Token
```

## 本课实现内容

### 1. 扩展 regex.go

添加新的 Token 类型和解析逻辑：

```go
// 解析转义序列
func parseEscape(pattern string, i int) (Token, int, error) {
    if i+1 >= len(pattern) {
        return Token{}, i, fmt.Errorf("incomplete escape sequence")
    }

    next := pattern[i+1]
    switch next {
    case 'd':
        return Token{Type: TokenDigit}, i+2, nil
    case 'w':
        return Token{Type: TokenWord}, i+2, nil
    case 's':
        return Token{Type: TokenSpace}, i+2, nil
    case '.', '*', '+', '?', '\\':
        // 字面转义
        return Token{Type: TokenLiteral, Value: next}, i+2, nil
    default:
        return Token{}, i, fmt.Errorf("unknown escape: \\%c", next)
    }
}
```

### 2. 扩展 engine.go

添加字符类匹配逻辑：

```go
func (m *RegexMatcher) matchCharClass(token Token, char byte) bool {
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

### 3. 辅助函数

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

## 测试用例

### 测试 1: \d 匹配数字

```bash
echo -e "abc\n123\na1b" | go run . "\\d"
# 输出:
# 123
# a1b

echo "ID: 12345" | go run . "ID: \\d\\d\\d\\d\\d"
# 输出: ID: 12345
```

### 测试 2: \w 匹配单词字符

```bash
echo -e "hello\nhello_world\nhello-world" | go run . "hello\\wworld"
# 输出:
# hello_world

echo "var_123" | go run . "\\w\\w\\w_\\d\\d\\d"
# 输出: var_123
```

### 测试 3: \s 匹配空白

```bash
echo -e "hello world\nhelloworld" | go run . "hello\\sworld"
# 输出: hello world

printf "tab\there\nnotab" | go run . "tab\\shere"
# 输出: tab    here
```

### 测试 4: 字面转义

```bash
echo -e "3.14\n3a14" | go run . "3\\.14"
# 输出: 3.14

echo "a.b.c" | go run . "a\\.b\\.c"
# 输出: a.b.c
```

### 测试 5: 混合使用

```bash
# 匹配日期格式: \d\d/\d\d/\d\d\d\d
echo "Date: 12/31/2024" | go run . "\\d\\d/\\d\\d/\\d\\d\\d\\d"
# 输出: Date: 12/31/2024

# 匹配邮箱格式: \w+@\w+\.\w+
# (注: + 量词在后续课程实现)
echo "user@example.com" | go run . "\\w\\w\\w\\w@\\w\\w\\w\\w\\w\\w\\w\\.\\w\\w\\w"
# 输出: user@example.com
```

### 测试 6: 与选项组合

```bash
# \d + -n: 显示包含数字的行号
echo -e "no digits\nline 123\nmore text" | go run . -n "\\d"
# 输出:
# 2:line 123

# \w + -i: 忽略大小写匹配单词字符
echo -e "ABC\n123\nxyz" | go run . -i "\\w\\w\\w"
# 输出:
# ABC
# 123
# xyz
```

## 关键知识点

### 1. 转义字符的双重含义

在命令行中，反斜杠本身需要转义：

```bash
# Shell 层面的转义
grep "\\d"   # Shell 解析为: \d，传给程序
grep "\d"    # 某些 shell 可能解析为: d

# 在程序中
pattern := "\\d"  // Go 字符串中: \d
```

### 2. 字符类的否定形式

大写字母表示否定：

```go
case TokenNotDigit:  // \D
    return !isDigit(char)

case TokenNotWord:   // \W
    return !isWordChar(char)

case TokenNotSpace:  // \S
    return !isSpace(char)
```

### 3. 空白字符的完整定义

```go
func isSpace(c byte) bool {
    switch c {
    case ' ':   // 空格
    case '\t':  // 制表符
    case '\n':  // 换行
    case '\r':  // 回车
    case '\f':  // 换页
    case '\v':  // 垂直制表符
        return true
    }
    return false
}
```

### 4. Unicode 支持

当前实现仅支持 ASCII 字符。对于 Unicode：

```go
// ASCII 版本（Day 4）
func isWordChar(c byte) bool {
    return (c >= 'a' && c <= 'z') || ...
}

// Unicode 版本（后续扩展）
import "unicode"

func isWordChar(r rune) bool {
    return unicode.IsLetter(r) ||
           unicode.IsDigit(r) ||
           r == '_'
}
```

## 常见错误和调试

### 1. 转义字符未正确解析

```bash
# 错误：忘记转义反斜杠
echo "123" | grep "\d"
# 可能不工作，取决于 shell

# 正确
echo "123" | grep "\\d"
```

### 2. 字符类拼写错误

```go
// 常见错误
case 'D':  // 应该是 'd'
    return Token{Type: TokenDigit}

// 正确
case 'd':
    return Token{Type: TokenDigit}
```

### 3. 匹配逻辑错误

```go
// 错误：\w 应该包含下划线
func isWordChar(c byte) bool {
    return (c >= 'a' && c <= 'z') ||
           (c >= 'A' && c <= 'Z') ||
           (c >= '0' && c <= '9')
    // 缺少: || (c == '_')
}
```

## 性能优化

### 1. 内联函数

```go
// 使用内联提高性能
//go:inline
func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}
```

### 2. 查表法

对于复杂的字符类，可以使用查表：

```go
var charClassTable [256]bool

func init() {
    // 预计算字符类
    for c := '0'; c <= '9'; c++ {
        charClassTable[c] = true
    }
    // ...
}

func isDigit(c byte) bool {
    return charClassTable[c]
}
```

## 与标准正则表达式的差异

| 特性 | 标准 Regex | 我们的实现 |
|------|-----------|------------|
| `\d` | Unicode 数字 | ASCII 0-9 |
| `\w` | Unicode 字母数字 | ASCII a-zA-Z0-9_ |
| `\s` | 所有空白 | 常见空白字符 |
| 转义字符 | 完整支持 | 基础支持 |

## 扩展思考

1. **更多字符类**：如何添加 `\h`（十六进制）、`\a`（字母）？
2. **自定义字符类**：如何实现 `[abc]`（字符组）？
3. **Unicode 支持**：如何支持 Unicode 字符类？
4. **性能优化**：大量字符类匹配时如何优化？

## 下一步

完成本课后，你应该能够：
- ✓ 理解字符类的概念和用途
- ✓ 实现转义字符解析
- ✓ 支持 `\d`、`\w`、`\s` 字符类
- ✓ 实现字面转义（`\.`、`\\` 等）
- ✓ 将字符类与已有功能组合使用

**下一课预告**：第五天我们将实现字符组 `[abc]` 和否定字符组 `[^abc]`，支持更灵活的字符匹配。

## 参考资料

- [Regular Expressions - Character Classes](https://www.regular-expressions.info/charclass.html)
- [MDN - Character Classes](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions/Character_Classes)
- [PCRE - Escape Sequences](https://www.pcre.org/original/doc/html/pcrepattern.html)
