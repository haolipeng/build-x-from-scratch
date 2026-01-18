# 第八天：量词（一）

## 学习目标

在这个阶段，你将学习如何：
1. 理解量词（Quantifiers）的概念
2. 实现 `?`（零或一次）
3. 实现 `+`（一或多次）
4. 实现 `*`（零或多次）
5. 理解贪婪匹配的基本原理

## 什么是量词？

量词用于指定前面的元素可以出现的次数。

### 基本量词

| 量词 | 含义 | 示例 |
|------|------|------|
| `?` | 零次或一次 | `colou?r` 匹配 "color" 或 "colour" |
| `+` | 一次或多次 | `a+` 匹配 "a", "aa", "aaa"... |
| `*` | 零次或多次 | `a*` 匹配 "", "a", "aa", "aaa"... |

### 使用示例

```bash
# ? - 零或一次
echo -e "color\ncolour" | grep "colou?r"
# 输出:
# color
# colour

# + - 一或多次
echo -e "a\naa\naaa\nb" | grep "a+"
# 输出:
# a
# aa
# aaa

# * - 零或多次
echo -e "ac\nabc\nabbc\nabbbc" | grep "ab*c"
# 输出:
# ac
# abc
# abbc
# abbbc
```

## 量词的工作原理

### 1. `?` - 零或一次

`?` 使前面的元素变为可选：

```
模式: "colou?r"

匹配 "color":
  c-o-l-o-r
  c o l o u? r
  ↓ ↓ ↓ ↓    ↓
  ✓ ✓ ✓ ✓ (u出现0次) ✓

匹配 "colour":
  c-o-l-o-u-r
  c o l o u? r
  ↓ ↓ ↓ ↓ ↓  ↓
  ✓ ✓ ✓ ✓ ✓  ✓
```

### 2. `+` - 一或多次

`+` 要求前面的元素至少出现一次：

```
模式: "a+"

匹配 "aaa":
  a-a-a
  a+
  ↓↓↓
  ✓✓✓ (a出现3次，满足1+)

不匹配 "b":
  b
  a+
  ↓
  ✗ (a出现0次，不满足1+)
```

### 3. `*` - 零或多次

`*` 允许前面的元素出现任意次（包括零次）：

```
模式: "ab*c"

匹配 "ac":
  a-c
  a b* c
  ↓    ↓
  ✓ (b出现0次) ✓

匹配 "abc":
  a-b-c
  a b* c
  ↓ ↓  ↓
  ✓ ✓  ✓

匹配 "abbbc":
  a-b-b-b-c
  a b*    c
  ↓ ↓↓↓   ↓
  ✓ ✓✓✓   ✓
```

## 贪婪匹配

默认情况下，量词是**贪婪**的，会尽可能多地匹配：

```
文本: "aaaaab"
模式: "a+"

贪婪匹配过程:
1. 尝试匹配尽可能多的 'a'
2. 匹配 "aaaaa"（5个a）
3. 然后模式结束

结果: 匹配 "aaaaa"
```

## 本课实现内容

### 1. Token 类型扩展

```go
const (
    // ... 已有类型
    TokenQuestion  // ? 量词
    TokenPlus      // + 量词
    TokenStar      // * 量词
)
```

### 2. 量词与前一个 Token 的关联

量词修饰它前面的 token。我们需要在解析时处理这种关系：

```go
type Token struct {
    Type      TokenType
    Value     byte
    CharGroup *CharGroup
    Quantifier *Quantifier  // 新增：量词信息
}

type Quantifier struct {
    Min int  // 最小次数
    Max int  // 最大次数（-1 表示无限）
}
```

### 3. 匹配算法

量词匹配需要回溯（backtracking）：

```go
func matchQuantified(token Token, tokens []Token, text string, patIdx, textIdx int) bool {
    min := token.Quantifier.Min
    max := token.Quantifier.Max

    // 贪婪匹配：先尝试最多次数
    for count := max; count >= min; count-- {
        if tryMatchN(token, text, textIdx, count) {
            newTextIdx := textIdx + count
            if matchTokens(tokens, text, patIdx+1, newTextIdx) {
                return true
            }
        }
    }
    return false
}
```

## 测试用例

### 测试 1: `?` 零或一次

```bash
echo -e "color\ncolour\ncolouur" | go run . "colou?r"
# 输出:
# color
# colour

echo -e "ac\nabc\nabbc" | go run . "ab?c"
# 输出:
# ac
# abc
```

### 测试 2: `+` 一或多次

```bash
echo -e "a\naa\naaa\nb" | go run . "a+"
# 输出:
# a
# aa
# aaa

echo -e "abc\naabc\naaabc\nbc" | go run . "a+bc"
# 输出:
# abc
# aabc
# aaabc
```

### 测试 3: `*` 零或多次

```bash
echo -e "ac\nabc\nabbc\nabbbc" | go run . "ab*c"
# 输出: 所有行

echo -e "bc\nabc\naabc\naaabc" | go run . "a*bc"
# 输出: 所有行
```

### 测试 4: 字符类 + 量词

```bash
echo -e "a\na1\na12\na123" | go run . "a\\d+"
# 输出:
# a1
# a12
# a123

echo -e "abc\nabc123\n123" | go run . "\\w+"
# 输出: 所有行（包含单词字符）
```

### 测试 5: 通配符 + 量词

```bash
echo -e "ac\naXc\naXXc\naXXXc" | go run . "a.+c"
# 输出:
# aXc
# aXXc
# aXXXc

echo -e "ac\naXc\naXXc" | go run . "a.*c"
# 输出: 所有行
```

### 测试 6: 与锚点组合

```bash
echo -e "hello\nhello world\nworld hello" | go run . "^hello.*"
# 输出:
# hello
# hello world

echo -e "123\nabc123\n123abc" | go run . "\\d+$"
# 输出:
# 123
# abc123
```

## 关键知识点

### 1. 量词的优先级

量词只修饰紧前面的一个元素：

```
a+b     → (a+)b      一个或多个 a，后跟 b
[abc]+  → ([abc])+   一个或多个 [abc] 中的字符
\d+     → (\d)+      一个或多个数字
```

### 2. 贪婪 vs 非贪婪

```
文本: "aaaaab"
模式: "a+"

贪婪（默认）: 匹配 "aaaaa"（尽可能多）
非贪婪: 匹配 "a"（尽可能少）- 在 Day 10 实现
```

### 3. 回溯机制

当贪婪匹配失败时，需要回溯减少匹配数量：

```
文本: "aaab"
模式: "a+b"

1. a+ 贪婪匹配 "aaa"
2. 剩下 "b"，期望匹配 "b"
3. 成功！

文本: "aaac"
模式: "a+b"

1. a+ 贪婪匹配 "aaa"
2. 剩下 "c"，期望匹配 "b"
3. 失败，回溯...
4. a+ 匹配 "aa"
5. 剩下 "ac"，期望匹配 "b"
6. 失败，继续回溯...
7. 最终失败（无论如何都匹配不到 "b"）
```

### 4. 空匹配

`*` 和 `?` 可以匹配零次，需要特别处理：

```
模式: "a*"
文本: "bbb"

结果: 匹配成功（a 出现 0 次）
```

## 实际应用场景

### 1. 匹配整数

```bash
echo "count: 123" | go run . "\\d+"
# 输出: count: 123
```

### 2. 匹配标识符

```bash
echo "var_name_123" | go run . "[a-zA-Z_][a-zA-Z0-9_]*"
# 输出: var_name_123
```

### 3. 匹配文件扩展名

```bash
echo -e "file.txt\nfile.tar.gz\nfile" | go run . "\\..+"
# 输出:
# file.txt
# file.tar.gz
```

### 4. 匹配可选前缀

```bash
echo -e "http://example.com\nhttps://example.com" | go run . "https?://"
# 输出: 两行都匹配
```

## 性能考虑

### 1. 回溯爆炸

避免嵌套量词导致的指数级回溯：

```
# 危险模式
(a+)+      # 嵌套量词
a*a*a*     # 多个相邻量词

# 对于某些输入可能导致性能问题
```

### 2. 优化策略

- 首字符优化：如果模式首字符固定，可以快速跳过
- 锚点优化：`^` 开头的模式只需在位置 0 匹配

## 常见错误

### 1. 忘记转义

```bash
# 错误：? 被当作量词
echo "what?" | grep "what?"

# 正确：转义
echo "what?" | grep "what\\?"
```

### 2. 贪婪匹配导致意外结果

```
文本: "<b>bold</b> and <i>italic</i>"
模式: "<.+>"

贪婪结果: "<b>bold</b> and <i>italic</i>"（整个文本）
期望结果: "<b>", "</b>", "<i>", "</i>"

解决: 使用非贪婪量词 <.+?>（Day 10 实现）
```

## 下一步

完成本课后，你应该能够：
- ✓ 理解量词的概念和工作原理
- ✓ 实现 `?`（零或一次）
- ✓ 实现 `+`（一或多次）
- ✓ 实现 `*`（零或多次）
- ✓ 理解贪婪匹配和回溯机制

**下一课预告**：第九天我们将实现更精确的量词 `{n}`、`{n,}` 和 `{n,m}`。

## 参考资料

- [Regular Expressions - Quantifiers](https://www.regular-expressions.info/repeat.html)
- [Backtracking](https://www.regular-expressions.info/catastrophic.html)
- [Greedy vs Lazy](https://www.regular-expressions.info/repeat.html#greedy)
