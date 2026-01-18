# Day 5: 字符组 - 示例代码

## 项目结构

```
day5-character-groups/
├── main.go       # 主程序入口
├── args.go       # 命令行参数解析
├── regex.go      # 正则表达式解析器（支持字符组）
├── engine.go     # 匹配引擎（支持字符组匹配）
├── output.go     # 输出格式化器
├── test.sh       # 测试脚本
├── go.mod        # Go 模块文件
└── README.md     # 本文档
```

## 新增功能

### 1. 字符枚举 `[abc]`

匹配方括号内的任意一个字符：

```bash
echo -e "cat\nhat\nbat\nrat" | go run . "[ch]at"
# 输出:
# cat
# hat
```

### 2. 字符范围 `[a-z]`

匹配指定范围内的任意字符：

```bash
# 小写字母
echo -e "hello\nHELLO\n12345" | go run . "[a-z]"
# 输出: hello

# 大写字母
echo -e "ABC\nabc" | go run . "[A-Z]"
# 输出: ABC

# 数字
echo -e "abc\n123" | go run . "[0-9]"
# 输出: 123
```

### 3. 否定字符组 `[^abc]`

匹配不在方括号内的任意字符：

```bash
echo -e "abc\nxyz\naxy" | go run . "[^abc]"
# 输出:
# xyz
# axy

echo -e "123\nabc" | go run . "[^0-9]"
# 输出: abc
```

### 4. 组合字符组

可以组合多个范围和字符：

```bash
# 字母（大小写）
echo -e "abc\nABC\n123" | go run . "[a-zA-Z]"
# 输出:
# abc
# ABC

# 字母数字
echo -e "abc\n123\n@#$" | go run . "[a-zA-Z0-9]"
# 输出:
# abc
# 123

# 十六进制
echo -e "0xff\n0xGG\n0x1a" | go run . "0x[0-9a-fA-F][0-9a-fA-F]"
# 输出:
# 0xff
# 0x1a
```

## 使用示例

### 基础用法

```bash
# 匹配元音字母
echo -e "apple\nsky\norange" | go run . "[aeiou]"
# 输出: apple, orange

# 匹配辅音字母（非元音）
echo -e "aaa\nbcd\neee" | go run . "[^aeiou]"
# 输出: bcd

# 匹配数字开头的行
echo -e "1abc\nabc\n2xyz" | go run . "[0-9]"
# 输出: 1abc, 2xyz
```

### 实际应用场景

#### 1. 匹配 ID 格式

```bash
echo "ID:123" | go run . "ID:[0-9][0-9][0-9]"
# 输出: ID:123
```

#### 2. 匹配十六进制颜色

```bash
echo "#ff0000" | go run . "#[0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F][0-9a-fA-F]"
# 输出: #ff0000
```

#### 3. 匹配简单用户名

```bash
echo -e "user_123\nuser@name\nuser.name" | go run . "[a-zA-Z_][a-zA-Z0-9_]"
# 输出: user_123
```

#### 4. 过滤非数字行

```bash
echo -e "abc\n123\nxyz\n456" | go run . -v "[a-zA-Z]"
# 输出:
# 123
# 456
```

### 与命令行选项组合

```bash
# -n: 显示行号
echo -e "abc\n123\nxyz" | go run . -n "[0-9]"
# 输出: 2:123

# -v: 反向匹配
echo -e "abc\n123\nxyz" | go run . -v "[0-9]"
# 输出:
# abc
# xyz

# -c: 统计数量
echo -e "a1\nb2\nc3\nde" | go run . -c "[0-9]"
# 输出: 3
```

## 运行测试

```bash
# 运行完整测试套件
bash test.sh

# 手动测试
echo -e "cat\nhat\nbat" | go run . "[ch]at"
```

## 实现细节

### 1. Token 类型扩展

```go
const (
    // ...
    TokenCharGroup  // 字符组 [abc] 或 [a-z]
)

type CharGroup struct {
    Negated bool    // 是否为否定字符组
    Chars   []byte  // 单个字符列表
    Ranges  []Range // 字符范围列表
}

type Range struct {
    Start byte
    End   byte
}
```

### 2. 字符组解析

解析 `[a-zA-Z0-9]`：

```go
func parseCharGroup(pattern string, pos int) (Token, int, error) {
    // 1. 检查是否为否定 [^...]
    // 2. 解析字符和范围
    // 3. 返回 Token
}
```

### 3. 字符组匹配

```go
func (cg *CharGroup) Match(char byte) bool {
    contains := cg.Contains(char)
    if cg.Negated {
        return !contains
    }
    return contains
}

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
```

## 字符组语法规则

### 1. 基本语法

```
[abc]       匹配 a、b 或 c
[a-z]       匹配 a 到 z
[A-Z]       匹配 A 到 Z
[0-9]       匹配 0 到 9
[^abc]      匹配除 a、b、c 外的任意字符
[^a-z]      匹配除小写字母外的任意字符
```

### 2. 组合规则

```
[a-zA-Z]        字母（大小写）
[a-zA-Z0-9]     字母和数字
[a-zA-Z0-9_]    等价于 \w
[aeiou]         元音字母
[^aeiou]        辅音字母
```

### 3. 特殊字符处理

```
[.]         匹配字面点号（不需要转义）
[-a]        - 在开头是字面字符
[a-]        - 在结尾是字面字符
[a\-z]      转义的 -
[\]]        转义的 ]
```

## 与 Day 4 的区别

| 特性 | Day 4 | Day 5 |
|------|-------|-------|
| 预定义字符类 | ✅ `\d`, `\w`, `\s` | ✅ |
| 字符枚举 | ❌ | ✅ `[abc]` |
| 字符范围 | ❌ | ✅ `[a-z]` |
| 否定字符组 | ❌ | ✅ `[^abc]` |
| 自定义字符集 | ❌ | ✅ |

### 对比示例

```bash
# Day 4: 使用预定义字符类
echo "abc" | day4/go run . "\\d"  # 匹配数字

# Day 5: 可以自定义匹配范围
echo "abc" | day5/go run . "[0-5]"  # 只匹配 0-5
echo "abc" | day5/go run . "[aeiou]"  # 只匹配元音
```

## 字符组 vs 预定义字符类

| 字符组 | 预定义字符类 | 说明 |
|--------|-------------|------|
| `[0-9]` | `\d` | 数字 |
| `[a-zA-Z0-9_]` | `\w` | 单词字符 |
| `[ \t\n\r]` | `\s` | 空白字符 |
| `[^0-9]` | `\D` | 非数字 |
| `[^a-zA-Z0-9_]` | `\W` | 非单词字符 |

**字符组的优势**：
- 可以自定义字符集合
- 更灵活的匹配范围
- 可以组合多个范围

## 调试技巧

### 1. 测试单个字符组

```bash
# 测试 [a-z]
echo "a" | go run . "[a-z]"  # 应该匹配
echo "A" | go run . "[a-z]"  # 不应该匹配
echo "1" | go run . "[a-z]"  # 不应该匹配

# 测试 [^a-z]
echo "1" | go run . "[^a-z]"  # 应该匹配
echo "a" | go run . "[^a-z]"  # 不应该匹配
```

### 2. 查看解析结果

在 `regex.go` 中添加调试输出：

```go
pattern, _ := ParsePattern("[a-zA-Z]")
for _, token := range pattern.Tokens() {
    fmt.Println(token)
}
// 输出: CharGroup([...], [{a z} {A Z}])
```

## 限制和注意事项

1. **ASCII 限制**: 只支持 ASCII 字符范围
2. **POSIX 类**: 不支持 `[[:alpha:]]` 等 POSIX 字符类
3. **嵌套**: 不支持字符组嵌套
4. **Unicode**: 不支持 Unicode 字符范围

这些限制将在后续课程中改进。

## 性能考虑

### 当前实现

```go
// O(n) 查找
func (cg *CharGroup) Contains(char byte) bool {
    for _, c := range cg.Chars { ... }
    for _, r := range cg.Ranges { ... }
}
```

### 可选优化：位图

```go
// O(1) 查找
type CharGroup struct {
    bitmap [256]bool
}

func (cg *CharGroup) Contains(char byte) bool {
    return cg.bitmap[char]
}
```

## 下一步

完成本课后，你应该掌握了：
- ✓ 字符组的概念和语法
- ✓ 字符枚举和范围的解析
- ✓ 否定字符组的实现
- ✓ 组合字符组的使用
- ✓ 字符组与已有功能的组合

**下一课预告**: Day 6 将学习组合字符类，实现更复杂的匹配模式。

## 常见问题

### Q: `[a-z]` 和 `\w` 有什么区别？
A: `[a-z]` 只匹配小写字母，而 `\w` 匹配所有字母、数字和下划线。

### Q: 如何匹配字面的 `-`？
A: 将 `-` 放在开头或结尾：`[-a]` 或 `[a-]`，或者转义：`[a\-z]`。

### Q: 如何匹配字面的 `]`？
A: 使用转义：`[\]]`，或将 `]` 放在开头：`[]a]`（某些实现支持）。

### Q: 字符组可以嵌套吗？
A: 当前实现不支持嵌套，如 `[[a-z]0-9]`。
