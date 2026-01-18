# 第五天：字符组

## 学习目标

在这个阶段，你将学习如何：
1. 理解字符组（Character Groups）的概念
2. 实现字符枚举 `[abc]`
3. 实现字符范围 `[a-z]`、`[0-9]`
4. 实现否定字符组 `[^abc]`
5. 支持组合使用 `[a-zA-Z0-9_]`

## 什么是字符组？

字符组（也称为字符类或字符集）允许你定义一组字符，匹配其中的任意一个。

### 基本语法

```
[abc]     匹�� a、b 或 c 中的任意一个
[a-z]     匹配 a 到 z 之间的任意小写字母
[A-Z]     匹配 A 到 Z 之间的任意大写字母
[0-9]     匹配 0 到 9 之间的任意数字
[^abc]    匹配除 a、b、c 之外的任意字符
```

### 组合使用

```
[a-zA-Z]        匹配任意字母
[a-zA-Z0-9]     匹配任意字母或数字
[a-zA-Z0-9_]    等价于 \w
[aeiou]         匹配元音字母
```

## 使用示例

### 示例 1: 字符枚举

```bash
# 匹配元音字母
echo -e "apple\nbanana\norange" | grep "[aeo]"
# 输出: 所有行（都包含元音）

# 匹配特定字符
echo -e "cat\nhat\nbat\nrat" | grep "[ch]at"
# 输出:
# cat
# hat
```

### 示例 2: 字符范围

```bash
# 匹配小写字母
echo -e "hello\nHELLO\n12345" | grep "[a-z]"
# 输出:
# hello

# 匹配数字
echo -e "abc\na1b\n123" | grep "[0-9]"
# 输出:
# a1b
# 123
```

### 示例 3: 否定字符组

```bash
# 匹配非数字字符
echo -e "abc\n123\na1b" | grep "[^0-9]"
# 输出:
# abc
# a1b

# 匹配非元音字母
echo -e "hello\nsky\nworld" | grep "[^aeiou]"
# 输出: 所有行（都包含辅音）
```

### 示例 4: 组合字符组

```bash
# 匹配字母或数字
echo -e "abc\n123\n@#$" | grep "[a-zA-Z0-9]"
# 输出:
# abc
# 123

# 匹配十六进制字符
echo -e "1a2b\nGHIJ\nffff" | grep "[0-9a-fA-F]"
# 输出:
# 1a2b
# ffff
```

## 实现架构

### 1. Token 类型扩展

```go
const (
    TokenLiteral    TokenType = iota
    TokenWildcard
    TokenDigit
    // ...
    TokenCharGroup     // 字符组 [abc]
    TokenNegCharGroup  // 否定字符组 [^abc]
)

// CharGroup 存储字���组的内容
type CharGroup struct {
    Chars   []byte    // 单个字符列表
    Ranges  []Range   // 字符范围列表
}

type Range struct {
    Start byte
    End   byte
}
```

### 2. 解析流程

```
输入: "[a-zA-Z0-9_]"

┌─────────────────────┐
│  遇到 '[' 开始解析   │
└──────────┬──────────┘
           │
           ▼
    是否 '^'？
    ┌──────┴──────┐
    │             │
   是            否
    │             │
    ▼             ▼
  否定组        正向组
    │             │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│  解析字符和范围      │
│  a-z, A-Z, 0-9, _   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  遇到 ']' 结束解析   │
└──────────┬──────────┘
           │
           ▼
   创建 Token
```

### 3. 字符范围识别

如何判断 `-` 是范围符号还是字面字符？

```
[a-z]   → a 到 z 的范围
[a-]    → 字符 'a' 和 '-'
[-z]    → 字符 '-' 和 'z'
[a-z-]  → a 到 z 的范围，加上 '-'
```

规则：
- 如果 `-` 在开头或结尾，它是字面字符
- 如果 `-` 两边都有字符，它是范围符号

## 本课实现内容

### 1. 扩展 regex.go

添加字符组解析：

```go
// CharGroup 表示字符组的内容
type CharGroup struct {
    Negated bool      // 是否为否定字符组
    Chars   []byte    // 单个字符
    Ranges  []Range   // 字符范围
}

type Range struct {
    Start byte
    End   byte
}

// Token 扩展
type Token struct {
    Type      TokenType
    Value     byte       // 字面字符
    CharGroup *CharGroup // 字符组（仅当 Type 为 TokenCharGroup 时使用）
}
```

### 2. 解析字符组

```go
func parseCharGroup(pattern string, pos int) (Token, int, error) {
    // 跳过 '['
    pos++

    // 检查是否为否定字符组
    negated := false
    if pos < len(pattern) && pattern[pos] == '^' {
        negated = true
        pos++
    }

    // 解析字符和范围
    chars := []byte{}
    ranges := []Range{}

    for pos < len(pattern) && pattern[pos] != ']' {
        // 检查是否为范围
        if pos+2 < len(pattern) &&
           pattern[pos+1] == '-' &&
           pattern[pos+2] != ']' {
            // 是范围
            ranges = append(ranges, Range{
                Start: pattern[pos],
                End:   pattern[pos+2],
            })
            pos += 3
        } else {
            // 单个字符
            chars = append(chars, pattern[pos])
            pos++
        }
    }

    // 跳过 ']'
    pos++

    return Token{
        Type: TokenCharGroup,
        CharGroup: &CharGroup{
            Negated: negated,
            Chars:   chars,
            Ranges:  ranges,
        },
    }, pos, nil
}
```

### 3. 匹配字符组

```go
func (cg *CharGroup) Contains(char byte) bool {
    // 检查单个字符
    for _, c := range cg.Chars {
        if c == char {
            return true
        }
    }

    // 检查范围
    for _, r := range cg.Ranges {
        if char >= r.Start && char <= r.End {
            return true
        }
    }

    return false
}

func (cg *CharGroup) Match(char byte) bool {
    contains := cg.Contains(char)
    if cg.Negated {
        return !contains
    }
    return contains
}
```

## 测试用例

### 测试 1: 字符枚举 [abc]

```bash
echo -e "cat\nhat\nbat\nrat" | go run . "[ch]at"
# 输出:
# cat
# hat

echo -e "apple\nbanana\norange" | go run . "[aeo]"
# 输出: 所有行
```

### 测试 2: 字符范围 [a-z]

```bash
echo -e "hello\nHELLO\n12345" | go run . "[a-z]"
# 输出:
# hello

echo -e "ABC\nabc\n123" | go run . "[A-Z]"
# 输出:
# ABC
```

### 测试 3: 数字范围 [0-9]

```bash
echo -e "abc\na1b\n123" | go run . "[0-9]"
# 输出:
# a1b
# 123

echo "ID: 12345" | go run . "ID: [0-9][0-9][0-9][0-9][0-9]"
# 输出: ID: 12345
```

### 测试 4: 否定字符组 [^abc]

```bash
echo -e "abc\nxyz\naxy" | go run . "[^abc]"
# 输出:
# xyz
# axy

echo -e "123\nabc\na1b" | go run . "[^0-9]"
# 输出:
# abc
# a1b
```

### 测试 5: 组合字符组

```bash
echo -e "abc\nABC\n123\n@#$" | go run . "[a-zA-Z]"
# 输出:
# abc
# ABC

echo -e "abc\n123\nabc123" | go run . "[a-zA-Z0-9]"
# 输出: 所有行
```

### 测试 6: 十六进制匹配

```bash
echo -e "0xff\n0xGG\n0x1a" | go run . "0x[0-9a-fA-F][0-9a-fA-F]"
# 输出:
# 0xff
# 0x1a
```

### 测试 7: 字面 `-` 和 `]`

```bash
# '-' 在开头或结尾是字面字符
echo "a-b" | go run . "[a-]"
# 输出: a-b

echo "a-b" | go run . "[-b]"
# 输出: a-b
```

### 测试 8: 与选项组合

```bash
# -n: 显示行号
echo -e "abc\n123\nxyz" | go run . -n "[0-9]"
# 输出:
# 2:123

# -v: 反向匹配
echo -e "abc\n123\nxyz" | go run . -v "[0-9]"
# 输出:
# abc
# xyz

# -i: 忽略大小写（对字符组不适用）
echo -e "ABC\nabc" | go run . -i "[a-z]"
# 输出: 两行都匹配（因为 -i 会将文本转为小写）
```

## 关键知识点

### 1. 字符组中的特殊字符

在字符组内部，大部分元字符失去特殊含义：

```
[.]     匹配字面点号（不需要转义）
[*]     匹配字面星号
[+]     匹配字面加号
[\]]    匹配字面右括号（需要转义）
[\\]    匹配字面反斜杠
[^]     如果不在开头，匹配字面 ^
[-]     如果在开头或结尾，匹配字面 -
```

### 2. 范围的有效性

```go
// 有效范围
[a-z]   // OK: a (97) < z (122)
[A-Z]   // OK: A (65) < Z (90)
[0-9]   // OK: 0 (48) < 9 (57)

// 无效范围（后续可添加错误检查）
[z-a]   // 错误: z > a
[9-0]   // 错误: 9 > 0
```

### 3. 预定义字符类 vs 字符组

```
\d      等价于 [0-9]
\w      等价于 [a-zA-Z0-9_]
\s      等价于 [ \t\n\r]

# 但字符组更灵活
[aeiou]         只匹配元音
[0-9a-fA-F]     只匹配十六进制
[^aeiou]        只匹配辅音
```

### 4. 性能考虑

对于大型字符组，可以使用位图优化：

```go
// 使用位图存储字符集
type CharGroup struct {
    bitmap [256]bool  // 256 个 ASCII 字符
    Negated bool
}

func (cg *CharGroup) Contains(char byte) bool {
    return cg.bitmap[char]
}
```

## 常见错误和调试

### 1. 未闭合的字符组

```
[abc     错误: 缺少 ]
```

### 2. 空字符组

```
[]       某些实现允许，某些不允许
[^]      同上
```

### 3. 范围顺序错误

```
[z-a]    错误: 结束字符小于开始字符
```

### 4. 忘记转义

```
[\]]     正确: 匹配 ]
[]]      某些实现允许（] 紧跟 [）
```

## 与标准正则表达式的差异

| 特性 | 标准 Regex | 我们的实现 |
|------|-----------|------------|
| POSIX 类 | `[[:alpha:]]` | 不支持 |
| Unicode | 支持 | 仅 ASCII |
| 交集/差集 | `[a-z&&[^aeiou]]` | 不支持 |
| 嵌套字符类 | 某些支持 | 不支持 |

## 扩展思考

1. **POSIX 字符类**：如何添加 `[[:alpha:]]`、`[[:digit:]]`？
2. **Unicode 范围**：如何支持 `[\u4e00-\u9fff]`（中文字符）？
3. **字符组简写**：如何实现 `[\w]` 在字符组内等价于 `a-zA-Z0-9_`？
4. **性能优化**：大型字符组如何高效匹配？

## 下一步

完成本课后，你应该能够：
- ✓ 理解字符组的概念和语法
- ✓ 实现字符枚举和字符范围解析
- ✓ 实现否定字符组
- ✓ 处理字符组内的特殊字符
- ✓ 将字符组与已有功能组合使用

**下一课预告**：第六天我们将学习组合字符类，实现更复杂的匹配模式和字符类的高级用法。

## 参考资料

- [Regular Expressions - Character Classes](https://www.regular-expressions.info/charclass.html)
- [MDN - Character Classes](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions/Character_Classes)
- [POSIX Character Classes](https://www.regular-expressions.info/posixbrackets.html)
