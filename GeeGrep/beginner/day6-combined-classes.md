# 第六天：组合字符类

## 学习目标

在这个阶段，你将学习如何：
1. 在字符组中使用预定义字符类（如 `[\d\s]`）
2. 组合多种匹配模式
3. 理解字符类的优先级和组合规则
4. 实现更复杂的匹配模式

## 什么是组合字符类？

组合字符类是指将预定义字符类（`\d`、`\w`、`\s`）与字符组（`[...]`）结合使用，创建更灵活的匹配模式。

### 基本语法

```
[\d\s]      匹配数字或空白字符
[\w-]       匹配单词字符或连字符
[a-z\d]     匹配小写字母或数字
[\d.]       匹配数字或点号
[^\d]       匹配非数字（等价于 \D）
[^\w\s]     匹配非单词字符且非空白的字符
```

### 使用示例

```bash
# 匹配数字或空格
echo "hello 123 world" | grep "[\d\s]"

# 匹配标识符（字母、数字、下划线、连字符）
echo "my-var_name123" | grep "[\w-]"

# 匹配非字母数字字符
echo "hello@world!" | grep "[^\w]"
```

## 本课实现内容

### 1. 字符组内支持字符类

扩展字符组解析，支持在 `[...]` 内使用 `\d`、`\w`、`\s` 等：

```
输入: "[\d\s]"

解析为:
CharGroup {
    CharClasses: [TokenDigit, TokenSpace]
    Chars: []
    Ranges: []
}
```

### 2. 扩展 CharGroup 结构

```go
type CharGroup struct {
    Negated     bool        // 是否为否定字符组
    Chars       []byte      // 单个字符
    Ranges      []Range     // 字符范围
    CharClasses []TokenType // 字符类（\d, \w, \s 等）
}
```

### 3. 匹配逻辑扩展

```go
func (cg *CharGroup) Contains(char byte) bool {
    // 检查单个字符
    for _, c := range cg.Chars { ... }

    // 检查范围
    for _, r := range cg.Ranges { ... }

    // 检查字符类
    for _, class := range cg.CharClasses {
        if matchCharClass(class, char) {
            return true
        }
    }

    return false
}
```

## 测试用例

### 测试 1: [\d\s] 匹配数字或空白

```bash
echo -e "abc\n123\n \t" | go run . "[\d\s]"
# 输出:
# 123
#  	（空白行）
```

### 测试 2: [\w-] 匹配单词字符或连字符

```bash
echo -e "hello-world\nhello@world\nhello_world" | go run . "[\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-][\w-]"
# 输出:
# hello-world
# hello_world
```

### 测试 3: [a-z\d] 匹配小写字母或数字

```bash
echo -e "abc\nABC\n123\na1b" | go run . "[a-z\d]"
# 输出:
# abc
# 123
# a1b
```

### 测试 4: [^\d] 匹配非数字

```bash
echo -e "123\nabc\na1b" | go run . "[^\d]"
# 输出:
# abc
# a1b
```

### 测试 5: [\d.] 匹配数字或点号

```bash
echo -e "3.14\n3a14\n3.14.15" | go run . "[\d.][\d.][\d.][\d.]"
# 输出:
# 3.14
# 3.14.15
```

### 测试 6: 复杂组合

```bash
# 匹配简单的浮点数模式
echo "pi = 3.14159" | go run . "[\d]+\.[\d]+"
# 注意：量词 + 在后续课程实现
```

## 实现细节

### 1. 字符组内的转义处理

在字符组内部，`\d`、`\w`、`\s` 需要特殊处理：

```go
func parseCharGroupContent(pattern string, pos int) (...) {
    if ch == '\\' && pos+1 < len(pattern) {
        next := pattern[pos+1]
        switch next {
        case 'd':
            charClasses = append(charClasses, TokenDigit)
            pos += 2
        case 'w':
            charClasses = append(charClasses, TokenWord)
            pos += 2
        case 's':
            charClasses = append(charClasses, TokenSpace)
            pos += 2
        case 'D':
            charClasses = append(charClasses, TokenNotDigit)
            pos += 2
        // ...
        default:
            // 其他转义当作字面字符
            chars = append(chars, next)
            pos += 2
        }
    }
}
```

### 2. 否定组合字符类

`[^\d\s]` 匹配既不是数字也不是空白的字符：

```go
func (cg *CharGroup) Match(char byte) bool {
    contains := cg.Contains(char)
    if cg.Negated {
        return !contains
    }
    return contains
}
```

### 3. 字符类匹配辅助函数

```go
func matchCharClass(tokenType TokenType, char byte) bool {
    switch tokenType {
    case TokenDigit:
        return isDigit(char)
    case TokenWord:
        return isWordChar(char)
    case TokenSpace:
        return isSpace(char)
    case TokenNotDigit:
        return !isDigit(char)
    case TokenNotWord:
        return !isWordChar(char)
    case TokenNotSpace:
        return !isSpace(char)
    }
    return false
}
```

## 关键知识点

### 1. 字符类在字符组中的行为

在字符组 `[...]` 内部：
- `\d` 展开为 `0-9`
- `\w` 展开为 `a-zA-Z0-9_`
- `\s` 展开为空白字符集

```
[\da-f]   等价于 [0-9a-f]
[\w-]     等价于 [a-zA-Z0-9_-]
```

### 2. 否定字符类的组合

```
[^\d]     匹配非数字（等价于 \D）
[^\w]     匹配非单词字符（等价于 \W）
[^\d\s]   匹配既不是数字也不是空白的字符
```

### 3. 组合顺序不影响结果

```
[\d\s]    等价于 [\s\d]
[a-z\d]   等价于 [\da-z]
```

### 4. 避免重复定义

```
[\d0-9]   有效但冗余（\d 已包含 0-9）
[\w\d]    有效但冗余（\w 已包含 \d）
```

## 实际应用场景

### 1. 匹配文件名

```bash
# 文件名可以包含字母、数字、下划线、连字符、点号
echo "my-file_v1.0.txt" | grep "[\w.-]+"
```

### 2. 匹配 URL 路径

```bash
# URL 路径可以包含字母、数字、连字符、斜杠
echo "/api/v1/users-list" | grep "[/\w-]+"
```

### 3. 匹配配置值

```bash
# 配置值可以是数字或布尔值
echo "enabled = true" | grep "[\w]+ = [\w\d]+"
```

### 4. 过滤特殊字符

```bash
# 过滤非字母数字字符
echo "hello@world!" | grep "[^\w]"
```

## 与标准正则表达式的对比

| 特性 | 标准 Regex | 我们的实现 |
|------|-----------|------------|
| `[\d\s]` | ✅ | ✅ |
| `[\w-]` | ✅ | ✅ |
| `[^\d\w]` | ✅ | ✅ |
| 嵌套字符类 | 某些支持 | 不支持 |
| Unicode 属性 | `\p{L}` | 不支持 |

## 扩展思考

1. **性能优化**：如何高效匹配包含多个字符类的字符组？
2. **嵌套字符组**：如何支持 `[[a-z][0-9]]`？
3. **Unicode 支持**：如何支持 Unicode 字符类？
4. **字符类运算**：如何支持交集、差集？

## 下一步

完成本课后，你应该能够：
- ✓ 在字符组内使用预定义字符类
- ✓ 创建复杂的匹配模式
- ✓ 理解字符类组合的规则
- ✓ 应用组合字符类解决实际问题

**下一课预告**：第七天我们将实现锚点，支持 `^`（行首）、`$`（行尾）、`\b`（单词边界）等。

## 参考资料

- [Regular Expressions - Character Classes](https://www.regular-expressions.info/charclass.html)
- [Combining Character Classes](https://www.regular-expressions.info/charclassintersect.html)
