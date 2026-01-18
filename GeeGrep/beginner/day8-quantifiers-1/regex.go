package main

import (
	"fmt"
)

// TokenType 定义 token 的类型
type TokenType int

const (
	TokenLiteral         TokenType = iota // 字面字符
	TokenWildcard                         // 通配符 .
	TokenDigit                            // \d 数字 [0-9]
	TokenWord                             // \w 单词字符 [a-zA-Z0-9_]
	TokenSpace                            // \s 空白字符
	TokenNotDigit                         // \D 非数字
	TokenNotWord                          // \W 非单词字符
	TokenNotSpace                         // \S 非空白
	TokenCharGroup                        // 字符组 [abc] 或 [a-z]
	TokenStartAnchor                      // ^ 行首锚点
	TokenEndAnchor                        // $ 行尾锚点
	TokenWordBoundary                     // \b 单词边界
	TokenNotWordBoundary                  // \B 非单词边界
)

// Quantifier 表示量词信息
type Quantifier struct {
	Min    int  // 最小匹配次数
	Max    int  // 最大匹配次数（-1 表示无限）
	Greedy bool // 是否贪婪（默认 true）
}

// Range 表示字符范围（如 a-z）
type Range struct {
	Start byte
	End   byte
}

// CharGroup 表示字符组的内容
type CharGroup struct {
	Negated     bool        // 是否为否定字符组 [^...]
	Chars       []byte      // 单个字符列表
	Ranges      []Range     // 字符范围列表
	CharClasses []TokenType // 字符类列表（\d, \w, \s 等）
}

// Contains 判断字符是否在字符组中
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

	// 检查字符类
	for _, class := range cg.CharClasses {
		if matchCharClass(class, char) {
			return true
		}
	}

	return false
}

// matchCharClass 判断字符是否匹配指定的字符类
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

// Match 判断字符是否匹配字符组（考虑否定）
func (cg *CharGroup) Match(char byte) bool {
	contains := cg.Contains(char)
	if cg.Negated {
		return !contains
	}
	return contains
}

// Token 表示正则表达式中的一个元素
type Token struct {
	Type       TokenType
	Value      byte       // 对于字面字符，存储其值
	CharGroup  *CharGroup // 对于字符组，存储其内容
	Quantifier *Quantifier // 量词信息（如果有）
}

// String 返回 Token 的字符串表示（用于调试）
func (t Token) String() string {
	quantStr := ""
	if t.Quantifier != nil {
		if t.Quantifier.Max == -1 {
			quantStr = fmt.Sprintf("{%d,∞}", t.Quantifier.Min)
		} else {
			quantStr = fmt.Sprintf("{%d,%d}", t.Quantifier.Min, t.Quantifier.Max)
		}
	}

	switch t.Type {
	case TokenLiteral:
		return fmt.Sprintf("Literal('%c')%s", t.Value, quantStr)
	case TokenWildcard:
		return fmt.Sprintf("Wildcard(.)%s", quantStr)
	case TokenDigit:
		return fmt.Sprintf("Digit(\\d)%s", quantStr)
	case TokenWord:
		return fmt.Sprintf("Word(\\w)%s", quantStr)
	case TokenSpace:
		return fmt.Sprintf("Space(\\s)%s", quantStr)
	case TokenNotDigit:
		return fmt.Sprintf("NotDigit(\\D)%s", quantStr)
	case TokenNotWord:
		return fmt.Sprintf("NotWord(\\W)%s", quantStr)
	case TokenNotSpace:
		return fmt.Sprintf("NotSpace(\\S)%s", quantStr)
	case TokenCharGroup:
		if t.CharGroup.Negated {
			return fmt.Sprintf("NegCharGroup(...)%s", quantStr)
		}
		return fmt.Sprintf("CharGroup(...)%s", quantStr)
	case TokenStartAnchor:
		return "StartAnchor(^)"
	case TokenEndAnchor:
		return "EndAnchor($)"
	case TokenWordBoundary:
		return "WordBoundary(\\b)"
	case TokenNotWordBoundary:
		return "NotWordBoundary(\\B)"
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

// StartsWithAnchor 检查模式是否以行首锚点开始
func (p *Pattern) StartsWithAnchor() bool {
	return len(p.tokens) > 0 && p.tokens[0].Type == TokenStartAnchor
}

// EndsWithAnchor 检查模式是否以行尾锚点结束
func (p *Pattern) EndsWithAnchor() bool {
	return len(p.tokens) > 0 && p.tokens[len(p.tokens)-1].Type == TokenEndAnchor
}

// ParsePattern 解析正则表达式字符串为 Pattern
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

		case '^':
			// 行首锚点
			tokens = append(tokens, Token{
				Type: TokenStartAnchor,
			})

		case '$':
			// 行尾锚点
			tokens = append(tokens, Token{
				Type: TokenEndAnchor,
			})

		case '?', '+', '*':
			// 量词 - 修饰前一个 token
			if len(tokens) == 0 {
				return nil, fmt.Errorf("quantifier '%c' at position %d has nothing to repeat", ch, i)
			}

			// 检查前一个 token 是否已有量词
			lastToken := &tokens[len(tokens)-1]
			if lastToken.Quantifier != nil {
				return nil, fmt.Errorf("multiple quantifiers at position %d", i)
			}

			// 检查前一个 token 是否为锚点（锚点不能有量词）
			if lastToken.Type == TokenStartAnchor || lastToken.Type == TokenEndAnchor ||
				lastToken.Type == TokenWordBoundary || lastToken.Type == TokenNotWordBoundary {
				return nil, fmt.Errorf("quantifier '%c' cannot follow anchor at position %d", ch, i)
			}

			// 设置量词
			switch ch {
			case '?':
				lastToken.Quantifier = &Quantifier{Min: 0, Max: 1, Greedy: true}
			case '+':
				lastToken.Quantifier = &Quantifier{Min: 1, Max: -1, Greedy: true}
			case '*':
				lastToken.Quantifier = &Quantifier{Min: 0, Max: -1, Greedy: true}
			}

		case '\\':
			// 转义序列
			token, newPos, err := parseEscape(pattern, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i = newPos - 1 // -1 因为 for 循环会 +1

		case '[':
			// 字符组
			token, newPos, err := parseCharGroup(pattern, i)
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
func parseEscape(pattern string, pos int) (Token, int, error) {
	if pos+1 >= len(pattern) {
		return Token{}, pos, fmt.Errorf("incomplete escape sequence at position %d", pos)
	}

	next := pattern[pos+1]

	switch next {
	case 'd':
		return Token{Type: TokenDigit}, pos + 2, nil
	case 'D':
		return Token{Type: TokenNotDigit}, pos + 2, nil
	case 'w':
		return Token{Type: TokenWord}, pos + 2, nil
	case 'W':
		return Token{Type: TokenNotWord}, pos + 2, nil
	case 's':
		return Token{Type: TokenSpace}, pos + 2, nil
	case 'S':
		return Token{Type: TokenNotSpace}, pos + 2, nil
	case 'b':
		return Token{Type: TokenWordBoundary}, pos + 2, nil
	case 'B':
		return Token{Type: TokenNotWordBoundary}, pos + 2, nil
	case '.', '*', '+', '?', '\\', '(', ')', '[', ']', '{', '}', '|', '^', '$':
		// 字面转义
		return Token{Type: TokenLiteral, Value: next}, pos + 2, nil
	default:
		return Token{}, pos, fmt.Errorf("unknown escape sequence: \\%c at position %d", next, pos)
	}
}

// parseCharGroup 解析字符组 [...]
func parseCharGroup(pattern string, pos int) (Token, int, error) {
	startPos := pos
	// 跳过 '['
	pos++

	if pos >= len(pattern) {
		return Token{}, startPos, fmt.Errorf("unclosed character group at position %d", startPos)
	}

	// 检查是否为否定字符组
	negated := false
	if pattern[pos] == '^' {
		negated = true
		pos++
	}

	// 解析字符、范围和字符类
	chars := []byte{}
	ranges := []Range{}
	charClasses := []TokenType{}

	// 特殊处理：] 在开头是字面字符
	if pos < len(pattern) && pattern[pos] == ']' {
		chars = append(chars, ']')
		pos++
	}

	for pos < len(pattern) && pattern[pos] != ']' {
		ch := pattern[pos]

		// 检查是否为转义字符
		if ch == '\\' && pos+1 < len(pattern) {
			next := pattern[pos+1]
			switch next {
			case 'd':
				charClasses = append(charClasses, TokenDigit)
				pos += 2
			case 'D':
				charClasses = append(charClasses, TokenNotDigit)
				pos += 2
			case 'w':
				charClasses = append(charClasses, TokenWord)
				pos += 2
			case 'W':
				charClasses = append(charClasses, TokenNotWord)
				pos += 2
			case 's':
				charClasses = append(charClasses, TokenSpace)
				pos += 2
			case 'S':
				charClasses = append(charClasses, TokenNotSpace)
				pos += 2
			case ']', '\\', '-', '^':
				// 转义的特殊字符
				chars = append(chars, next)
				pos += 2
			default:
				// 其他转义当作字面字符
				chars = append(chars, next)
				pos += 2
			}
			continue
		}

		// 检查是否为范围 (a-z)
		if pos+2 < len(pattern) &&
			pattern[pos+1] == '-' &&
			pattern[pos+2] != ']' {
			// 是范围
			start := ch
			end := pattern[pos+2]
			if start <= end {
				ranges = append(ranges, Range{Start: start, End: end})
			} else {
				// 无效范围，当作三个独立字符
				chars = append(chars, start, '-', end)
			}
			pos += 3
		} else {
			// 单个字符
			chars = append(chars, ch)
			pos++
		}
	}

	// 检查是否找到闭合的 ']'
	if pos >= len(pattern) {
		return Token{}, startPos, fmt.Errorf("unclosed character group at position %d", startPos)
	}

	// 跳过 ']'
	pos++

	return Token{
		Type: TokenCharGroup,
		CharGroup: &CharGroup{
			Negated:     negated,
			Chars:       chars,
			Ranges:      ranges,
			CharClasses: charClasses,
		},
	}, pos, nil
}

// Compile 是 ParsePattern 的别名
func Compile(pattern string) (*Pattern, error) {
	return ParsePattern(pattern)
}

// 辅助函数
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		(c == '_')
}

func isSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}
