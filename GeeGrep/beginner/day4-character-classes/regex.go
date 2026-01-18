package main

import (
	"fmt"
)

// TokenType 定义 token 的类型
type TokenType int

const (
	TokenLiteral  TokenType = iota // 字面字符
	TokenWildcard                   // 通配符 .
	TokenDigit                      // \d 数字 [0-9]
	TokenWord                       // \w 单词字符 [a-zA-Z0-9_]
	TokenSpace                      // \s 空白字符
	TokenNotDigit                   // \D 非数字
	TokenNotWord                    // \W 非单词字符
	TokenNotSpace                   // \S 非空白
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
	case TokenDigit:
		return "Digit(\\d)"
	case TokenWord:
		return "Word(\\w)"
	case TokenSpace:
		return "Space(\\s)"
	case TokenNotDigit:
		return "NotDigit(\\D)"
	case TokenNotWord:
		return "NotWord(\\W)"
	case TokenNotSpace:
		return "NotSpace(\\S)"
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
// 支持：
// - 字面字符：直接匹配
// - 通配符 .：匹配任意单个字符
// - \d：匹配数字 [0-9]
// - \w：匹配单词字符 [a-zA-Z0-9_]
// - \s：匹配空白字符
// - \D, \W, \S：对应的否定形式
// - 字面转义：\., \\, \*, \+, \? 等
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

		case '\\':
			// 转义序列
			token, newPos, err := parseEscape(pattern, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i = newPos - 1 // -1 因为 for 循环会 +1

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

// parseEscape 解析转义序列
// 返回：token, 新位置, 错误
func parseEscape(pattern string, pos int) (Token, int, error) {
	if pos+1 >= len(pattern) {
		return Token{}, pos, fmt.Errorf("incomplete escape sequence at position %d", pos)
	}

	next := pattern[pos+1]

	switch next {
	case 'd':
		// \d: 数字
		return Token{Type: TokenDigit}, pos + 2, nil

	case 'D':
		// \D: 非数字
		return Token{Type: TokenNotDigit}, pos + 2, nil

	case 'w':
		// \w: 单词字符
		return Token{Type: TokenWord}, pos + 2, nil

	case 'W':
		// \W: 非单词字符
		return Token{Type: TokenNotWord}, pos + 2, nil

	case 's':
		// \s: 空白字符
		return Token{Type: TokenSpace}, pos + 2, nil

	case 'S':
		// \S: 非空白字符
		return Token{Type: TokenNotSpace}, pos + 2, nil

	case '.', '*', '+', '?', '\\', '(', ')', '[', ']', '{', '}', '|', '^', '$':
		// 字面转义：匹配这些特殊字符的字面值
		return Token{Type: TokenLiteral, Value: next}, pos + 2, nil

	default:
		return Token{}, pos, fmt.Errorf("unknown escape sequence: \\%c at position %d", next, pos)
	}
}

// Compile 是 ParsePattern 的别名，提供更符合习惯的 API
func Compile(pattern string) (*Pattern, error) {
	return ParsePattern(pattern)
}
