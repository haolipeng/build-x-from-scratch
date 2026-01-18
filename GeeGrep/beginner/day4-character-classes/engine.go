package main

import (
	"strings"
)

// RegexMatcher 实现基于正则表达式的匹配器
type RegexMatcher struct {
	pattern    *Pattern
	ignoreCase bool // 是否忽略大小写
}

// NewRegexMatcher 创建一个新的正则表达式匹配器
func NewRegexMatcher(pattern *Pattern) *RegexMatcher {
	return &RegexMatcher{
		pattern:    pattern,
		ignoreCase: false,
	}
}

// SetIgnoreCase 设置是否忽略大小写
func (m *RegexMatcher) SetIgnoreCase(ignore bool) {
	m.ignoreCase = ignore
}

// Match 判断文本中是否包含匹配模式的子串
func (m *RegexMatcher) Match(text string) bool {
	// 如果忽略大小写，将文本转为小写
	searchText := text
	if m.ignoreCase {
		searchText = strings.ToLower(text)
	}

	// 尝试从文本的每个位置开始匹配
	for i := 0; i <= len(searchText); i++ {
		if m.matchAt(searchText, i) {
			return true
		}
	}
	return false
}

// matchAt 从文本的指定位置开始尝试匹配模式
func (m *RegexMatcher) matchAt(text string, startPos int) bool {
	return m.matchTokens(m.pattern.tokens, text, 0, startPos)
}

// matchTokens 递归匹配 token 序列
func (m *RegexMatcher) matchTokens(tokens []Token, text string, patIdx int, textIdx int) bool {
	// 基础情况 1: 模式已全部匹配
	if patIdx >= len(tokens) {
		return true
	}

	// 基础情况 2: 文本耗尽但模式未完
	if textIdx >= len(text) {
		return false
	}

	// 获取当前要匹配的 token 和字符
	token := tokens[patIdx]
	char := text[textIdx]

	// 根据 token 类型进行匹配
	if m.matchChar(token, char) {
		return m.matchTokens(tokens, text, patIdx+1, textIdx+1)
	}

	return false
}

// matchChar 判断单个字符是否匹配 token
func (m *RegexMatcher) matchChar(token Token, char byte) bool {
	switch token.Type {
	case TokenWildcard:
		// 通配符 . 匹配任意单个字符
		return true

	case TokenLiteral:
		// 字面字符匹配
		patternChar := token.Value
		if m.ignoreCase {
			patternChar = toLowerByte(patternChar)
		}
		return char == patternChar

	case TokenDigit:
		// \d: 匹配数字 [0-9]
		return isDigit(char)

	case TokenNotDigit:
		// \D: 匹配非数字
		return !isDigit(char)

	case TokenWord:
		// \w: 匹配单词字符 [a-zA-Z0-9_]
		return isWordChar(char)

	case TokenNotWord:
		// \W: 匹配非单词字符
		return !isWordChar(char)

	case TokenSpace:
		// \s: 匹配空白字符
		return isSpace(char)

	case TokenNotSpace:
		// \S: 匹配非空白字符
		return !isSpace(char)

	default:
		return false
	}
}

// Pattern 返回匹配器使用的模式
func (m *RegexMatcher) Pattern() string {
	return m.pattern.Raw()
}

// isDigit 判断字符是否为数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isWordChar 判断字符是否为单词字符
// 单词字符：a-z, A-Z, 0-9, _
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		(c == '_')
}

// isSpace 判断字符是否为空白字符
// 空白字符：空格、制表符、换行符、回车符
func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

// toLowerByte 将单个字节转为小写（仅处理 ASCII）
func toLowerByte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
