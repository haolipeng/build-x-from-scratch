# 第七天：锚点

## 学习目标

在这个阶段，你将学习如何：
1. 理解锚点（Anchors）的概念
2. 实现行首锚点 `^`
3. 实现行尾锚点 `$`
4. 实现单词边界 `\b`
5. 理解锚点如何影响匹配行为

## 什么是锚点？

锚点是正则表达式中的特殊元字符，它们不匹配具体的字符，而是匹配**位置**。

### 常用锚点

| 锚点 | 含义 | 示例 |
|------|------|------|
| `^` | 行首（字符串开头） | `^hello` 匹配以 "hello" 开头的行 |
| `$` | 行尾（字符串结尾） | `world$` 匹配以 "world" 结尾的行 |
| `\b` | 单词边界 | `\bword\b` 匹配独立的 "word" |
| `\B` | 非单词边界 | `\Bword` 匹配不在单词开头的 "word" |

### 使用示例

```bash
# 匹配以 "hello" 开头的行
echo -e "hello world\nworld hello\nhello" | grep "^hello"
# 输出:
# hello world
# hello

# 匹配以 "world" 结尾的行
echo -e "hello world\nworld hello\nworld" | grep "world$"
# 输出:
# hello world
# world

# 匹配独立的单词 "the"
echo -e "the cat\nthere\nother" | grep "\bthe\b"
# 输出:
# the cat
```

## 锚点的工作原理

### 1. 行首锚点 `^`

`^` 匹配字符串的开始位置：

```
文本: "hello world"
模式: "^hello"

位置: 0 1 2 3 4 5 ...
      h e l l o   w ...
      ^
      └── ^ 匹配位置 0（开始）

结果: 匹配成功
```

```
文本: "world hello"
模式: "^hello"

位置: 0 1 2 3 4 5 6 7 8 9 10
      w o r l d   h e l l o
      ^
      └── 在位置 0，下一个字符是 'w'，不是 'h'

结果: 匹配失败
```

### 2. 行尾锚点 `$`

`$` 匹配字符串的结束位置：

```
文本: "hello world"
模式: "world$"

位置: 0 1 2 3 4 5 6 7 8 9 10 [END]
      h e l l o   w o r l d  $
                            ^
                            └── $ 匹配位置 11（结束）

结果: 匹配成功
```

### 3. 单词边界 `\b`

`\b` 匹配单词字符和非单词字符之间的位置：

```
文本: "the cat"
模式: "\bthe\b"

位置: [START] t h e   c a t [END]
      ^       ^   ^
      |       |   └── \b（单词结束：'e' 后是空格）
      |       └── 't', 'h', 'e' 匹配
      └── \b（单词开始：开头是单词字符）

结果: 匹配成功
```

```
文本: "there"
模式: "\bthe\b"

位置: [START] t h e r e [END]
      ^       ^   ^
      |       |   └── 'r' 不是 \b（都是单词字符）
      |       └── 't', 'h', 'e' 匹配
      └── \b（单词开始）

结果: 匹配失败（"the" 后面的 \b 不匹配）
```

## 本课实现内容

### 1. Token 类型扩展

```go
const (
    // ... 已有类型
    TokenStartAnchor    // ^ 行首
    TokenEndAnchor      // $ 行尾
    TokenWordBoundary   // \b 单词边界
    TokenNotWordBoundary // \B 非单词边界
)
```

### 2. 锚点匹配逻辑

```go
// 检查行首锚点
func matchStartAnchor(textIdx int) bool {
    return textIdx == 0
}

// 检查行尾锚点
func matchEndAnchor(text string, textIdx int) bool {
    return textIdx == len(text)
}

// 检查单词边界
func matchWordBoundary(text string, textIdx int) bool {
    atStart := textIdx == 0
    atEnd := textIdx == len(text)

    prevIsWord := !atStart && isWordChar(text[textIdx-1])
    currIsWord := !atEnd && isWordChar(text[textIdx])

    // 单词边界：一边是单词字符，另一边不是
    return prevIsWord != currIsWord
}
```

### 3. 匹配算法调整

锚点不消耗字符，只检查位置：

```go
func matchTokens(tokens, text, patIdx, textIdx) bool {
    token := tokens[patIdx]

    switch token.Type {
    case TokenStartAnchor:
        if textIdx == 0 {
            return matchTokens(tokens, text, patIdx+1, textIdx) // 不移动 textIdx
        }
        return false

    case TokenEndAnchor:
        if textIdx == len(text) {
            return matchTokens(tokens, text, patIdx+1, textIdx)
        }
        return false

    case TokenWordBoundary:
        if isWordBoundary(text, textIdx) {
            return matchTokens(tokens, text, patIdx+1, textIdx)
        }
        return false

    default:
        // 普通字符匹配
        if matchChar(token, text[textIdx]) {
            return matchTokens(tokens, text, patIdx+1, textIdx+1)
        }
        return false
    }
}
```

## 测试用例

### 测试 1: 行首锚点 `^`

```bash
echo -e "hello world\nworld hello\nhello" | go run . "^hello"
# 输出:
# hello world
# hello

echo -e "  hello\nhello" | go run . "^hello"
# 输出:
# hello
```

### 测试 2: 行尾锚点 `$`

```bash
echo -e "hello world\nworld hello\nworld" | go run . "world$"
# 输出:
# hello world
# world

echo -e "hello\nhello  " | go run . "hello$"
# 输出:
# hello
```

### 测试 3: 同时使用 `^` 和 `$`

```bash
echo -e "hello\nhello world\nworld" | go run . "^hello$"
# 输出:
# hello

echo -e "a\nab\nabc" | go run . "^..$"
# 输出:
# ab
```

### 测试 4: 单词边界 `\b`

```bash
echo -e "the cat\nthere\nother" | go run . "\\bthe\\b"
# 输出:
# the cat

echo -e "cat\ncatalog\ncat food" | go run . "\\bcat\\b"
# 输出:
# cat
# cat food
```

### 测试 5: 非单词边界 `\B`

```bash
echo -e "there\nthe\nother" | go run . "\\Bthe"
# 输出:
# other (匹配 "other" 中的 "the")

echo -e "catalog\ncat\nconcat" | go run . "cat\\B"
# 输出:
# catalog (匹配 "catalog" 开头的 "cat")
# concat (匹配 "concat" 结尾的 "cat")
```

### 测试 6: 与选项组合

```bash
# ^pattern + -n
echo -e "hello\nworld\nhello world" | go run . -n "^hello"
# 输出:
# 1:hello
# 3:hello world

# pattern$ + -v
echo -e "hello\nworld\nhello world" | go run . -v "world$"
# 输出:
# hello
```

## 关键知识点

### 1. 锚点不消耗字符

锚点只检查位置，不消耗（移动）文本指针：

```go
// 普通字符：消耗一个字符
return matchTokens(tokens, text, patIdx+1, textIdx+1)

// 锚点：不消耗字符
return matchTokens(tokens, text, patIdx+1, textIdx)
```

### 2. `^` 在不同位置的含义

```
^hello     行首的 hello
[^abc]     否定字符组（不是 a、b、c）
```

规则：只有在模式开头时，`^` 才是行首锚点。

### 3. 单词边界的定义

单词边界出现在：
- 字符串开头，如果第一个字符是单词字符
- 字符串结尾，如果最后一个字符是单词字符
- 单词字符和非单词字符之间

```go
func isWordBoundary(text string, pos int) bool {
    prevIsWord := pos > 0 && isWordChar(text[pos-1])
    currIsWord := pos < len(text) && isWordChar(text[pos])
    return prevIsWord != currIsWord
}
```

### 4. 锚点与包含匹配

使用锚点时，匹配行为会改变：

```bash
# 无锚点：在任意位置匹配
echo "hello world" | grep "hello"  # 匹配

# 有行首锚点：只在开头匹配
echo "world hello" | grep "^hello"  # 不匹配
```

## 实际应用场景

### 1. 验证整行格式

```bash
# 验证是否为有效的数字
echo -e "123\nabc\n456" | go run . "^[0-9]+$"
# 输出: 123, 456（注：+ 量词在后续课程实现）
```

### 2. 查找独立单词

```bash
# 查找独立的 "is"（不是 "this" 或 "island" 中的）
echo -e "this is\nisland\nit is here" | go run . "\\bis\\b"
# 输出:
# this is
# it is here
```

### 3. 查找行首的注释

```bash
# 查找以 # 开头的行（注释）
echo -e "# comment\ncode # comment\n  # indented" | go run . "^#"
# 输出:
# # comment
```

### 4. 查找空行

```bash
# 匹配空行（行首紧接行尾）
echo -e "hello\n\nworld" | go run . "^$"
# 输出: (空行)
```

## 性能考虑

### 1. 行首锚点优化

当模式以 `^` 开头时，只需在位置 0 尝试匹配：

```go
func (m *RegexMatcher) Match(text string) bool {
    if startsWithAnchor(m.pattern) {
        // 只在位置 0 尝试
        return m.matchAt(text, 0)
    }

    // 否则尝试所有位置
    for i := 0; i <= len(text); i++ {
        if m.matchAt(text, i) {
            return true
        }
    }
    return false
}
```

### 2. 行尾锚点优化

类似地，可以从结尾开始匹配：

```go
if endsWithAnchor(m.pattern) {
    // 从接近结尾的位置开始尝试
    startPos := len(text) - patternLength + 1
    if startPos < 0 {
        startPos = 0
    }
    // ...
}
```

## 常见错误和调试

### 1. 混淆 `^` 的含义

```bash
# 行首锚点
echo "hello" | grep "^h"     # 匹配

# 否定字符组
echo "hello" | grep "[^h]"   # 匹配 'e', 'l', 'o'
```

### 2. 忘记转义 `\b`

```bash
# 错误：\b 被 shell 解释
echo "the cat" | grep "\bthe"

# 正确：使用双反斜杠或单引号
echo "the cat" | grep "\\bthe"
echo "the cat" | grep '\bthe'
```

### 3. 空模式边界情况

```bash
# 空模式 + 锚点
echo "hello" | grep "^$"   # 只匹配空行
echo "" | grep "^$"        # 匹配
```

## 扩展思考

1. **多行模式**：如何让 `^` 和 `$` 匹配每行而不是整个字符串？
2. **字符串边界**：`\A`（字符串开头）和 `\Z`（字符串结尾）与 `^`、`$` 的区别？
3. **Unicode 单词边界**：如何处理中文等非 ASCII 字符的单词边界？

## 下一步

完成本课后，你应该能够：
- ✓ 理解锚点的概念和工作原理
- ✓ 实现行首锚点 `^`
- ✓ 实现行尾锚点 `$`
- ✓ 实现单词边界 `\b` 和 `\B`
- ✓ 理解锚点如何影响匹配行为

**下一课预告**：第八天我们将实现量词，支持 `?`（零或一次）、`+`（一或多次）、`*`（零或多次）。

## 参考资料

- [Regular Expressions - Anchors](https://www.regular-expressions.info/anchors.html)
- [Word Boundaries](https://www.regular-expressions.info/wordboundaries.html)
- [GNU Grep - Anchoring](https://www.gnu.org/software/grep/manual/grep.html#Anchoring)
