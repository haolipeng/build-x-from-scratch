package main

import (
	"fmt"
)

// TokenType 定义 token 的类型
type TokenType int

const (
	TokenLiteral  TokenType = iota // 字面字符
	TokenWildcard                   // 通配符 .
)

// Token 表示正则表达式中的一个元素
type Token struct {
	Type  TokenType
	Value byte // 对于字面字符，存储其值
}

// String 返回 Token 的字符串表示（用于调试）
func (t Token) String() string {
	switch t.Type {
	case TokenLiteral:
		return fmt.Sprintf("Literal('%c')", t.Value)
	case TokenWildcard:
		return "Wildcard(.)"
	default:
		return "Unknown"
	}
}

// Pattern 表示一个已解析的正则表达式模式
type Pattern struct {
	raw    string  // 原始模式字符串
	tokens []Token // 解析后的 token 序列
}

// NewPattern 创建一个新的 Pattern
func NewPattern(raw string, tokens []Token) *Pattern {
	return &Pattern{
		raw:    raw,
		tokens: tokens,
	}
}

// Raw 返回原始模式字符串
func (p *Pattern) Raw() string {
	return p.raw
}

// Tokens 返回 token 序列
func (p *Pattern) Tokens() []Token {
	return p.tokens
}

// String 返回 Pattern 的字符串表示
func (p *Pattern) String() string {
	return fmt.Sprintf("Pattern{raw=%q, tokens=%v}", p.raw, p.tokens)
}

// ParsePattern 解析正则表达式字符串为 Pattern
// 目前支持：
// - 字面字符：直接匹配
// - 通配符 .：匹配任意单个字符
func ParsePattern(pattern string) (*Pattern, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	tokens := make([]Token, 0, len(pattern))

	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]

		switch ch {
		case '.':
			// 通配符
			tokens = append(tokens, Token{
				Type: TokenWildcard,
			})
		default:
			// 字面字符
			tokens = append(tokens, Token{
				Type:  TokenLiteral,
				Value: ch,
			})
		}
	}

	return NewPattern(pattern, tokens), nil
}

// Compile 是 ParsePattern 的别名，提供更符合习惯的 API
func Compile(pattern string) (*Pattern, error) {
	return ParsePattern(pattern)
}
