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

	// 获取当前要匹配的 token
	token := tokens[patIdx]

	// 根据 token 类型进行匹配
	switch token.Type {
	case TokenWildcard:
		// 通配符 . 匹配任意单个字符
		return m.matchTokens(tokens, text, patIdx+1, textIdx+1)

	case TokenLiteral:
		// 字面字符匹配
		textChar := text[textIdx]
		patternChar := token.Value

		// 如果忽略大小写，将模式字符也转为小写
		if m.ignoreCase {
			patternChar = toLowerByte(patternChar)
		}

		if textChar == patternChar {
			return m.matchTokens(tokens, text, patIdx+1, textIdx+1)
		}
		return false

	default:
		return false
	}
}

// Pattern 返回匹配器使用的模式
func (m *RegexMatcher) Pattern() string {
	return m.pattern.Raw()
}

// toLowerByte 将单个字节转为小写（仅处理 ASCII）
func toLowerByte(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
