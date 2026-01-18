# Day 6: 组合字符类 - 示例代码

## 项目结构

```
day6-combined-classes/
├── main.go       # 主程序入口
├── args.go       # 命令行参数解析
├── regex.go      # 正则表达式解析器（支持组合字符类）
├── engine.go     # 匹配引擎
├── output.go     # 输出格式化器
├── test.sh       # 测试脚本
├── go.mod        # Go 模块文件
└── README.md     # 本文档
```

## 新增功能

### 1. 字符组内支持字符类

Day 6 支持在字符组 `[...]` 内使用预定义字符类：

```bash
# [\d\s]: 匹配数字或空白
echo -e "abc\n123\n   " | go run . "[\d\s]"
# 输出: 123, (空格行)

# [\w-]: 匹配单词字符或连字符
echo -e "hello-world\nhello@world" | go run . "hello[\w-]world"
# 输出: hello-world

# [a-z\d]: 匹配小写字母或数字
echo -e "abc\n123\nABC" | go run . "[a-z\d]"
# 输出: abc, 123
```

### 2. 否定组合字符类

```bash
# [^\d]: 匹配非数字
echo -e "123\nabc" | go run . "[^\d]"
# 输出: abc

# [^\w\s]: 匹配非单词字符且非空白
echo -e "hello\nhello@world" | go run . "[^\w\s]"
# 输出: hello@world
```

### 3. 复杂组合

```bash
# 匹配十六进制
echo "0xff" | go run . "0x[\da-fA-F][\da-fA-F]"
# 输出: 0xff

# 匹配数字或点号（浮点数）
echo "3.14" | go run . "[\d.][\d.][\d.][\d.]"
# 输出: 3.14
```

## 使用示例

### 基础用法

```bash
# 数字或空白
echo -e "abc\n123\n   " | go run . "[\d\s]"

# 单词字符或连字符
echo "my-var_name" | go run . "[\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-]"

# 字母或数字
echo "abc123" | go run . "[a-zA-Z\d]"
```

### 实际应用

#### 1. 匹配标识符（带连字符）

```bash
echo -e "my-var\nmy_var\nmy var" | go run . "my[\w-]var"
# 输出:
# my-var
# my_var
```

#### 2. 匹配浮点数

```bash
echo "pi = 3.14159" | go run . "[\d.][\d.][\d.][\d.][\d.][\d.][\d.]"
# 输出: pi = 3.14159
```

#### 3. 过滤特殊字符

```bash
echo -e "hello\nhello@world\nhello_world" | go run . "[^\w\s]"
# 输出: hello@world
```

### 与命令行选项组合

```bash
# -n: 显示行号
echo -e "abc\n123\nxyz" | go run . -n "[\d]"
# 输出: 2:123

# -v: 反向匹配
echo -e "abc\n123\nxyz" | go run . -v "[\d]"
# 输出:
# abc
# xyz

# -c: 统计数量
echo -e "a1\nb2\nc3\nde" | go run . -c "[\d]"
# 输出: 3
```

## 运行测试

```bash
# 运行完整测试套件
bash test.sh

# 手动测试
echo -e "abc\n123" | go run . "[\d\s]"
```

## 实现细节

### 1. CharGroup 结构扩展

```go
type CharGroup struct {
    Negated     bool        // 是否为否定字符组
    Chars       []byte      // 单个字符
    Ranges      []Range     // 字符范围
    CharClasses []TokenType // 字符类列表（新增）
}
```

### 2. 字符组解析扩展

```go
func parseCharGroup(pattern string, pos int) (...) {
    // ...
    if ch == '\\' && pos+1 < len(pattern) {
        next := pattern[pos+1]
        switch next {
        case 'd':
            charClasses = append(charClasses, TokenDigit)
        case 'w':
            charClasses = append(charClasses, TokenWord)
        case 's':
            charClasses = append(charClasses, TokenSpace)
        // ...
        }
    }
    // ...
}
```

### 3. 匹配逻辑扩展

```go
func (cg *CharGroup) Contains(char byte) bool {
    // 检查单个字符
    for _, c := range cg.Chars { ... }

    // 检查范围
    for _, r := range cg.Ranges { ... }

    // 检查字符类（新增）
    for _, class := range cg.CharClasses {
        if matchCharClass(class, char) {
            return true
        }
    }

    return false
}
```

## 支持的组合模式

| 模式 | 说明 | 示例 |
|------|------|------|
| `[\d\s]` | 数字或空白 | 匹配 `1`, ` `, `\t` |
| `[\w-]` | 单词字符或连字符 | 匹配 `a`, `1`, `-` |
| `[a-z\d]` | 小写字母或数字 | 匹配 `a`, `1` |
| `[\d.]` | 数字或点号 | 匹配 `1`, `.` |
| `[^\d]` | 非数字 | 匹配 `a`, `@` |
| `[^\w\s]` | 非单词非空白 | 匹配 `@`, `!` |

## 与 Day 5 的区别

| 特性 | Day 5 | Day 6 |
|------|-------|-------|
| 字符枚举 `[abc]` | ✅ | ✅ |
| 字符范围 `[a-z]` | ✅ | ✅ |
| 否定字符组 `[^abc]` | ✅ | ✅ |
| 字符组内字符类 `[\d\s]` | ❌ | ✅ |
| 组合模式 `[a-z\d]` | ❌ | ✅ |

### 对比示例

```bash
# Day 5: 只能手动枚举
echo "123" | day5/go run . "[0-9 \t]"  # 需要手动列出空白字符

# Day 6: 可以使用字符类
echo "123" | day6/go run . "[\d\s]"  # 更简洁
```

## 字符类在字符组中的等价关系

| 字符类 | 等价于 |
|--------|--------|
| `[\d]` | `[0-9]` |
| `[\w]` | `[a-zA-Z0-9_]` |
| `[\s]` | `[ \t\n\r]` |
| `[^\d]` | `\D` |
| `[^\w]` | `\W` |
| `[^\s]` | `\S` |

## 调试技巧

### 查看解析结果

```go
pattern, _ := ParsePattern("[\d\s]")
for _, token := range pattern.Tokens() {
    fmt.Println(token)
}
// 输出: CharGroup(chars=[], ranges=[], classes=[TokenDigit TokenSpace])
```

### 测试单个组合

```bash
# 测试 [\d\s]
echo "1" | go run . "[\d\s]"  # 应该匹配
echo " " | go run . "[\d\s]"  # 应该匹配
echo "a" | go run . "[\d\s]"  # 不应该匹配
```

## 限制和注意事项

1. **嵌套字符类**: 不支持 `[[\d][\w]]` 这样的嵌套
2. **字符类运算**: 不支持交集 `[a-z&&[^aeiou]]`
3. **Unicode**: 只支持 ASCII 字符类
4. **POSIX 类**: 不支持 `[[:digit:]]`

## 性能考虑

### 字符类匹配

```go
// 每个字符类都是 O(1) 的检查
func matchCharClass(tokenType TokenType, char byte) bool {
    switch tokenType {
    case TokenDigit:
        return char >= '0' && char <= '9'  // O(1)
    // ...
    }
}
```

### 总体复杂度

- 检查单个字符: O(c + r + k)
  - c = 单个字符数量
  - r = 范围数量
  - k = 字符类数量

## 下一步

完成本课后，你应该掌握了：
- ✓ 在字符组内使用预定义字符类
- ✓ 创建复杂的组合匹配模式
- ✓ 理解字符类在字符组中的行为
- ✓ 应用组合字符类解决实际问题

**下一课预告**: Day 7 将实现锚点：
- `^`: 行首锚点
- `$`: 行尾锚点
- `\b`: 单词边界

## 常见问题

### Q: `[\d\w]` 和 `[\w]` 有什么区别？
A: 没有区别，因为 `\w` 已经包含了 `\d`。`[\d\w]` 是冗余的写法。

### Q: 如何匹配字面的 `\d`？
A: 使用 `\\d` 或在字符组中使用 `[\\d]`（但这会被解析为字符类）。要匹配字面的反斜杠和 d，需要 `\\\\d`。

### Q: 字符类的顺序重要吗？
A: 不重要。`[\d\s]` 和 `[\s\d]` 是等价的。

### Q: 为什么 `[^\d\s]` 匹配字母？
A: 因为它匹配"既不是数字也不是空白"的字符，字母满足这个条件。
