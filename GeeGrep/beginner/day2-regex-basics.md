# 第二天：正则表达式基础

## 学习目标

在这个阶段，你将学习如何：
1. 理解正则表达式的基本概念
2. 实现一个简单的正则表达式匹配引擎
3. 支持字面字符匹配
4. 支持通配符 `.` (匹配任意字符)
5. 理解正则表达式匹配的工作原理

## 什么是正则表达式？

正则表达式（Regular Expression，简称 regex）是一种强大的文本模式匹配工具。它使用特殊的语法来描述字符串的模式，可以用于：

- 搜索文本
- 验证输入
- 提取信息
- 替换文本

### 基本概念

**字面字符（Literal Characters）**：直接匹配自己
```
模式 "hello" 匹配字符串 "hello"
```

**元字符（Metacharacters）**：具有特殊含义的字符
```
. ^ $ * + ? { } [ ] \ | ( )
```

**通配符 `.`**：匹配任意单个字符（除换行符外）
```
模式 "h.llo" 可以匹配 "hello", "hallo", "hxllo" 等
```

## 本课实现内容

### 1. 正则表达式匹配引擎架构

```
输入: pattern = "h.llo", text = "hello world"

┌─────────────────────┐
│  解析正则表达式      │
│  (Parse Pattern)    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  生成匹配规则       │
│  [h] [.] [l] [l] [o]│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  在文本中查找匹配   │
│  (Find Match)       │
└──────────┬──────────┘
           │
           ▼
    ┌──────┴──────┐
    │  匹配成功？  │
    └──────┬──────┘
           │
    ┌──────┴──────┐
    │             │
   是            否
    │             │
    ▼             ▼
 返回true      继续搜索
```

### 2. 匹配算法

我们将实现一个简单但强大的匹配算法：

```go
// 伪代码
func match(pattern, text):
    if pattern is empty:
        return true

    if text is empty:
        return pattern matches empty string

    first_match = (pattern[0] matches text[0])

    if first_match:
        return match(pattern[1:], text[1:])
    else:
        return false
```

### 3. 支持的模式

**Day 2 支持的模式**：
- 字面字符：`a`, `b`, `1`, `@` 等
- 通配符：`.` 匹配任意单个字符

**示例**：
```
模式 "cat"    匹配 "cat"
模式 "c.t"    匹配 "cat", "cot", "c9t" 等
模式 "..."    匹配任意三个字符
```

## 实现细节

### 1. 正则表达式模式解析

第一步是将正则表达式字符串解析为内部表示：

```go
// Pattern 表示一个正则表达式模式
type Pattern struct {
    tokens []Token
}

// Token 表示模式中的一个元素
type Token struct {
    Type  TokenType  // 字面字符或通配符
    Value byte       // 如果是字面字符，存储其值
}

type TokenType int

const (
    TokenLiteral  TokenType = iota  // 字面字符
    TokenWildcard                   // 通配符 .
)
```

### 2. 匹配算法实现

我们实现两种匹配方式：

**完整匹配**（Full Match）：
- 整个文本必须完全匹配模式
- 例如：模式 "cat" 只匹配 "cat"，不匹配 "cats" 或 "scat"

**包含匹配**（Contains Match）：
- 文本中包含匹配模式的子串即可
- 例如：模式 "cat" 可以匹配 "my cat is cute" 中的 "cat"
- 这是 grep 的默认行为

### 3. 匹配流程

```go
// 包含匹配：在文本中查找是否包含匹配的子串
func containsMatch(pattern Pattern, text string) bool {
    // 尝试从文本的每个位置开始匹配
    for i := 0; i <= len(text); i++ {
        if fullMatch(pattern, text[i:]) {
            return true
        }
    }
    return false
}

// 完整匹配：从当前位置开始完整匹配模式
func fullMatch(pattern Pattern, text string) bool {
    // 递归匹配每个 token
    return matchTokens(pattern.tokens, text, 0, 0)
}
```

## 测试用例

### 测试 1：字面字符匹配
```bash
echo "hello" | go run . "hello"
# 输出: hello

echo "hello world" | go run . "hello"
# 输出: hello world

echo "goodbye" | go run . "hello"
# 无输出（未匹配）
```

### 测试 2：通配符匹配
```bash
echo "cat" | go run . "c.t"
# 输出: cat

echo "cot" | go run . "c.t"
# 输出: cot

echo "coat" | go run . "c.t"
# 无输出（模式是3个字符，"coat" 有4个字符，但会匹配 "coa"）
```

### 测试 3：多个通配符
```bash
echo "hello" | go run . "h...o"
# 输出: hello

echo "hallo" | go run . "h...o"
# 输出: hallo

echo "halo" | go run . "h...o"
# 无输出（只有4个字符）
```

### 测试 4：混合模式
```bash
echo "pattern123test" | go run . "pattern...test"
# 输出: pattern123test

echo -e "cat\ncut\ncot\ncaught" | go run . "c.t"
# 输出:
# cat
# cut
# cot
```

## 关键知识点

### 1. 正则表达式的两种匹配模式

```go
// Full match: 整个字符串必须匹配
"cat" matches "cat"          ✓
"cat" matches "cats"         ✗
"cat" matches "my cat"       ✗

// Contains match (grep 默认)
"cat" in "cat"               ✓
"cat" in "cats"              ✓
"cat" in "my cat is cute"    ✓
```

### 2. 递归匹配算法

```go
func matchAt(pattern []Token, patIdx int, text string, textIdx int) bool {
    // 基础情况：模式匹配完成
    if patIdx >= len(pattern) {
        return true  // 模式已全部匹配
    }

    // 基础情况：文本耗尽但模式未完
    if textIdx >= len(text) {
        return false
    }

    token := pattern[patIdx]

    // 匹配当前字符
    if token.Type == TokenWildcard {
        // 通配符匹配任意字符
        return matchAt(pattern, patIdx+1, text, textIdx+1)
    } else {
        // 字面字符必须精确匹配
        if text[textIdx] == token.Value {
            return matchAt(pattern, patIdx+1, text, textIdx+1)
        }
        return false
    }
}
```

### 3. 通配符的特殊性

通配符 `.` 匹配**任意单个字符**，但有例外：
- 在标准正则表达式中，`.` 不匹配换行符 `\n`
- 在我们的实现中，暂时简化为匹配任意字符

### 4. 转义字符

如果想匹配字面意义的 `.`，需要转义：`\.`
```
模式 "\\."  匹配字符 "."
模式 "3\\.14" 匹配 "3.14"
```

在 Day 2 中，我们暂时不实现转义，将在后续课程中添加。

## 代码结构

### regex.go
- `Token` 和 `TokenType`：表示模式的元素
- `Pattern`：表示整个正则表达式
- `ParsePattern()`：解析正则表达式字符串
- `Compile()`：编译正则表达式

### engine.go
- `Match()`：判断文本是否包含匹配
- `MatchFull()`：判断文本是否完全匹配
- `matchAt()`：从指定位置开始匹配

### main.go
- 集成新的正则表达式引擎
- 保持与 Day 1 相同的命令行接口

## 性能考虑

### 时间复杂度

对于简单的模式（无量词），匹配算法的时间复杂度：
- 最坏情况：O(n * m)，其中 n 是文本长度，m 是模式长度
- 我们需要在文本的每个位置尝试匹配

### 优化方向

后续课程中我们将学习：
1. **提前终止**：如果模式的第一个字符不匹配，跳过
2. **编译优化**：将模式编译为状态机
3. **Boyer-Moore 算法**：智能跳过不可能匹配的位置

## 与 Day 1 的区别

| 特性 | Day 1 | Day 2 |
|------|-------|-------|
| 匹配类型 | 字符串包含检查 | 正则表达式匹配 |
| 支持通配符 | ✗ | ✓ (`.`) |
| 模式解析 | 无需解析 | 需要解析 Token |
| 匹配算法 | `strings.Contains()` | 递归匹配算法 |
| 扩展性 | 难以扩展 | 易于添加新特性 |

## 常见问题

### Q: 为什么要设计 Token 结构？
A: Token 是正则表达式内部表示的基本单元。通过将字符串解析为 Token 序列，我们可以：
- 区分字面字符和元字符
- 为后续添加量词、字符类等复杂特性做准备
- 实现更高效的匹配算法

### Q: 递归算法会不会栈溢出？
A: 对于简单模式（Day 2 的实现），递归深度等于模式长度，通常不会有问题。在后续课程中，我们会：
- 添加尾递归优化
- 或改用迭代实现
- 或使用状态机

### Q: `.` 为什么不匹配换行符？
A: 这是正则表达式的传统约定。原因是：
- 文本通常是按行处理的
- grep 默认逐行匹配
- 如果需要匹配换行符，可以使用特殊标志（如 Perl 的 `/s` 修饰符）

### Q: 这个实现和标准库的 regexp 包有什么区别？
A: 标准库使用了更复杂的算法（如 Thompson NFA），支持更多特性。我们的实现：
- 更简单，易于理解
- 性能较低，但对学习足够了
- 逐步添加特性，循序渐进

## 扩展思考

1. **贪婪匹配**：如果文本是 "hello"，模式 "h.llo" 有多个匹配方式吗？
2. **多行模式**：如何让 `.` 也匹配换行符？
3. **Unicode 支持**：如何处理多字节字符？
4. **性能优化**：如何避免不必要的回溯？

## 下一步

完成本课后，你应该能够：
- ✓ 理解正则表达式的基本工作原理
- ✓ 实现简单的模式解析器
- ✓ 实现字面字符和通配符匹配
- ✓ 理解递归匹配算法

**下一课预告**：第三天我们将添加基础命令行选项，包括 `-n`（显示行号）、`-i`（忽略大小写）、`-v`（反向匹配）等。

## 参考资料

- [正则表达式30分钟入门教程](https://deerchao.cn/tutorials/regex/regex.htm)
- [Regular Expression Matching Can Be Simple And Fast](https://swtch.com/~rsc/regexp/regexp1.html)
- [Implementing Regular Expressions](https://swtch.com/~rsc/regexp/)
- [Go regexp 包源码](https://github.com/golang/go/tree/master/src/regexp)
