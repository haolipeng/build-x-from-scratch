# Day 2: 正则表达式基础 - 示例代码

## 项目结构

```
day2-regex-basics/
├── main.go       # 主程序入口
├── args.go       # 命令行参数解析
├── regex.go      # 正则表达式解析器
├── engine.go     # 匹配引擎
├── test.sh       # 测试脚本
├── go.mod        # Go 模块文件
└── README.md     # 本文档
```

## 核心组件

### regex.go - 正则表达式解析器
定义了正则表达式的内部表示：
- `TokenType`：定义 token 类型（字面字符、通配符）
- `Token`：表示模式中的一个元素
- `Pattern`：表示完整的正则表达式模式
- `ParsePattern()`：将字符串解析为 Pattern 对象

### engine.go - 匹配引擎
实现核心匹配算法：
- `RegexMatcher`：正则表达式匹配器
- `Match()`：在文本中查找匹配（包含匹配）
- `MatchFull()`：完整匹配整个文本
- `matchAt()`：从指定位置开始匹配
- `matchTokens()`：递归匹配 token 序列

### args.go
命令行参数解析（与 Day 1 相同）

### main.go
程序入口，集成正则表达式引擎

## 支持的正则表达式特性

### 1. 字面字符
直接匹配字符本身：
```bash
echo "hello" | go run . "hello"
# 输出: hello

echo "hello world" | go run . "world"
# 输出: hello world
```

### 2. 通配符 `.`
匹配任意单个字符：
```bash
echo -e "cat\ncot\ncut" | go run . "c.t"
# 输出:
# cat
# cot
# cut

echo "hello" | go run . "h.llo"
# 输出: hello
```

### 3. 混合使用
字面字符和通配符可以组合：
```bash
echo "pattern123test" | go run . "pattern...test"
# 输出: pattern123test

echo -e "hello\nhallo\nhxllo" | go run . "h.llo"
# 输出:
# hello
# hallo
# hxllo
```

## 运行示例

### 基础用法

```bash
# 从标准输入搜索
echo "hello world" | go run . "h.llo"

# 从文件搜索
go run . "pattern" file.txt

# 多文件搜索
go run . "pattern" file1.txt file2.txt
```

### 示例 1：匹配三字母单词

```bash
echo -e "cat\ndog\nelephant\nfox\nbird" | go run . "..."
```
输出:
```
cat
dog
elephant  # 包含3个连续字符
fox
```

### 示例 2：匹配电话号码模式

```bash
echo -e "123-4567\n456-7890\nabcd-efgh" | go run . "...-......"
```
输出:
```
123-4567
456-7890
abcd-efgh
```

### 示例 3：匹配邮箱前缀

```bash
echo -e "user@example.com\nadmin@test.com\ntest@domain.org" | go run . "....@"
```
输出:
```
user@example.com
admin@test.com  # "admin" 有5个字符，也包含4个字符的子串
```

### 示例 4：验证长度

```bash
# 匹配至少5个字符的行
echo -e "ab\nabc\nabcde\nabcdef" | go run . "....."
```
输出:
```
abcde
abcdef
```

## 运行测试

```bash
# 运行完整测试套件
bash test.sh

# 或者手动测试各个功能
echo -e "cat\ncot\ncut" | go run . "c.t"
```

## 工作原理

### 1. 模式解析

输入: `"h.llo"`

解析为 Token 序列:
```
[
  Token{Type: Literal, Value: 'h'},
  Token{Type: Wildcard},
  Token{Type: Literal, Value: 'l'},
  Token{Type: Literal, Value: 'l'},
  Token{Type: Literal, Value: 'o'}
]
```

### 2. 匹配过程

文本: `"hello world"`
模式: `"h.llo"`

```
尝试位置 0: "hello"
  h == h ✓
  e == . ✓ (通配符匹配任意字符)
  l == l ✓
  l == l ✓
  o == o ✓
  匹配成功！
```

### 3. 包含匹配 vs 完整匹配

**包含匹配**（grep 默认）:
```go
// 尝试从每个位置开始匹配
for i := 0; i <= len(text); i++ {
    if matchAt(text, i) {
        return true
    }
}
```

**完整匹配**:
```go
// 从位置 0 开始，且必须匹配整个文本
return matchAt(text, 0) && len(pattern) == len(text)
```

## 与 Day 1 的对比

| 特性 | Day 1 | Day 2 |
|------|-------|-------|
| 匹配类型 | 字符串包含 | 正则表达式 |
| 通配符 | 不支持 | 支持 `.` |
| 实现方式 | `strings.Contains()` | 递归匹配算法 |
| 灵活性 | 低 | 高 |
| 可扩展性 | 难 | 易 |

### 示例对比

```bash
# Day 1: 只能精确匹配字面 "h.llo"
echo "h.llo" | day1/go run . "h.llo"  # ✓ 匹配
echo "hello" | day1/go run . "h.llo"  # ✗ 不匹配

# Day 2: "." 被解释为通配符
echo "h.llo" | day2/go run . "h.llo"  # ✓ 匹配
echo "hello" | day2/go run . "h.llo"  # ✓ 匹配（. 匹配 e）
```

## 算法复杂度

### 时间复杂度
- **最坏情况**: O(n × m)
  - n = 文本长度
  - m = 模式长度
  - 需要在文本的每个位置尝试匹配

- **最好情况**: O(n)
  - 第一个字符不匹配时可以快速跳过

### 空间复杂度
- **递归深度**: O(m)（模式长度）
- **Token 数组**: O(m)

## 限制和注意事项

1. **换行符处理**: 当前实现中，`.` 匹配任意字符包括换行符（简化实现）
2. **转义未实现**: 无法匹配字面意义的 `.`
3. **Unicode**: 当前只处理 ASCII 字符，多字节字符可能有问题
4. **性能**: 对于长文本和复杂模式，性能不是最优的

这些限制将在后续课程中逐步改进。

## 下一步

完成本课后，你应该掌握了：
- ✓ 正则表达式的基本概念
- ✓ 如何设计 Token 系统
- ✓ 递归匹配算法的实现
- ✓ 包含匹配 vs 完整匹配的区别

**下一课预告**: Day 3 将添加命令行选项支持，包括：
- `-n`: 显示行号
- `-i`: 忽略大小写
- `-v`: 反向匹配（显示不匹配的行）
- `-c`: 统计匹配数量

## 扩展练习

1. **添加转义支持**: 实现 `\.` 来匹配字面的点
2. **优化性能**: 实现首字符快速跳过
3. **添加调试模式**: 显示匹配过程
4. **完整匹配模式**: 添加一个选项来启用完整匹配

示例调试输出:
```
Pattern: "h.llo"
Text: "hello world"

Position 0: "hello"
  Match 'h' with 'h' ✓
  Match '.' with 'e' ✓
  Match 'l' with 'l' ✓
  Match 'l' with 'l' ✓
  Match 'o' with 'o' ✓
  ✓ MATCH at position 0-5
```

## 常见问题

### Q: 为什么使用递归而不是循环？
A: 递归更容易理解和实现。在后续课程中，我们会优化为迭代实现以提高性能。

### Q: 通配符能匹配换行符吗？
A: 在我们的简化实现中可以。标准 grep 中，`.` 默认不匹配换行符。

### Q: 如何匹配字面的点号？
A: 当前版本不支持。Day 4 会添加转义字符支持，使用 `\.` 匹配字面点号。

### Q: 性能如何？
A: 对于简单模式已经足够快。对于复杂模式和大文件，后续会优化。
