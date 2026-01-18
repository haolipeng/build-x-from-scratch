# Day 7: 锚点 - 示例代码

## 项目结构

```
day7-anchors/
├── main.go       # 主程序入口
├── args.go       # 命令行参数解析
├── regex.go      # 正则表达式解析器（支持锚点）
├── engine.go     # 匹配引擎（锚点匹配逻辑）
├── output.go     # 输出格式化器
├── test.sh       # 测试脚本
├── go.mod        # Go 模块文件
└── README.md     # 本文档
```

## 新增功能

### 1. 行首锚点 `^`

匹配行的开始位置：

```bash
echo -e "hello world\nworld hello\nhello" | go run . "^hello"
# 输出:
# hello world
# hello
```

### 2. 行尾锚点 `$`

匹配行的结束位置：

```bash
echo -e "hello world\nworld hello\nworld" | go run . "world$"
# 输出:
# hello world
# world
```

### 3. 同时使用 `^` 和 `$`

匹配整行：

```bash
echo -e "hello\nhello world\nworld" | go run . "^hello$"
# 输出:
# hello

echo -e "a\nab\nabc" | go run . "^..$"
# 输出:
# ab
```

### 4. 单词边界 `\b`

匹配单词的开始或结束位置：

```bash
echo -e "the cat\nthere\nother" | go run . "\\bthe\\b"
# 输出:
# the cat

echo -e "cat\ncatalog\ncat food" | go run . "\\bcat\\b"
# 输出:
# cat
# cat food
```

### 5. 非单词边界 `\B`

匹配非单词边界位置：

```bash
echo -e "there\nthe\nother" | go run . "\\Bthe"
# 输出:
# other

echo -e "catalog\ncat\nconcat" | go run . "cat\\B"
# 输出:
# catalog
```

## 使用示例

### 基础用法

```bash
# 匹配以 hello 开头的行
echo -e "hello world\nworld hello" | go run . "^hello"

# 匹配以 world 结尾的行
echo -e "hello world\nworld hello" | go run . "world$"

# 匹配空行
echo -e "hello\n\nworld" | go run . "^$"

# 匹配独立的单词 "is"
echo -e "this is\nisland" | go run . "\\bis\\b"
```

### 实际应用

#### 1. 查找以 # 开头的注释行

```bash
echo -e "# comment\ncode # comment\n  # indented" | go run . "^#"
# 输出:
# # comment
```

#### 2. 查找以数字结尾的行

```bash
echo -e "abc\nabc123\n123abc" | go run . "\\d$"
# 输出:
# abc123
```

#### 3. 验证整行格式

```bash
echo -e "hello\nhello world" | go run . "^hello$"
# 输出: hello（只匹配完整的 "hello"）
```

#### 4. 查找独立单词

```bash
echo -e "the cat\nthere\nother" | go run . "\\bthe\\b"
# 输出: the cat
```

### 与命令行选项组合

```bash
# -n: 显示行号
echo -e "hello\nworld\nhello world" | go run . -n "^hello"
# 输出:
# 1:hello
# 3:hello world

# -v: 反向匹配
echo -e "hello\nworld\nhello world" | go run . -v "world$"
# 输出:
# hello

# -c: 统计数量
echo -e "hello\nhello world\nworld hello" | go run . -c "^hello"
# 输出:
# 2
```

## 运行测试

```bash
# 运行完整测试套件
bash test.sh

# 手动测试
echo -e "hello\nworld" | go run . "^hello"
```

## 实现细节

### 1. Token 类型扩展

```go
const (
    // ... 已有类型
    TokenStartAnchor     // ^ 行首
    TokenEndAnchor       // $ 行尾
    TokenWordBoundary    // \b 单词边界
    TokenNotWordBoundary // \B 非单词边界
)
```

### 2. 锚点解析

```go
func ParsePattern(pattern string) (*Pattern, error) {
    // ...
    case '^':
        tokens = append(tokens, Token{Type: TokenStartAnchor})
    case '$':
        tokens = append(tokens, Token{Type: TokenEndAnchor})
    // ...
}

func parseEscape(pattern string, pos int) (Token, int, error) {
    // ...
    case 'b':
        return Token{Type: TokenWordBoundary}, pos + 2, nil
    case 'B':
        return Token{Type: TokenNotWordBoundary}, pos + 2, nil
    // ...
}
```

### 3. 锚点匹配逻辑

锚点的关键特点是**不消耗字符**，只检查位置：

```go
func matchTokens(tokens, text, patIdx, textIdx) bool {
    token := tokens[patIdx]

    switch token.Type {
    case TokenStartAnchor:
        // ^ 只在位置 0 匹配
        if textIdx == 0 {
            return matchTokens(tokens, text, patIdx+1, textIdx) // 不移动 textIdx
        }
        return false

    case TokenEndAnchor:
        // $ 只在文本末尾匹配
        if textIdx == len(text) {
            return matchTokens(tokens, text, patIdx+1, textIdx)
        }
        return false

    case TokenWordBoundary:
        if isWordBoundary(text, textIdx) {
            return matchTokens(tokens, text, patIdx+1, textIdx)
        }
        return false
    // ...
    }
}
```

### 4. 单词边界判断

```go
func isWordBoundary(text string, pos int) bool {
    prevIsWord := pos > 0 && isWordChar(text[pos-1])
    currIsWord := pos < len(text) && isWordChar(text[pos])
    // 单词边界：一边是单词字符，另一边不是
    return prevIsWord != currIsWord
}
```

### 5. 性能优化

当模式以 `^` 开头时，只需在位置 0 尝试匹配：

```go
func (m *RegexMatcher) Match(text string) bool {
    if m.pattern.StartsWithAnchor() {
        return m.matchAt(text, 0) // 只在位置 0 尝试
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

## 锚点工作原理

### 行首锚点 `^`

```
文本: "hello world"
模式: "^hello"

位置: 0 1 2 3 4 5 ...
      h e l l o   w ...
      ^
      └── ^ 匹配位置 0（开始）

结果: 匹配成功
```

### 行尾锚点 `$`

```
文本: "hello world"
模式: "world$"

位置: 0 1 2 3 4 5 6 7 8 9 10 [END]
      h e l l o   w o r l d   $
                              ^
                              └── $ 匹配位置 11（结束）

结果: 匹配成功
```

### 单词边界 `\b`

```
文本: "the cat"
模式: "\bthe\b"

位置: [START] t h e   c a t [END]
      ^       ^   ^
      |       |   └── \b（'e' 后是空格）
      |       └── 匹配 'the'
      └── \b（开头是单词字符）

结果: 匹配成功
```

## 与 Day 6 的区别

| 特性 | Day 6 | Day 7 |
|------|-------|-------|
| 组合字符类 | ✅ | ✅ |
| 行首锚点 `^` | ❌ | ✅ |
| 行尾锚点 `$` | ❌ | ✅ |
| 单词边界 `\b` | ❌ | ✅ |
| 非单词边界 `\B` | ❌ | ✅ |
| 位置匹配 | ❌ | ✅ |

### 对比示例

```bash
# Day 6: 匹配任意位置的 hello
echo "world hello" | day6/go run . "hello"  # 匹配

# Day 7: 只匹配行首的 hello
echo "world hello" | day7/go run . "^hello"  # 不匹配
echo "hello world" | day7/go run . "^hello"  # 匹配
```

## `^` 的两种含义

```bash
# 行首锚点（在模式开头）
echo "hello" | go run . "^h"     # 匹配行首的 h

# 否定字符组（在 [...] 内部）
echo "hello" | go run . "[^aeiou]"  # 匹配非元音字符
```

## 调试技巧

### 测试锚点

```bash
# 测试 ^
echo "hello" | go run . "^hello"   # 应该匹配
echo " hello" | go run . "^hello"  # 不应该匹配

# 测试 $
echo "hello" | go run . "hello$"   # 应该匹配
echo "hello " | go run . "hello$"  # 不应该匹配

# 测试 \b
echo "the" | go run . "\\bthe\\b"    # 应该匹配
echo "there" | go run . "\\bthe\\b"  # 不应该匹配
```

## 限制和注意事项

1. **多行模式**: 当前实现按行处理，`^` 和 `$` 匹配每行的开始和结束
2. **Unicode**: 单词边界只考虑 ASCII 字符
3. **性能**: `^` 开头的模式已优化，其他锚点暂未优化

## 下一步

完成本课后，你应该掌握了：
- ✓ 锚点的概念和工作原理
- ✓ 行首锚点 `^` 的实现
- ✓ 行尾锚点 `$` 的实现
- ✓ 单词边界 `\b` 和 `\B` 的实现
- ✓ 锚点不消耗字符的特性

**下一课预告**: Day 8 将实现量词：
- `?`: 零或一次
- `+`: 一或多次
- `*`: 零或多次

## 常见问题

### Q: `^` 和 `$` 是否匹配换行符？
A: 不，它们匹配位置而不是字符。grep 按行处理，所以 `^` 匹配行首，`$` 匹配行尾。

### Q: `\b` 如何判断单词边界？
A: 当一边是单词字符（`\w`），另一边不是时，就是单词边界。

### Q: 为什么 `^hello$` 只匹配 "hello" 而不匹配 "hello world"？
A: 因为 `^` 要求在行首，`$` 要求在行尾，整个模式必须匹配完整的行。

### Q: `\B` 有什么用？
A: 用于匹配在单词中间的位置，如 `\Bcat` 匹配 "scatter" 中的 "cat"。
